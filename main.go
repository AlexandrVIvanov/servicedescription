package main

import (
	"flag"
	//_ "github.com/microsoft/go-mssqldb"
	"log"
	"main/src/description"
	"main/src/searchsn"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	port *int
)

func init() {
	port = flag.Int("p", 8134, "Port service")
}

func main() {

	println("Help comandline arguments run: \n\tservicedescription -p PORT")

	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("/description", description.ShowDescription)
	mux.HandleFunc("/writedesription", description.WriteDescription)
	mux.HandleFunc("/search", searchsn.Searchsn)

	strport := strconv.Itoa(*port)

	text, _ := description.ReadLines(filepath.Join("src", "welcome.txt"))

	log.Printf(strings.Join(text, "\n"), strport)
	err := http.ListenAndServe(":"+strport, mux)
	log.Fatal(err)
}
