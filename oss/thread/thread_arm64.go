//go:build linux && arm64

package thread

import (
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"runtime"
	"unsafe"
)

type schedParam struct {
	SchedPriority int32
	_             [4]byte // padding (确保与 C struct 对齐)
}

// SetupRealtime 创建并绑定当前线程为实时任务
func SetupRealtime(cpuID int, priority int) error {
	// 1. 绑定当前线程到指定 CPU
	runtime.LockOSThread()

	var cpuset unix.CPUSet
	cpuset.Zero()
	cpuset.Set(cpuID)

	if err := unix.SchedSetaffinity(0, &cpuset); err != nil {
		return fmt.Errorf("failed to set CPU affinity: %v", err)
	}

	param := schedParam{SchedPriority: int32(priority)}
	_, _, errno := unix.Syscall(
		unix.SYS_SCHED_SETSCHEDULER,
		0,                        // 0 = current thread
		uintptr(unix.SCHED_FIFO), // policy
		uintptr(unsafe.Pointer(&param)),
	)
	if errno != 0 {
		return fmt.Errorf("sched_setscheduler failed: %v", errno)
	}

	// 3. 锁住内存，防止 swap
	if err := unix.Mlockall(unix.MCL_CURRENT | unix.MCL_FUTURE); err != nil {
		return fmt.Errorf("failed to lock memory: %v", err)
	}

	log.Printf("✅ Realtime thread initialized on CPU %d with priority %d", cpuID, priority)
	return nil
}
