package lifecycle

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"runtime"
)

const defaultRegistryDir = ".harness-services"

// DefaultRegistryRoot returns the shared registry directory for tool-level services.
// Windows: C:\Users\<user>\.harness-services
// Unix: ~/.harness-services
func DefaultRegistryRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, defaultRegistryDir), nil
}

// InstanceKeyFromWorkDir derives a stable key from an absolute worktree/startup path.
func InstanceKeyFromWorkDir(workDir string) (string, error) {
	abs, err := filepath.Abs(workDir)
	if err != nil {
		return "", err
	}
	clean := filepath.Clean(abs)
	sum := sha256.Sum256([]byte(clean))
	return hex.EncodeToString(sum[:6]), nil
}

func registryFileName(branchName, instanceKey string) string {
	return branchName + "+" + instanceKey + ".json"
}

func registryFilePath(registryRoot, appName, branchName, instanceKey string) string {
	return filepath.Join(registryRoot, appName, registryFileName(branchName, instanceKey))
}

func executablePath() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	resolved, err := filepath.EvalSymlinks(exe)
	if err != nil {
		return exe
	}
	return resolved
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}
