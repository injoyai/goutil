package ip

import (
	"github.com/injoyai/conv"
	"net"
)

// RangeFunc 遍历ip地址
func RangeFunc(start, end net.IP, f func(ip net.IP) bool) {
	for i := conv.Uint32([]byte(start.To4())); i <= conv.Uint32([]byte(end.To4())); i++ {
		ip := net.IP(conv.Bytes(i))
		if !f(ip) {
			break
		}
	}
}
