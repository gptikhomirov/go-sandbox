/*
Пакет os — файлы, переменные окружения, аргументы, потоки процесса.

Зачем:   конфиги/секреты из env, работа с файлами, аргументы CLI, коды выхода.
Когда:   12-factor конфигурация (env), чтение/запись файлов, скрипты.
Грабли:  os.ReadFile грузит файл ЦЕЛИКОМ в память — для больших файлов используйте Open + потоковое
         чтение (io.Copy/bufio). Проверка «файл существует» делается через ошибку os.IsNotExist(err),
         а не отдельной функцией. os.Exit немедленно завершает процесс — defer'ы не отработают.
*/
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// LookupEnv vs Getenv: LookupEnv отличает «пусто» от «не задано» — важно для дефолтов.
	if port, ok := os.LookupEnv("PORT"); ok {
		fmt.Println("PORT:", port)
	} else {
		fmt.Println("PORT not set, using 8080") // => PORT not set, using 8080
	}

	// Args — аргументы запуска. Args[0] — путь к программе, остальное — параметры.
	fmt.Println("program:", filepath.Base(os.Args[0])) // => program: main (или имя бинарника)

	// WriteFile/ReadFile — записать и прочитать файл целиком одним вызовом (для НЕбольших файлов).
	// 0o644 — права rw-r--r--.
	tmp := filepath.Join(os.TempDir(), "go-base-demo.txt")
	defer os.Remove(tmp) // подчистим за собой

	if err := os.WriteFile(tmp, []byte("hello file\n"), 0o644); err != nil {
		fmt.Println("write error:", err)
		os.Exit(1) // фатальная ошибка скрипта → ненулевой код выхода
	}
	data, _ := os.ReadFile(tmp)
	fmt.Printf("%q\n", string(data)) // => "hello file\n"

	// Проверка существования файла через os.IsNotExist.
	_, err := os.Stat("no-such-file")
	fmt.Println(os.IsNotExist(err)) // => true

	// Stdout/Stderr — стандартные потоки как io.Writer.
	fmt.Fprintln(os.Stdout, "done") // => done
}

/*
Что запомнить (что чаще и почему):
  • ReadFile/WriteFile vs Open/Create:
      ReadFile/WriteFile — одна строка, весь файл сразу. Для конфигов, шаблонов, мелких данных. Чаще всего.
      Open/Create        — возвращают *os.File (io.Reader/Writer) для ПОТОКОВОЙ работы: большие файлы,
                           докачка, копирование через io.Copy. Берут, когда файл может не влезть в память.
  • Getenv vs LookupEnv: Getenv вернёт "" и для пустой, и для незаданной переменной. LookupEnv даёт bool —
    используйте его, когда «не задано» и «пусто» — разные случаи (выбор дефолта).
  • os.Exit vs return/panic: Exit немедленный, без defer — только в main на фатальных ошибках старта.
    В библиотеках/хендлерах возвращайте ошибку, а не вызывайте Exit.
  • Права файла: 0o644 (rw-r--r--) для обычных, 0o600 для секретов, 0o755 для исполняемых/директорий.

Типичные сценарии:
  1) Конфиг из env:  port := cmp.Or(os.Getenv("PORT"), "8080")
  2) Чтение секрета: key, err := os.ReadFile("/run/secrets/api_key")
  3) Большой файл:   f, _ := os.Open(path); defer f.Close(); io.Copy(dst, f)
*/
