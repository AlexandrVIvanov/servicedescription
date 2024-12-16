package searchsn

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"main/src/readconfig"
	"net/http"
	"strings"
	"time"
)

// var db *sql.DB

// Searchsn добавил функцию для поиска серийных номеров. На вход подается /searchsn?sn=...
func Searchsn(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	if r.Method == "GET" {
		log.Println(r.RemoteAddr, r.RequestURI)
		sn := r.URL.Query().Get("sn")
		if sn == "" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		ret, err := searchIntoBase(sn)
		if err != nil {
			http.Error(w, "Internal error "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println(string(ret))

		log.Printf("Время выполнения: %s", time.Since(start))

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

func searchIntoBase(sn string) ([]byte, error) {
	var err error
	var db *sql.DB

	var conf *readconfig.TypeSqlConfiguration

	conf, err = readconfig.Getconfigsqlserver()

	if err != nil {
		return []byte(""), err
	}
	server := conf.ServerName
	portdb := conf.Port
	user := conf.UserName
	password := conf.Password
	database := conf.Database

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, portdb, database)

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Println("Error creating connection pool: ", err.Error())
		return []byte(""), err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		conf = nil
		cancel()
	}()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		conf = nil
		return []byte(""), err
	}
	defer func() {
		_ = db.Close()
		cancel()
	}()

	log.Println(sn)

	answer, err := readsnFromBase(db, sn)

	ctx.Done()

	if err != nil {
		log.Println(err.Error())
		return []byte(""), err
	}
	return answer, nil
}

func readsnFromBase(db *sql.DB, sn string) ([]byte, error) {
	var findsn string
	var finddate time.Time
	var exportdate time.Time
	var retaildate time.Time
	var repairdate time.Time
	var warrantydate time.Time
	var productname string
	var customer string
	var code string

	type retType struct {
		Id           string
		IdFound      bool
		DateImport   string
		DateExport   string
		RetailDate   string
		Productname  string
		Customer     string
		Code         string
		DateRepair   string
		WarrantyDate string
	}

	var ret retType

	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return []byte("error"), err
	}

	//tsql := "SELECT [sn], " +
	//	"[importdate], " +
	//	"[exportdate], " +
	//	"trim([productname]), " +
	//	"trim([customer]), " +
	//	"trim([code]), " +
	//	"[retaildate], " +
	//	"[repairdate], [warrantydate] FROM [DBQlik-log-xml].[dbo].[sntable] where sn=@sn;"

	tsql := `Select 
	s2.sn sn,
	s2.importdate importdate,
	s2.exportdate,
	 isnull(s3.productname,'') productname,
	 s2.customer,
	 isnull(s3.productcode,'') productcode,
	 s2.retaildate,
	 s2.repairedate,
	 s2.warrantydate from
	(Select 
		@sn as sn, 
		max(idnom) idnom, 
		isnull(MIN(importdate),		'2001-01-01 00:00:00.000') importdate, 
		isnull(max(exportdate),		'2001-01-01 00:00:00.000') exportdate, 
		isnull(max(retaildate),		'2001-01-01 00:00:00.000') retaildate, 
		isnull(max(warrantydate),	'2001-01-01 00:00:00.000') warrantydate, 
		isnull(max(repairedate),	'2001-01-01 00:00:00.000') repairedate, 
		isnull(max(customer),'')	customer from 
	(
		Select idnom as idnom , importdate as importdate, null as exportdate, null as retaildate ,null as warrantydate, null as repairedate, null as customer from [dbo].[sn_importtable] with (NOLOCK) where sn = @sn
		union
		Select idnom, null, exportdate, null, null, null, null from [dbo].[sn_exporttable] with (NOLOCK) where sn = @sn
		union
		Select idnom, null, null, [retaildate], [warrantydate], null, customer from [dbo].[sn_retailtable] with (NOLOCK) where sn = @sn
		union
		Select idnom,  null, null,  null, null, [repairedate], null from [dbo].[sn_repairtable] with (NOLOCK) where sn = @sn
	) as s1) as s2
left join [dbo].[sn_nom] s3 on (s2.idnom = s3.idnom)

`

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
	retDateExport := ""
	retDateRetail := ""
	retDateRepair := ""
	retDateWarranty := ""
	retIdFound := false

	if rows.Next() {

		// Get values from row.

		err := rows.Scan(
			&findsn,
			&finddate,
			&exportdate,
			&productname,
			&customer,
			&code,
			&retaildate,
			&repairdate,
			&warrantydate)
		if err != nil {
			return []byte("error"), err
		}

		ss, _ := finddate.MarshalJSON()
		retDateImport = string(ss)
		retDateImport = strings.Replace(retDateImport, "\"", "", -1)

		ss, _ = exportdate.MarshalJSON()
		retDateExport = string(ss)
		retDateExport = strings.Replace(retDateExport, "\"", "", -1)

		ss, _ = retaildate.MarshalJSON()
		retDateRetail = string(ss)
		retDateRetail = strings.Replace(retDateRetail, "\"", "", -1)

		ss, _ = repairdate.MarshalJSON()
		retDateRepair = string(ss)
		retDateRepair = strings.Replace(retDateRepair, "\"", "", -1)

		ss, _ = warrantydate.MarshalJSON()
		retDateWarranty = string(ss)
		retDateWarranty = strings.Replace(retDateWarranty, "\"", "", -1)

		retIdFound = true
	}

	ret.Id = sn
	ret.IdFound = retIdFound
	ret.DateImport = retDateImport
	ret.DateExport = retDateExport
	ret.Customer = customer
	ret.Productname = productname
	ret.Code = code
	ret.RetailDate = retDateRetail
	ret.DateRepair = retDateRepair
	ret.WarrantyDate = retDateWarranty

	retjson, _ := json.Marshal(ret)

	err = rows.Close()
	if err != nil {
		log.Println("Error close connection: ", err.Error())
	}

	return retjson, nil
}
