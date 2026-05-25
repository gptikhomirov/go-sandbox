/*
Пакет net/url — разбор и сборка веб-адресов (URL).

URL — это адрес в интернете, например:
    https://site.com/users?page=2&sort=name
Разберём по частям:
    https        — схема (протокол)
    site.com     — хост (сам сайт)
    /users       — путь (что запрашиваем)
    page=2&sort=name — параметры запроса (после знака ?), пары "ключ=значение" через &

Зачем нужно: достать параметры из адреса (page, sort) и безопасно собрать адрес самому.

Как запустить:  go run main.go
*/
package main

import (
	"fmt"
	"net/url"
)

func main() {
	// Parse разбирает строку-адрес на части.
	u, err := url.Parse("https://site.com/users?page=2&sort=name")
	if err != nil {
		fmt.Println("плохой адрес:", err)
		return
	}
	fmt.Println(u.Scheme) // => https
	fmt.Println(u.Host)   // => site.com
	fmt.Println(u.Path)   // => /users

	// Query() достаёт параметры запроса. Get(имя) возвращает значение по имени.
	q := u.Query()
	fmt.Println(q.Get("page")) // => 2
	fmt.Println(q.Get("sort")) // => name

	// Сборка адреса. Параметры лучше задавать через url.Values, а не клеить строку вручную:
	// тогда спецсимволы (пробелы, &) автоматически «экранируются» (заменяются на безопасные коды).
	params := url.Values{}
	params.Set("q", "go лучший")  // пробел станет +
	params.Set("page", "1")
	fmt.Println(params.Encode()) // => page=1&q=go+%D0%BB%D1%83%D1%87%D1%88%D0%B8%D0%B9

	// Соберём полный адрес.
	link := "https://site.com/search?" + params.Encode()
	fmt.Println(link)
}

/*
Что важно запомнить:
  • URL состоит из частей: схема, хост, путь и параметры после "?". Parse раскладывает их по полочкам.
  • Параметры читают так: u.Query().Get("имя"). Это самый частый приём в веб-программах.
  • Адрес собирай через url.Values{} + Encode(), а не склейкой строк: спецсимволы экранируются сами,
    иначе пробел или & сломают адрес.

Маленькая задача:
  1) Разбери адрес "https://shop.com/items?id=42&color=red" и напечатай значения id и color.
*/
