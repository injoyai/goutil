package main

import (
	"fmt"
	"github.com/injoyai/base/oss/shell"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/net/ip"
	"github.com/injoyai/logs"
	"github.com/spf13/cobra"
	"go.bug.st/serial"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func handlerScanICMP(cmd *cobra.Command, args []string, flags *Flags) {
	gateIPv4 := []byte(net.ParseIP(ip.GetLocal())[12:15])
	wg := sync.WaitGroup{}
	for i := conv.Uint32(append(gateIPv4, 0)); i <= conv.Uint32(append(gateIPv4, 255)); i++ {
		ipv4 := net.IPv4(uint8(i>>24), uint8(i>>16), uint8(i>>8), uint8(i))
		wg.Add(1)
		go func(ipv4 net.IP) {
			defer wg.Done()
			used, err := ip.Ping(ipv4.String(), time.Second)
			if err == nil {
				fmt.Printf("%s: %s\n", ipv4, used.String())
			}
		}(ipv4)
	}
	wg.Wait()
}

func handlerScanSSH(cmd *cobra.Command, args []string, flags *Flags) {
	handlerScanPort(cmd, []string{"22"}, flags)
}

func handlerScanPort(cmd *cobra.Command, args []string, flags *Flags) {
	if len(args) == 0 {
		log.Println("[错误]", "缺少端口")
	}
	gateIPv4 := []byte(net.ParseIP(ip.GetLocal())[12:15])
	wg := sync.WaitGroup{}
	for i := conv.Uint32(append(gateIPv4, 0)); i <= conv.Uint32(append(gateIPv4, 255)); i++ {
		ipv4 := net.IPv4(uint8(i>>24), uint8(i>>16), uint8(i>>8), uint8(i))
		wg.Add(1)
		go func(ipv4 net.IP) {
			defer wg.Done()
			addr := fmt.Sprintf("%s:%s", ipv4, args[0])
			c, err := net.DialTimeout("tcp", addr, time.Second)
			if err == nil {
				c.Close()
				fmt.Printf("%s   开启\n", addr)
			}
		}(ipv4)
	}
	wg.Wait()
}

func handlerScanSerial(cmd *cobra.Command, args []string, flags *Flags) {
	list, err := serial.GetPortsList()
	if err != nil {
		logs.Err(err)
		return
	}
	fmt.Println(strings.Join(list, "\n"))
}

func handlerScanEdge(cmd *cobra.Command, args []string, flags *Flags) {
	ipv4 := ip.GetLocal()
	startIP := append(net.ParseIP(ipv4)[:15], 0)
	endIP := append(net.ParseIP(ipv4)[:15], 255)
	ch, ctx := qlScanEdge(startIP, endIP)
	for i := 1; ; i++ {
		select {
		case <-ctx.Done():
			return
		case data := <-ch:
			fmt.Printf("%v: %v\n", data.IP, data.SN)
			if flags.GetBool("open") {
				logs.PrintErr(shell.OpenBrowser(fmt.Sprintf("http://%s:10001", data.IP)))
			}
		}
	}
}

func handlerScan(cmd *cobra.Command, args []string, flags *Flags) {
	switch true {
	case len(args) == 0:
		log.Println("[错误]", "缺少扫描类型(icmp,serial...)")
	default:

		number := flags.GetInt("number")

		switch args[0] {
		case "icmp":

			gateIPv4 := []byte(net.ParseIP(ip.GetLocal())[12:15])
			wg := sync.WaitGroup{}
			for i := conv.Uint32(append(gateIPv4, 0)); i <= conv.Uint32(append(gateIPv4, 255)); i++ {
				ipv4 := net.IPv4(uint8(i>>24), uint8(i>>16), uint8(i>>8), uint8(i))
				wg.Add(1)
				go func(ipv4 net.IP) {
					defer wg.Done()
					used, err := ip.Ping(ipv4.String(), time.Second)
					if err == nil {
						fmt.Printf("%s: %s\n", ipv4, used.String())
					}
				}(ipv4)
			}
			wg.Wait()

		case "ssh":

			gateIPv4 := []byte(net.ParseIP(ip.GetLocal())[12:15])
			wg := sync.WaitGroup{}
			for i := conv.Uint32(append(gateIPv4, 0)); i <= conv.Uint32(append(gateIPv4, 255)); i++ {
				ipv4 := net.IPv4(uint8(i>>24), uint8(i>>16), uint8(i>>8), uint8(i))
				wg.Add(1)
				go func(ipv4 net.IP) {
					defer wg.Done()
					c, err := net.Dial("tcp", ipv4.String()+":22")
					if err == nil {
						c.Close()
						fmt.Printf("%s\n", ipv4)
					}
				}(ipv4)
			}
			wg.Wait()

		case "serial":

			list, err := serial.GetPortsList()
			if err != nil {
				logs.Err(err)
				return
			}
			fmt.Println(strings.Join(list, "\n"))

		case "edge":

			ipv4 := ip.GetLocal()
			startIP := append(net.ParseIP(ipv4)[:15], 0)
			endIP := append(net.ParseIP(ipv4)[:15], 255)
			ch, ctx := qlScanEdge(startIP, endIP)
			for i := 1; ; i++ {
				select {
				case <-ctx.Done():
					return
				case data := <-ch:
					fmt.Printf("%v: %v\n", data.IP, data.SN)
					if flags.GetBool("open") {
						logs.PrintErr(shell.OpenBrowser(fmt.Sprintf("http://%s:10001", data.IP)))
					}
					if number > 0 && i >= number {
						return
					}
				}
			}

		}
	}
}
