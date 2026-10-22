//go:build windows

package utils

import (
	"encoding/binary"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

const (
	HKCU uintptr = 0x80000001
	HKLM uintptr = 0x80000002
)

func RegBytes(root uintptr, path, name string) ([]byte, bool) {
	_, data, ok := regQuery(root, path, name)
	return data, ok
}

func RegString(root uintptr, path, name string) (string, bool) {
	kind, data, ok := regQuery(root, path, name)
	if !ok || kind != 1 || len(data) < 2 {
		return "", false
	}
	u := make([]uint16, 0, len(data)/2)
	for i := 0; i+1 < len(data); i += 2 {
		c := binary.LittleEndian.Uint16(data[i:])
		if c == 0 {
			break
		}
		u = append(u, c)
	}
	return string(utf16.Decode(u)), true
}

func RegKeyExists(root uintptr, path string) bool {
	adv := syscall.NewLazyDLL("advapi32.dll")
	open := adv.NewProc("RegOpenKeyExW")
	closeKey := adv.NewProc("RegCloseKey")
	sub, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false
	}
	var h syscall.Handle
	r, _, _ := open.Call(root, uintptr(unsafe.Pointer(sub)), 0, 0x20019, uintptr(unsafe.Pointer(&h)))
	if r != 0 {
		return false
	}
	closeKey.Call(uintptr(h))
	return true
}

func RegUint(root uintptr, path, name string) (uint32, bool) {
	kind, data, ok := regQuery(root, path, name)
	if !ok || (kind != 4 && kind != 11) || len(data) < 4 {
		return 0, false
	}
	return binary.LittleEndian.Uint32(data[:4]), true
}

func regQuery(root uintptr, path, name string) (uint32, []byte, bool) {
	adv := syscall.NewLazyDLL("advapi32.dll")
	open := adv.NewProc("RegOpenKeyExW")
	query := adv.NewProc("RegQueryValueExW")
	closeKey := adv.NewProc("RegCloseKey")
	sub, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return 0, nil, false
	}
	var h syscall.Handle
	r, _, _ := open.Call(root, uintptr(unsafe.Pointer(sub)), 0, 0x20019, uintptr(unsafe.Pointer(&h)))
	if r != 0 {
		return 0, nil, false
	}
	defer closeKey.Call(uintptr(h))

	var namePtr *uint16
	if name != "" {
		namePtr, err = syscall.UTF16PtrFromString(name)
		if err != nil {
			return 0, nil, false
		}
	}
	var kind, n uint32
	r, _, _ = query.Call(uintptr(h), uintptr(unsafe.Pointer(namePtr)), 0, uintptr(unsafe.Pointer(&kind)), 0, uintptr(unsafe.Pointer(&n)))
	if r != 0 && r != 234 {
		return 0, nil, false
	}
	if n == 0 {
		return kind, nil, true
	}
	buf := make([]byte, n)
	r, _, _ = query.Call(uintptr(h), uintptr(unsafe.Pointer(namePtr)), 0, uintptr(unsafe.Pointer(&kind)), uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&n)))
	if r != 0 || int(n) > len(buf) {
		return 0, nil, false
	}
	return kind, buf[:n], true
}
