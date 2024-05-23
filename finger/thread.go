package finger

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"test/requestpoc"
	"time"
)

var (
	TimestampStr string
	VulnMapMutex sync.Mutex
	DomainMutex  sync.Mutex
)

const maxConcurrentIPScans = 30

func init() {
	timestamp := time.Now().Unix()
	TimestampStr = fmt.Sprintf("%d", timestamp)
}

func processDomain(domain string, wg *sync.WaitGroup, semaphore chan struct{}) {
	defer wg.Done()
	defer func() { <-semaphore }() // 释放一个信号量
	semaphore <- struct{}{}        // 占用一个信号量

	title, htmlContent, htmlheaders, code, length, icoHash, req := Httprequest(domain)
	fmt.Printf("[%s%s |title: \033[31m%s\033[0m |响应码: %d |返回长度: %d ]\n", req, domain, title, code, length)

	Finger(domain, title, htmlContent, htmlheaders, code, length, icoHash, req)
}

func Http_thread(domains []string) {
	semaphore := make(chan struct{}, maxConcurrentIPScans)
	var wg sync.WaitGroup
	for _, domain := range domains {
		wg.Add(1)
		go processDomain(domain, &wg, semaphore)
	}
	wg.Wait()
	domainresults()
}

func domainresults() {
	// 创建文件夹
	err := os.MkdirAll("results/"+TimestampStr, 0755)
	if err != nil {
		fmt.Println(err)
	}

	// 输出poc结果到txt
	file_poc, err_poc := os.OpenFile("results/"+TimestampStr+"/PocResults.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err_poc != nil {
		fmt.Println(err_poc)
		fmt.Println("创建文件失败")
	}
	defer file_poc.Close() // 确保文件关闭
	writer_poc := bufio.NewWriter(file_poc)
	VulnMapMutex.Lock()
	for url, vulnerabilities := range requestpoc.VulnMap {
		writer_poc.WriteString("URL: " + url + "\n")
		writer_poc.WriteString("漏洞列表: ")
		for i, vulnerability := range vulnerabilities {
			writer_poc.WriteString(vulnerability)
			// 如果不是最后一个漏洞，添加逗号分隔符
			if i < len(vulnerabilities)-1 {
				writer_poc.WriteString(", ")
			}
		}
		writer_poc.WriteString("\n\n")
	}
	VulnMapMutex.Unlock()
	writer_poc.Flush()

	// 对NoFingerdomain按照返回长度排序
	sort.Slice(NoFingerdomain, func(i, j int) bool {
		// 解析返回长度
		lengthI, _ := strconv.Atoi(strings.Split(strings.Split(NoFingerdomain[i], "|返回长度: ")[1], " ")[0])
		lengthJ, _ := strconv.Atoi(strings.Split(strings.Split(NoFingerdomain[j], "|返回长度: ")[1], " ")[0])
		return lengthI > lengthJ
	})

	// 输出排序后的结果到文件
	file, err := os.OpenFile("results/"+TimestampStr+"/NoFinger.json", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		fmt.Println("创建文件失败")
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	DomainMutex.Lock()
	for _, url := range NoFingerdomain {
		writer.WriteString(url)
	}
	DomainMutex.Unlock()
	writer.Flush()

	// 对Fingerdomain按照返回长度排序
	sort.Slice(Fingerdomain, func(i, j int) bool {
		// 解析返回长度
		lengthI, _ := strconv.Atoi(strings.Split(strings.Split(Fingerdomain[i], "|返回长度: ")[1], " ")[0])
		lengthJ, _ := strconv.Atoi(strings.Split(strings.Split(Fingerdomain[j], "|返回长度: ")[1], " ")[0])
		return lengthI > lengthJ
	})

	// 输出排序后的结果到文件
	file_finger, err_finger := os.OpenFile("results/"+TimestampStr+"/Finger.json", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err_finger != nil {
		fmt.Println(err_finger)
		fmt.Println("创建文件失败")
	}
	defer file_finger.Close()
	writer_finger := bufio.NewWriter(file_finger)
	DomainMutex.Lock()
	for _, url := range Fingerdomain {
		writer_finger.WriteString(url)
	}
	DomainMutex.Unlock()
	writer_finger.Flush()
}
