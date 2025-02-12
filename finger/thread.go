package finger

import (
	"bufio"
	"fmt"
	"github.com/tealeg/xlsx"
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

const maxConcurrentIPScans = 25

func init() {
	timestamp := time.Now().Unix()
	TimestampStr = fmt.Sprintf("%d", timestamp)
}

func processDomain(domain string, wg *sync.WaitGroup, semaphore chan struct{}) {
	defer wg.Done()
	defer func() { <-semaphore }() // 释放一个信号量
	semaphore <- struct{}{}        // 占用一个信号量

	title, htmlContent, htmlheaders, code, length, icoHash, req := Httprequest(domain)

	// 替换标题中的 "|" 符号
	title = strings.ReplaceAll(title, "|", "&#124;")
	fmt.Printf("[%s%s |title: \033[31m%s\033[0m |响应码: %d |返回长度: %d ]\n", req, domain, title, code, length)
	if code == 0 {
		return
	}
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

func writeToExcel(filename string, data []string, includeFingerprint bool) error {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Results")
	if err != nil {
		return err
	}

	// 设置样式
	normalStyle := xlsx.NewStyle()
	normalStyle.Font = *xlsx.NewFont(10, "Arial")

	redStyle := xlsx.NewStyle()
	redStyle.Font = *xlsx.NewFont(10, "Arial")
	redStyle.Font.Color = "FF0000" // 红色字体

	row := sheet.AddRow()
	row.AddCell().Value = "URL"
	row.AddCell().Value = "响应码"
	row.AddCell().Value = "返回长度"
	row.AddCell().Value = "Title"
	if includeFingerprint {
		row.AddCell().Value = "指纹"
	}

	for _, result := range data {
		row := sheet.AddRow()
		parts := strings.Split(result, "|")
		if len(parts) < 4 {
			continue
		}
		row.AddCell().Value = strings.TrimSpace(parts[0]) // URL
		row.AddCell().Value = strings.TrimSpace(parts[1]) // 响应码
		row.AddCell().Value = strings.TrimSpace(parts[2]) // 返回长度
		// 还原标题中的 "|" 符号
		title := strings.ReplaceAll(strings.TrimSpace(parts[3]), "&#124;", "|")
		row.AddCell().Value = title // Title

		// 判断指纹内容并设置字体颜色
		if includeFingerprint && len(parts) > 4 {
			fingerprint := strings.TrimSpace(parts[4])
			cell := row.AddCell()
			cell.Value = fingerprint

			if fingerprint != "疑似登陆页面" {
				// 如果不是 "疑似登陆页面"，则使用红色字体
				cell.SetStyle(redStyle)
			} else {
				// 否则使用正常字体
				cell.SetStyle(normalStyle)
			}
		}
	}
	return file.Save(filename)
}

func domainresults() {
	err := os.MkdirAll("results/"+TimestampStr, 0755)
	if err != nil {
		fmt.Println(err)
	}

	file_poc, err_poc := os.OpenFile("results/"+TimestampStr+"/PocResults.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err_poc != nil {
		fmt.Println(err_poc)
		fmt.Println("创建文件失败")
	}
	defer file_poc.Close()
	writer_poc := bufio.NewWriter(file_poc)

	VulnMapMutex.Lock()
	for url, vulnerabilities := range requestpoc.VulnMap {
		writer_poc.WriteString("URL: " + url + "\n")
		writer_poc.WriteString("漏洞列表: ")
		for i, vulnerability := range vulnerabilities {
			writer_poc.WriteString(vulnerability)
			if i < len(vulnerabilities)-1 {
				writer_poc.WriteString(", ")
			}
		}
		writer_poc.WriteString("\n\n")
	}
	VulnMapMutex.Unlock()
	writer_poc.Flush()

	sort.Slice(NoFingerdomain, func(i, j int) bool {
		partsI := strings.Split(NoFingerdomain[i], "|")
		partsJ := strings.Split(NoFingerdomain[j], "|")
		if len(partsI) < 4 || len(partsJ) < 4 {
			return false
		}
		lengthI, errI := strconv.Atoi(partsI[3])
		lengthJ, errJ := strconv.Atoi(partsJ[3])
		if errI != nil || errJ != nil {
			return false
		}
		return lengthI > lengthJ
	})

	err = writeToExcel("results/"+TimestampStr+"/NoFinger.xlsx", NoFingerdomain, false)
	if err != nil {
		fmt.Println("写入 NoFinger.xlsx 错误:", err)
		return
	}

	sort.Slice(Fingerdomain, func(i, j int) bool {
		partsI := strings.Split(Fingerdomain[i], "|")
		partsJ := strings.Split(Fingerdomain[j], "|")
		if len(partsI) < 4 || len(partsJ) < 4 {
			return false
		}
		lengthI, errI := strconv.Atoi(partsI[3])
		lengthJ, errJ := strconv.Atoi(partsJ[3])
		if errI != nil || errJ != nil {
			return false
		}
		return lengthI > lengthJ
	})

	err = writeToExcel("results/"+TimestampStr+"/Finger.xlsx", Fingerdomain, true)
	if err != nil {
		fmt.Println("写入 Finger.xlsx 错误:", err)
		return
	}
}
