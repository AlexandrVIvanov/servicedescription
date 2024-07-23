package main

import (
	"flag"
	_ "github.com/microsoft/go-mssqldb"
	"log"
	"main/src/description"
	"main/src/searchsn"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	port *int
)

func init() {
	//port = 8431
	port = flag.Int("p", 8431, "Port service")
}

func main() {

	//println("Help comandline arguments run: \n\tservicedescription -p PORT")
	go func() {
		http.ListenAndServe("localhost:8080", nil)
	}()

	flag.Parse()

	strport := strconv.Itoa(*port)

	text, _ := description.ReadLines(filepath.Join("src", "welcome.txt"))
	log.Printf(strings.Join(text, "\n"), strport)

	mux := http.NewServeMux()
	mux.HandleFunc("/description", description.ShowDescription)
	mux.HandleFunc("/writedesription", description.WriteDescription)
	mux.HandleFunc("/search", searchsn.Searchsn)

	err := http.ListenAndServe(":"+strport, mux)

	log.Fatal(err)

}
