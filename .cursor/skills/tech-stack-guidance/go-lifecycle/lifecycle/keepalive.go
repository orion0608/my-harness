package lifecycle

import (
	"context"
	"time"
)

const (
	// KeepaliveInterval is how often a running instance refreshes lastKeepalive in the registry.
	KeepaliveInterval = time.Minute
	// KeepaliveStaleAfter marks instances whose lastKeepalive is older than this as likely abnormally exited.
	KeepaliveStaleAfter = 2 * KeepaliveInterval
)

// IsKeepaliveStale reports whether lastKeepalive indicates the instance probably exited abnormally.
func IsKeepaliveStale(lastKeepalive time.Time, now time.Time) bool {
	if lastKeepalive.IsZero() {
		return true
	}
	return now.Sub(lastKeepalive) > KeepaliveStaleAfter
}

func (m *Manager) nowUTC() time.Time {
	return time.Now().UTC()
}

func (m *Manager) touchSelfKeepalive() error {
	now := m.nowUTC()
	return m.mutateRegistry(func(reg *RegistryFile) (*RegistryFile, error) {
		if reg == nil {
			return nil, nil
		}

		found := false
		for i := range reg.Instances {
			if reg.Instances[i].PID != m.pid {
				continue
			}
			reg.Instances[i].LastKeepalive = now
			found = true
		}
		if !found {
			return reg, nil
		}

		m.mu.Lock()
		m.lastKeepalive = now
		m.mu.Unlock()
		return reg, nil
	})
}

func (m *Manager) startKeepalive(ctx context.Context) {
	keepaliveCtx, cancel := context.WithCancel(ctx)
	m.keepaliveMu.Lock()
	if m.keepaliveCancel != nil {
		m.keepaliveCancel()
	}
	m.keepaliveCancel = cancel
	m.keepaliveMu.Unlock()

	go m.runKeepalive(keepaliveCtx)
}

func (m *Manager) stopKeepalive() {
	m.keepaliveMu.Lock()
	if m.keepaliveCancel != nil {
		m.keepaliveCancel()
		m.keepaliveCancel = nil
	}
	m.keepaliveMu.Unlock()
}

func (m *Manager) runKeepalive(ctx context.Context) {
	ticker := time.NewTicker(KeepaliveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = m.touchSelfKeepalive()
		}
	}
}
