/*
Пакет crypto/rand — КРИПТОСТОЙКИЕ случайные числа.

Важно отличать от math/rand: math/rand предсказуем (для игр, перемешивания, тестов).
crypto/rand непредсказуем — для ВСЕГО, что связано с безопасностью: токены, пароли, соли, ключи, session id.

Зачем: сгенерировать токен сброса пароля, API-ключ, соль для хеша — так, чтобы его нельзя было угадать.

Как запустить:  go run main.go
*/
package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

func main() {
	// Read заполняет байты криптослучайными данными. Это основа генерации токенов.
	token := make([]byte, 16)
	if _, err := rand.Read(token); err != nil {
		fmt.Println("ошибка:", err) // на проде ошибку crypto/rand нельзя игнорировать
		return
	}
	// Сырые байты обычно кодируют в hex или base64, чтобы получить строку-токен.
	fmt.Println("длина токена:", len(hex.EncodeToString(token))) // => длина токена: 32 (16 байт = 32 hex-символа)

	// Int — криптослучайное число в диапазоне [0, max). Например, для одноразового кода.
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	code := fmt.Sprintf("%06d", n) // 6-значный код
	fmt.Println("код длиной:", len(code)) // => код длиной: 6

	// rand.Reader — глобальный источник; его передают в Int, в tls, в генераторы ключей.
	buf := make([]byte, 4)
	rand.Reader.Read(buf)
	fmt.Println("случайных байт:", len(buf)) // => случайных байт: 4
}

/*
Что важно запомнить:
  • crypto/rand vs math/rand — ГЛАВНОЕ различие по безопасности:
      math/rand   — быстрый, ПРЕДСКАЗУЕМЫЙ. Игры, перемешивание, тесты, семплирование.
      crypto/rand — непредсказуемый. Токены, пароли, соли, ключи, session id, CSRF, OTP.
    Спутать — серьёзная уязвимость (предсказуемые токены можно подобрать).
  • Ошибку rand.Read НЕЛЬЗЯ игнорировать: если источник энтропии недоступен, токен будет небезопасен.
  • Сырые байты -> строка через hex.EncodeToString или base64.RawURLEncoding (для URL-safe токенов).
  • Для случайного числа в диапазоне — rand.Int(rand.Reader, max), а не «math/rand по модулю».

Задача:
  1) Сгенерируй URL-safe токен из 32 байт (crypto/rand.Read + base64.RawURLEncoding.EncodeToString).
*/
