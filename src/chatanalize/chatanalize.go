package chatanalize

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
)

type TypeIncomingText struct {
	Id   string `json:"Id"`
	Text string `json:"Text"`
}
type Response struct {
	Answer string `json:"Answer"`
}

//Gigachat analize text

func GetTokenGigachat() string {
	return ""
}

func AnalizeText(Text []byte) string {

	gigatoken := GetTokenGigachat()

	return gigatoken
}

func Chatanalize(w http.ResponseWriter, r *http.Request) {

	var d TypeIncomingText

	//Check authorization header
	token := "A1B2C3D1E2F3"

	auth := r.Header.Get("Authorization")
	if auth != "Bearer "+token {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == "POST" {

		//read body request 1MB max
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		err := dec.Decode(&d)
		if err != nil {
			msg := "Error request body"
			http.Error(w, msg, http.StatusBadRequest)
			log.Println("Error", err.Error())
			return
		}

		//id := d.id
		text, err := base64.StdEncoding.DecodeString(d.Text)

		if err != nil {
			msg := "Error Decode Base64 field test"
			http.Error(w, msg, http.StatusBadRequest)
			log.Println("Error", err.Error())
			return
		}
		retString := AnalizeText(text)

		response := Response{retString}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}
