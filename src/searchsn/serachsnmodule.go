package searchsn

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var db *sql.DB

var server = "app01"
var portdb = 1433
var user = "DBQlik"
var password = "Yfcnhjqrf48"
var database = "DBQlik-log-xml"
var err error

// добавил функцию для поиска серийных номеров. На вход подается /searchsn?sn=...
func Searchsn(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Println(r.RemoteAddr, r.RequestURI)
		sn := r.URL.Query().Get("sn")
		if sn == "" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		ret, err := searchIntoBase(sn)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		log.Println(string(ret))

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(ret)
		if err != nil {
			return
		}
		return

	} else {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

func readsnFromBase(sn string) ([]byte, error) {
	var findsn string
	var finddate time.Time

	type retType struct {
		Id         string
		DateImport string
	}

	var ret retType

	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return []byte("error"), err
	}

	tsql := fmt.Sprintf("SELECT [sn], [importdate] FROM [DBQlik-log-xml].[dbo].[sntable] where sn=@sn;")

	// Execute query
	rows, err := db.QueryContext(ctx, tsql, sql.Named("sn", sn))
	if err != nil {
		return []byte("error"), err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Error close connection: ", err.Error())
		}
	}(rows)

	retDateImport := ""

	for rows.Next() {

		// Get values from row.
		err := rows.Scan(&findsn, &finddate)
		if err != nil {
			return []byte("error"), err
		}
		ss, _ := finddate.MarshalJSON()
		retDateImport = string(ss)

	}

	ret.Id = sn
	ret.DateImport = retDateImport

	retjson, _ := json.Marshal(ret)

	return retjson, nil
}

func searchIntoBase(sn string) ([]byte, error) {

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, portdb, database)

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Println("Error creating connection pool: ", err.Error())
		return []byte(""), err
	}

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Println(err.Error())
		return []byte(""), err
	}
	fmt.Printf("Connected!\n")

	log.Println(sn)

	answer, err := readsnFromBase(sn)

	if err != nil {
		log.Println(err.Error())
		return []byte(""), err
	}
	return answer, nil
}