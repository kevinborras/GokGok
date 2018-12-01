package dnsdumpster

import (
	"github.com/kevinborras/GokGok/modules/parser/auxiliary"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var client = &http.Client{Timeout: time.Second * 10}

//regex for csrftoken
var re = regexp.MustCompile(`^csrftoken=([\S\s]{32}\;)`)

func GetMapFromDumpster(domain string) (dnsD auxiliary.Domain) {
	csrfToken, cookie := getCSRFToken()
	subdomains := make(map[string]bool)

	data := url.Values{}
	data.Set("csrfmiddlewaretoken", csrfToken)
	data.Set("targetip", domain)

	reqPost, err := http.NewRequest("POST", "https://dnsdumpster.com/", strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	reqPost.Header.Set("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0`)
	reqPost.Header.Set("content-type", `application/x-www-form-urlencoded`)
	reqPost.Header.Set("referer", "https://dnsdumpster.com/")
	reqPost.Header.Set("cookie", cookie)

	respPost, err := client.Do(reqPost)
	if err != nil {
		log.Fatal(err)
	}
	defer respPost.Body.Close()

	// Here we have the response
	body, err := ioutil.ReadAll(respPost.Body)
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(`"col-md-4">[A-Za-z0-9]+\.(` + domain + `)<br>`)
	re2 := regexp.MustCompile(`"col-md-4">(.*?)<br>`)
	subdomainsList := re.FindAllString(string(body), -1)
	for _, v := range subdomainsList {
		subdomains[re2.FindStringSubmatch(v)[1]] = true
		dnsD.Subdomains = subdomains
	}
	// fmt.Println(subdomainsList)
	dnsD.Domain = domain
	return dnsD
}

//Get the CSRF token from dnsdumpster.com
func getCSRFToken() (csrftoken, cookie string) {

	req, err := http.NewRequest("GET", "https://dnsdumpster.com/", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	cookie = resp.Header.Get("Set-Cookie")
	csrftoken = re.FindStringSubmatch(cookie)[1]
	return csrftoken[:len(csrftoken)-1], cookie

}
