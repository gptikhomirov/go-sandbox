/*
Пакет os/exec — запуск ВНЕШНИХ программ из своего кода.

Зачем: вызвать системную утилиту (git, ffmpeg, ls), скрипт или другой бинарник и получить его вывод.
Это частая задача в CI-инструментах, деплое, обработке файлов.

Грабли: НЕ передавай команду одной строкой через шелл — это инъекция. Command(name, args...) запускает
программу напрямую, без шелла, аргументы передаются списком (безопасно).

Как запустить:  go run main.go
*/
package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func main() {
	// LookPath — найти путь к программе в PATH (проверить, установлена ли).
	if path, err := exec.LookPath("echo"); err == nil {
		fmt.Println("echo найден:", path != "") // => echo найден: true
	}

	// Command + Output — запустить и получить stdout одним вызовом (самый частый случай).
	out, err := exec.Command("echo", "привет", "мир").Output()
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Printf("вывод: %q\n", strings.TrimSpace(string(out))) // => вывод: "привет мир"

	// CommandContext — с таймаутом/отменой: если команда зависнет, context её прибьёт.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	out2, _ := exec.CommandContext(ctx, "echo", "с таймаутом").Output()
	fmt.Printf("вывод2: %q\n", strings.TrimSpace(string(out2))) // => вывод2: "с таймаутом"

	// Тонкая настройка: можно задать рабочую папку (Dir), окружение (Env), перенаправить потоки.
	cmd := exec.Command("pwd")
	cmd.Dir = "/tmp"
	dir, _ := cmd.Output()
	fmt.Printf("pwd в /tmp: %q\n", strings.TrimSpace(string(dir))) // => pwd в /tmp: "/tmp" (или /private/tmp на macOS)

	// CombinedOutput — stdout + stderr вместе (удобно для отладки/логов команды).
	combined, _ := exec.Command("echo", "stdout+stderr").CombinedOutput()
	fmt.Printf("combined: %q\n", strings.TrimSpace(string(combined)))

	// Run = Start + Wait: запустить и дождаться. Вывод можно направить в свой буфер (cmd.Stdout).
	var sb strings.Builder
	runCmd := exec.Command("echo", "через Run")
	runCmd.Stdout = &sb
	runCmd.Run()
	fmt.Printf("run: %q\n", strings.TrimSpace(sb.String())) // => run: "через Run"

	// Start (не ждёт) + Wait (дожидается) — для асинхронного запуска / параллельных команд.
	asyncCmd := exec.Command("echo", "async")
	asyncCmd.Start()
	asyncCmd.Wait()
	fmt.Println("async-команда завершена")
}

/*
Что важно запомнить:
  • Command(name, args...) запускает программу НАПРЯМУЮ (без шелла) — аргументы списком, не одной строкой.
    Это защищает от инъекций. Хочешь пайпы/шелл-фичи — явно зови "sh", "-c", "cmd" (и осознавай риск).
  • Способы получить результат: Output() — только stdout; CombinedOutput() — stdout+stderr; либо назначить
    cmd.Stdout/Stderr своими io.Writer и звать Run() (запустить и дождаться) или Start()+Wait() (асинхронно).
  • CommandContext — ВСЕГДА предпочтительно для внешних вызовов: даёт таймаут/отмену, иначе зависшая команда повесит сервис.
  • Настройки до запуска: cmd.Dir (рабочая папка), cmd.Env (окружение), cmd.Stdin (подать данные на вход).

Задача:
  1) Запусти "go version" через exec.Command(...).Output() и выведи результат.
*/
