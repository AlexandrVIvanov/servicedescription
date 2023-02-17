package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"path/filepath"
	//"fmt"
	"log"
	"net/http"
	"os"
	"strings" // сплитим адрес для айдишников
)

type TypeDescription struct {
	IdText string
	Text   string
}

var template1, template2 []string

// Description: readLines -
// в Го файлы читаются в []byte
// чтобы перевести []bytes->[]string добавил эту функцию
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Обработчик для отображения содержимого заметки.
func showDescription(w http.ResponseWriter, r *http.Request) {
	// Извлекаем значение параметра id из URL и попытаемся
	// конвертировать строку в integer используя функцию strconv.Atoi(). Если его нельзя
	// конвертировать в integer, или значение меньше 1, возвращаем ответ
	var outputstrings []string
	var outputbyte []byte

	id := r.URL.Query().Get("id") // получаем строку с айдишками

	//fmt.Println(id)

	strsplit := strings.Split(id, ",") // разделяем по запятой

	//fmt.Println(strsplit)
	// проходимся циклом по массиву, чтобы выцепить айди
	outputstrings = append(template1)
	for number := range strsplit {

		filename := filepath.Join("service", strsplit[number]+".txt")
		content, err := readLines(filename)
		if err != nil {
			content = []string{"Описание услуги не найдено"}
		}
		for i := range content {
			if i == 0 {
				content[i] = "<h2>" + content[i] + "</h2>"
			} else {
				content[i] = "<p>" + content[i] + "</p>"
			}
		}

		outputstrings = append(outputstrings, "<section>")
		outputstrings = append(outputstrings, content...)
		outputstrings = append(outputstrings, "</section>")
	}

	outputstrings = append(outputstrings, template2...)
	for strindex := range outputstrings {
		strline := outputstrings[strindex]
		outputbyte = append(outputbyte, []byte(strline)...)
	}

	_, err := w.Write(outputbyte)
	if err != nil {
		return
	}
}

func WriteDescriptionFile(id string, text []byte) {
	filename := filepath.Join("service", id+".txt")
	err := os.WriteFile(filename, text, 0666)
	if err != nil {
		log.Println("Error writing file " + filename)
		return
	}
}

// Обработчик для записи заметки в вебсервис
// формат тело запроса JSON
// {id :  number id service must by int
//  text: base64 string }

func writeDescription(w http.ResponseWriter, r *http.Request) {
	var d TypeDescription
	if r.Method == "POST" {

		//read body request 1MB max
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)
		dec := json.NewDecoder(r.Body)
		//dec.DisallowUnknownFields()

		err := dec.Decode(&d)
		if err != nil {
			msg := "Error request body"
			http.Error(w, msg, http.StatusBadRequest)
		}
		id := d.IdText
		text, err := base64.StdEncoding.DecodeString(d.Text)
		if err != nil {
			msg := "Error Decode Base64 field test"
			http.Error(w, msg, http.StatusBadRequest)
		}

		WriteDescriptionFile(id, text)

	} else {

		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

}

func main() {

	template1, _ = readLines("template.html")
	template2, _ = readLines("template2.html")

	mux := http.NewServeMux()
	mux.HandleFunc("/description", showDescription)
	mux.HandleFunc("/writedesription", writeDescription)

	log.Println("Запуск веб-сервера на http://127.0.0.1:8080(locallhost)")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
