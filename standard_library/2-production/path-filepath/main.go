/*
Пакет path/filepath (уровень 2) — только новое поверх base.
В base уже: Join, Base/Dir/Ext, Abs, WalkDir.
Здесь: Clean/Split/Rel/IsAbs, Glob/Match, ToSlash/FromSlash.

Зачем: нормализовать пути (Clean — против "../" атак и грязного ввода), строить относительные пути,
искать файлы по маске (Glob), приводить разделители к "/" для URL/хранения.

Как запустить:  go run main.go
*/
package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	// Clean нормализует путь: убирает лишние "/", "." и схлопывает "..".
	// ВАЖНО для безопасности: чистить пользовательский путь перед доступом к файлу.
	fmt.Println(filepath.Clean("/var/../etc/./passwd")) // => /etc/passwd

	// Split делит путь на папку и имя (в отличие от Dir/Base — отдаёт оба сразу).
	dir, file := filepath.Split("/var/log/app.txt")
	fmt.Printf("%q %q\n", dir, file) // => "/var/log/" "app.txt"

	// Rel строит путь ОТ одного места К другому (относительный).
	rel, _ := filepath.Rel("/var/log", "/var/log/app/x.txt")
	fmt.Println(rel) // => app/x.txt

	// IsAbs — абсолютный ли путь.
	fmt.Println(filepath.IsAbs("/etc"), filepath.IsAbs("etc")) // => true false

	// Match проверяет имя по маске (* ? [..]). Glob ищет реальные файлы по маске на диске.
	ok, _ := filepath.Match("*.txt", "app.txt")
	fmt.Println(ok) // => true
	matches, _ := filepath.Glob("/tmp/*") // вернёт реальные пути в /tmp (или пусто)
	fmt.Println(matches != nil || matches == nil) // => true (просто показываем, что вызов отработал)

	// ToSlash приводит разделители к "/" (для URL, хранения в БД, кроссплатформенности).
	fmt.Println(filepath.ToSlash("a\\b\\c")) // на Windows => a/b/c (на Unix без изменений)
}

/*
Что важно запомнить:
  • Clean — обязателен для путей из пользовательского ввода (загрузки, параметры): схлопывает "..",
    защищая от выхода за пределы папки (path traversal). Часто комбинируют: filepath.Clean(filepath.Join(base, userPath)).
  • Match (по строке-маске) vs Glob (реальный поиск на диске). Match не ходит в ФС, Glob ходит.
  • Split отдаёт (dir, file) сразу; Dir/Base — по отдельности. Выбор по удобству.
  • ToSlash/FromSlash — когда путь хранится/передаётся в едином формате "/" независимо от ОС.

Задача:
  1) Возьми base="/data" и пользовательский userPath="../secret" — собери безопасный путь через
     filepath.Clean(filepath.Join(base, userPath)) и посмотри, куда он ведёт.
*/
