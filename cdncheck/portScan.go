package cdncheck

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sort"
	"sync"
	"test/finger"
	"time"
)

func UniqueSortedIPs(ips []net.IP) []net.IP {
	seen := make(map[string]struct{})
	var unique []net.IP

	for _, ip := range ips {
		ipStr := ip.String()
		if _, found := seen[ipStr]; !found {
			unique = append(unique, ip)
			seen[ipStr] = struct{}{}
		}
	}

	// 对IP地址进行排序
	sort.Slice(unique, func(i, j int) bool {
		return unique[i].String() < unique[j].String()
	})

	// 限制同时扫描的IP地址数量
	const maxConcurrentIPScans = 10
	semaphore := make(chan struct{}, maxConcurrentIPScans)

	var wg sync.WaitGroup
	for _, ip := range unique {
		wg.Add(1)
		semaphore <- struct{}{} // 占用一个信号量
		go func(ip net.IP) {
			defer func() { <-semaphore }() // 释放一个信号量
			defer wg.Done()
			scanIP(ip)
		}(ip)
	}
	wg.Wait()

	return unique
}

func scanIP(ip net.IP) {
	file, err := os.OpenFile("results/"+finger.TimestampStr+"/Ports.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		fmt.Println("创建文件失败")
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// 示例需要扫描的端口列表
	ports := []int{21, 22, 80, 81, 82, 83, 85, 88, 89, 90, 91, 92, 99,
		135, 139, 443, 445, 800, 808, 888, 1000, 1010, 1080, 1081, 1082,
		1099, 1433, 1521, 1888, 2008, 2020, 2375, 27017, 3000, 3008, 3306, 3389, 3505, 5432, 6080, 6379, 7000, 7001, 7002, 7003,
		7070, 7071, 7074, 7078, 7080, 7088, 7200, 7680, 7687, 7688, 7777, 7890, 8000, 8001, 8002, 8003, 8004, 8006, 8008, 8009, 8010, 8011,
		8012, 8016, 8018, 8020, 8028, 8030, 8038, 8042, 8044, 8046, 8048, 8053, 8060,
		8069, 8070, 8080, 8081, 8082, 8083, 8084, 8085, 8086, 8087, 8088, 8089, 8090,
		8091, 8092, 8093, 8094, 8095, 8096, 8097, 8098, 8099, 8100, 8101, 8108, 8118,
		8161, 8172, 8180, 8181, 8200, 8244, 8280, 8288, 8300, 8360, 8443, 8448, 8800,
		8848, 8858, 8868, 8879, 8880, 8881, 8888, 8899, 8983, 8989, 9000, 9001, 9002, 9008,
		9010, 9043, 9060, 9080, 9081, 9082, 9083, 9084, 9085, 9086, 9087, 9088, 9089, 9090,
		9091, 9092, 9093, 9094, 9095, 9096, 9097, 9098, 9099, 9100, 9200, 9443, 9448, 9800,
		9981, 9986, 9988, 9998, 9999, 10000, 10001, 10002, 10004, 10008, 10010, 10250, 11211,
		12018, 12443, 14000, 16080, 18000, 18001, 18002, 18004, 18008, 18080, 18082, 18088,
		18090, 18098, 19001, 20000, 20720, 21000, 21501, 21502, 28018, 20880, 27017, 61616} // 修改为需要扫描的端口列表
	var wg sync.WaitGroup
	for _, port := range ports {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			scanPort(ip, port, writer)
		}(port)
	}
	wg.Wait()
}

func isOpen(ip net.IP, port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip.String(), port), 15*time.Second) // 设置超时时间为15秒
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func scanPort(ip net.IP, port int, writer *bufio.Writer) {
	if isOpen(ip, port) {
		result := fmt.Sprintf("%s:%d open\n", ip.String(), port)
		fmt.Print(result)
		writer.WriteString(result)
	}
}
