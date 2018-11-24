package webapp

import (
	
	"github.com/kevinborras/GokGok/modules/parser"
	//"fmt"
	"html/template"
	"net/http"
	"github.com/op/go-logging"

)

var NmapResults parser.Hosts

//Format
var log = logging.MustGetLogger("example")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func mainpage(w http.ResponseWriter, r *http.Request ) {
	log.Info(" -  Method:", r.Method, " - /")
	//fmt.Println(NmapResults.List)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("modules/webapp/html/index.html")
		t.Execute(w, NmapResults)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

//Init starts the web server 
func Init() {
	http.HandleFunc("/", mainpage)
	fileServer := http.FileServer(http.Dir("html/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	log.Info(" -  Listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(" -  ListenAndServe: ", err)
	}
}