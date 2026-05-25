/*
Пакет encoding/base64 — кодирование произвольных байтов в безопасную текстовую строку.

Зачем: передать «сырые» байты (картинку, токен, ключ) там, где можно только текст — в JSON, URL,
заголовках, cookie. Base64 превращает любые байты в строку из букв/цифр/+/-/= и обратно.

Это НЕ шифрование! Base64 легко раскодировать. Это просто способ представить байты текстом.

Как запустить:  go run main.go
*/
package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	data := []byte("привет мир")

	// StdEncoding — обычный base64 (с символами + / и паддингом =).
	enc := base64.StdEncoding.EncodeToString(data)
	fmt.Println(enc) // => 0L/RgNC40LLQtdGCINC80LjRgA==

	dec, _ := base64.StdEncoding.DecodeString(enc)
	fmt.Println(string(dec)) // => привет мир

	// URLEncoding — URL-safe вариант: вместо + / используются - _ (чтобы не ломать URL/имена файлов).
	urlEnc := base64.URLEncoding.EncodeToString([]byte{0xfb, 0xff})
	fmt.Println(urlEnc) // => -_8= (а не +/8= как в StdEncoding)

	// RawURLEncoding — URL-safe и БЕЗ паддинга '='. Идеально для токенов/JWT.
	token := base64.RawURLEncoding.EncodeToString([]byte("session-id-123"))
	fmt.Println(token) // => без '=' на конце
}

/*
Что важно запомнить:
  • Base64 — это КОДИРОВАНИЕ, не шифрование. Любой раскодирует обратно. Не клади туда секреты «для защиты».
  • Какой вариант брать:
      StdEncoding     — общий случай (вложения, бинарь в JSON).
      URLEncoding     — когда строка идёт в URL/имя файла (+ и / там опасны, заменены на - и _).
      RawURLEncoding  — то же, но без '=' на конце: стандарт для токенов, JWT, query-параметров.
  • Частая связка: crypto/rand.Read -> base64.RawURLEncoding.EncodeToString = компактный безопасный токен.
  • Для бинарного отпечатка хеша иногда удобнее hex (см. encoding/hex) — он короче читается «по парам».

Задача:
  1) Закодируй случайные 24 байта (crypto/rand) в RawURLEncoding и убедись, что в результате нет '='.
*/
