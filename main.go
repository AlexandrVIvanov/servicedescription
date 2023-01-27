package main
 
import (
	//"fmt"
	"log"
	"net/http"
	//"strconv" 
	"strings" // сплитим адрес для айдишников 
	"os"
)
 
func home(w http.ResponseWriter, r *http.Request) {
	// Проверяется, если текущий путь URL запроса точно совпадает с шаблоном "/". Если нет, вызывается
	// функция http.NotFound() для возвращения клиенту ошибки 404.
	// Важно, чтобы мы завершили работу обработчика через return. Если мы забудем про "return", то обработчик
	// продолжит работу и выведет сообщение  как ни в чем не бывало.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
 
	w.Write([]byte("Просто Отсканируйте Код"))
}
 
// Обработчик для отображения содержимого заметки.
func showDescription(w http.ResponseWriter, r *http.Request) {
	// Извлекаем значение параметра id из URL и попытаемся
	// конвертировать строку в integer используя функцию strconv.Atoi(). Если его нельзя
	// конвертировать в integer, или значение меньше 1, возвращаем ответ

	id := r.URL.Query().Get("id") // получаем строку с айдишками

	//fmt.Println(id)

	strsplit := strings.Split(id,",") // разделяем по запятой 

	//fmt.Println(strsplit)


	// проходимся циклом по массиву, чтобы выцепить айди 
	for docnumber := range strsplit {
  	 	//fmt.Println(strsplit[docnumber])

  	 	filename :=string("service/"+strsplit[docnumber]+".txt")

  	 	file, err := os.Open(filename)
    	if err != nil {
        	log.Fatal(err)
        	//w.Write([]byte(""))
   		 }
    	defer func() {
        	if err = file.Close(); err != nil {
            	log.Fatal(err)
       		}
   		}()
 		b, err := os.ReadFile(file)
  		w.Write([]byte(b))
	 }



	

	 // if err != nil || id < 1 
	 // {
	 // 	http.NotFound(w, r)
	 // 	return
	 // }
 	
	
	// Используем функцию fmt.Fprintf() для вставки значения из id в строку ответа
	// и записываем его в http.ResponseWriter.
	//fmt.Fprintf(w, "Отображение выбранной заметки с ID %d...",id)
}
 

 
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/description", showDescription)
	
	log.Println("Запуск веб-сервера на http://127.0.0.1:8080(locallhost)")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
