/*
Пакет text/template — шаблоны для ЛЮБОГО текста (НЕ HTML): конфиги, письма, код, скрипты.

Тот же движок шаблонов, что и html/template, но БЕЗ авто-экранирования HTML. Поэтому его используют
для текста, где экранирование не нужно или мешает: генерация конфигов, SQL, кода, текстовых писем, YAML.

⚠️ Для HTML бери html/template (см. соседний пример) — text/template не защищает от XSS!

Как запустить:  go run main.go
*/
package main

import (
	"os"
	"strings"
	"text/template"
)

type Service struct {
	Name  string
	Port  int
	Hosts []string
}

func main() {
	// Шаблон конфига. {{.Field}}, {{range}}, {{if}} — как в html/template.
	const tpl = `service: {{.Name}}
port: {{.Port}}
upstreams:
{{- range .Hosts}}
  - {{.}}
{{- end}}
`
	// Funcs добавляет свои функции в шаблон. Здесь — upper для верхнего регистра.
	t := template.Must(template.New("cfg").
		Funcs(template.FuncMap{"upper": strings.ToUpper}).
		Parse(tpl + "name_upper: {{upper .Name}}\n"))

	cfg := Service{Name: "api", Port: 8080, Hosts: []string{"10.0.0.1", "10.0.0.2"}}
	t.Execute(os.Stdout, cfg)
	// => service: api
	//    port: 8080
	//    upstreams:
	//      - 10.0.0.1
	//      - 10.0.0.2
	//    name_upper: API
}

/*
Что важно запомнить:
  • text/template — для НЕ-HTML текста (конфиги, письма, кодогенерация, SQL-скрипты). НЕ экранирует.
    Для HTML — только html/template (тот же синтаксис, но с защитой от XSS).
  • {{- и -}} убирают пробелы/переносы вокруг (важно для аккуратного отступа в YAML/конфигах).
  • Funcs(FuncMap{...}) — подключить свои функции (форматирование, upper, default). Вызываются как {{upper .X}}.
  • Execute пишет в io.Writer; шаблоны грузят из строк (Parse) или файлов (ParseFiles/ParseGlob).

Задача:
  1) Сгенерируй .env-файл из map[string]string: строки вида KEY=value через {{range $k, $v := .}}.
*/
