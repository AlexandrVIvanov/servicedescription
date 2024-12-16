package certificates

/*
 TODO Добавить функцию создания сертификата
 TODO добавить функцию авторизации,
 TODO добавить функцию разрешенных ip адресов

*/

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"main/src/readconfig"
	"net/http"
	"time"
)

// TypeCertificate структура для сертификата
type TypeCertificate struct {
	Certbarcode string `json:"certbarcode"`
	Certprice   int    `json:"certprice"`
	Certurl     string `json:"certurl"`
}

// TypeAddCertificate структура для добавления сертификата
type TypeAddCertificate struct {
	Payuuid      string            `json:"payuuid"`
	Paytimestamp string            `json:"paytimestamp"`
	Paysendtel   string            `json:"paysendtel"`
	Paysendemail string            `json:"paysendemail"`
	Certs        []TypeCertificate `json:"certs"`
}
type TypeAddCertificates struct {
	Certificates []TypeAddCertificate `json:"certificates"`
}

// функция keyTruth проверки корректности X-API-Authorization
func keyTruth(msg []byte, key string) (bool, error) {

	secret := "zwtvl-v^))%tcw#&p(a%jax4rt%dg!qpw9c6wo+ljc$j32#v1d"

	h := hmac.New(md5.New, []byte(secret))
	h.Write(msg)
	signature := hex.EncodeToString(h.Sum(nil))

	if signature == key {
		return true, nil
	} else {
		return false, errors.New("invalid key")
	}
}

// insertpaycheck в базу данных

func Insertpaycheck(db *sql.DB,
	payuuid string,
	paytimestamp string,
	paysendtel string,
	paysendemail string) error {

	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return err
	}

	tsql := "INSERT INTO [dbo].[cert_paytable] ([payuuid] ,[paytimestamp] ,[paysendtel] ,[paysendemail]) VALUES ( @p1, @p2, @p3, @p4)"

	_, err = db.ExecContext(ctx, tsql,
		sql.Named("p1", payuuid),
		sql.Named("p2", paytimestamp),
		sql.Named("p3", paysendtel),
		sql.Named("p4", paysendemail))

	if err != nil {
		return err
	}
	return nil
}

func Insertcert(db *sql.DB,
	payuuid string,
	certbarcode string,
	certprice int,
	certurl string) error {

	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return err
	}

	tsql := `INSERT INTO [dbo].[cert_certtable] (payuuid, certbarcode, certprice, certurl) 
			VALUES (@p1, @p2, @p3, @p4)`

	_, err = db.ExecContext(ctx, tsql,
		sql.Named("p1", payuuid),
		sql.Named("p2", certbarcode),
		sql.Named("p3", certprice),
		sql.Named("p4", certurl))

	if err != nil {
		return err
	}

	return nil
}

// CertificateAddDB Функция добавления нового сертификата в базу данных
func CertificateAddDB(d TypeAddCertificates) error {

	var db *sql.DB
	var ctx context.Context
	var cancel context.CancelFunc

	conf, err := readconfig.Getconfigsqlserver()

	if err != nil {
		return err
	}
	server := conf.ServerName
	portdb := conf.Port
	user := conf.UserName
	password := conf.Password
	database := conf.Database

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, portdb, database)

	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Println("Error creating connection pool: ", err.Error())
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	err = db.PingContext(ctx)

	if err != nil {
		return err
	}

	defer func() {
		_ = db.Close()
		conf = nil
		cancel()
	}()

	for _, certificate := range d.Certificates {

		payuuid := certificate.Payuuid
		paytimestamp := certificate.Paytimestamp
		paysendtel := certificate.Paysendtel
		paysendemail := certificate.Paysendemail

		err = Insertpaycheck(db, payuuid, paytimestamp, paysendtel, paysendemail)
		if err != nil {
			return err
		}

		for _, cert := range certificate.Certs {
			certbarcode := cert.Certbarcode
			certprice := cert.Certprice
			certurl := cert.Certurl

			err = Insertcert(db, payuuid, certbarcode, certprice, certurl)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

// CertificateAddHttp Функция записи сертификата через http
func CertificateAddHttp(bodytext []byte) error {

	conf, err := readconfig.Getconfighttpclient()

	if err != nil {
		return err
	}
	urlserver := conf.URLServerName
	urlpath := conf.URLPath

	resp, err := http.Post(urlserver+urlpath, "application/json", bytes.NewBuffer(bodytext))
	if err != nil {
		log.Println("Error creating http client: ", err.Error())
		return err
	}

	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	return nil
}

// функция CertificateAdd добавления сертификата

func CertificateAdd(w http.ResponseWriter, r *http.Request) {

	var d TypeAddCertificates

	defer func() {
		_ = r.Body.Close()
	}()

	if r.Method == "POST" {

		// Проверяем X-API-Authorization
		xapikey := r.Header.Get("X-API-Authorization")

		if xapikey == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		//читаем r.Body в bytes[]

		body, _ := io.ReadAll(r.Body)

		valid, err := keyTruth(body, xapikey)

		if err != nil || !valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// читаем json из тела проверяем валидность
		bodyjson, err := base64.StdEncoding.DecodeString(string(body))
		if err != nil {
			log.Println("Error decoding Base64", body)
			http.Error(w, "Error decoding Base64", http.StatusInternalServerError)
			return
		}

		// прочитать json из bodyjson
		var decoder = json.NewDecoder(bytes.NewReader(bodyjson))
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&d)

		if err != nil {
			log.Println("Error decoding JSON", string(bodyjson))
			return
		}

		go func() {
			err := CertificateAddHttp(bodyjson)
			if err != nil {

			}
		}()

		// Прочитали json необходимо записать в базу sql и отправить http запрос в 1с
		err = CertificateAddDB(d)
		if err != nil {
			log.Println("Error add certificate to DB "+err.Error(), string(bodyjson))
			http.Error(w, "Error add certificate to DB "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println(r.RemoteAddr, r.RequestURI)

	} else {

		http.Error(w, "Bad request", http.StatusBadRequest)
		return

	}
}
