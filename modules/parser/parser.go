package parser

import (
	"strings"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/fatih/color"
	nmap "github.com/tomsteele/go-nmap"
)

// Color support
var yellow = color.New(color.Bold, color.FgYellow).SprintFunc()
var red = color.New(color.Bold, color.FgRed).SprintFunc()
var cyan = color.New(color.Bold, color.FgCyan).SprintFunc()
var green = color.New(color.Bold, color.FgGreen).SprintFunc()

//Host contains the fields with the required information
type Host struct {
 	Hostnames []nmap.Hostname
 	CVElist   []CVE
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
func GetNmapData(path string) (result Hosts){

	var h Host
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
			h.Hostnames = host.Hostnames
			fmt.Fprintf(color.Output, "%v Host: %s IP: %s \n", cyan(" [i] INFO: "), host.Hostnames[i].Name, host.Addresses[i].Addr)
			h.Hostnames = host.Hostnames
			for _, port := range host.Ports {
	
				fmt.Fprintf(color.Output, "%v Port: %d Service: %s Version: %s\n", cyan(" [i] INFO: "), port.PortId, port.Service.Name, port.Service.Product+" "+port.Service.Version)
				if len(port.Scripts) >0 {
					aux := strings.Split(port.Scripts[0].Output,"\n")
					for i,value:= range aux{
						if i == 1{
							cve.CPE = value
							fmt.Fprintf(color.Output, "%v CPE %v \n", cyan(" [i] INFO: "), value)
						} else if i !=0 || i >1{
							fmt.Fprintf(color.Output, "%v CVE %v \n", cyan(" [i] INFO: "), value)
					}
					}
					
				}
	
			}
			//fmt.Println(h)
			result.List = append(result.List, host)
		}
	}
	return result
}

/*X returns the open ports and services of the host
func extractor(nmapFile []byte) (result Hosts) {

	scan, err := nmap.Parse(nmapFile)
	if err != nil {
		log.Fatal(err)
	}

	for i, host := range scan.Hosts {

		fmt.Fprintf(color.Output, "%v Host: %s IP: %s \n", cyan(" [i] INFO: "), host.Hostnames[i].Name, host.Addresses[i].Addr)

		for _, port := range host.Ports {

			fmt.Fprintf(color.Output, "%v Port: %d Service: %s Version: %s\n", cyan(" [i] INFO: "), port.PortId, port.Service.Name, port.Service.Product+" "+port.Service.Version)
			
			if len(port.Scripts) >0 {
				
				fmt.Fprintf(color.Output, "%v CVE's %v \n", cyan(" [i] INFO: "), port.Scripts[0].Output)
			}

		}
		result.List = append(result.List, host)
		//return host
	}
	fmt.Println(len(result.List))
	return result
}
*/