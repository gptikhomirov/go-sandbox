/*
Пакет crypto/tls — шифрованные соединения (HTTPS работает поверх TLS).

В большинстве случаев TLS настраивается «сам»: http.ListenAndServeTLS на сервере и http.Get("https://...")
на клиенте уже шифруют трафик. Пакет crypto/tls нужен, когда требуется ТОНКАЯ настройка:
свой сертификат, минимальная версия протокола, mTLS (взаимная аутентификация), TLS поверх своего протокола.

Этот пример НЕ ходит в сеть и не требует сертификатов — он показывает, как СОБИРАЮТ конфигурацию.

Как запустить:  go run main.go
*/
package main

import (
	"crypto/tls"
	"fmt"
)

func main() {
	// tls.Config — центральный объект настройки. Его передают в сервер/клиент.
	clientCfg := &tls.Config{
		MinVersion: tls.VersionTLS12, // не принимать устаревшие небезопасные версии
		ServerName: "api.example.com", // имя для проверки сертификата сервера
	}
	fmt.Println(clientCfg.MinVersion == tls.VersionTLS12) // => true

	// Серверный конфиг обычно содержит сертификат. Его грузят из файлов через LoadX509KeyPair:
	//   cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	//   serverCfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	// Здесь файлов нет — показываем, что вызов существует и корректно вернёт ошибку.
	_, err := tls.LoadX509KeyPair("нет.crt", "нет.key")
	fmt.Println(err != nil) // => true (файлов нет — это ожидаемо)

	// Подключение к TLS-серверу (требует сеть, поэтому только в комментарии):
	//   conn, err := tls.Dial("tcp", "api.example.com:443", clientCfg)
	//   defer conn.Close()
	//
	// А обычный HTTPS вообще не требует ручного tls: http.Get("https://...") сделает всё сам.

	fmt.Println("TLS-конфиг собран") // => TLS-конфиг собран
}

/*
Что важно запомнить:
  • Чаще всего ручной crypto/tls НЕ нужен: http.ListenAndServeTLS(cert, key, h) и http.Get("https://...")
    уже шифруют. Пакет берут для тонкой настройки.
  • tls.Config — что важно задавать: MinVersion (>= TLS 1.2), ServerName (проверка имени), Certificates (свой серт),
    ClientAuth (для mTLS — требовать сертификат и от клиента).
  • LoadX509KeyPair(certFile, keyFile) — загрузить пару серт+ключ для сервера/клиента.
  • НИКОГДА не ставь InsecureSkipVerify: true в проде — это отключает проверку сертификата (man-in-the-middle).
    Допустимо только в локальных тестах с самоподписанными сертами.
  • tls.Dial — TLS поверх голого TCP (свои протоколы); для HTTP всё делает net/http.

Задача:
  1) Собери tls.Config с MinVersion TLS 1.3 (tls.VersionTLS13) и распечатай выбранную версию.
*/
