/*
Пакет strings — операции над строками.

Зачем:   разбор и нормализация текста, проверки, эффективная сборка строк.
Когда:   email/логин перед записью в БД, разбор путей и заголовков, генерация текста.
Грабли:  строки в Go неизменяемы — каждая операция возвращает НОВУЮ строку; склейка через
         += в цикле порождает мусор, для этого есть Builder. Индексы строк — это БАЙТЫ,
         а не символы (для Unicode см. unicode/utf8 на уровне production).
*/
package main

import (
	"fmt"
	"strings"
)

func main() {
	// Нормализация ввода: TrimSpace + ToLower. Стандарт перед сравнением/записью email.
	email := strings.ToLower(strings.TrimSpace("  Alice@Example.COM  "))
	fmt.Println(email) // => alice@example.com

	// EqualFold — сравнение без учёта регистра (методы HTTP, имена заголовков).
	fmt.Println(strings.EqualFold("GET", "get")) // => true

	// HasPrefix/HasSuffix/Contains — проверки для роутинга и фильтрации.
	fmt.Println(strings.HasPrefix("/api/users", "/api")) // => true

	// Split — разбить по разделителю; Fields — по любым пробелам (и схлопывает повторы).
	fmt.Printf("%q\n", strings.Split("/api/users/42", "/")) // => ["" "api" "users" "42"]
	fmt.Printf("%q\n", strings.Fields("a   b\tc"))          // => ["a" "b" "c"]

	// Join — обратная операция: срез → строка через разделитель.
	fmt.Println(strings.Join([]string{"id", "name"}, ",")) // => id,name

	// ReplaceAll + Count — массовая замена и подсчёт.
	tpl := "Hi {n}, code {n}"
	fmt.Println(strings.ReplaceAll(tpl, "{n}", "Ann")) // => Hi Ann, code Ann
	fmt.Println(strings.Count(tpl, "{n}"))             // => 2

	// Builder — эффективная сборка в цикле (вместо s += ...). Для SQL, отчётов, больших текстов.
	var b strings.Builder
	for i := 0; i < 3; i++ {
		fmt.Fprintf(&b, "row%d;", i)
	}
	fmt.Println(b.String()) // => row0;row1;row2;
}

/*
Что запомнить (что чаще и почему):
  • Split vs Fields vs SplitN:
      Split  — фиксированный разделитель (CSV, путь). Пустые части сохраняются.
      Fields — разбить по пробелам, схлопывая повторы. Идеально для разбора пользовательского ввода.
      SplitN — когда нужно ограничить число частей (напр. "key=value=x" → 2 части).
  • Contains vs HasPrefix/HasSuffix: Contains — «есть где-то»; Has* — «в начале/конце».
    Для роутинга и проверки расширений почти всегда Has*.
  • += vs Builder: одна-две склейки — можно +=; цикл/много кусков — только Builder (нет лишних аллокаций).
  • strings vs bytes: тот же API, но bytes для []byte (см. уровень production) — без конвертаций.

Типичные сценарии:
  1) Нормализация:  email = strings.ToLower(strings.TrimSpace(raw))
  2) Разбор пути:   parts := strings.Split(r.URL.Path, "/")
  3) Сборка SQL:    b.WriteString("INSERT ..."); ... ; q := b.String()
*/
