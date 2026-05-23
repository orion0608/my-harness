package lifecycle

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func processAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	if isWindows() {
		// FindProcess always succeeds on Windows; use tasklist-style check via OpenProcess.
		return windowsProcessAlive(pid)
	}
	err = proc.Signal(syscall.Signal(0))
	return err == nil
}

func killProcess(pid int) error {
	if pid <= 0 {
		return nil
	}
	if isWindows() {
		cmd := exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/F")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("taskkill pid %d: %w", pid, err)
		}
		return nil
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return proc.Kill()
}
