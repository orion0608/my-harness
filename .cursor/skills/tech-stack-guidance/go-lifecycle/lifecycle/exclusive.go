package lifecycle

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	shutdownPath       = "/__service/shutdown"
	shutdownTimeout    = 5 * time.Second
	shutdownWaitPeriod = 100 * time.Millisecond
	shutdownWaitMax    = 8 * time.Second
)

// EnsureExclusive stops every registered instance for the same app + worktree key.
// Order per instance: POST /__service/shutdown, wait for exit, then kill if needed.
func (m *Manager) EnsureExclusive(ctx context.Context) error {
	instances, err := m.loadRegistryInstances()
	if err != nil {
		return fmt.Errorf("load registry: %w", err)
	}
	if len(instances) == 0 {
		return nil
	}

	remaining := make([]InstanceRecord, 0, len(instances))
	for _, inst := range instances {
		if inst.PID == m.pid {
			continue
		}
		if !processAlive(inst.PID) {
			continue
		}

		if err := requestShutdown(ctx, inst); err != nil {
			// Shutdown API failed; fall through to kill.
		} else if waitProcessExit(inst.PID, shutdownWaitMax) {
			continue
		}

		if err := killProcess(inst.PID); err != nil {
			remaining = append(remaining, inst)
			continue
		}
		if !waitProcessExit(inst.PID, shutdownWaitMax) {
			remaining = append(remaining, inst)
		}
	}

	if err := m.replaceRegistryInstances(remaining); err != nil {
		return fmt.Errorf("update registry after exclusive start: %w", err)
	}
	return nil
}

func requestShutdown(ctx context.Context, inst InstanceRecord) error {
	url := fmt.Sprintf("http://%s:%d%s", inst.Host, inst.Port, shutdownPath)
	reqCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("shutdown returned status %d", resp.StatusCode)
	}
	return nil
}

func waitProcessExit(pid int, maxWait time.Duration) bool {
	deadline := time.Now().Add(maxWait)
	for time.Now().Before(deadline) {
		if !processAlive(pid) {
			return true
		}
		time.Sleep(shutdownWaitPeriod)
	}
	return !processAlive(pid)
}
