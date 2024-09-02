package chatanalize

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type TypeIncomingText struct {
	Id   string `json:"Id"`
	Text string `json:"Text"`
}

type Response struct {
	Answer string `json:"Answer"`
}

//Gigachat analize text

func GetAnswerGigachat(accesstoken string, Text []byte) (string, error) {

	type TypeIncomingText struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	type TypeBodyRequest struct {
		Model             string              `json:"model"`
		Temperature       float32             `json:"temperature"`
		N                 int                 `json:"n"`
		MaxTokens         int                 `json:"max_tokens"`
		Stream            bool                `json:"stream"`
		Updateinterval    int                 `json:"update_interval"`
		RepetitionPenalty float32             `json:"repetition_penalty"`
		Messages          [2]TypeIncomingText `json:"messages"`
	}

	type TypeIncomingChoice struct {
		Message struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"message"`
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
	}

	type TypeBodyResponse struct {
		Choices []TypeIncomingChoice `json:"choices"`
		Created int                  `json:"created"`
		Model   string               `json:"model"`
		Object  string               `json:"object"`
		Usage   struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}

	var b TypeBodyRequest
	var err error
	var BodyResponse TypeBodyResponse

	firstPrompt := []byte("Эмоциональная оценка текста если Положительная: 1, если отрицательная: 0")
	firstPrompt = append(firstPrompt, Text...)

	b.Model = "GigaChat"
	b.Temperature = 0.87
	b.MaxTokens = 512
	b.Stream = false
	b.N = 1
	b.Updateinterval = 0
	b.RepetitionPenalty = 1.07

	b.Messages[0].Role = "system"
	b.Messages[0].Content = "Отвечай одним словом"
	b.Messages[1].Role = "user"
	b.Messages[1].Content = string(firstPrompt)

	requestServer := "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"

	client := &http.Client{Timeout: 10 * time.Second}

	byteBody, _ := json.Marshal(b)
	log.Println(string(byteBody))

	req, err := http.NewRequest("POST", requestServer, bytes.NewReader(byteBody))
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+accesstoken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&BodyResponse)
	if err != nil {
		return "", err
	}

	return BodyResponse.Choices[0].Message.Content, err
}

func GetTokenGigachat() (string, error) {

	type TypeAccessToken struct {
		AccessToken string `json:"access_token"`
		ExpiresAt   int    `json:"expires_at"`
	}

	var d TypeAccessToken

	credentials := "NzM0NWQ0ODUtNzBhZC00YjA2LWE3MWMtZDZmOThjZTU3ZjM4OmZiZmZkNTJjLTMzYWItNGQ4ZS1iYzA5LTJhNjhiNTkxMTAyMQ=="
	authserver := "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"

	client := &http.Client{Timeout: 10 * time.Second}

	jsonBody := []byte(`scope=GIGACHAT_API_PERS`)
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest("POST", authserver, bodyReader)
	if err != nil {
		return "", err
	}

	req.Header.Add("RqUID", "6f0b1291-c7f3-43c6-bb2e-9f3efb2dc98e")
	req.Header.Add("Authorization", "Bearer "+credentials)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", nil
	}

	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&d)
	if err != nil {
		return "", err
	}

	return d.AccessToken, err
}

func AnalyzeText(Text []byte) (string, error) {
	var (
		err         error  = nil
		accesstoken string = ""
		retString   string = ""
	)

	accesstoken, err = GetTokenGigachat()
	if err != nil {
		return accesstoken, err
	}

	log.Println(accesstoken)

	retString, err = GetAnswerGigachat(accesstoken, Text)
	if err != nil {
		return retString, err
	}

	return retString, nil
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

		text, err := base64.StdEncoding.DecodeString(d.Text)
		if err != nil {
			msg := "Error Decode Base64 field test"
			http.Error(w, msg, http.StatusBadRequest)
			log.Println("Error", err.Error())
			return
		}

		retString, err := AnalyzeText(text)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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

		log.Println("Request:", string(text), ", Answer:", retString)

	} else {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

}
