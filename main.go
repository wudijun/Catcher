package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"test/cdncheck"
	"test/finger"
)

// 自定义帮助信息
func customUsage() {
	// 打印参数说明
	fmt.Println("Usage:")
	fmt.Println("  -f  是否执行指纹识别")
	fmt.Println("  -p  是否进行poc测试")
	fmt.Println("  -d  指定读取的文件路径，默认为domain.txt")
	fmt.Println("  example : -f -p -d domain.txt")
}

func main() {
	// 打印启动信息
	fmt.Println("   ___             _              _                      ")
	fmt.Println("  / __|   __ _    | |_     __    | |_      ___      _ _  ")
	fmt.Println(" | (__   / _` |   |  _|   / _|   | ' \\    / -_)    | '_| ")
	fmt.Println("  \\___|  \\__,_|   _|_|_   \\__|_  |_||_|   \\___|   _|_|_  ")
	fmt.Println("_|\"\"\"\"\"|_|\"\"\"\"\"|_|\"\"\"\"\"|_|\"\"\"\"\"|_|\"\"\"\"\"|_|\"\"\"\"\"|_|\"\"\"\"\"| ")
	fmt.Print("\033[32m") // 设置颜色为绿色
	fmt.Println("                                 Verson1.0 --- by:island")
	fmt.Print("\033[0m") // 重置颜色

	// 定义命令行参数
	fFlag := flag.Bool("f", false, "是否执行指纹识别")
	dFlag := flag.String("d", "domain.txt", "指定读取的文件路径")
	pFlag := flag.Bool("p", false, "是否要进行poc测试") // 使用 bool 类型

	flag.Parse()

	// 如果没有传入任何命令行参数，打印帮助信息
	if flag.NFlag() == 0 {
		customUsage()
		return
	}

	// 根据 -p 参数设置 ProcessValue 的值
	if *pFlag {
		finger.SetProcessValue(1) // 传入了 -p 参数，设置为 1
	} else {
		finger.SetProcessValue(0) // 没有传入 -p 参数，设置为 0
	}

	// 打开文件
	file, err := os.Open(*dFlag)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	// 使用 Scanner 逐行读取文件内容
	scanner := bufio.NewScanner(file)
	var domains []string
	for scanner.Scan() {
		line := scanner.Text()
		domains = append(domains, strings.Fields(line)...)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error scanning file:", err)
	}

	// 根据 -f 标志执行指纹识别
	if *fFlag {
		fmt.Print("\033[31m")
		fmt.Println("执行指纹识别...")
		fmt.Print("\033[0m")

		finger.Http_thread(domains)
		if *pFlag {
			fmt.Print("\033[31m")
			fmt.Println("执行POC...")
			fmt.Print("\033[0m")
		}
		//判断是否cdn
		cdndomains := cdncheck.CdnCheck(domains)
		//将未使用cdn进行ip获取
		iplist := cdncheck.Getip(cdndomains)
		//将获取到的ip进行去重
		cdncheck.UniqueSortedIPs(iplist)
	}

	// 输出结果
	red := "\033[31m"
	reset := "\033[0m"
	text := "此次结果保存在/results/" + finger.TimestampStr + "下"
	fmt.Println(red + text + reset)
}
