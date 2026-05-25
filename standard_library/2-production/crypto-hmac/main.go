/*
Пакет crypto/hmac — ПОДПИСЬ данных секретным ключом (HMAC).

HMAC отвечает на вопрос «эти данные точно от того, у кого есть секрет, и их не подменили?».
В отличие от простого хеша, HMAC использует секретный ключ — подделать подпись без ключа нельзя.

Зачем: подпись webhook-ов (GitHub/Stripe шлют X-Signature), API-токенов, cookie, защита от подмены.

Как запустить:  go run main.go
*/
package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// sign считает HMAC-SHA256 подпись сообщения секретным ключом.
func sign(message, key []byte) []byte {
	h := hmac.New(sha256.New, key) // HMAC поверх SHA-256
	h.Write(message)
	return h.Sum(nil)
}

func main() {
	key := []byte("секретный-ключ")
	msg := []byte(`{"event":"payment","amount":100}`)

	// Отправитель считает подпись и прикладывает её к запросу.
	sig := sign(msg, key)
	fmt.Println("подпись:", hex.EncodeToString(sig)[:16], "...") // => подпись: ... (первые символы)

	// Получатель считает подпись сам и сравнивает с присланной — через hmac.Equal (constant-time!).
	got := sign(msg, key)
	fmt.Println(hmac.Equal(sig, got)) // => true (подпись верна, данные не подменены)

	// Если данные подменили — подпись не сойдётся.
	tampered := sign([]byte(`{"event":"payment","amount":999999}`), key)
	fmt.Println(hmac.Equal(sig, tampered)) // => false

	// С чужим ключом тоже не сойдётся.
	wrongKey := sign(msg, []byte("другой-ключ"))
	fmt.Println(hmac.Equal(sig, wrongKey)) // => false
}

/*
Что важно запомнить:
  • HMAC = хеш + СЕКРЕТНЫЙ КЛЮЧ. Доказывает и целостность (не подменили), и подлинность (от владельца ключа).
    Простой sha256 целостность проверяет, но подделать его может кто угодно — для подписи нужен именно HMAC.
  • Сравнивай подписи ТОЛЬКО через hmac.Equal — это constant-time сравнение (см. crypto/subtle).
    Обычное == сравнивает побайтово с ранним выходом и открывает timing-атаку на подбор подписи.
  • Типичный сценарий: проверка входящего webhook — посчитать HMAC тела присланным секретом и сравнить
    с заголовком-подписью через hmac.Equal.

Задача:
  1) Сделай функцию verify(msg, key, sig) bool, которая считает HMAC и сравнивает через hmac.Equal.
*/
