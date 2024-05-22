package requestpoc

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func rndua() string {
	ua := []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 YaBrowser/22.1.0.2517 Yowser/2.5 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
		"Mozilla/5.0 (X11; Linux x86_64; rv:96.0) Gecko/20100101 Firefox/96.0",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64; rv:95.0) Gecko/20100101 Firefox/95.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:96.0) Gecko/20100101 Firefox/96.0",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.3 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 YaBrowser/22.1.0.2517 Yowser/2.5 Safari/537.36"}
	n := rand.Intn(13) + 1
	return ua[n]
}

type RequestData struct {
	ID   string `yaml:"id"`
	Info struct {
		Name        string `yaml:"name"`
		Severity    int    `yaml:"severity"`
		Description string `yaml:"description"`
		Details     string `yaml:"details"`
		Repair      string `yaml:"repair"`
	} `yaml:"info"`
	Request1 struct {
		Method  string            `yaml:"method"`
		Path    string            `yaml:"path"`
		Headers map[string]string `yaml:"headers"`
		Data    string            `yaml:"data"`
	} `yaml:"requests1"`
	Response1 struct {
		Headers map[string]string `yaml:"headers"`
		Body    string            `yaml:"body"`
		Code    int               `yaml:"code"`
	} `yaml:"response1"`
	Request2 struct {
		Method  string            `yaml:"method"`
		Path    string            `yaml:"path"`
		Headers map[string]string `yaml:"headers"`
		Data    string            `yaml:"data"`
	} `yaml:"requests2"`
	Response2 struct {
		Headers map[string]string `yaml:"headers"`
		Body    string            `yaml:"body"`
		Code    int               `yaml:"code"`
	} `yaml:"response2"`
}

var VulnMap = make(map[string][]string)

// 打开yaml文件
func File_poc(url string, cms string) {
	dir, err := os.Open("poc\\" + cms)
	if err != nil {
		//fmt.Println("无 \"" + cms + "\" poc文件，不做poc测试")
		return
	}
	defer dir.Close()
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println("读取文件夹内容失败:", err)
		return
	}
	files := make([]string, 0)
	for _, fileInfo := range fileInfos {
		files = append(files, fileInfo.Name())
	}
	var wg sync.WaitGroup
	//多线程读取
	for _, file := range files {
		wg.Add(1)
		go processYaml(url, file, cms, &wg)
	}
	wg.Wait()
}

// 多线程读取yaml
func processYaml(url string, file string, cms string, wg *sync.WaitGroup) {
	defer wg.Done()
	yamlData, err := ioutil.ReadFile("poc\\" + cms + "\\" + file)
	if err != nil {
		fmt.Println("读取 YAML 文件失败:", err)
		return
	}
	// 解析 YAML 数据
	var requestData RequestData
	err = yaml.Unmarshal(yamlData, &requestData)
	if err != nil {
		fmt.Println("解析 YAML 数据失败:", err)
		return
	}
	//进行请求
	request1_poc(url, requestData)
	if requestData.Request2.Method != "" {
		request2_poc(url, requestData)
	}

}

// 这里构造http请求的函数，直接将请求url和从yaml获取值后的结构体传入
func request1_poc(url string, requestData RequestData) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 忽略证书验证
	}
	client := http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}
	// 将字符串转换为 Reader 接口类型
	body := strings.NewReader(requestData.Request1.Data)

	req, err := http.NewRequest(requestData.Request1.Method, url+requestData.Request1.Path, body)
	for key, value := range requestData.Request1.Headers {
		req.Header.Set(key, value)
	}
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", rndua())
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送 HTTP 请求失败：%v\n", err)
		return
	}

	htmlContent, err := ioutil.ReadAll(response.Body)
	//fmt.Println(string(htmlContent))

	if err != nil {
		fmt.Printf("获取响应体失败：%v\n", err)
		return
	}

	//判断yaml文件中是否有第一次请求返回的特征，如果没有则就是要去请求第二次，就直接去第二次请求。如果有一个有值那么就进if
	if requestData.Response1.Body != "" || requestData.Response1.Code != 0 || requestData.Response1.Headers != nil {
		//当yaml文件只有返回体有值时
		if requestData.Response1.Body != "" && requestData.Response1.Code == 0 && requestData.Response1.Headers == nil {
			if strings.Contains(string(htmlContent), requestData.Response1.Body) {
				VulnMap[url] = append(VulnMap[url], requestData.Info.Name)
			}
			//当yaml文件只有返回值有值时
		} else if requestData.Response1.Body == "" && requestData.Response1.Code != 0 && requestData.Response1.Headers == nil {
			if response.StatusCode == requestData.Response1.Code {
				VulnMap[url] = append(VulnMap[url], requestData.Info.Name)
			}
			//当yaml文件只有返回头有值时
		} else if requestData.Response1.Body == "" && requestData.Response1.Code == 0 && requestData.Response1.Headers != nil {
			if HeaderYaml(response.Header, requestData.Response1.Headers) {
				VulnMap[url] = append(VulnMap[url], requestData.Info.Name)
			}
			//当三个都有值
		} else if requestData.Response1.Body != "" && requestData.Response1.Code != 0 && requestData.Response1.Headers != nil {
			if (strings.Contains(string(htmlContent), requestData.Response1.Body)) && (response.StatusCode == requestData.Response1.Code) && (HeaderYaml(response.Header, requestData.Response1.Headers)) {
				VulnMap[url] = append(VulnMap[url], requestData.Info.Name)
			}
			//当yaml文件返回体和返回码有值时
		} else if requestData.Response1.Body != "" && requestData.Response1.Code != 0 && requestData.Response1.Headers == nil {
			if (strings.Contains(string(htmlContent), requestData.Response1.Body)) && (response.StatusCode == requestData.Response1.Code) {
				VulnMap[url] = append(VulnMap[url], requestData.Info.Name)
			}
			//当yaml文件返回体和返回头有值时
		} else if requestData.Response1.Body != "" && requestData.Response1.Headers != nil && requestData.Response1.Code == 0 {
			if (strings.Contains(string(htmlContent), requestData.Response1.Body)) && (HeaderYaml(response.Header, requestData.Response1.Headers)) {
				VulnMap[url] = append(VulnMap[url], requestData.Info.Name)
			}
			//当yaml文件返回头和响应码有值时
		} else if requestData.Response1.Headers != nil && requestData.Response1.Code != 0 && requestData.Response1.Body == "" {
			if (response.StatusCode == requestData.Response1.Code) && (HeaderYaml(response.Header, requestData.Response1.Headers)) {
				VulnMap[url] = append(VulnMap[url], requestData.Info.Name)
			}
		}

	}
	defer response.Body.Close()

}

