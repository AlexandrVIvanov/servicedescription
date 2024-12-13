package certificates

/*
 TODO Добавить функцию создания сертификата
 TODO добавить функцию авторизации,
 TODO добавить функцию разрешенных ip адресов

*/

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

// структура для сертификата
type TypeCertificate struct {
	Certbarcode string `json:"certbarcode"`
	Certprice   int    `json:"certprice"`
	Certurl     string `json:"certurl"`
}

// структура для добавления сертификата
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
		return false, errors.New("Invalid key")
	}
}

// Функция добавления нового сертификата в базу данных
func CertificateAddDB(d TypeAddCertificate) error {
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

		// Прочитали json необходимо записать в базу sql и отправить http запрос в 1с

		log.Println(r.RemoteAddr, r.RequestURI)

	} else {

		http.Error(w, "Bad request", http.StatusBadRequest)
		return

	}
}
