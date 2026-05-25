/*
Пакет net — сеть на уровне TCP/UDP/DNS (ниже, чем net/http).

net/http покрывает 95% веба. Пакет net нужен, когда работаешь НИЖЕ HTTP: свой протокол поверх TCP,
UDP, ручной разбор адресов, DNS-резолвинг. Также его типы (net.Conn, net.Listener) лежат в основе http.

Зачем: TCP-серверы/клиенты своих протоколов, разбор "host:port", проверка IP, DNS-запросы.

Как запустить:  go run main.go
*/
package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	// Разбор адресов (без сети): SplitHostPort / JoinHostPort / ParseIP.
	host, port, _ := net.SplitHostPort("127.0.0.1:8080")
	fmt.Println(host, port)                       // => 127.0.0.1 8080
	fmt.Println(net.JoinHostPort("::1", "443"))   // => [::1]:443 (IPv6 берётся в скобки)
	fmt.Println(net.ParseIP("10.0.0.1") != nil)   // => true
	fmt.Println(net.ParseIP("не-ip") == nil)      // => true (некорректный IP -> nil)

	// --- Локальный TCP: Listen (сервер) + Dial (клиент) ---
	// Listen на :0 — система сама выберет свободный порт.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("listen:", err)
		return
	}
	defer ln.Close()

	// Сервер в горутине: принять одно соединение и ответить «эхом» с приставкой.
	go func() {
		conn, err := ln.Accept() // ждём подключения; conn — это net.Conn (io.Reader+Writer)
		if err != nil {
			return
		}
		defer conn.Close()
		line, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Fprintf(conn, "эхо: %s", line)
	}()

	// Клиент: Dial подключается к адресу слушателя.
	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		fmt.Println("dial:", err)
		return
	}
	defer conn.Close()
	fmt.Fprintf(conn, "привет\n")              // отправить
	reply, _ := bufio.NewReader(conn).ReadString('\n') // получить
	fmt.Print(reply)                            // => эхо: привет

	// DNS-резолвинг (требует доступности резолвера; localhost обычно резолвится из /etc/hosts).
	if ips, err := net.LookupHost("localhost"); err == nil {
		fmt.Println("localhost ->", len(ips) > 0) // => localhost -> true
	}
}

/*
Что важно запомнить:
  • Для HTTP бери net/http; net нужен для своих протоколов поверх TCP/UDP и низкоуровневых задач.
  • net.Conn — это io.Reader + io.Writer + Close: с соединением работаешь как с потоком (bufio, io.Copy и т.д.).
  • Listen("tcp", "127.0.0.1:0") -> система выберет порт; реальный адрес узнаёшь через ln.Addr().
  • SplitHostPort/JoinHostPort корректно обрабатывают IPv6 (скобки) — не парси "host:port" вручную.
  • LookupHost/LookupIP — DNS; в проде давай им context (net.Resolver) и таймаут, чтобы не зависнуть.

Задача:
  1) Подними TCP-сервер, который отвечает текущим временем на любое подключение, и обратись к нему клиентом.
*/
