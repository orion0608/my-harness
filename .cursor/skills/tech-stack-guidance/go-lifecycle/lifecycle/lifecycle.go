package lifecycle

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	InfoPath     = "/__service/info"
	ShutdownPath = "/__service/shutdown"
	defaultHost  = "127.0.0.1"
)

// Config configures a tool-level service lifecycle manager.
// AppName, version, and workingDirectory are detected automatically (git + os.Getwd).
type Config struct {
	// RegistryRoot overrides ~/.harness-services.
	RegistryRoot string
	// PortMin/PortMax define the tool-level port range; default 50000-60000.
	PortMin int
	PortMax int
	// OnShutdown runs before the HTTP server exits (optional cleanup hook).
	OnShutdown func(context.Context) error
}

// InfoResponse is returned by GET /__service/info.
type InfoResponse struct {
	AppName          string    `json:"appName"`
	BranchName       string    `json:"branchName"`
	InstanceKey      string    `json:"instanceKey"`
	PID              int       `json:"pid"`
	Version          string    `json:"version"`
	StartedAt        time.Time `json:"startedAt"`
	ExecutablePath   string    `json:"executablePath"`
	WorkingDirectory string    `json:"workingDirectory"`
	Port             int       `json:"port"`
}

// Manager implements registry, exclusive startup, and standard HTTP endpoints.
type Manager struct {
	appName        string
	branchName     string
	version        string
	workDir        string
	instanceKey    string
	registryRoot   string
	host           string
	portMin        int
	portMax        int
	onShutdown     func(context.Context) error
	executablePath string
	pid            int
	port           int
	startedAt      time.Time

	mu     sync.Mutex
	server *http.Server
}

// New validates config and prepares a Manager. Call EnsureExclusive before ListenAndServe.
func New(cfg Config) (*Manager, error) {
	startedAt := time.Now().UTC()

	workDir, err := currentWorkingDirectory()
	if err != nil {
		return nil, fmt.Errorf("lifecycle: work dir: %w", err)
	}

	appName, err := detectAppName(workDir)
	if err != nil {
		return nil, fmt.Errorf("lifecycle: app name: %w", err)
	}

	instanceKey, err := InstanceKeyFromWorkDir(workDir)
	if err != nil {
		return nil, fmt.Errorf("lifecycle: instance key: %w", err)
	}

	registryRoot := cfg.RegistryRoot
	if registryRoot == "" {
		registryRoot, err = DefaultRegistryRoot()
		if err != nil {
			return nil, fmt.Errorf("lifecycle: registry root: %w", err)
		}
	}

	return &Manager{
		appName:        appName,
		branchName:     detectBranchName(workDir),
		version:        detectVersion(workDir, startedAt),
		workDir:        workDir,
		instanceKey:    instanceKey,
		registryRoot:   registryRoot,
		host:           defaultHost,
		portMin:        cfg.PortMin,
		portMax:        cfg.PortMax,
		onShutdown:     cfg.OnShutdown,
		executablePath: executablePath(),
		pid:            os.Getpid(),
		startedAt:      startedAt,
	}, nil
}

// BranchName returns the detected git branch used in the registry filename.
func (m *Manager) BranchName() string {
	return m.branchName
}

// InstanceKey returns the derived worktree key.
func (m *Manager) InstanceKey() string {
	return m.instanceKey
}

// RegistryPath returns the JSON registry file path for this app + branch + worktree.
func (m *Manager) RegistryPath() string {
	return m.registryPath()
}

// Info returns the current service metadata.
func (m *Manager) Info() InfoResponse {
	return InfoResponse{
		AppName:          m.appName,
		BranchName:       m.branchName,
		InstanceKey:      m.instanceKey,
		PID:              m.pid,
		Version:          m.version,
		StartedAt:        m.startedAt,
		ExecutablePath:   m.executablePath,
		WorkingDirectory: m.workDir,
		Port:             m.port,
	}
}

// Addr returns host:port after ListenAndServe picks a port.
func (m *Manager) Addr() string {
	if m.port == 0 {
		return ""
	}
	return fmt.Sprintf("%s:%d", m.host, m.port)
}

// Port returns the bound port, or 0 before ListenAndServe.
func (m *Manager) Port() int {
	return m.port
}

// ListenAndServe binds a random port in range, registers routes on mux, and blocks.
// Standard routes are always mounted at /__service/info and /__service/shutdown.
func (m *Manager) ListenAndServe(ctx context.Context, mux *http.ServeMux) error {
	if mux == nil {
		mux = http.NewServeMux()
	}
	m.mountRoutes(mux)

	port, err := pickPort(m.host, m.portMin, m.portMax)
	if err != nil {
		return err
	}
	m.port = port

	addr := fmt.Sprintf("%s:%d", m.host, m.port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	if err := m.appendSelfRecord(); err != nil {
		_ = ln.Close()
		return fmt.Errorf("register instance: %w", err)
	}
	defer func() { _ = m.removeSelfRecord() }()

	m.server = &http.Server{Handler: mux}
	errCh := make(chan error, 1)
	go func() {
		errCh <- m.server.Serve(ln)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = m.gracefulStop(shutdownCtx)
		<-errCh
		return ctx.Err()
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

func (m *Manager) mountRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET "+InfoPath, m.handleInfo)
	mux.HandleFunc("POST "+ShutdownPath, m.handleShutdown)
}

// gracefulStop unregisters this process (matched by pid) then shuts down the server.
func (m *Manager) gracefulStop(ctx context.Context) error {
	if err := m.removeSelfRecord(); err != nil {
		return err
	}
	return m.shutdown(ctx)
}

func (m *Manager) shutdown(ctx context.Context) error {
	if m.onShutdown != nil {
		if err := m.onShutdown(ctx); err != nil {
			return err
		}
	}
	if m.server != nil {
		return m.server.Shutdown(ctx)
	}
	return nil
}
