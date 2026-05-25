/*
Пакет io (уровень 2) — только новое поверх base.
В base уже: Reader/Writer, ReadAll, Copy, EOF, Closer/ReadCloser.
Здесь: CopyN, ReadFull/ReadAtLeast, LimitReader, MultiReader/MultiWriter, TeeReader, Pipe,
NopCloser, Discard, доп. интерфейсы (Seeker/ReaderAt/WriterAt/ReadWriter).

Зачем: ограничивать чтение (защита от огромных тел), писать сразу в несколько мест, считать хеш
на лету (TeeReader), отбрасывать вывод (Discard).

Как запустить:  go run main.go
*/
package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
)

func main() {
	// LimitReader ограничивает чтение N байтами — защита от «бесконечного» источника.
	limited := io.LimitReader(strings.NewReader("очень длинные данные"), 6)
	data, _ := io.ReadAll(limited)
	fmt.Printf("%q\n", string(data)) // => "оче" (6 байт = 3 кириллические буквы)

	// CopyN копирует РОВНО n байт. Discard — «мусорка»: принимает данные и выбрасывает.
	n, _ := io.CopyN(io.Discard, strings.NewReader("abcdef"), 3)
	fmt.Println(n) // => 3 (скопировано 3 байта в никуда)

	// MultiWriter пишет сразу в НЕСКОЛЬКО приёмников (например, в файл И в лог).
	var a, b bytes.Buffer
	mw := io.MultiWriter(&a, &b)
	io.WriteString(mw, "оба")
	fmt.Println(a.String(), b.String()) // => оба оба

	// MultiReader склеивает несколько источников в один последовательный поток.
	mr := io.MultiReader(strings.NewReader("раз "), strings.NewReader("два"))
	all, _ := io.ReadAll(mr)
	fmt.Println(string(all)) // => раз два

	// TeeReader: читаем данные И попутно пишем их копию в другой приёмник.
	// Классика: посчитать хеш файла, не читая его дважды.
	h := sha256.New()
	tee := io.TeeReader(strings.NewReader("payload"), h)
	io.ReadAll(tee)                       // читаем данные (попутно они утекли в h)
	fmt.Printf("%x\n", h.Sum(nil)[:4])    // => первые 4 байта хеша

	// ReadFull читает РОВНО len(buf) байт или вернёт ошибку (в отличие от обычного Read).
	buf := make([]byte, 3)
	io.ReadFull(strings.NewReader("xyz"), buf)
	fmt.Println(string(buf)) // => xyz

	// Pipe — связанная пара Reader/Writer: что пишут в writer, читается из reader (как труба между горутинами).
	pr, pw := io.Pipe()
	go func() {
		io.WriteString(pw, "через трубу")
		pw.Close() // закрыть, чтобы читатель получил EOF
	}()
	piped, _ := io.ReadAll(pr)
	fmt.Println(string(piped)) // => через трубу

	// NopCloser оборачивает Reader в ReadCloser с пустым Close (когда API требует Closer, а закрывать нечего).
	rc := io.NopCloser(strings.NewReader("данные"))
	rc.Close() // ничего не делает, но удовлетворяет интерфейс
	fmt.Println("NopCloser ок")
}

/*
Что важно запомнить:
  • LimitReader — главный приём безопасности при чтении из сети/файла: «не больше N байт».
    В HTTP для тела запроса есть готовый http.MaxBytesReader.
  • TeeReader — «прочитать и одновременно скопировать»: хеширование, логирование, кэширование на лету.
  • MultiWriter — один Write -> в несколько мест (файл + stdout + метрика). MultiReader — наоборот, склейка источников.
  • io.Discard — приёмник-«мусорка»: когда данные надо прочитать/пропустить, но не сохранять.
  • Read vs ReadFull: обычный Read может вернуть МЕНЬШЕ запрошенного; ReadFull гарантирует полный буфер или ошибку.

Задача:
  1) Посчитай sha256 строки через TeeReader, одновременно скопировав её в io.Discard.
*/
