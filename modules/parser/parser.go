package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

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
	Hostname []string
	IP       string
	Ports    []map[int]string
	Status   string
}

//Hosts contains a slice of Host
type Hosts struct {
	Info []Host
}

//GetNmapData returns the open ports and services of the host
func GetNmapData(path string) {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		// Open nmap file
		nmapFile, err := os.Open(f.Name())
		if err != nil {
			fmt.Fprintf(color.Output, "%v Opening: %s \n", red(" [-] ERROR: "), f.Name())
		}
		byteNmap, _ := ioutil.ReadAll(nmapFile)
		println(string(byteNmap))
		extractor(byteNmap)
	}

}

//X returns the open ports and services of the host
func extractor(nmapFile []byte) (h Host) {

	//var aux []map[int]string
	//var portInfo map[int]string
	scan, err := nmap.Parse(nmapFile)
	if err != nil {
		log.Fatal(err)
	}
	for i, host := range scan.Hosts {
		h.Hostname = append(h.Hostname, host.Hostnames[i].Name)
		for _, ip := range host.Addresses {
			h.IP = ip.Addr
		}
		for _, port := range host.Ports {
			fmt.Println(port.PortId)
			fmt.Println(port.Service.Name)
			//portInfo[port.PortId] = port.Service.Name + " " + port.Service.Product + " " + port.Service.Version
			//aux = append(aux, portInfo)

		}

	}
	//h.Ports = aux
	return h
}
