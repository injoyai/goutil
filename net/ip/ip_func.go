package ip

import (
	"errors"
	"github.com/injoyai/conv"
	"net"
	"strings"
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

func ParseV4(s string) (net.IP, uint16, error) {
	list := strings.Split(s, ":")
	ip := net.ParseIP(list[0])
	if ip == nil {
		return nil, 0, errors.New("无效地址: " + s)
	}
	port := uint16(80)
	if len(list) > 1 {
		port = conv.Uint16(list[1])
	}
	return ip, port, nil
}
