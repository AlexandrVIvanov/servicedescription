# servicedescription

Сomandline arguments run:
servicedescription -p PORT"
PORT - порт запуска сервиса по умолчанию 8134

Запуск веб-сервера на http://127.0.0.1:%s
Сервисы
GET:
/descrption?id=xx,yy - Возвращает страницу с описанием услуг
xx,yy - id (int) вида услуги
/search?sn=серийный номер
- Возвращает json
{"Id": SN,
"DateImport": Дата производства}
POST:
/writedesription  - Добавление или обновление описания услуги
BODY request (json):
{"IdText" : id вида услуги,
"Text": " текст описания услуги закодированые в BASE64 "}

Исходники URL: https://github.com/AlexandrVIvanov/servicedescription


ВАЖНО!!!!!
Обратить внимание на подключение к рабочей базе!!!


