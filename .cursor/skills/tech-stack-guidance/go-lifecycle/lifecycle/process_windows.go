//go:build windows

package lifecycle

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	modKernel32         = syscall.NewLazyDLL("kernel32.dll")
	procOpenProcess     = modKernel32.NewProc("OpenProcess")
	procGetExitCodeProc = modKernel32.NewProc("GetExitCodeProcess")
	procLockFileEx      = modKernel32.NewProc("LockFileEx")
	procUnlockFileEx    = modKernel32.NewProc("UnlockFileEx")
)

const (
	processQueryLimitedInformation = 0x1000
	stillActive                    = 259
	lockFileExclusiveLock          = 0x00000002
)

type lockOverlapped struct {
	Internal     uintptr
	InternalHigh uintptr
	Offset       uint32
	OffsetHigh   uint32
	hEvent       uintptr
}

func windowsProcessAlive(pid int) bool {
	handle, _, _ := procOpenProcess.Call(
		uintptr(processQueryLimitedInformation),
		0,
		uintptr(pid),
	)
	if handle == 0 {
		return false
	}
	defer syscall.CloseHandle(syscall.Handle(handle))

	var exitCode uint32
	ok, _, _ := procGetExitCodeProc.Call(handle, uintptr(unsafe.Pointer(&exitCode)))
	if ok == 0 {
		return false
	}
	return exitCode == stillActive
}

func acquireFileLock(f *os.File) error {
	handle := syscall.Handle(f.Fd())
	var overlapped lockOverlapped
	ok, _, err := procLockFileEx.Call(
		uintptr(handle),
		uintptr(lockFileExclusiveLock),
		0,
		1,
		0,
		uintptr(unsafe.Pointer(&overlapped)),
	)
	if ok == 0 {
		if err != syscall.Errno(0) {
			return err
		}
		return syscall.EINVAL
	}
	return nil
}

func releaseFileLock(f *os.File) error {
	handle := syscall.Handle(f.Fd())
	var overlapped lockOverlapped
	ok, _, err := procUnlockFileEx.Call(
		uintptr(handle),
		0,
		1,
		0,
		uintptr(unsafe.Pointer(&overlapped)),
	)
	if ok == 0 {
		if err != syscall.Errno(0) {
			return err
		}
		return syscall.EINVAL
	}
	return nil
}
