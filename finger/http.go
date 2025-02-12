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

func xegexpjs(reg string, resp string) (result1 [][]string) {
	reg1 := regexp.MustCompile(reg)
	if reg1 == nil {
		log.Println("regexp err")
		return nil
	}
	result1 = reg1.FindAllStringSubmatch(resp, -1)
	return result1
}

func rndua() string {
	ua := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 YaBrowser/22.1.0.2517 Yowser/2.5 Safari/537.36",
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
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 YaBrowser/22.1.0.2517 Yowser/2.5 Safari/537.36",
	}
	n := rand.Intn(len(ua))
	return ua[n]
}

func Gettitle(httpbody string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(httpbody))
	if err != nil {
		return "Not found"
	}
	title := doc.Find("title").Text()
	title = strings.ReplaceAll(title, "\n", "")
	title = strings.TrimSpace(title)
	return title
}

func Getfavicon(httpBody string, turl string) string {
	faviconPaths := xegexpjs(`href="(.*?favicon....)"`, httpBody)
	var faviconPath string
	_, err := url.Parse(turl)
	if err != nil {
		//panic(err)
	}
	turl = strings.TrimRight(turl, "/") // 移除末尾的斜杠
	if len(faviconPaths) > 0 {
		fav := faviconPaths[0][1]
		if strings.HasPrefix(fav, "//") {
			faviconPath = "http:" + fav // 或者 "https:"，取决于需求
		} else if strings.HasPrefix(fav, "http") {
			faviconPath = fav
		} else {
			// 确保拼接时没有多余的斜杠
			faviconPath = strings.TrimRight(turl, "/") + "/" + strings.TrimLeft(fav, "/")
		}
	} else {
		faviconPath = strings.TrimRight(turl, "/") + "/favicon.ico"
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

	makeRequest := func(url string) (string, http.Header, int, int) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return "", nil, 0, 0
		}
		req.Header.Set("User-Agent", rndua())
		req.Header.Set("Accept", "*/*;q=0.8")
		req.Header.Set("Connection", "close")
		//cookie := &http.Cookie{
		//	Name:  "rememberMe",
		//	Value: "me",
		//}
		//req.AddCookie(cookie)
		httpResponse, httpErr := client.Do(req)
		if httpErr != nil {
			return "", nil, 0, 0
		}
		htmlContent, _ := ioutil.ReadAll(httpResponse.Body)
		defer httpResponse.Body.Close()
		return string(htmlContent), httpResponse.Header, httpResponse.StatusCode, len(htmlContent)
	}

	parseAndFollowRedirect := func(htmlContent, baseURL string) (string, http.Header, int, int) {
		redirectPatterns := []string{
			`window\.parent\.location\.href\s*=\s*['"](.*?)['"];`,
			`window\.location\.href\s*=\s*['"](.*?)['"];`,
			`window\.top\.location\s*=\s*['"](.*?)['"];`,
			`window\.location\s*=\s*['"](.*?)['"];`,
			`location\.href\s*=\s*['"](.*?)['"];`,
			`location\s*=\s*['"](.*?)['"];`,
			`eval\("window\.".*?\.location\s*=\s*['"](.*?)['"]\);`,
			`<meta\s+http-equiv=["']refresh["']\s+content=["']\d*;\s*url=(.*?)["']`,
		}

		for _, pattern := range redirectPatterns {
			redirectRegex := regexp.MustCompile(pattern)
			matches := redirectRegex.FindStringSubmatch(htmlContent)
			if len(matches) > 1 {
				redirectURL := matches[1]
				if !strings.HasPrefix(redirectURL, "http") {
					base, err := url.Parse(baseURL)
					if err != nil {
						return "", nil, 0, 0
					}
					redirectURL = base.ResolveReference(&url.URL{Path: redirectURL}).String()
				}
				return makeRequest(redirectURL)
			}
		}
		return htmlContent, nil, 0, 0
	}

	var htmlContent string
	var htmlHeader http.Header
	var statusCode, contentLength int

	if strings.HasPrefix(domain, "http://") || strings.HasPrefix(domain, "https://") {
		htmlContent, htmlHeader, statusCode, contentLength = makeRequest(domain)
	} else {
		htmlContent, htmlHeader, statusCode, contentLength = makeRequest("https://" + domain)
		if statusCode == 0 {
			htmlContent, htmlHeader, statusCode, contentLength = makeRequest("http://" + domain)
		}
	}

	title := Gettitle(htmlContent)
	favicon := Getfavicon(htmlContent, domain)

	if statusCode == 0 {
		// 处理跳转
		htmlContent, htmlHeader, statusCode, contentLength = parseAndFollowRedirect(htmlContent, domain)
		title = Gettitle(htmlContent)
		favicon = Getfavicon(htmlContent, domain)
	}
	//print(htmlContent)
	return title, htmlContent, htmlHeader, statusCode, contentLength, favicon, ""
}
