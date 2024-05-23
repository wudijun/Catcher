package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"test/finger"
)

func main() {
	fmt.Println("   ___             _              _                      ")
	fmt.Println("  / __|   __ _    | |_     __    | |_      ___      _ _  ")
	fmt.Println(" | (__   / _` |   |  _|   / _|   | ' \\    / -_)    | '_| ")
	fmt.Println("  \\___|  \\__,_|   _|_|_   \\__|_  |_||_|   \\___|   _|_|_  ")
	fmt.Println("_|\"\"\"\"\"|_|\"\"\"\"\"|_|\"\"\"\"\"|_|\"\"\"\"\"|_|\"\"\"\"\"|_|\"\"\"\"\"|_|\"\"\"\"\"| ")
	fmt.Print("\033[32m") // 设置颜色为绿色
	fmt.Println(" Verson1.0 --- by:1sl4nd ")
	fmt.Print("\033[0m") // 重置颜色
	fmt.Println("\n")

	// 打开文本文件
	file, err := os.Open("domain.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// 创建一个 Scanner 来扫描文件中的内容
	scanner := bufio.NewScanner(file)

	var domains []string

	// 逐行扫描文件内容
	for scanner.Scan() {
		line := scanner.Text()
		domains = append(domains, strings.Fields(line)...)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return
	}

	//指纹识别，poc测试
	finger.Http_thread(domains)
	//判断是否cdn
	cdndomains := cdncheck.CdnCheck(domains)
	//将未使用cdn进行ip获取
	iplist := cdncheck.Getip(cdndomains)
	//将获取到的ip进行去重
	cdncheck.UniqueSortedIPs(iplist)
	red := "\033[31m"
	reset := "\033[0m"
	text := "此次结果保存在/results/" + finger.TimestampStr + "下"
	fmt.Println(red + text + reset)
}
