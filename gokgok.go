package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/integrii/flaggy"
	"github.com/kevinborras/GokGok/modules/parser/auxiliary"
	crtsh "github.com/kevinborras/GokGok/modules/parser/crtsh"
	dnsDumpster "github.com/kevinborras/GokGok/modules/parser/dnsdumpster"
	parser "github.com/kevinborras/GokGok/modules/parser/nmap"
	"github.com/kevinborras/GokGok/modules/utils"
	"github.com/kevinborras/GokGok/modules/webapp"
	"os"
)

// Color support
var yellow = color.New(color.Bold, color.FgYellow).SprintFunc()
var red = color.New(color.Bold, color.FgRed).SprintFunc()
var cyan = color.New(color.Bold, color.FgCyan).SprintFunc()
var green = color.New(color.Bold, color.FgGreen).SprintFunc()

//VERSION of the program
var VERSION = `v0.0.9b`

var targetList, parseFiles, domain string
var version, scanThem, html, subdomains = false, false, false, false

func init() {

	//gokgok Flags
	flaggy.SetName("Gok-Gok")
	flaggy.SetDescription("Is something there?")
	flaggy.DefaultParser.ShowVersionWithVersionFlag = false

	flaggy.Bool(&version, "v", "version", "Print version")
	flaggy.String(&targetList, "t", "targetList", "File with targets to be checked")
	flaggy.String(&domain, "d", "domain", "Domain is used with -d")
	flaggy.Bool(&subdomains, "sd", "subdomains", "Enumerate subdomains")
	flaggy.Bool(&scanThem, "s", "scanThem", "Scan the the targets with Nmap")
	flaggy.String(&parseFiles, "p", "parseFiles", "Parse the nmap resut files, needs the path of the scans")
	flaggy.Bool(&html, "o", "html", "HTML output")
}

func main() {

	flaggy.Parse()

	if version == true {
		fmt.Fprintf(color.Output, "%v %v\n", cyan(" [i] Gok-Gok "), VERSION)
		os.Exit(0)
	}

	if targetList != "" && !scanThem && !subdomains {
		utils.IsAlive(targetList)

	} else if targetList != "" && scanThem && !subdomains {
		ipList := utils.IsAlive(targetList)
		utils.RunNmap(ipList)

	} else if targetList != "" && subdomains {
		resCrt := make(chan auxiliary.Domain)
		resDnsD := make(chan auxiliary.Domain)
		// var aux auxiliary.Domain
		file, err := os.Open(targetList)
		if err != nil {
			fmt.Fprintf(color.Output, red(" [-] ERROR: "), err)
		}
		defer file.Close()
		fmt.Fprintf(color.Output, "%v Checking for subdomains \n", cyan(" [i] INFO: "))
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			domain := scanner.Text()
			fmt.Fprintf(color.Output, "%v %s \n", cyan(" [i] DOMAIN: "), domain)

			go crtsh.GetMapfromCRT(domain, resCrt)
			go dnsDumpster.GetMapFromDumpster(domain, resDnsD)
			subCRT, subDNSD := <-resCrt, <-resDnsD

			for k, _ := range subCRT.Subdomains {
				fmt.Fprintf(color.Output, "%v %s  \n", green(" [+] "+subCRT.Source+":"), k)
			}
			for k, _ := range subDNSD.Subdomains {
				fmt.Fprintf(color.Output, "%v %s  \n", green(" [+] "+subDNSD.Source+":"), k)
			}
		}

	}
	if domain != "" && subdomains {
		resCrt := make(chan auxiliary.Domain)
		resDnsD := make(chan auxiliary.Domain)
		fmt.Fprintf(color.Output, "%v Checking for subdomains of %s \n", cyan(" [i] INFO: "), domain)

		go crtsh.GetMapfromCRT(domain, resCrt)
		go dnsDumpster.GetMapFromDumpster(domain, resDnsD)

		subCRT, subDNSD := <-resCrt, <-resDnsD
		fmt.Fprintf(color.Output, "%v %s \n", cyan(" [i] DOMAIN: "), subCRT.Domain)
		for k, _ := range subCRT.Subdomains {
			fmt.Fprintf(color.Output, "%v %s  \n", green(" [+] "+subCRT.Source+":"), k)
		}
		for k, _ := range subDNSD.Subdomains {
			fmt.Fprintf(color.Output, "%v %s  \n", green(" [+] "+subDNSD.Source+":"), k)
		}

	}
	if parseFiles != "" && !html {
		parser.GetNmapData(parseFiles)
	} else if parseFiles != "" && html {
		webapp.NmapResults, webapp.CVEHost = parser.GetNmapData(parseFiles)
		webapp.Init()
	}

}
