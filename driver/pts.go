package main

import (
	"fmt"
	"os"
	"path"
	"syscall"
	"unsafe"
)

func CreatePTS(slavePath string) (masterPTS *os.File, err error) {
	masterPTS, err = os.OpenFile("/dev/ptmx", syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK, 0666)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			masterPTS.Close()
			masterPTS = nil
		}
	}()

	setupTermios(masterPTS.Fd())

	// create pts
	var data int
	err = ioctl(masterPTS.Fd(), syscall.TIOCGPTN, unsafe.Pointer(&data))
	if err != nil {
		return
	}
	slaverName := fmt.Sprintf("/dev/pts/%d", data)

	os.Remove(slavePath)
	os.MkdirAll(path.Dir(slavePath), 0755)
	err = os.Symlink(slaverName, slavePath)
	if err != nil {
		return
	}
	// set owner
	if err = os.Chown(slaverName, os.Getuid(), os.Getgid()); err != nil {
		return
	}
	if err = os.Chmod(slaverName, 0666); err != nil {
		return
	}

	// unlock master
	data = 0
	err = ioctl(masterPTS.Fd(), syscall.TIOCSPTLCK, unsafe.Pointer(&data))
	return
}

func setupTermios(fd uintptr) (*syscall.Termios, error) {
	var oldState syscall.Termios

	err := ioctl(fd, syscall.TCGETS, unsafe.Pointer(&oldState))
	if err != nil {
		return nil, err
	}

	newState := oldState
	// newState.Cflag ^= (IGNBRK | BRKINT | ICRNL |
	// 	INLCR | PARMRK | INPCK | ISTRIP | IXON)

	newState.Cflag &= syscall.B9600 | syscall.CLOCAL | syscall.CS8 | syscall.CREAD
	newState.Iflag &^= (syscall.ISTRIP | syscall.INLCR | syscall.ICRNL | syscall.IGNCR | syscall.IXON | syscall.IXOFF)
	newState.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.ISIG | syscall.ECHONL | syscall.IEXTEN
	newState.Iflag = 0
	newState.Oflag = 0
	newState.Ispeed = syscall.B9600
	newState.Ospeed = syscall.B9600

	// Base settings
	// cflagToUse := syscall.CREAD | syscall.CLOCAL | rate
	// cflagToUse |= syscall.CS8

	// t := syscall.Termios{
	// 	Iflag:  syscall.IGNPAR,
	// 	Cflag:  cflagToUse,
	// 	Cc:     [32]uint8{syscall.VMIN: 0, syscall.VTIME: 0},
	// 	Ispeed: rate,
	// 	Ospeed: rate,
	//	}

	err = ioctl(fd, syscall.TCSETS, unsafe.Pointer(&newState))
	if err != nil {
		return nil, err
	}

	{ //debug
		var s syscall.Termios
		ioctl(fd, syscall.TCGETS, unsafe.Pointer(&s))
		//		fmt.Printf("Old: %#v \n New: %#v\n Real: %v\n", oldState, newState, s)
	}

	return &oldState, nil
}

func makePTSLink(fpath string, num int) error {
	return os.Symlink(fmt.Sprintf("/dev/pts/%d", num), fpath)
}

func ioctl(fd uintptr, command uint, data unsafe.Pointer) error {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd),
		uintptr(command), uintptr(data))
	if err != 0 {
		return err
	}
	return nil
}
