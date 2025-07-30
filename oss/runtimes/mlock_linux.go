package runtimes

import (
	"errors"
)

/*
#include <sys/mman.h>
*/
import "C"

/*
MLock linux有效
锁定当前进程已分配的所有内存
锁定未来可能新分配的内存
所有锁定内存将不会被换出到 swap
通常配合 ulimit 或 systemd 设置 LimitMEMLOCK=infinity 才能生效

必须有权限
在 shell 中运行前使用 ulimit -l unlimited
或 systemd 配置 LimitMEMLOCK=infinity
*/
func MLock() error {
	ret := C.mlockall(C.MCL_CURRENT | C.MCL_FUTURE)
	if ret != 0 {
		return errors.New("mlockall failed")
	}
	return nil
}
