package thread

import (
	"github.com/injoyai/logs"
)

// StartRealtime 创建并绑定当前线程为实时任务,绑定到指定CPU,设置优先级(1-99),越大优先级越高
func StartRealtime(cpuID int, priority int, task func()) {
	go func() {
		//绑定当前线程到指定CPU上执行
		err := BindCPU(cpuID, priority)
		if err != nil {
			logs.Fatalf("Failed to setup realtime thread on CPU %d: %v", cpuID, err)
		}
		//执行实时任务
		task()
	}()
}
