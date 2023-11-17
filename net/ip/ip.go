package ip

import (
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
)

// IsLocal 是否是局域网ip
func IsLocal() bool {
	return regexp.MustCompile(
		`^(\[::1\]:\d+)|(127\.0\.0\.1)|(localhost)|(10\.\d{1,3}\.\d{1,3}\.\d{1,3})|(172\.((1[6-9])|(2\d)|(3[01]))\.\d{1,3}\.\d{1,3})|(192\.168\.\d{1,3}\.\d{1,3})`,
	).MatchString(GetLocal())
}

func GetLocal() string {
	ipv4 := ""
	for _, v := range GetLocalAll() {
		ipv4 = v
		if strings.Contains(v, "192.168.") {
			return v
		}
	}
	return ipv4
}

func GetLocalAll() []string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	ips := []string(nil)
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ips = append(ips, ipnet.IP.String())
					}
				}
			}
		}
	}
	return ips
}

func GetRemote(r *http.Request) string {
	if len(r.RemoteAddr) == 0 {
		realIps := r.Header.Get("X-Forwarded-For")
		if realIps != "" && len(realIps) != 0 && !strings.EqualFold("unknown", realIps) {
			ipArray := strings.Split(realIps, ",")
			r.RemoteAddr = ipArray[0]
		}
		if r.RemoteAddr == "" || strings.EqualFold("unknown", realIps) {
			r.RemoteAddr = r.Header.Get("Proxy-Client-IP")
		}
		if r.RemoteAddr == "" || strings.EqualFold("unknown", realIps) {
			r.RemoteAddr = r.Header.Get("WL-Proxy-Client-IP")
		}
		if r.RemoteAddr == "" || strings.EqualFold("unknown", realIps) {
			r.RemoteAddr = r.Header.Get("HTTP_CLIENT_IP")
		}
		if r.RemoteAddr == "" || strings.EqualFold("unknown", realIps) {
			r.RemoteAddr = r.Header.Get("HTTP_X_FORWARDED_FOR")
		}
		if r.RemoteAddr == "" || strings.EqualFold("unknown", realIps) {
			r.RemoteAddr = r.Header.Get("X-Real-IP")
		}
	}
	return r.RemoteAddr
}

func GetNetV4() string {
	response, errClient := http.Get("https://ipv4.netarm.com")
	if errClient != nil {
		return ""
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	return string(body)
}

//func GetIPv4() {
//	// Parse a STUN URI
//	u, err := stun.ParseURI("stun:stun.l.google.com:19302")
//	if err != nil {
//		panic(err)
//	}
//
//	// Creating a "connection" to STUN server.
//	c, err := stun.DialURI(u, &stun.DialConfig{})
//	if err != nil {
//		panic(err)
//	}
//	// Building binding request with random transaction id.
//	message := stun.MustBuild(stun.TransactionID, stun.BindingRequest)
//	// Sending request to STUN server, waiting for response message.
//	if err := c.Do(message, func(res stun.Event) {
//		if res.Error != nil {
//			panic(res.Error)
//		}
//		// Decoding XOR-MAPPED-ADDRESS attribute from message.
//		var xorAddr stun.XORMappedAddress
//		if err := xorAddr.GetFrom(res.Message); err != nil {
//			panic(err)
//		}
//		fmt.Println("your IP is", xorAddr)
//	}); err != nil {
//		panic(err)
//	}
//}

func GetNetV6() string {
	response, errClient := http.Get("https://ipv6.netarm.com")
	if errClient != nil {
		return ""
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	return string(body)
}
