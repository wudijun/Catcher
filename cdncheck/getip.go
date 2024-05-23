package cdncheck

import (
	"fmt"
	"net"
)

func Getip(domains []string) []net.IP {
	var iplist []net.IP // 在函数内部声明 iplist 变量

	for _, domain := range domains {
		ips, err := net.LookupIP(domain)
		fmt.Println("正在解析" + domain + "的ip")
		if err != nil {
			//fmt.Printf("Failed to lookup IP for domain %s: %v\n", domain, err)
			continue
		}

		// 只打印第一个IP地址
		if len(ips) > 0 {
			// fmt.Printf("First IP address for domain %s: %s\n", domain, ips[0])
			iplist = append(iplist, ips[0])
		}
	}
	return iplist
}
