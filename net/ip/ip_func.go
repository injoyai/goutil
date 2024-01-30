package ip

import (
	"errors"
	"github.com/injoyai/conv"
	"net"
	"strconv"
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
	host, portStr, err := net.SplitHostPort(s)
	if err != nil {
		if strings.HasSuffix(err.Error(), "missing port in address") {
			host, portStr, err = net.SplitHostPort(s + ":80")
		}
	}
	if err != nil {
		return nil, 0, err
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return nil, 0, errors.New("无效地址: " + s)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, 0, err
	}
	return ip, uint16(port), nil
}
