/*
Пакет os (уровень 2) — только новое поверх base.
В base уже: ReadFile/WriteFile, Open/Create, Getenv/LookupEnv, Remove, Args, Stdin/out/err, Exit.
Здесь: Stat/Lstat, классификация ошибок (IsNotExist/IsExist), Mkdir/MkdirAll/RemoveAll/Rename, ReadDir,
временные файлы/папки (CreateTemp/MkdirTemp), Setenv/Unsetenv/Environ, Getwd/Chdir,
OpenFile с флагами, методы файла (Seek/Sync/WriteString).

Как запустить:  go run main.go
*/
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// MkdirTemp/CreateTemp — временные папка/файл (для тестов, кэшей, промежуточных данных).
	dir, _ := os.MkdirTemp("", "demo")
	defer os.RemoveAll(dir) // RemoveAll удаляет папку со всем содержимым

	// MkdirAll создаёт всю цепочку папок (Mkdir — только одну, упадёт если родителя нет).
	os.Mkdir(filepath.Join(dir, "one"), 0o755) // одна папка (родитель dir уже есть)
	sub := filepath.Join(dir, "a", "b")
	os.MkdirAll(sub, 0o755)

	// OpenFile с флагами: открыть на дозапись (O_APPEND), создать если нет (O_CREATE).
	path := filepath.Join(dir, "log.txt")
	f, _ := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	f.WriteString("строка 1\n")
	f.Sync() // сбросить на диск принудительно
	f.Seek(0, 0) // вернуть позицию в начало (для чтения/перезаписи)
	f.Close()

	// Stat — информация о файле (размер, имя, права). Lstat не идёт по symlink.
	info, _ := os.Stat(path)
	fmt.Println(info.Name(), info.Size() > 0) // => log.txt true

	// Классификация ошибок: «файла нет» / «уже есть».
	_, err := os.Stat(filepath.Join(dir, "no.txt"))
	fmt.Println(os.IsNotExist(err)) // => true

	// ReadDir — список содержимого папки (что внутри dir).
	entries, _ := os.ReadDir(dir)
	fmt.Println(len(entries) >= 2) // => true (есть папка "a" и файл log.txt)

	// Rename — переименовать/переместить.
	os.Rename(path, filepath.Join(dir, "renamed.txt"))

	// Окружение целиком: Environ() — все переменные; Setenv/Unsetenv — задать/убрать.
	os.Setenv("DEMO_KEY", "1")
	fmt.Println(len(os.Environ()) > 0) // => true
	os.Unsetenv("DEMO_KEY")

	// Рабочая директория.
	wd, _ := os.Getwd()
	fmt.Println(wd != "") // => true
	_ = os.Chdir(wd)      // сменить рабочую папку (здесь — на ту же)

	// Lstat для полноты (на обычном файле ведёт себя как Stat).
	_, _ = os.Lstat(dir)
	fmt.Println(os.IsExist(err)) // => false (ошибка была «не существует», а не «уже есть»)
}

/*
Что важно запомнить:
  • Временные файлы/папки — CreateTemp/MkdirTemp (+ defer RemoveAll): безопаснее, чем хардкодить /tmp/xxx,
    т.к. имя уникально и нет гонок.
  • Mkdir vs MkdirAll: MkdirAll создаёт всю цепочку и не ругается, если папка уже есть. Чаще нужен он.
  • OpenFile + флаги (O_CREATE|O_APPEND|O_WRONLY) — когда нужен режим, которого нет у Open/Create:
    дозапись в лог, открытие без перезаписи и т.п.
  • Проверка существования — через ошибку Stat + os.IsNotExist(err), отдельной функции «Exists» в Go нет.
  • ReadDir отдаёт os.DirEntry (быстрее старого Readdir): имя и тип без лишнего Stat.

Задача:
  1) Создай временную папку, запиши в неё файл через OpenFile с O_APPEND (дважды), прочитай и проверь, что строки накопились.
*/
