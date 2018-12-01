package webapp

import (
	parser "github.com/kevinborras/GokGok/modules/parser/nmap"
	"github.com/op/go-logging"
	"html/template"
	"net/http"
)

var NmapResults parser.Hosts
var CVEHost parser.CVEHost

//Format
var log = logging.MustGetLogger("example")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func mainpage(w http.ResponseWriter, r *http.Request) {
	log.Info(" -  Method:", r.Method, " - /")
	if r.Method == "GET" {
		t, _ := template.ParseFiles("modules/webapp/html/index.html")
		t.Execute(w, NmapResults)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func cvePage(w http.ResponseWriter, r *http.Request) {
	log.Info(" -  Method:", r.Method, " - /cve")
	if r.Method == "GET" {
		t, _ := template.ParseFiles("modules/webapp/html/cve.html")
		t.Execute(w, CVEHost)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

//Init starts the web server
func Init() {
	http.HandleFunc("/", mainpage)
	http.HandleFunc("/cve", cvePage)
	fileServer := http.FileServer(http.Dir("html/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	log.Info(" -  Listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(" -  ListenAndServe: ", err)
	}
}
