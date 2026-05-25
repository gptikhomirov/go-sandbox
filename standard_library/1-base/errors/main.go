/*
Пакет errors — создание, оборачивание и проверка ошибок.

Зачем:   пробрасывать ошибку наверх с контекстом, но уметь узнать её первопричину.
Когда:   репозиторий вернул «не найдено» → хендлер должен отдать 404; валидация → 400.
Грабли:  сравнивать ошибки через == работает только для sentinel-значений и ломается,
         как только ошибку обернули. Поэтому всегда errors.Is, а не err == ErrNotFound.
*/
package main

import (
	"errors"
	"fmt"
)

// Sentinel-ошибка — заранее объявленное значение для сравнения через errors.Is.
// Типичный приём в репозиториях/сервисах.
var ErrNotFound = errors.New("not found")

// Свой ТИП ошибки — когда нужно передать данные (поле, код). Достаётся через errors.As.
type ValidationError struct{ Field string }

func (e *ValidationError) Error() string { return "invalid field: " + e.Field }

func repoGet(id int) error {
	if id <= 0 {
		// %w кладёт ErrNotFound в цепочку — её увидит errors.Is выше по стеку.
		return fmt.Errorf("get user %d: %w", id, ErrNotFound)
	}
	return &ValidationError{Field: "id"}
}

func main() {
	err := repoGet(0)
	fmt.Println(err) // => get user 0: not found  (контекст + причина)

	// errors.Is — проверка на КОНКРЕТНОЕ значение сквозь все обёртки.
	if errors.Is(err, ErrNotFound) {
		fmt.Println("Is ErrNotFound -> true") // => Is ErrNotFound -> true  → отдать 404
	}

	// errors.As — достать ошибку нужного ТИПА и прочитать её поля.
	var ve *ValidationError
	if errors.As(repoGet(5), &ve) {
		fmt.Println("As field:", ve.Field) // => As field: id  → отдать 400
	}

	// errors.Unwrap — снять один слой обёртки. Напрямую нужен редко (обычно хватает Is/As).
	fmt.Println(errors.Unwrap(err)) // => not found
}

/*
Что запомнить (что чаще и почему):
  • errors.Is vs errors.As:
      Is  — «это та самая ошибка?» (сравнение со ЗНАЧЕНИЕМ, напр. sql.ErrNoRows). Чаще всего.
      As  — «это ошибка такого ТИПА?» — когда нужны её поля (код, имя поля, HTTP-статус).
  • errors.New vs fmt.Errorf:
      New    — для постоянных sentinel-значений (var ErrX = errors.New(...)).
      Errorf — когда добавляете контекст к существующей ошибке; обязателен глагол %w,
               иначе цепочка рвётся и Is/As перестают находить причину.
  • Никогда не сравнивайте ошибки через ==: после первой же обёртки сравнение сломается.

Типичные сценарии:
  1) Репозиторий:  if rows == 0 { return ErrNotFound }
  2) Хендлер:      if errors.Is(err, ErrNotFound) { http.Error(w, "", 404); return }
  3) Валидация:    var ve *ValidationError; if errors.As(err, &ve) { ...badRequest(ve.Field) }
*/
