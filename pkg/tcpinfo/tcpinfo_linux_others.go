//go:build linux && !386

package tcpinfo

import (
	"syscall"
	"unsafe"
)

// GetRawTCPInfo calls getsockopt(2) on Linux to retrieve tcp_info and unpacks that into the golang-friendly TCPInfo.
// This variant is for all non-x86 (386) architectures.
func GetRawTCPInfo(fd uintptr) (*RawTCPInfo, error) {
	var value RawTCPInfo
	length := uint32(sizeOfRawTCPInfo)
	_, _, errNo := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		uintptr(fd),
		uintptr(syscall.SOL_TCP),
		uintptr(syscall.TCP_INFO),
		uintptr(unsafe.Pointer(&value)),
		uintptr(unsafe.Pointer(&length)),
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
