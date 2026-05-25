/*
Пакет errors (уровень 2) — только то, чего не было в base.
В base уже: New, Is, As, Unwrap. Здесь — Join.

Зачем Join: собрать НЕСКОЛЬКО ошибок в одну (например, накопить все ошибки валидации формы
и вернуть разом, а не падать на первой).

Как запустить:  go run main.go
*/
package main

import (
	"errors"
	"fmt"
)

func validate(age int, name string) error {
	var errs []error // копим ошибки в срез
	if age < 0 {
		errs = append(errs, errors.New("возраст отрицательный"))
	}
	if name == "" {
		errs = append(errs, errors.New("имя пустое"))
	}
	// Join объединяет срез ошибок в одну. Если errs пуст — вернёт nil.
	return errors.Join(errs...)
}

func main() {
	err := validate(-1, "")
	fmt.Println(err)
	// => возраст отрицательный
	//    имя пустое

	// errors.Is по-прежнему работает: найдёт ЛЮБУЮ из объединённых ошибок.
	notFound := errors.New("не найдено")
	combined := errors.Join(notFound, errors.New("ещё одна"))
	fmt.Println(errors.Is(combined, notFound)) // => true

	fmt.Println(validate(20, "Аня") == nil) // => true (ошибок нет -> Join вернул nil)
}

/*
Что важно запомнить:
  • errors.Join(a, b, ...) — склеить несколько ошибок в одну. Пустой набор -> nil.
  • errors.Is умеет искать конкретную ошибку и внутри Join (а не только внутри %w-обёртки).
  • Типичный сценарий: валидация — собрать ВСЕ проблемы и показать пользователю списком,
    вместо «упасть на первой ошибке».

Задача:
  1) Сделай функцию проверки пароля: накопи ошибки «короче 8 символов» и «нет цифры», верни через Join.
*/
