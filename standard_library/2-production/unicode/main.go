/*
Пакет unicode — классификация символов (рун): буква, цифра, пробел, регистр.

Руна (rune) — это один символ Unicode (в Go это тип rune = int32). unicode отвечает на вопросы
«это буква?», «это цифра?», «это пробел?» — правильно для ЛЮБОГО языка, а не только латиницы.

Зачем: валидация ввода, разбор текста, фильтрация символов — без хрупких сравнений вроде c >= 'a' && c <= 'z'.

Как запустить:  go run main.go
*/
package main

import (
	"fmt"
	"unicode"
)

func main() {
	// Классификация рун. Работает для всех алфавитов (кириллица, китайский и т.д.).
	fmt.Println(unicode.IsLetter('я')) // => true
	fmt.Println(unicode.IsDigit('7'))  // => true
	fmt.Println(unicode.IsSpace(' '))  // => true
	fmt.Println(unicode.IsLetter('!')) // => false
	// IsNumber шире, чем IsDigit: ловит не только 0-9, но и числовые символы вроде ½, Ⅷ.
	fmt.Println(unicode.IsNumber('5'), unicode.IsNumber('½')) // => true true

	// IsUpper/IsLower и смена регистра ОДНОЙ руны.
	fmt.Println(unicode.IsUpper('Я'))      // => true
	fmt.Println(string(unicode.ToLower('Я'))) // => я
	fmt.Println(string(unicode.ToUpper('a'))) // => A

	// Пример применения: оставить в строке только буквы и цифры (перебор по рунам через range).
	clean := ""
	for _, r := range "Привет, мир! 123" {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			clean += string(r)
		}
	}
	fmt.Println(clean) // => Приветмир123

	// Подсчёт цифр в строке.
	digits := 0
	for _, r := range "тел: +7 999 123" {
		if unicode.IsDigit(r) {
			digits++
		}
	}
	fmt.Println(digits) // => 8
}

/*
Что важно запомнить:
  • Руна (rune) — один символ Unicode. Перебор строки через `for _, r := range s` даёт руны (а не байты!).
  • unicode.IsLetter/IsDigit/IsSpace/IsUpper/IsLower — правильная проверка для ЛЮБОГО языка.
    НЕ пиши c >= 'a' && c <= 'z' — это сломается на кириллице, диакритике и т.п.
  • ToUpper/ToLower здесь — для ОДНОЙ руны; для целой строки используй strings.ToUpper/ToLower.
  • Часто работает в паре с unicode/utf8 (длина в рунах) и со strings.Map / strings.TrimFunc.

Задача:
  1) Напиши функцию isValidLogin(s): только буквы/цифры/'_' и длиной 3..16. Проверь через unicode.IsLetter/IsDigit.
*/
