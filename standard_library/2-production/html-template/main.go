/*
Пакет html/template — HTML-шаблоны с АВТОМАТИЧЕСКИМ экранированием (защита от XSS).

Зачем: генерировать HTML-страницы, подставляя данные. Ключевое отличие от text/template:
html/template САМ экранирует опасные символы (<, >, &, кавычки), поэтому пользовательский ввод
не сможет внедрить <script> на страницу. Это защита от XSS «из коробки».

Как запустить:  go run main.go
*/
package main

import (
	"html/template"
	"os"
)

type Page struct {
	Title string
	User  string
	Items []string
}

func main() {
	// Шаблон: {{.Field}} — подстановка поля, {{range}} — цикл, {{if}} — условие.
	const tpl = `<h1>{{.Title}}</h1>
<p>Привет, {{.User}}!</p>
<ul>{{range .Items}}<li>{{.}}</li>{{end}}</ul>
`
	// New + Parse создают шаблон. Must — обёртка: паникует, если шаблон кривой (удобно для зашитых шаблонов).
	t := template.Must(template.New("page").Parse(tpl))

	data := Page{
		Title: "Каталог",
		User:  "<script>alert(1)</script>", // опасный ввод!
		Items: []string{"яблоко", "банан"},
	}

	// Execute рендерит шаблон с данными в io.Writer (здесь stdout; в вебе — http.ResponseWriter).
	t.Execute(os.Stdout, data)
	// => <h1>Каталог</h1>
	//    <p>Привет, &lt;script&gt;alert(1)&lt;/script&gt;!</p>   <-- скрипт ОБЕЗВРЕЖЕН экранированием
	//    <ul><li>яблоко</li><li>банан</li></ul>
}

/*
Что важно запомнить:
  • html/template ЭКРАНИРУЕТ автоматически и контекстно (в HTML, в атрибуте, в URL, в JS — по-разному).
    Поэтому для HTML бери ИМЕННО html/template, а НЕ text/template — иначе открываешь XSS.
  • Синтаксис: {{.Field}} — поле, {{range .Slice}}...{{end}} — цикл, {{if .Cond}}...{{else}}...{{end}} — условие.
  • template.Must(... Parse(...)) — для шаблонов, зашитых в код (паника при ошибке = баг находится сразу).
    Для шаблонов из файлов в рантайме лучше Parse с обработкой ошибки.
  • Реальные шаблоны грузят из файлов: ParseFiles(...)/ParseGlob("templates/*.html") или ParseFS(embedFS, ...).
  • Funcs(funcMap) добавляет свои функции (форматирование дат, плюрализация) в шаблоны.

Задача:
  1) Сделай шаблон письма "Здравствуйте, {{.Name}}" и отрендери его с именем, содержащим символ '&'.
*/
