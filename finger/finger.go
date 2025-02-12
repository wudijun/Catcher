package finger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"test/requestpoc"
)

// 公开的 ProcessValue
var ProcessValue int

// 设置 ProcessValue 的函数
func SetProcessValue(value int) {
	ProcessValue = value
}

// 获取 ProcessValue 的函数
func GetProcessValue() int {
	return ProcessValue
}

type Fingerprint struct {
	Cms      string   `json:"cms"`
	Location string   `json:"location"`
	Keyword  []string `json:"keyword"`
}

type FingerprintData struct {
	Fingerprint []Fingerprint `json:"fingerprint"`
}

var Fingerdomain []string
var NoFingerdomain []string
var count int

func processFingerprint(fp Fingerprint, domain string, title string, htmlContent string, htmlheaders http.Header, code int, length int, icoHash string, req string) {
	// 如果 title 为空，设置为 "空"
	if title == "" {
		title = "空"
	}
	for _, kw := range fp.Keyword {
		// 按逗号拆分关键字，将其视为“或”条件
		keywords := strings.Split(kw, ",")
		matched := false
		for _, subKw := range keywords {
			subKw = strings.TrimSpace(subKw) // 去除多余空格

			// 根据指纹位置进行匹配
			switch fp.Location {
			case "body":
				if strings.Contains(htmlContent, subKw) {
					matched = true
				}
			case "faviconhash":
				if icoHash == subKw {
					matched = true
				}
			case "header":
				if HeaderContains(htmlheaders, subKw) {
					matched = true
				}
			case "title":
				if strings.Contains(title, subKw) {
					matched = true
				}
			}
			// 如果匹配成功，输出结果并停止进一步的关键字检查
			if matched {
				count++
				s := fmt.Sprintf("%s%s|%d|%d|%s|%s\n", req, domain, code, length, title, fp.Cms)
				Fingerdomain = append(Fingerdomain, s)
				//请求poc
				if ProcessValue != 0 {
					requestpoc.File_poc(req+domain, fp.Cms)
				}
				return
			}
		}
	}
}

func Finger(domain string, title string, htmlContent string, htmlheaders http.Header, code int, length int, icoHash string, req string) {
	if title == "" {
		title = "空"
	}
	data, err := ioutil.ReadFile("finger.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	var fingerprintData FingerprintData
	if err := json.Unmarshal(data, &fingerprintData); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	count = 0
	for _, fp := range fingerprintData.Fingerprint {
		processFingerprint(fp, domain, title, htmlContent, htmlheaders, code, length, icoHash, req)
	}
	if count == 0 {
		s := fmt.Sprintf("%s%s|%d|%d|%s \n", req, domain, code, length, title)
		NoFingerdomain = append(NoFingerdomain, s)
	}
}

func HeaderContains(response http.Header, searchString string) bool {
	// 遍历响应头部
	for key, values := range response {
		// 查找特定的字符串
		if strings.Contains(strings.ToLower(key), strings.ToLower(searchString)) {
			return true
		}
		// 查找头部值中是否包含特定的字符串
		for _, value := range values {
			if strings.Contains(strings.ToLower(value), strings.ToLower(searchString)) {
				return true
			}
		}
	}
	return false
}
