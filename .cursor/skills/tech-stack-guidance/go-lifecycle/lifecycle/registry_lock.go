package lifecycle

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func registryLockPath(registryPath string) string {
	return registryPath + ".lock"
}

// withFileLock acquires an OS file lock on lockPath, blocking until available.
func withFileLock(lockPath string, fn func() error) error {
	if err := os.MkdirAll(filepath.Dir(lockPath), 0o755); err != nil {
		return err
	}

	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return fmt.Errorf("open lock file: %w", err)
	}
	defer f.Close()

	if err := acquireFileLock(f); err != nil {
		return fmt.Errorf("acquire file lock: %w", err)
	}
	defer func() { _ = releaseFileLock(f) }()

	return fn()
}

func (m *Manager) withRegistryLock(fn func() error) error {
	return withFileLock(registryLockPath(m.registryPath()), fn)
}

func saveRegistryAtomic(path string, reg *RegistryFile) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(reg, "", "  ")
	if err != nil {
		return err
	}

	tmpPath := path + ".tmp." + strconv.Itoa(os.Getpid())
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return nil
}

func (m *Manager) mutateRegistry(update func(reg *RegistryFile) (*RegistryFile, error)) error {
	path := m.registryPath()
	return m.withRegistryLock(func() error {
		reg, err := loadRegistry(path)
		if err != nil {
			return err
		}

		newReg, err := update(reg)
		if err != nil {
			return err
		}
		if newReg == nil || len(newReg.Instances) == 0 {
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return err
			}
			return nil
		}
		return saveRegistryAtomic(path, newReg)
	})
}
