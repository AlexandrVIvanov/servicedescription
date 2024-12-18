package main

import (
	"flag"
	_ "github.com/microsoft/go-mssqldb"
	"log"
	"main/src/certificates"
	"main/src/chatanalize"
	"main/src/description"
	"main/src/searchsn"
	"net/http"
	"os"
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

	logfilename := filepath.Join("", "servicedescription.log")

	logFile, err := os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = logFile.Close()
	}()

	flag.Parse()

	strport := strconv.Itoa(*port)

	text, _ := description.ReadLines(filepath.Join("src", "welcome.txt"))
	log.Printf(strings.Join(text, "\n"), strport)

	mux := http.NewServeMux()
	mux.HandleFunc("/description", description.ShowDescription)
	mux.HandleFunc("/writedesription", description.WriteDescription)
	mux.HandleFunc("/search", searchsn.Searchsn)
	mux.HandleFunc("/chatanalize", chatanalize.Chatanalize)
	mux.HandleFunc("/certificates/add", certificates.CertificateAdd)
	mux.HandleFunc("/certificates/status", certificates.CertificateGetStatus)

	// TODO Проверить утечку памяти
	// FIXME Память куда то утекает

	err = http.ListenAndServe(":"+strport, mux)

	log.Fatal(err)

}
