package lifecycle

import (
	"testing"
	"time"
)

func TestIsKeepaliveStale(t *testing.T) {
	now := time.Date(2026, 5, 24, 12, 0, 0, 0, time.UTC)

	if !IsKeepaliveStale(time.Time{}, now) {
		t.Fatal("zero time should be stale")
	}
	if IsKeepaliveStale(now.Add(-time.Minute), now) {
		t.Fatal("1 minute ago should not be stale")
	}
	if !IsKeepaliveStale(now.Add(-3*time.Minute), now) {
		t.Fatal("3 minutes ago should be stale")
	}
}

func TestTouchSelfKeepalive(t *testing.T) {
	root := t.TempDir()
	startedAt := time.Date(2026, 5, 24, 10, 0, 0, 0, time.UTC)
	m := &Manager{
		appName:      "testapp",
		branchName:   "main",
		instanceKey:  "abc123",
		workDir:      "/tmp/work",
		registryRoot: root,
		pid:          4242,
		host:         "127.0.0.1",
		port:         50123,
		version:      "1.0.0",
		startedAt:    startedAt,
		lastKeepalive: startedAt,
	}

	if err := m.appendSelfRecord(); err != nil {
		t.Fatal(err)
	}

	reg, err := loadRegistry(m.registryPath())
	if err != nil {
		t.Fatal(err)
	}
	before := reg.Instances[0].LastKeepalive
	if err := m.touchSelfKeepalive(); err != nil {
		t.Fatal(err)
	}

	reg, err = loadRegistry(m.registryPath())
	if err != nil {
		t.Fatal(err)
	}
	if !reg.Instances[0].LastKeepalive.After(before) {
		t.Fatalf("expected lastKeepalive to advance from %v, got %v", before, reg.Instances[0].LastKeepalive)
	}
}
