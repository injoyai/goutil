//go:build !linux || !arm64

package thread

// BindCPU 创建并绑定当前线程为实时任务,windows 不支持
func BindCPU(cpuID int, priority int) error {
	return nil
}