func request2_poc(url string, requestData RequestData) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 忽略证书验证
	}
	client := http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}
	// 将字符串转换为 Reader 接口类型
	body := strings.NewReader(requestData.Request2.Data)

	req, err := http.NewRequest(requestData.Request2.Method, url+requestData.Request2.Path, body)
	for key, value := range requestData.Request2.Headers {
		req.Header.Set(key, value)
	}
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", rndua())
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送 HTTP 请求失败：%v\n", err)
		return
	}
	htmlContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("获取响应体失败：%v\n", err)
		return
	}
	//fmt.Println(string(htmlContent))
	if requestData.Response2.Body != "" || requestData.Response2.Code != 0 || requestData.Response2.Headers != nil {
		//当yaml文件只有返回体有值时
		if requestData.Response2.Body != "" && requestData.Response2.Code == 0 && requestData.Response2.Headers == nil {
			if strings.Contains(string(htmlContent), requestData.Response2.Body) {
				VulnMap[url] = append(VulnMap[url], requestData.Info.Name)
			}
			//当yaml文件只有返回值有值时
		} else if requestData.Response2.Body == "" && requestData.Response2.Code != 0 && requestData.Response2.Headers == nil {
			if response.StatusCode == requestData.Response2.Code {
				VulnMap[url] = append(VulnMap[url], requestData.Info.Name)
			}
			//当yaml文件只有返回头有值时
		} else if requestData.Response2.Body == "" && requestData.Response2.Code == 0 && requestData.Response2.Headers != nil {
			if HeaderYaml(response.Header, requestData.Response2.Headers) {
				VulnMap[url] = append(VulnMap[url], requestData.Info.Name)
			}
			//当三个都有值
		} else if requestData.Response2.Body != "" && requestData.Response2.Code != 0 && requestData.Response2.Headers != nil {
			if (strings.Contains(string(htmlContent), requestData.Response2.Body)) && (response.StatusCode == requestData.Response2.Code) && (HeaderYaml(response.Header, requestData.Response2.Headers)) {
				VulnMap[url] = append(VulnMap[url], requestData.Info.Name)
			}
			//当yaml文件返回体和返回码有值时
		} else if requestData.Response2.Body != "" && requestData.Response2.Code != 0 && requestData.Response2.Headers == nil {
			if (strings.Contains(string(htmlContent), requestData.Response2.Body)) && (response.StatusCode == requestData.Response2.Code) {
				VulnMap[url] = append(VulnMap[url], requestData.Info.Name)
			}
			//当yaml文件返回体和返回头有值时
		} else if requestData.Response2.Body != "" && requestData.Response2.Headers != nil && requestData.Response2.Code == 0 {
			if (strings.Contains(string(htmlContent), requestData.Response2.Body)) && (HeaderYaml(response.Header, requestData.Response2.Headers)) {
				VulnMap[url] = append(VulnMap[url], requestData.Info.Name)
			}
			//当yaml文件返回头和响应码有值时
		} else if requestData.Response2.Headers != nil && requestData.Response2.Code != 0 && requestData.Response2.Body == "" {
			if (response.StatusCode == requestData.Response2.Code) && (HeaderYaml(response.Header, requestData.Response2.Headers)) {
				VulnMap[url] = append(VulnMap[url], requestData.Info.Name)
			}
		}

	}
	defer response.Body.Close()

}

// 检查响应头中是否包含yaml文件自定义头的内容
func HeaderYaml(responseHeader http.Header, yamlHeader map[string]string) bool {
	for key, value := range yamlHeader {
		// 如果返回头中有这个键那么就把对应的值赋给headerValue，并且headerValue不为空就进入if
		if headerValue := responseHeader.Get(key); headerValue != "" {
			//当返回头和yaml文件中都有这个键的时候，如果返回头中这个键的值包含yaml文件中这个键的值那么就匹配成功
			if strings.Contains(headerValue, value) {
				// 如果任何一个值包含了YAML格式头部中的值，返回true
				return true
			}
		}
	}
	// 如果没有找到匹配，返回false
	return false
}
