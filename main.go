package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"flag"
	"path/filepath"
	"strconv"

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

var (
	template1, template2 []string
	port                 *int
)

// Description: readLines -
// в Го файлы читаются в []byte
// чтобы перевести []bytes->[]string добавил эту функцию
func readLines(path string) ([]string, error) {
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
func showDescription(w http.ResponseWriter, r *http.Request) {
	// Извлекаем значение параметра id из URL и попытаемся
	// конвертировать строку в integer используя функцию strconv.Atoi(). Если его нельзя
	// конвертировать в integer, или значение меньше 1, возвращаем ответ
	var outputstrings []string
	var outputbyte []byte
	//var itsBe map[string]bool

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

func writeDescription(w http.ResponseWriter, r *http.Request) {
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

func init() {
	port = flag.Int("p", 8134, "Port service")
}

func main() {

	println("Help comandline arguments run: \n\tservicedescription -p PORT")

	flag.Parse()

	template1, _ = readLines(filepath.Join("templates", "template.html"))
	template2, _ = readLines(filepath.Join("templates", "template2.html"))

	mux := http.NewServeMux()
	mux.HandleFunc("/description", showDescription)
	mux.HandleFunc("/writedesription", writeDescription)

	strport := strconv.Itoa(*port)

	text := "\nЗапуск веб-сервера на http://127.0.0.1:" + strport + "\n" +
		"Сервисы\n" +
		" GET: /descrption?id=xx,yy - Возвращает страницу с описанием услуг\n" +
		" xx,yy - id (int) вида услуги\n" +
		"\n" +
		" POST: /writedesription  - Добавление или обновление описания услуги \n" +
		"  BODY request (json): \n" +
		"	{\"IdText\" : \"id вида услуги \", \n" +
		"	\"Text\": \" текст описания услуги закодированые в BASE64 \"}\n" +
		"\nsource URL: https://github.com/AlexandrVIvanov/servicedescription"

	log.Println(text)
	err := http.ListenAndServe(":"+strport, mux)
	log.Fatal(err)
}
