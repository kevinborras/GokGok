package parser

import (
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
// type Host struct {
// 	Hostname []string
// 	IP       string
// 	Ports    []map[int]string
// 	Status   string
// }

//Hosts contains a slice of Host
type Hosts struct {
	List []nmap.Host
}

//GetNmapData returns the open ports and services of the host
func GetNmapData(path string) {

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

		extractor(nmapFile)
	}

}

//X returns the open ports and services of the host
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

	}

	return result
}
