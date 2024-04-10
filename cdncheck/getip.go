package cdncheck

import (
	"fmt"
	"net"
)

func Getip(domains []string) []net.IP {
	var iplist []net.IP

	for _, domain := range domains {
		ips, err := net.LookupIP(domain)
		fmt.Println("正在解析" + domain + "的ip")
		if err != nil {
			continue
		}

		// 只打印第一个IP地址
		if len(ips) > 0 {
			iplist = append(iplist, ips[0])
		}
	}
	return iplist
}
