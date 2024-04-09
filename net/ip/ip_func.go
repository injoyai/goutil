package ip

import (
	"errors"
	"fmt"
	"github.com/injoyai/conv"
	"net"
	"strconv"
	"strings"
)

//// RangeFunc 遍历ip地址
//func RangeFunc(start, end net.IP, f func(ip net.IP) bool) {
//	for i := conv.Uint32([]byte(start.To4())); i <= conv.Uint32([]byte(end.To4())); i++ {
//		ip := net.IP(conv.Bytes(i))
//		if !f(ip) {
//			break
//		}
//	}
//}

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

func RangeNetwork(network string, fn func(inter *Inter)) error {
	inters, err := net.Interfaces()
	if err != nil {
		return err
	}
	for i, inter := range inters {
		if inter.Flags&(1<<net.FlagLoopback) == 1 || inter.Flags&(1<<net.FlagUp) == 0 {
			continue
		}
		if len(network) > 0 && !strings.Contains(inter.Name, network) {
			continue
		}
		fn(&Inter{
			Index:     i,
			Interface: inter,
		})
	}
	return nil
}

type Inter struct {
	Index int
	net.Interface
}

func (this *Inter) Print() {
	fmt.Printf("\n%d: %s (%s):\n", this.Index, this.Name, this.HardwareAddr)
}

// RangeSegment 遍历网段
func (this *Inter) RangeSegment(fn func(ipv4 net.IP, self bool) bool) error {
	return this.RangeV4(func(ipv4 net.IP) bool {
		for i := conv.Uint32([]byte{ipv4[0], ipv4[1], ipv4[2], 0}); i <= conv.Uint32([]byte{ipv4[0], ipv4[1], ipv4[2], 255}); i++ {
			ip := net.IP(conv.Bytes(i))
			if !fn(ip, ip.String() == ipv4.String()) {
				return false
			}
		}
		return true
	})
}

// RangeV4 遍历ipv4
func (this *Inter) RangeV4(fn func(ipv4 net.IP) bool) error {
	addrs, err := this.Addrs()
	if err != nil {
		return err
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			fn(ipNet.IP.To4())
		}
	}
	return nil
}
