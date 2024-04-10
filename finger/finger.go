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
	if fp.Location == "body" {
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

		}
	}

}

func Finger(domain string, title string, htmlContent string, htmlheaders http.Header, code int, length int, icoHash string, req string) {
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
		s := fmt.Sprintf("url: \"%s%s\" |响应码: %d |返回长度: %d |title: \"%s\" \n", req, domain, code, length, title)
		NoFingerdomain = append(NoFingerdomain, s)
	}
}

func HeaderContains(response http.Header, searchString string) bool {
	for key, values := range response {
		if strings.Contains(strings.ToLower(key), strings.ToLower(searchString)) {
			return true
		}
		for _, value := range values {
			if strings.Contains(strings.ToLower(value), strings.ToLower(searchString)) {
				return true
			}
		}
	}
	return false
}
