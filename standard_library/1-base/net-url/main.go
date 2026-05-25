/*
Пакет net/url — разбор и сборка URL и query-параметров.

Зачем:   достать параметры запроса, безопасно собрать URL к внешнему API с экранированием.
Когда:   чтение ?page=2&sort=name в хендлере, построение ссылок, разбор адреса из конфига.
Грабли:  не собирайте URL и query вручную через конкатенацию — спецсимволы (пробел, &, =, /)
         сломают запрос или откроют инъекцию. Только Values.Encode() / QueryEscape.
*/
package main

import (
	"fmt"
	"net/url"
)

func main() {
	// Parse — разобрать строку URL на части.
	u, err := url.Parse("https://api.example.com/users?page=2&sort=name&q=hello+world")
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}
	fmt.Println(u.Scheme, u.Host, u.Path) // => https api.example.com /users

	// Query() → url.Values (map[string][]string). Get берёт первое значение и ДЕКОДИРУЕТ его.
	// Так читают параметры в хендлере: r.URL.Query().Get("page").
	q := u.Query()
	fmt.Println(q.Get("page")) // => 2
	fmt.Println(q.Get("q"))    // => hello world  (+ декодирован в пробел)

	// Сборка URL: задаём параметры через Values и Encode() — значения экранируются сами.
	params := url.Values{}
	params.Set("token", "a b&c") // пробел и & будут закодированы
	params.Add("ids", "1")       // Add позволяет несколько значений одного ключа
	params.Add("ids", "2")
	built := url.URL{Scheme: "https", Host: "example.com", Path: "/search", RawQuery: params.Encode()}
	fmt.Println(built.String())
	// => https://example.com/search?ids=1&ids=2&token=a+b%26c

	// QueryEscape — экранировать одиночное значение для подстановки в шаблон URL.
	fmt.Println(url.QueryEscape("a/b c")) // => a%2Fb+c
}

/*
Что запомнить (что чаще и почему):
  • Чтение параметров — почти всегда r.URL.Query().Get("key"). Query() парсит строку в Values;
    если читаете несколько ключей, вызовите Query() один раз и переиспользуйте результат.
  • Get vs Values["key"]: Get возвращает первое значение (или ""), Values["key"] — весь []string
    (когда ключ повторяется, как ids=1&ids=2). В 90% случаев нужен Get.
  • Сборка URL: Values.Encode() для query + url.URL{} для всего адреса. Ручная склейка строк —
    источник багов с экранированием. Encode ещё и сортирует ключи (детерминированный вывод).
  • url.Parse не валидирует «доступность» — только синтаксис. Проверяйте Scheme/Host, если важно.

Типичные сценарии:
  1) Параметр в хендлере: page := r.URL.Query().Get("page")
  2) Запрос к API:        u := url.URL{Scheme:"https", Host:host, Path:p, RawQuery:params.Encode()}
  3) Экранирование:       link := "/items?id=" + url.QueryEscape(rawID)
*/
