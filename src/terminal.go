package src

import (
	"syscall"
	"unsafe"
)

type Winsize struct {
	Rows   uint16
	Cols   uint16
	Xpixel uint16
	Ypixel uint16
}

// ws, err := IoctlGetWinsize(int(os.Stdout.Fd()))
func IoctlGetWinsize(fd int) (*Winsize, error) {
	var value Winsize
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&value)),
	)
	if errno != 0 {
		return nil, errno
	}
	return &value, nil
}
