package crtsh

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

var client = &http.Client{Timeout: time.Second * 15}

var re = regexp.MustCompile(`\?id=[0-9]+`)

func GetMapfromCRT(domain string) map[string]bool {
	subdomains := make(map[string]bool)
	req, err := http.NewRequest("GET", "https://crt.sh/?q="+domain, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// we need to get ?id=XXXXX
	res := re.FindAllString(string(body), -1)

	//once we have this id, we iterate over the results and extract only the subdomains
	for i, id := range res {
		req, err := http.NewRequest("GET", "https://crt.sh/"+id+"&output=json", nil)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		//extracting the subdomains
		subD := regexp.MustCompile(`[A-Za-z0-9]+\.(` + domain + `)`)
		subdomainsList := subD.FindAllString(string(body), -1)
		//if we have something proceed to check for duplicities
		if i == 0 {
			for i := 0; i < len(subdomainsList); i++ {
				subdomains[subdomainsList[i]] = true
			}
		} else if len(subdomainsList) > 0 {
			for _, d := range subdomainsList {
				if _, ok := subdomains[d]; !ok {
					subdomains[d] = true
				}
			}
		}

	}
	return subdomains
}
