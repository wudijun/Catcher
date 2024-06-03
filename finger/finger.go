package finger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"test/requestpoc"
)

type Fingerprint struct {
	Cms      string   `json:"cms"`
	Method   string   `json:"method"`
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
	//defer wg.Done()
	if title == "" {
		title = "空"
	}
	if fp.Location == "body" {
		// 检查指纹对象关键字是否与 HTML 内容匹配
		for _, kw := range fp.Keyword {
			if strings.Contains(htmlContent, kw) {
				count++
				s := fmt.Sprintf("url: \"%s%s\" |响应码: %d |返回长度: %d |title: \"%s\" | 指纹: %s]\n", req, domain, code, length, title, fp.Cms)
				Fingerdomain = append(Fingerdomain, s)
				requestpoc.File_poc(req+domain, fp.Cms)
				break
			}
		}
	} else if fp.Location == "faviconhash" {
		if icoHash == fp.Keyword[0] {
			count++
			s := fmt.Sprintf("url: \"%s%s\" |响应码: %d |返回长度: %d |title: \"%s\" | 指纹: %s]\n", req, domain, code, length, title, fp.Cms)
			Fingerdomain = append(Fingerdomain, s)
			requestpoc.File_poc(req+domain, fp.Cms)
		}

	} else if fp.Location == "header" {
		if HeaderContains(htmlheaders, fp.Keyword[0]) {
			count++
			s := fmt.Sprintf("url: \"%s%s\" |响应码: %d |返回长度: %d |title: \"%s\" | 指纹: %s]\n", req, domain, code, length, title, fp.Cms)
			Fingerdomain = append(Fingerdomain, s)
			requestpoc.File_poc(req+domain, fp.Cms)
		}
	} else if fp.Location == "title" {
		if strings.Contains(title, fp.Keyword[0]) {
			count++
			s := fmt.Sprintf("url: \"%s%s\" |响应码: %d |返回长度: %d |title: \"%s\" | 指纹: %s]\n", req, domain, code, length, title, fp.Cms)
			Fingerdomain = append(Fingerdomain, s)
			requestpoc.File_poc(req+domain, fp.Cms)
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

	//var wg sync.WaitGroup
	count = 0
	for _, fp := range fingerprintData.Fingerprint {
		//wg.Add(1)
		processFingerprint(fp, domain, title, htmlContent, htmlheaders, code, length, icoHash, req)
	}
	if count == 0 {
		s := fmt.Sprintf("url: \"%s%s\" |响应码: %d |返回长度: %d |title: \"%s\" \n", req, domain, code, length, title)
		NoFingerdomain = append(NoFingerdomain, s)
	}
	//wg.Wait()
	//等待所有协程结束，也就是所有指纹全都匹配结束如果计数器没有增加则一个都没匹配成功
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
