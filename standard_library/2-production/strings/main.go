/*
Пакет strings (уровень 2) — только новое поверх base.
В base уже: Contains, HasPrefix/Suffix, TrimSpace/Prefix/Suffix, Split, Fields, Join, ReplaceAll,
ToLower/Upper, EqualFold, Count, Index, NewReader, Builder.
Здесь: Cut/CutPrefix/CutSuffix, гибкая обрезка (Trim/TrimFunc), SplitN, Replace, LastIndex,
Repeat, Map, NewReplacer, точная сборка Builder.

Как запустить:  go run main.go
*/
package main

import (
	"fmt"
	"strings"
	"unicode"
)

func main() {
	// Cut — современный идиоматичный способ разрезать строку по первому разделителю.
	// Возвращает «до», «после» и флаг «нашёлся ли разделитель». Заменяет связку Index+срез.
	key, val, found := strings.Cut("name=Alice", "=")
	fmt.Println(key, val, found) // => name Alice true

	// CutPrefix/CutSuffix — срезать начало/конец И сразу узнать, было ли что срезать.
	rest, ok := strings.CutPrefix("Bearer abc123", "Bearer ")
	fmt.Println(rest, ok) // => abc123 true

	// Trim с набором символов / TrimFunc с условием — гибкая обрезка по краям.
	fmt.Println(strings.Trim("***важно***", "*"))                       // => важно
	fmt.Println(strings.TrimFunc("123abc456", unicode.IsDigit))          // => abc

	// SplitN ограничивает число частей; Replace заменяет только первые n вхождений.
	fmt.Printf("%q\n", strings.SplitN("a=b=c", "=", 2)) // => ["a" "b=c"]
	fmt.Println(strings.Replace("ля-ля-ля", "ля", "ла", 1)) // => ла-ля-ля (только первое)

	// LastIndex — позиция ПОСЛЕДНЕГО вхождения. Repeat — повторить строку.
	fmt.Println(strings.LastIndex("a/b/c", "/")) // => 3
	fmt.Println(strings.Repeat("=", 5))          // => =====

	// Map преобразует каждый символ функцией (например, «зашифровать» сдвигом, убрать символы).
	rot := strings.Map(func(r rune) rune { return r + 1 }, "abc")
	fmt.Println(rot) // => bcd

	// NewReplacer — много замен за один проход (эффективнее цепочки ReplaceAll).
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;")
	fmt.Println(r.Replace("a < b & c")) // => a &lt; b &amp; c

	// Builder: точная запись по байту/руне + Grow (заранее выделить память под N байт).
	var b strings.Builder
	b.Grow(16)
	b.WriteByte('>')
	b.WriteRune('я')
	fmt.Println(b.String()) // => >я
}

/*
Что важно запомнить:
  • Cut — теперь главный способ «разрезать по разделителю» (Go 1.18+): читается лучше, чем Index+slice.
    Split — когда частей много; Cut — когда нужно «ключ и остаток».
  • Replace(s, old, new, n) vs ReplaceAll: ReplaceAll = Replace c n = -1. Указывай n, если нужно ограничить.
  • NewReplacer выгоднее нескольких ReplaceAll подряд (один проход, заранее построенная таблица).
  • Builder.Grow(n) заранее резервирует память — ускоряет сборку, если примерный размер известен.

Задача:
  1) Разбери заголовок "Authorization: Bearer xyz" через Cut по ": ", затем CutPrefix "Bearer ".
*/
