package nmap

import (
	"regexp"
	"strings"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/fatih/color"
	nmap "github.com/tomsteele/go-nmap"
)

//regex vars
var reCVE, _ = regexp.Compile(`([A-Z]+[\-]+[\d{4}]+[\-]+\d{4,5})`)
var reScore, _ = regexp.Compile(`(?:0|[1-9][0-9]*)\.[0-9]`)
var reURL, _ = regexp.Compile(`(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$`)

// Color support
var yellow = color.New(color.Bold, color.FgYellow).SprintFunc()
var red = color.New(color.Bold, color.FgRed).SprintFunc()
var cyan = color.New(color.Bold, color.FgCyan).SprintFunc()
var green = color.New(color.Bold, color.FgGreen).SprintFunc()

type CVEHost struct{
	Hostnames []nmap.Hostname
	ListCVE	  []CVE
}

 //CVE -- CPE, CVE number, CVSS base Score, and URL
 type CVE struct {
	CPE string
	ID	string
	Score string
	URL string
 }

//Hosts contains a slice of Host
type Hosts struct {
	List []nmap.Host
}

//GetNmapData returns the open ports and services of the host
func GetNmapData(path string) (result Hosts,c CVEHost){

	var cve CVE

	if path[len(path)-1:] != "/" {
		path = path + "/"
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		// Open nmap file
		nmapFile, err := ioutil.ReadFile(path + f.Name())
		if err != nil {
			fmt.Fprintf(color.Output, "%v Opening %s \n", red(" [-] ERROR: "), f.Name())
		}
		scan, err := nmap.Parse(nmapFile)
		if err != nil {
			log.Fatal(err)
		}
		for i, host := range scan.Hosts {
			c.Hostnames = host.Hostnames
			fmt.Fprintf(color.Output, "%v Host: %s IP: %s \n", cyan(" [i] INFO: "), host.Hostnames[i].Name, host.Addresses[i].Addr)
			for _, port := range host.Ports {
	
				fmt.Fprintf(color.Output, "%v Port: %d Service: %s Version: %s\n", cyan(" [i] INFO: "), port.PortId, port.Service.Name, port.Service.Product+" "+port.Service.Version)
				if len(port.Scripts) >0 {
					aux := strings.Split(port.Scripts[0].Output,"\n")
					for i,value:= range aux{
						if i == 1{
							cve.CPE = value[2:len(value)-2]
							fmt.Fprintf(color.Output, "%v CPE %v \n", cyan(" [i] INFO: "), value)
						} else if i !=0 || i >1{
							if reCVE.MatchString(value) {
								cve.ID = reCVE.FindAllString(value, 1)[0]
							}
							if reURL.MatchString(value) {
								cve.URL = reURL.FindAllString(value, -1)[0]
							}
							if reScore.MatchString(value) {
								cve.Score = reScore.FindAllString(value,-1)[0]
							}

							fmt.Fprintf(color.Output, "%v CVE %v \n", cyan(" [i] INFO: "), value)
					}
					if cve.ID !="" {
						c.ListCVE = append(c.ListCVE, cve)
					}
					
					}
					
				}
				
	
			}
			result.List = append(result.List, host)
		}
	}
	return result,c
}