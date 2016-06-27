package runtime // import "blitznote.com/src/go.runtime"

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// NumCPU queries the system for the count of threads available
// for use to this process.
//
// Issues two syscalls.
// Returns 0 on errors. Use |runtime.NumCPU| in that case.
func NumCPU() int {
	// Gets the affinity mask for this process.
	pid, _, _ := syscall.RawSyscall(unix.SYS_GETPID, 0, 0, 0)

	var mask [1024 / 64]uintptr
	_, _, err := syscall.RawSyscall(unix.SYS_SCHED_GETAFFINITY, pid, uintptr(len(mask)*8), uintptr(unsafe.Pointer(&mask[0])))
	if err != 0 {
		return 0
	}

	// For every available thread a bit is set in the mask.
	ncpu := 0
	for _, e := range mask {
		if e == 0 {
			continue
		}
		ncpu += int(popcnt(uint64(e)))
	}
	return ncpu
}
