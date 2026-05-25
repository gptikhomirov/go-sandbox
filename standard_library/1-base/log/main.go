/*
Пакет log — простое логирование в stderr.

Зачем:   стартовые и фатальные сообщения, быстрый лог без зависимостей.
Когда:   CLI-утилиты, маленькие сервисы, скрипты, ранний этап проекта.
Грабли:  log.Fatal вызывает os.Exit(1) — defer'ы НЕ отработают; в хендлерах не использовать,
         иначе один битый запрос убьёт весь процесс. Для прод-сервисов берут log/slog.
*/
package main

import (
	"log"
	"os"
)

func main() {
	// Flags задают метаданные строки. Lshortfile (файл:строка) очень помогает в отладке.
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[app] ")

	// Print/Printf/Println — обычный лог уровня «инфо».
	log.Println("service starting")        // => [app] 2026/... main.go:NN: service starting
	log.Printf("listening on port %d", 8080)

	// Свой logger для отдельного назначения (лог доступа, отдельный файл) — через log.New.
	access := log.New(os.Stdout, "[access] ", log.LstdFlags)
	access.Printf("GET /users 200") // => [access] 2026/... GET /users 200

	// Fatal = лог + os.Exit(1). Только на СТАРТЕ, когда продолжать невозможно
	// (не открылась БД, не прочитан конфиг). Здесь — демонстрация удачного старта.
	if err := os.Setenv("OK", "1"); err != nil {
		log.Fatalf("cannot set env: %v", err) // прервало бы программу с кодом 1
	}
	log.Println("startup finished")
}

/*
Что запомнить (что чаще и почему):
  • Println/Printf — 90% использования: обычные информационные сообщения.
  • Fatal vs Panic vs обычный лог:
      обычный  — штатные события, программа продолжает работу;
      Fatal    — фатальная ошибка ИНИЦИАЛИЗАЦИИ (выход с кодом 1, без defer);
      Panic    — почти не нужен для логов; паника + лог. Используйте только для «не должно случиться».
  • log vs log/slog (уровень production):
      log  — простой текст, без уровней и полей. Хорош для CLI и прототипов.
      slog — структурные логи (JSON, уровни, поля), стандарт для сервисов в проде.

Типичные сценарии:
  1) Старт сервиса:  if err := db.Ping(); err != nil { log.Fatalf("db: %v", err) }
  2) Инфо-событие:   log.Printf("processed %d items", n)
  3) В хендлере:     логируйте и возвращайте ошибку клиенту — НЕ Fatal.
*/
