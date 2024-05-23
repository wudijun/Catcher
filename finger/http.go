package finger

import (
	"crypto/tls"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func xegexpjs(reg string, resp string) (reslut1 [][]string) {
	reg1 := regexp.MustCompile(reg)
	if reg1 == nil {
		log.Println("regexp err")
		return nil
	}
	result1 := reg1.FindAllStringSubmatch(resp, -1)
	return result1
}

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

func Gettitle(httpbody string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(httpbody))
	if err != nil {
		return "Not found"
	}
	title := doc.Find("title").Text()
	title = strings.Replace(title, "\n", "", -1)
	title = strings.Trim(title, " ")
	return title
}

//	func Getfavicon(httpbody string, turl string) string {
//		faviconpaths := xegexpjs(`href="(.*?favicon....)"`, httpbody)
//		var faviconpath string
//		_, err := url.Parse(turl)
//		if err != nil {
//			panic(err)
//		}
//		turl = turl + "/"
//		if len(faviconpaths) > 0 {
//			fav := faviconpaths[0][1]
//			if fav[:2] == "//" {
//				faviconpath = "http:" + fav
//			} else {
//				if fav[:4] == "http" {
//					faviconpath = fav
//				} else {
//					faviconpath = turl + fav
//				}
//
//			}
//		} else {
//			faviconpath = turl + "favicon.ico"
//		}
//		return favicohash(faviconpath)
//	}
func Getfavicon(httpBody string, turl string) string {
	faviconPaths := xegexpjs(`href="(.*?favicon....)"`, httpBody)
	var faviconPath string
	_, err := url.Parse(turl)
	if err != nil {
		panic(err)
	}
	turl = strings.TrimRight(turl, "/") // 移除末尾的斜杠

	if len(faviconPaths) > 0 {
		fav := faviconPaths[0][1]
		if strings.HasPrefix(fav, "//") {
			faviconPath = "http:" + fav
		} else if strings.HasPrefix(fav, "http") {
			faviconPath = fav
		} else {
			faviconPath = turl + "/" + fav
		}
	} else {
		faviconPath = turl + "/favicon.ico"
	}
	return favicohash(faviconPath)
}

func Httprequest(domain string) (string, string, http.Header, int, int, string, string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 忽略证书验证
	}
	client := http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}
	//看域名是否有协议，如果是http则直接请求，如果不是http则判断是否是https
	if strings.HasPrefix(domain, "http://") {
		req, err := http.NewRequest("GET", domain, nil)
		if err != nil {
			//fmt.Println("1Error creating request:", err)
			return "", "", nil, 0, 0, "", ""
		}
		req.Header.Set("User-Agent", rndua())
		cookie := &http.Cookie{
			Name:  "rememberMe",
			Value: "me",
		}
		req.Header.Set("Accept", "*/*;q=0.8")
		req.Header.Set("Connection", "close")
		req.AddCookie(cookie)
		req.Header.Set("User-Agent", rndua())
		http_response, http_err := client.Do(req)
		if http_err != nil {
			return "", "", nil, 0, 0, "", ""
		}
		htmlContent, _ := ioutil.ReadAll(http_response.Body)
		htmlheaders := http_response.Header
		defer http_response.Body.Close()
		title := Gettitle(string(htmlContent))
		return title, string(htmlContent), htmlheaders, http_response.StatusCode, len(string(htmlContent)), Getfavicon(string(htmlContent), domain), ""
	} else if strings.HasPrefix(domain, "https://") {
		//看域名是否以https开头，如果是则直接请求
		req, err := http.NewRequest("GET", domain, nil)
		if err != nil {
			//fmt.Println("1Error creating request:", err)
			return "", "", nil, 0, 0, "", ""
		}
		req.Header.Set("User-Agent", rndua())
		cookie := &http.Cookie{
			Name:  "rememberMe",
			Value: "me",
		}
		req.Header.Set("Accept", "*/*;q=0.8")
		req.Header.Set("Connection", "close")
		req.AddCookie(cookie)
		req.Header.Set("User-Agent", rndua())
		http_response, http_err := client.Do(req)
		if http_err != nil {
			return "", "", nil, 0, 0, "", ""
		}
		htmlContent, _ := ioutil.ReadAll(http_response.Body)
		htmlheaders := http_response.Header
		defer http_response.Body.Close()
		title := Gettitle(string(htmlContent))
		return title, string(htmlContent), htmlheaders, http_response.StatusCode, len(string(htmlContent)), Getfavicon(string(htmlContent), domain), ""
	} else {
		//如果既不是以http开头也不是https开头则证明没有写，那么先请求https再请求http
		req, err := http.NewRequest("GET", "https://"+domain, nil)
		if err != nil {
			//fmt.Println("1Error creating request:", err)
			return "", "", nil, 0, 0, "", ""
		}
		req.Header.Set("User-Agent", rndua())
		cookie := &http.Cookie{
			Name:  "rememberMe",
			Value: "me",
		}
		req.Header.Set("Accept", "*/*;q=0.8")
		req.Header.Set("Connection", "close")
		req.AddCookie(cookie)
		req.Header.Set("User-Agent", rndua())
		http_response, http_err := client.Do(req)

		if http_err != nil {
			req, err := http.NewRequest("GET", "http://"+domain, nil)
			if err != nil {
				//fmt.Println("2Error creating request:", err)
				return "", "", nil, 0, 0, "", ""
			}
			cookie := &http.Cookie{
				Name:  "rememberMe",
				Value: "me",
			}
			req.Header.Set("Accept", "*/*;q=0.8")
			req.Header.Set("Connection", "close")
			req.AddCookie(cookie)
			req.Header.Set("User-Agent", rndua())
			http_response, http_err := client.Do(req)

			if http_err != nil {
				//fmt.Println("3Error creating request:", http_err)
				return "", "", nil, 0, 0, "", ""
			}
			//获取响应体
			htmlContent, _ := ioutil.ReadAll(http_response.Body)
			//fmt.Println(string(htmlContent))
			//获取响应头
			htmlheaders := http_response.Header

			defer http_response.Body.Close()
			title := Gettitle(string(htmlContent))
			return title, string(htmlContent), htmlheaders, http_response.StatusCode, len(string(htmlContent)), Getfavicon(string(htmlContent), "http://"+domain), "http://"
		}
		//读取响应的body内容
		htmlContent, _ := ioutil.ReadAll(http_response.Body)
		//fmt.Println(string(htmlContent))
		//获取响应头
		htmlheaders := http_response.Header
		defer http_response.Body.Close()
		title := Gettitle(string(htmlContent))
		return title, string(htmlContent), htmlheaders, http_response.StatusCode, len(string(htmlContent)), Getfavicon(string(htmlContent), "https://"+domain), "https://"
	}

}
