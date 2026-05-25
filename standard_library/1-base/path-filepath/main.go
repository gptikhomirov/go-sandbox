/*
Пакет path/filepath — работа с путями к файлам и папкам.

Путь — это «адрес» файла в системе, например var/log/app.txt. Разные системы используют
разный разделитель: Linux/Mac — "/", Windows — "\". filepath сам подставляет правильный,
поэтому путь НИКОГДА не склеивают вручную через "+" — для этого есть filepath.Join.

Зачем нужно: собрать путь, вытащить имя файла или его расширение, обойти все файлы в папке.

Как запустить:  go run main.go
*/
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Join собирает путь из кусочков и сам ставит правильный разделитель.
	path := filepath.Join("var", "log", "app.txt")
	fmt.Println(path) // => var/log/app.txt (на Linux/Mac)

	// Разбор пути на части:
	fmt.Println(filepath.Base(path)) // => app.txt  (имя файла)
	fmt.Println(filepath.Dir(path))  // => var/log  (папка)
	fmt.Println(filepath.Ext(path))  // => .txt     (расширение)

	// Создадим временную папку с парой файлов, чтобы показать обход.
	dir, _ := os.MkdirTemp("", "demo")  // временная папка
	defer os.RemoveAll(dir)             // defer = «выполнить в конце функции»: удалим папку при выходе
	os.WriteFile(filepath.Join(dir, "a.txt"), nil, 0o644)
	os.WriteFile(filepath.Join(dir, "b.log"), nil, 0o644)

	// WalkDir обходит ВСЕ файлы и подпапки. Для каждого вызывает нашу функцию.
	// Здесь соберём только файлы с расширением .txt.
	var txtFiles []string
	filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(p) == ".txt" { // не папка И расширение .txt
			txtFiles = append(txtFiles, filepath.Base(p))
		}
		return nil // nil = «продолжаем обход без ошибок»
	})
	fmt.Println(txtFiles) // => [a.txt]
}

/*
Что важно запомнить:
  • Пути собирай через filepath.Join(...), а не склейкой строк "+": Join сам подставит "/" или "\".
  • Base — имя файла, Dir — папка, Ext — расширение. Частый приём при работе с загруженными файлами.
  • filepath.WalkDir обходит всю папку рекурсивно и вызывает твою функцию на каждый файл/папку.
  • defer ставит действие «на потом» — оно выполнится, когда функция завершится (удобно для очистки).

Маленькая задача:
  1) Собери путь "data/users/report.csv" через Join и напечатай имя файла и его расширение.
*/
