/*
Пакет path/filepath — пути файловой системы (с учётом разделителя ОС).

Зачем:   собирать пути кроссплатформенно, доставать имя/расширение, рекурсивно обходить каталоги.
Когда:   работа с файлами и загрузками, поиск миграций/ассетов, нормализация путей из ввода.
Грабли:  не путайте с пакетом path — тот всегда использует "/" и нужен для URL, а filepath
         учитывает ОС (\ на Windows). НИКОГДА не склеивайте пути через "+" или вручную — только Join.
*/
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Join — собирает путь правильно: подставит / или \, уберёт лишние разделители.
	p := filepath.Join("var", "log", "app", "server.log")
	fmt.Println(p) // => var/log/app/server.log (на Unix)

	// Base/Dir/Ext — разбор пути. Частый случай при обработке загруженных файлов.
	fmt.Println(filepath.Base(p)) // => server.log
	fmt.Println(filepath.Dir(p))  // => var/log/app
	fmt.Println(filepath.Ext(p))  // => .log

	// Abs — привести относительный путь к абсолютному (от рабочей директории).
	abs, _ := filepath.Abs("config.yaml")
	fmt.Println(filepath.IsAbs(abs)) // => true

	// WalkDir — рекурсивный обход дерева. Типичная задача: собрать список миграций.
	dir, _ := os.MkdirTemp("", "walk-demo")
	defer os.RemoveAll(dir)
	os.WriteFile(filepath.Join(dir, "001_init.sql"), nil, 0o644)
	os.WriteFile(filepath.Join(dir, "002_users.sql"), nil, 0o644)

	var sqls []string
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".sql" {
			sqls = append(sqls, filepath.Base(path))
		}
		return nil
	})
	fmt.Println(sqls) // => [001_init.sql 002_users.sql]
}

/*
Что запомнить (что чаще и почему):
  • filepath vs path: filepath — для ФАЙЛОВ (учитывает ОС-разделитель), path — для URL/слэш-путей.
    Работаете с диском → filepath; режете URL-путь → path. Их легко перепутать.
  • Join vs ручная склейка: Join сам разрулит разделители и ".." — склейка "+"/Sprintf
    ломается на Windows и при двойных слэшах. Всегда Join.
  • WalkDir vs устаревший Walk: WalkDir (Go 1.16+) быстрее — не делает Stat на каждый файл,
    отдаёт os.DirEntry. В новом коде берут WalkDir.
  • Base/Dir/Ext — стандартный разбор имени файла загрузки: имя + расширение для валидации типа.

Типичные сценарии:
  1) Путь к файлу:   p := filepath.Join(uploadDir, userID, fileName)
  2) Расширение:     if filepath.Ext(name) != ".csv" { reject() }
  3) Список миграций: filepath.WalkDir(migrationsDir, collectSQL)
*/
