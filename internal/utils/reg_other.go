//go:build !windows

package utils

const (
	HKCU uintptr = 0x80000001
	HKLM uintptr = 0x80000002
)

func RegBytes(root uintptr, path, name string) ([]byte, bool) {
	return nil, false
}

func RegString(root uintptr, path, name string) (string, bool) {
	return "", false
}

func RegKeyExists(root uintptr, path string) bool {
	return false
}

func RegUint(root uintptr, path, name string) (uint32, bool) {
	return 0, false
}

func RegUint64(root uintptr, path, name string) (uint64, bool) {
	return 0, false
}

func RegSubkeys(root uintptr, path string) ([]string, bool) {
	return nil, false
}
