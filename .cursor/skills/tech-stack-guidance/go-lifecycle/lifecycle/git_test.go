package lifecycle

import (
	"testing"
	"time"
)

func TestRegistryFileName(t *testing.T) {
	got := registryFileName("rp-admin-header", "a1b2c3d4e5f6")
	want := "rp-admin-header+a1b2c3d4e5f6.json"
	if got != want {
		t.Fatalf("registryFileName = %q, want %q", got, want)
	}
}

func TestSanitizeBranchFileName(t *testing.T) {
	got := sanitizeBranchFileName("rp/admin-header")
	want := "rp-admin-header"
	if got != want {
		t.Fatalf("sanitizeBranchFileName = %q, want %q", got, want)
	}
}

func TestRepoNameFromRemoteURL(t *testing.T) {
	cases := map[string]string{
		"https://github.com/org/my_harness.git": "my_harness",
		"git@github.com:org/my_harness.git":     "my_harness",
		"https://github.com/org/my_harness":     "my_harness",
		"/bare/repo.git":                        "repo",
	}
	for input, want := range cases {
		if got := repoNameFromRemoteURL(input); got != want {
			t.Fatalf("repoNameFromRemoteURL(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestDetectVersion(t *testing.T) {
	startedAt := time.Date(2026, 5, 24, 10, 0, 0, 0, time.UTC)
	got := detectVersion(t.TempDir(), startedAt)
	want := "unknown@2026-05-24T10:00:00Z"
	if got != want {
		t.Fatalf("detectVersion without git = %q, want %q", got, want)
	}
}

func TestDetectAppNameFallback(t *testing.T) {
	dir := t.TempDir()
	got, err := detectAppName(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got == "" {
		t.Fatal("expected non-empty app name fallback")
	}
}
