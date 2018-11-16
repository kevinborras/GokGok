package utils

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/fatih/color"
)

var timeout = time.Duration(3 * time.Second)
var client = http.Client{
	Timeout: timeout,
}
var path = "nmapResults"

// Color support
var yellow = color.New(color.Bold, color.FgYellow).SprintFunc()
var red = color.New(color.Bold, color.FgRed).SprintFunc()
var cyan = color.New(color.Bold, color.FgCyan).SprintFunc()
var green = color.New(color.Bold, color.FgGreen).SprintFunc()

//GetIP get the ip of the given Host
func getIP(host string) (ip string) {

	addr, err := net.LookupHost(host)
	if err != nil {
		log.Fatal(err)
	}
	ip = addr[0]
	return ip

}

// IsAlive checks if a domain is up or down
func IsAlive(domainList string) (ipList []string) {

	file, err := os.Open(domainList)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()
		resp, err := client.Get("http://" + url)
		if err != nil {
			fmt.Fprintf(color.Output, "%v No such host or it's down: %s \n", red(" [-] ERROR: "), url)

		} else {
			fmt.Fprintf(color.Output, "%v Host %s is alive with %s  \n", green(" [+] SUCCESS: "), url, resp.Status)

			ip := getIP(url)
			ipList = append(ipList, ip)

		}
	}
	return ipList

}

//RunNmap runs a nmap scan
func RunNmap(ipList []string) {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 666)
	}
	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	for _, ip := range ipList {
		if re.MatchString(ip) {
			fmt.Fprintf(color.Output, "%v Executing command %s \n", cyan(" [i] INFO: "), "nmap -sV -P0 "+ip+" -oX "+path+"/"+ip+".xml")
			err := exec.Command("nmap", "-sV", "-P0", ip, "-oX", path+"/"+ip+".xml").Run()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Fprintf(color.Output, "%v Scan done  \n", green(" [+] SUCCESS: "))
		} else {
			fmt.Println("IP chekc failed,trying with Host check")
			addr, err := net.LookupHost(ip)
			if err != nil {
				log.Fatal("Neither is an IP nor is a Host")
				os.Exit(1)
			}
			ip = addr[0]
			fmt.Fprintf(color.Output, "%v Executing command %s \n", cyan(" [i] INFO: "), "nmap -sV -P0 "+ip+" -oX "+path+"/"+ip+".xml")
			err = exec.Command("nmap", "-sV", "-P0", ip, "-oX", path+"/"+ip).Run()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Fprintf(color.Output, "%v Scan done  \n", green(" [+] SUCCESS: "))
		}

	}

}

// IsAlive checks a list of hosts hostList *parser.Hosts

// utils.RunNmap(ip)

// xmlFile, err := os.Open("result_" + ip + ".xml")
// if err != nil {
// 	fmt.Println(err)
// }
// defer xmlFile.Close()
// byteValue, err := ioutil.ReadAll(xmlFile)

// result := parser.GetNmapData(byteValue)
// result.Status = resp.Status
// hostList.Info = append(hostList.Info, result)
// fmt.Println(hostList.Info)
