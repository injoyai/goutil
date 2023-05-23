package ip

import (
	"net"
	"time"
)

func Ping(ip string, timeout time.Duration) (time.Duration, error) {
	conn, err := net.DialTimeout("ip:icmp", ip, timeout)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	t := time.Now()
	if err = conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return 0, err
	}
	if _, err = conn.Write([]byte{
		8, 0, 247, 253, 0, 1, 0, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0}); err != nil {
		return 0, err
	}
	buf := make([]byte, 65535)
	_, err = conn.Read(buf)
	return time.Since(t), err
}
