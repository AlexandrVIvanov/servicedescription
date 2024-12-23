package description

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type TypeDescription struct {
	IdText string
	Text   string
}

// Description: readLines - в Го файлы читаются в []byte чтобы перевести
// []bytes->[]string добавил эту функцию

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Println("Error", err.Error())
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("Error", err.Error())
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
// Description: showDescription - отображает содержимое заметки

func ShowDescription(w http.ResponseWriter, r *http.Request) {
	// Извлекаем значение параметра id из URL и попытаемся
	// конвертировать строку в integer используя функцию strconv.Atoi(). Если его нельзя
	// конвертировать в integer, или значение меньше 1, возвращаем ответ

	var outputstrings []string
	var outputbyte []byte
	var template1, template2 []string

	if (template1 == nil) || (template2 == nil) {
		template1, _ = ReadLines(filepath.Join("templates", "template.html"))
		template2, _ = ReadLines(filepath.Join("templates", "template2.html"))
	}

	log.Println(r.RemoteAddr, r.RequestURI)
	id := r.URL.Query().Get("id") // получаем строку с айдишками

	//fmt.Println(id)

	strsplit := strings.Split(id, ",") // разделяем по запятой

	itsBe := make(map[string]bool)
	//fmt.Println(strsplit)
	// проходимся циклом по массиву, чтобы выцепить айди
	outputstrings = template1
	for number := range strsplit {

		//ищем может уже выводили текст для этого номера?
		if itsBe[strsplit[number]] {
			continue
		}

		filename := filepath.Join("service", strsplit[number]+".txt")
		content, err := ReadLines(filename)
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

		itsBe[strsplit[number]] = true
	}

	outputstrings = append(outputstrings, template2...)
	for strindex := range outputstrings {
		strline := outputstrings[strindex]
		outputbyte = append(outputbyte, []byte(strline)...)
	}

	_, err := w.Write(outputbyte)
	if err != nil {
		log.Println("Error", err.Error())
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

func WriteDescription(w http.ResponseWriter, r *http.Request) {
	var d TypeDescription
	if r.Method == "POST" {
		log.Println(r.RemoteAddr, r.RequestURI)

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
		id := d.IdText
		text, err := base64.StdEncoding.DecodeString(d.Text)
		if err != nil {
			msg := "Error Decode Base64 field test"
			http.Error(w, msg, http.StatusBadRequest)
			log.Println("Error", err.Error())
			return
		}

		WriteDescriptionFile(id, text)

	} else {

		http.Error(w, "Bad request", http.StatusBadRequest)
		return

	}

}
