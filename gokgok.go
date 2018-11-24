package main

import (
	"github.com/kevinborras/GokGok/modules/webapp"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/integrii/flaggy"
	"github.com/kevinborras/GokGok/modules/parser"
	"github.com/kevinborras/GokGok/modules/utils"
	
)

// Color support
var yellow = color.New(color.Bold, color.FgYellow).SprintFunc()
var red = color.New(color.Bold, color.FgRed).SprintFunc()
var cyan = color.New(color.Bold, color.FgCyan).SprintFunc()
var green = color.New(color.Bold, color.FgGreen).SprintFunc()

//VERSION of the program
var VERSION = `v1.0.0`

var targetList, parseFiles string
var version, scanThem, html = false, false, false

func init() {

	//gokgok Flags
	flaggy.SetName("Gok-Gok")
	flaggy.SetDescription("Is something there?")
	flaggy.DefaultParser.ShowVersionWithVersionFlag = false

	flaggy.Bool(&version, "v", "version", "Print version")
	flaggy.String(&targetList, "t", "targetList", "File with targets to be checked")
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

	if targetList != "" && !scanThem {
		utils.IsAlive(targetList)

	} else if targetList != "" && scanThem {
		ipList := utils.IsAlive(targetList)
		utils.RunNmap(ipList)

	}
	if parseFiles != "" && !html {
		parser.GetNmapData(parseFiles)
	} else if  parseFiles != "" && html {
		webapp.NmapResults,webapp.CVEHost = parser.GetNmapData(parseFiles)
		webapp.Init()
	}

}
