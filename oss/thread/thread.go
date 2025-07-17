package thread

import (
	"log"
	"runtime"
)

// StartRealtime 创建并绑定当前线程为实时任务,绑定到指定CPU,设置优先级(1-99),越大优先级越高
func StartRealtime(cpuID int, priority int, task func()) {
	go func() {
		//锁定当前协程到当前(固定)系统线程,防止调度迁移
		runtime.LockOSThread()
		//绑定当前线程到指定CPU上执行
		err := SetupRealtime(cpuID, priority)
		if err != nil {
			log.Fatalf("Failed to setup realtime thread on CPU %d: %v", cpuID, err)
		}
		log.Printf("✅ Realtime thread initialized on CPU %d with priority %d", cpuID, priority)
		//执行实时任务
		task()
	}()
}
