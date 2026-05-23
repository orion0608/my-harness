package lifecycle

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInstanceKeyFromWorkDirStable(t *testing.T) {
	dir := t.TempDir()
	key1, err := InstanceKeyFromWorkDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	key2, err := InstanceKeyFromWorkDir(filepath.Join(dir, "."))
	if err != nil {
		t.Fatal(err)
	}
	if key1 != key2 {
		t.Fatalf("expected stable key, got %q vs %q", key1, key2)
	}
	if len(key1) != 12 {
		t.Fatalf("expected 12-char hex key, got %q", key1)
	}
}

func TestRegistryRoundTrip(t *testing.T) {
	root := t.TempDir()
	m := &Manager{
		appName:        "testapp",
		branchName:     "main",
		instanceKey:    "abc123",
		workDir:        "/tmp/work",
		registryRoot:   root,
		pid:            4242,
		host:           "127.0.0.1",
		port:           50123,
		version:        "1.0.0",
		startedAt:      mustParseTime(t, "2026-05-24T10:00:00Z"),
		executablePath: "/bin/app",
	}

	if err := m.appendSelfRecord(); err != nil {
		t.Fatal(err)
	}

	reg, err := loadRegistry(m.registryPath())
	if err != nil {
		t.Fatal(err)
	}
	if reg == nil || len(reg.Instances) != 1 {
		t.Fatalf("expected one instance, got %+v", reg)
	}
	if reg.Instances[0].PID != 4242 {
		t.Fatalf("unexpected pid: %+v", reg.Instances[0])
	}

	if err := m.removeSelfRecord(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(m.registryPath()); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected registry file removed, stat err=%v", err)
	}
}

func TestRemoveSelfRecordByPIDKeepsOthers(t *testing.T) {
	root := t.TempDir()
	path := registryFilePath(root, "testapp", "main", "abc123")
	reg := &RegistryFile{
		AppName:     "testapp",
		BranchName:  "main",
		InstanceKey: "abc123",
		WorkDir:     "/tmp/work",
		Instances: []InstanceRecord{
			{PID: 100, Host: "127.0.0.1", Port: 50001},
			{PID: 200, Host: "127.0.0.1", Port: 50002},
		},
	}
	if err := saveRegistry(path, reg); err != nil {
		t.Fatal(err)
	}

	m := &Manager{
		appName:      "testapp",
		branchName:   "main",
		instanceKey:  "abc123",
		workDir:      "/tmp/work",
		registryRoot: root,
		pid:          100,
	}

	if err := m.removeSelfRecord(); err != nil {
		t.Fatal(err)
	}

	updated, err := loadRegistry(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(updated.Instances) != 1 {
		t.Fatalf("expected one remaining instance, got %+v", updated.Instances)
	}
	if updated.Instances[0].PID != 200 {
		t.Fatalf("expected pid 200 to remain, got %+v", updated.Instances[0])
	}
}

func mustParseTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatal(err)
	}
	return parsed
}
