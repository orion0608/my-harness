//go:build !windows

package lifecycle

func windowsProcessAlive(pid int) bool {
	return false
}
