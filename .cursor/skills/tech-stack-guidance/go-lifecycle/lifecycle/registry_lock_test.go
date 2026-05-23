package lifecycle

import (
	"path/filepath"
	"testing"
	"time"
)

func TestWithFileLockBlocksUntilReleased(t *testing.T) {
	lockPath := filepath.Join(t.TempDir(), "test.lock")
	started := make(chan struct{})
	release := make(chan struct{})

	go func() {
		err := withFileLock(lockPath, func() error {
			close(started)
			<-release
			return nil
		})
		if err != nil {
			t.Errorf("holder lock: %v", err)
		}
	}()

	<-started

	done := make(chan error, 1)
	go func() {
		done <- withFileLock(lockPath, func() error { return nil })
	}()

	select {
	case err := <-done:
		t.Fatalf("expected waiter to block, got err=%v", err)
	case <-time.After(200 * time.Millisecond):
	}

	close(release)

	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("waiter did not acquire lock after release")
	}
}

func TestMutateRegistryConcurrentAppend(t *testing.T) {
	root := t.TempDir()
	base := &Manager{
		appName:      "testapp",
		branchName:   "main",
		instanceKey:  "abc123",
		workDir:      "/tmp/work",
		registryRoot: root,
		host:         "127.0.0.1",
		version:      "1.0.0",
		startedAt:    time.Now().UTC(),
	}

	errCh := make(chan error, 2)
	for i, pid := range []int{100, 200} {
		m := *base
		m.pid = pid
		m.port = 50000 + i
		go func(mgr Manager) {
			errCh <- mgr.appendSelfRecord()
		}(m)
	}

	for i := 0; i < 2; i++ {
		if err := <-errCh; err != nil {
			t.Fatal(err)
		}
	}

	reg, err := loadRegistry(base.registryPath())
	if err != nil {
		t.Fatal(err)
	}
	if reg == nil || len(reg.Instances) != 2 {
		t.Fatalf("expected two instances, got %+v", reg)
	}
}
