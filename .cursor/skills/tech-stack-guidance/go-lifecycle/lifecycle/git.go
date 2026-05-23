package lifecycle

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func currentWorkingDirectory() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	abs, err := filepath.Abs(wd)
	if err != nil {
		return "", err
	}
	return filepath.Clean(abs), nil
}

func detectAppName(workDir string) (string, error) {
	if url, err := gitOutput(workDir, "remote", "get-url", "origin"); err == nil {
		if name := repoNameFromRemoteURL(strings.TrimSpace(url)); name != "" {
			return name, nil
		}
	}
	if top, err := gitOutput(workDir, "rev-parse", "--show-toplevel"); err == nil {
		if name := filepath.Base(filepath.Clean(top)); name != "" && name != "." {
			return name, nil
		}
	}
	name := filepath.Base(workDir)
	if name == "" || name == "." {
		return "", fmt.Errorf("cannot derive app name from %q", workDir)
	}
	return name, nil
}

func detectBranchName(workDir string) string {
	name, err := gitOutput(workDir, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "unknown"
	}
	name = strings.TrimSpace(name)
	if name == "" || name == "HEAD" {
		if short, err := gitOutput(workDir, "rev-parse", "--short", "HEAD"); err == nil {
			short = strings.TrimSpace(short)
			if short != "" {
				return "detached-" + short
			}
		}
		return "unknown"
	}
	return sanitizeBranchFileName(name)
}

// sanitizeBranchFileName makes a git branch name safe for use in a filename.
func sanitizeBranchFileName(name string) string {
	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "-",
		"?", "-",
		"\"", "-",
		"<", "-",
		">", "-",
		"|", "-",
	)
	return replacer.Replace(name)
}

func detectVersion(workDir string, startedAt time.Time) string {
	shortID, err := gitOutput(workDir, "rev-parse", "--short", "HEAD")
	if err != nil || strings.TrimSpace(shortID) == "" {
		return fmt.Sprintf("unknown@%s", startedAt.Format(time.RFC3339))
	}
	return fmt.Sprintf("%s@%s", strings.TrimSpace(shortID), startedAt.Format(time.RFC3339))
}

func gitOutput(workDir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = workDir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func repoNameFromRemoteURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	raw = strings.TrimSuffix(raw, ".git")
	raw = strings.TrimSuffix(raw, "/")
	if i := strings.LastIndexAny(raw, "/:"); i >= 0 && i+1 < len(raw) {
		return raw[i+1:]
	}
	return raw
}
