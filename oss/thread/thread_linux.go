package thread

import (
	"fmt"
	"runtime"

	"golang.org/x/sys/unix"
)

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

	// 2. 设置实时调度策略 SCHED_FIFO + 优先级
	param := &unix.SchedParam{SchedPriority: priority}
	if err := unix.SchedSetscheduler(0, unix.SCHED_FIFO, param); err != nil {
		return fmt.Errorf("failed to set SCHED_FIFO: %v", err)
	}

	// 3. 锁住内存，防止 swap
	if err := unix.Mlockall(unix.MCL_CURRENT | unix.MCL_FUTURE); err != nil {
		return fmt.Errorf("failed to lock memory: %v", err)
	}

	return nil
}
