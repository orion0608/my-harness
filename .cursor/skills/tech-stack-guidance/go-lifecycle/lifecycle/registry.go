package lifecycle

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// InstanceRecord describes one running service instance.
type InstanceRecord struct {
	PID              int       `json:"pid"`
	Host             string    `json:"host"`
	Port             int       `json:"port"`
	StartedAt        time.Time `json:"startedAt"`
	LastKeepalive    time.Time `json:"lastKeepalive"`
	Version          string    `json:"version"`
	ExecutablePath   string    `json:"executablePath"`
	WorkingDirectory string    `json:"workingDirectory"`
}

// RegistryFile tracks all known instances for one app + worktree key.
type RegistryFile struct {
	AppName     string           `json:"appName"`
	BranchName  string           `json:"branchName"`
	InstanceKey string           `json:"instanceKey"`
	WorkDir     string           `json:"workDir"`
	Instances   []InstanceRecord `json:"instances"`
}

func (m *Manager) registryPath() string {
	return registryFilePath(m.registryRoot, m.appName, m.branchName, m.instanceKey)
}

func loadRegistry(path string) (*RegistryFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var reg RegistryFile
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, err
	}
	return &reg, nil
}

func (m *Manager) newInstanceRecord() InstanceRecord {
	now := m.nowUTC()
	m.mu.Lock()
	m.lastKeepalive = now
	m.mu.Unlock()
	return InstanceRecord{
		PID:              m.pid,
		Host:             m.host,
		Port:             m.port,
		StartedAt:        m.startedAt,
		LastKeepalive:    now,
		Version:          m.version,
		ExecutablePath:   m.executablePath,
		WorkingDirectory: m.workDir,
	}
}

func (m *Manager) newRegistryFile() *RegistryFile {
	return &RegistryFile{
		AppName:     m.appName,
		BranchName:  m.branchName,
		InstanceKey: m.instanceKey,
		WorkDir:     m.workDir,
		Instances:   nil,
	}
}

func (m *Manager) appendSelfRecord() error {
	record := m.newInstanceRecord()
	return m.mutateRegistry(func(reg *RegistryFile) (*RegistryFile, error) {
		if reg == nil {
			reg = m.newRegistryFile()
		}
		reg.Instances = append(reg.Instances, record)
		return reg, nil
	})
}

// removeSelfRecord removes the instance entry whose pid matches this process.
func (m *Manager) removeSelfRecord() error {
	return m.mutateRegistry(func(reg *RegistryFile) (*RegistryFile, error) {
		if reg == nil {
			return nil, nil
		}

		filtered := reg.Instances[:0]
		for _, inst := range reg.Instances {
			if inst.PID != m.pid {
				filtered = append(filtered, inst)
			}
		}
		reg.Instances = filtered
		if len(reg.Instances) == 0 {
			return nil, nil
		}
		return reg, nil
	})
}

func (m *Manager) replaceRegistryInstances(instances []InstanceRecord) error {
	if len(instances) == 0 {
		return m.mutateRegistry(func(reg *RegistryFile) (*RegistryFile, error) {
			return nil, nil
		})
	}
	reg := m.newRegistryFile()
	reg.Instances = instances
	return m.mutateRegistry(func(_ *RegistryFile) (*RegistryFile, error) {
		return reg, nil
	})
}

func (m *Manager) loadRegistryInstances() ([]InstanceRecord, error) {
	path := m.registryPath()
	var instances []InstanceRecord
	err := m.withRegistryLock(func() error {
		reg, err := loadRegistry(path)
		if err != nil {
			return err
		}
		if reg != nil {
			instances = append(instances, reg.Instances...)
		}
		return nil
	})
	return instances, err
}

// saveRegistry writes registry data for tests and setup helpers.
func saveRegistry(path string, reg *RegistryFile) error {
	return saveRegistryAtomic(path, reg)
}
