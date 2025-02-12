package cdncheck

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"test/finger"

	"github.com/miekg/dns"
)

var (
	NoCdndomains []string
)

const maxConcurrentIPScans = 10

func Cdn(domain string) (bool, string) {
	file_nocdn, err := os.OpenFile("results/"+finger.TimestampStr+"/NoCdn.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		fmt.Println("创建文件失败")
	}
	defer file_nocdn.Close()
	file_cdn, err := os.OpenFile("results/"+finger.TimestampStr+"/Cdn.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		fmt.Println("创建文件失败")
	}
	defer file_cdn.Close()
	file_errorcdn, err := os.OpenFile("results/"+finger.TimestampStr+"/ErrorCdn.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		fmt.Println("创建文件失败")
	}
	defer file_errorcdn.Close()

	dnsServer := "8.8.8.8" // 使用 Google 的 DNS 服务器地址

	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	m.RecursionDesired = true

	r, _, err := c.Exchange(m, dnsServer+":53")
	if err != nil {
		//log.Println("Error resolving domain:", err)  报错
		//fmt.Println(err)
		writer_error := bufio.NewWriter(file_errorcdn)
		defer writer_error.Flush() // 延迟刷新缓冲区到函数结束
		writer_error.WriteString(domain + "\n")
		return true, domain
	}

	ipCount := 0
	for _, ans := range r.Answer {
		if _, ok := ans.(*dns.A); ok {
			ipCount++
			if ipCount >= 2 {
				writer_cdn := bufio.NewWriter(file_cdn)
				defer writer_cdn.Flush() // 延迟刷新缓冲区到函数结束
				writer_cdn.WriteString(domain + "\n")
				return true, domain
			}
		} else if cname, ok := ans.(*dns.CNAME); ok {
			cnameStr := strings.TrimRight(cname.Target, ".")
			parts1 := strings.Split(cnameStr, ".")
			parts2 := strings.Split(domain, ".")
			if strings.Join(parts1[len(parts1)-2:], ".") == strings.Join(parts2[len(parts2)-2:], ".") {
				writer_nocdn := bufio.NewWriter(file_nocdn)
				defer writer_nocdn.Flush() // 延迟刷新缓冲区到函数结束
				writer_nocdn.WriteString(domain + "\n")
				return false, domain
			} else {
				writer_cdn := bufio.NewWriter(file_cdn)
				defer writer_cdn.Flush() // 延迟刷新缓冲区到函数结束
				writer_cdn.WriteString(domain + "\n")
				return true, domain
			}
		}
	}

	if ipCount == 1 {
		writer_nocdn := bufio.NewWriter(file_nocdn)
		defer writer_nocdn.Flush() // 延迟刷新缓冲区到函数结束
		writer_nocdn.WriteString(domain + "\n")
		return false, domain
	}
	writer_nocdn := bufio.NewWriter(file_nocdn)
	defer writer_nocdn.Flush() // 延迟刷新缓冲区到函数结束
	writer_nocdn.WriteString(domain + "\n")
	return false, domain
}

func CdnCheck(domains []string) []string {
	semaphore := make(chan struct{}, maxConcurrentIPScans)
	var wg sync.WaitGroup
	for _, domain := range domains {
		if strings.HasPrefix(domain, "http://") {
			domain = strings.TrimPrefix(domain, "http://")
		} else if strings.HasPrefix(domain, "https://") {
			domain = strings.TrimPrefix(domain, "https://")
		}
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			defer func() { <-semaphore }()
			semaphore <- struct{}{}
			processcdn(domain)
		}(domain)
	}
	wg.Wait()
	return NoCdndomains
}

func processcdn(domain string) {
	regex := regexp.MustCompile(`(?m)^([a-zA-Z0-9-_]+\.[a-zA-Z0-9-_]+\.[a-zA-Z0-9-_.]+)`)
	match := regex.FindStringSubmatch(domain)
	if len(match) > 1 {
		domain = match[1]
	}
	isCDN, domain := Cdn(domain)
	if !isCDN {
		NoCdndomains = append(NoCdndomains, domain)
		fmt.Println(domain, "is CDN:", isCDN)
	}
}
