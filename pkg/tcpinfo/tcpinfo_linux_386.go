//go:build linux && 386

package tcpinfo

import (
	"syscall"
	"unsafe"
)

const netGetSockOpt = 15

// GetTCPInfo calls socketcall(2) on Linux to retrieve tcp_info and unpacks that into the golang-friendly TCPInfo.
// This variant is for the 32-bit x86 (386) architecture.
func GetRawTCPInfo(fd uintptr) (*RawTCPInfo, error) {
	var value RawTCPInfo
	length := uint32(sizeOfRawTCPInfo)

	args := [5]uintptr{
		uintptr(fd),
		uintptr(syscall.SOL_TCP), uintptr(syscall.TCP_INFO),
		uintptr(unsafe.Pointer(&value)), uintptr(unsafe.Pointer(&length)),
	}

	_, _, errNo := syscall.RawSyscall(
		syscall.SYS_SOCKETCALL,
		netGetSockOpt,
		uintptr(unsafe.Pointer(&args)),
		0,
	)
	if errNo != 0 {
		switch errNo {
		case syscall.EAGAIN:
			return nil, EAGAIN
		case syscall.EINVAL:
			return nil, EINVAL
		case syscall.ENOENT:
			return nil, ENOENT
		}
		return nil, errNo
	}

	return &value, nil
}
