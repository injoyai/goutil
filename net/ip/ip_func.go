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

func RangeV4(search string, fn func(n net.Interface, ip net.IP, self bool) bool) error {
	is, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, network := range is {

		if network.Flags&(1<<net.FlagLoopback) == 1 || network.Flags&(1<<net.FlagUp) == 0 {
			continue
		}
		if len(search) > 0 && search != "all" && !strings.Contains(network.Name, search) {
			continue
		}

		addrs, err := network.Addrs()
		if err != nil {
			return err
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				selfIP := ipnet.IP.To4()
				for i := conv.Uint32([]byte{selfIP[0], selfIP[1], selfIP[2], 0}); i <= conv.Uint32([]byte{selfIP[0], selfIP[1], selfIP[2], 255}); i++ {
					ip := net.IP(conv.Bytes(i))
					if !fn(network, ip, selfIP.String() == ip.String()) {
						break
					}
				}
			}
		}
	}
	return nil
}
