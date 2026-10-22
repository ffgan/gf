//go:build windows

package modules

import (
	"sort"
	"strings"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"
)

type winProc struct {
	pid     uint32
	ppid    uint32
	threads uint32
	cpu     int64
	name    string
}

func queryProcesses() []winProc {
	ntdll := syscall.NewLazyDLL("ntdll.dll")
	proc := ntdll.NewProc("NtQuerySystemInformation")
	if proc.Find() != nil {
		return nil
	}
	size := uint32(1 << 20)
	var buf []byte
	for attempt := 0; attempt < 6; attempt++ {
		buf = make([]byte, size)
		var need uint32
		status, _, _ := proc.Call(5, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)), uintptr(unsafe.Pointer(&need)))
		if status == 0 {
			if need > 0 && int(need) <= len(buf) {
				buf = buf[:need]
			}
			break
		}
		if status != 0xC0000004 {
			return nil
		}
		if need > size {
			size = need + 65536
		} else {
			size *= 2
		}
		buf = nil
	}
	if len(buf) < 8 {
		return nil
	}
	return walkProcesses(buf)
}

func walkProcesses(buf []byte) []winProc {
	ptr := unsafe.Sizeof(uintptr(0))
	var out []winProc
	off := 0
	for off+8 <= len(buf) {
		next := int(uint32(buf[off]) | uint32(buf[off+1])<<8 | uint32(buf[off+2])<<16 | uint32(buf[off+3])<<24)
		threads := uint32(buf[off+4]) | uint32(buf[off+5])<<8 | uint32(buf[off+6])<<16 | uint32(buf[off+7])<<24
		var userOff, kernelOff, nameOff, pidOff, ppidOff int
		if ptr == 8 {
			userOff, kernelOff, nameOff, pidOff, ppidOff = 40, 48, 56, 80, 88
		} else {
			userOff, kernelOff, nameOff, pidOff, ppidOff = 40, 48, 56, 68, 72
		}
		if off+ppidOff+int(ptr) > len(buf) {
			break
		}
		user := readInt64(buf[off+userOff:])
		kernel := readInt64(buf[off+kernelOff:])
		nameLen := uint16(buf[off+nameOff]) | uint16(buf[off+nameOff+1])<<8
		var namePtr uintptr
		if ptr == 8 {
			namePtr = readUintptr(buf[off+nameOff+8:])
		} else {
			namePtr = uintptr(uint32(buf[off+nameOff+4]) | uint32(buf[off+nameOff+5])<<8 | uint32(buf[off+nameOff+6])<<16 | uint32(buf[off+nameOff+7])<<24)
		}
		pid := readPtrID(buf[off+pidOff:], ptr)
		ppid := readPtrID(buf[off+ppidOff:], ptr)
		out = append(out, winProc{
			pid:     pid,
			ppid:    ppid,
			threads: threads,
			cpu:     user + kernel,
			name:    utf16At(namePtr, nameLen),
		})
		if next == 0 {
			break
		}
		off += next
	}
	return out
}

func readInt64(b []byte) int64 {
	if len(b) < 8 {
		return 0
	}
	return int64(uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56)
}

func readUintptr(b []byte) uintptr {
	if len(b) < 8 {
		return 0
	}
	return uintptr(uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56)
}

func readPtrID(b []byte, ptr uintptr) uint32 {
	if ptr == 8 {
		return uint32(readUintptr(b))
	}
	if len(b) < 4 {
		return 0
	}
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func utf16At(ptr uintptr, nbytes uint16) string {
	if ptr == 0 || nbytes < 2 {
		return ""
	}
	n := int(nbytes) / 2
	raw := unsafe.Slice((*uint16)(unsafe.Pointer(ptr)), n)
	chars := make([]uint16, n)
	copy(chars, raw)
	return strings.TrimRight(string(utf16.Decode(chars)), "\x00")
}

func windowsProcesses() string {
	procs := queryProcesses()
	if len(procs) == 0 {
		return ""
	}
	var n, threads uint64
	for _, proc := range procs {
		if proc.ppid == 0 {
			continue
		}
		n++
		threads += uint64(proc.threads)
	}
	return formatProcesses(n, threads, true)
}

func windowsInit() string {
	for _, proc := range queryProcesses() {
		if proc.ppid != 4 {
			continue
		}
		name := strings.TrimSuffix(proc.name, ".exe")
		name = strings.TrimSuffix(name, ".EXE")
		if name != "" {
			return name
		}
	}
	return ""
}

func windowsTop() string {
	first := queryProcesses()
	if len(first) == 0 {
		return ""
	}
	if usagePause > 0 {
		time.Sleep(usagePause)
	}
	second := queryProcesses()
	prev := map[uint32]int64{}
	for _, proc := range first {
		prev[proc.pid] = proc.cpu
	}
	type rank struct {
		name  string
		delta int64
	}
	var ranks []rank
	use := second
	if len(use) == 0 {
		use = first
	}
	for _, proc := range use {
		if proc.ppid == 0 || proc.name == "" {
			continue
		}
		delta := proc.cpu
		if old, ok := prev[proc.pid]; ok && proc.cpu >= old {
			delta = proc.cpu - old
		}
		name := proc.name
		if i := strings.LastIndexAny(name, `\/`); i >= 0 {
			name = name[i+1:]
		}
		ranks = append(ranks, rank{name: name, delta: delta})
	}
	sort.Slice(ranks, func(i, j int) bool { return ranks[i].delta > ranks[j].delta })
	if len(ranks) > 5 {
		ranks = ranks[:5]
	}
	var lines []string
	for _, rank := range ranks {
		lines = append(lines, rank.name)
	}
	return joinLines(lines)
}

func windowsPlayer() string {
	var names []string
	for _, proc := range queryProcesses() {
		names = append(names, proc.name)
	}
	return playersFromNames(names)
}

func systemTimes() (idle, kernel, user uint64, ok bool) {
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("GetSystemTimes")
	if proc.Find() != nil {
		return 0, 0, 0, false
	}
	var i, k, u uint64
	r, _, _ := proc.Call(uintptr(unsafe.Pointer(&i)), uintptr(unsafe.Pointer(&k)), uintptr(unsafe.Pointer(&u)))
	if r == 0 {
		return 0, 0, 0, false
	}
	return i, k, u, true
}

func windowsCPUUsage() string {
	return usageFromSamples(func() (cpuSample, bool) {
		idle, kernel, user, ok := systemTimes()
		if !ok || kernel < idle {
			return cpuSample{}, false
		}
		inUse := user + (kernel - idle)
		return cpuSample{inUse: inUse, total: inUse + idle}, true
	})
}
