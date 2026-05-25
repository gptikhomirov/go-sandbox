/*
Пакет log/slog — СТРУКТУРНОЕ логирование (новый пакет, стандарт для прода).

Чем отличается от log (из base): обычный log пишет просто текст. slog пишет лог как набор
ПОЛЕЙ (ключ=значение) и поддерживает УРОВНИ (Debug/Info/Warn/Error). Такие логи легко искать
и фильтровать в системах вроде Grafana/Kibana, особенно в формате JSON.

Зачем: в проде логи читают не глазами, а машинами — им нужны поля, а не строки.

Как запустить:  go run main.go
*/
package main

import (
	"log/slog"
	"os"
)

func main() {
	// Логи уровня через стандартный logger (пишет в stderr в текстовом формате).
	// Аргументы после сообщения — ПАРЫ ключ-значение.
	slog.Info("сервер запущен", "port", 8080, "env", "dev")
	// => 2026/05/25 ... INFO сервер запущен port=8080 env=dev
	//    (пакетный slog.Info по умолчанию пишет через log; явный формат задаёт handler ниже)
	slog.Warn("мало памяти", "free_mb", 120)
	slog.Error("ошибка БД", "err", "connection refused")

	// JSON-формат — то, что нужно в проде (машиночитаемо). New создаёт logger с нужным handler.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("заказ создан", "order_id", 42, "amount", 1999)
	// => {"time":"...","level":"INFO","msg":"заказ создан","order_id":42,"amount":1999}

	// With — «приклеить» поля ко всем последующим записям logger'а (например, request_id).
	reqLog := logger.With("request_id", "abc-123")
	reqLog.Info("обработка началась")
	reqLog.Info("обработка завершена", "ms", 12)
	// у обеих записей будет request_id=abc-123

	// Типизированные атрибуты (slog.Int, slog.String, slog.Any, ...) — точный тип/производительность.
	logger.LogAttrs(nil, slog.LevelInfo, "метрика",
		slog.Int("count", 5), slog.String("unit", "rps"), slog.Any("extra", []int{1, 2}))

	// WithGroup/Group — вложенные поля (группировка под общим ключом).
	logger.WithGroup("http").Info("запрос", "method", "GET", "code", 200)
	logger.Info("событие", slog.Group("user", "id", 7, "role", "admin"))

	// SetDefault делает наш logger глобальным — тогда пакетные slog.Info(...) пойдут через него.
	slog.SetDefault(logger)
	slog.Info("теперь и это в JSON")

	// Enabled — проверить, включён ли уровень (чтобы не тратить силы на дорогой лог).
	if slog.Default().Enabled(nil, slog.LevelDebug) {
		slog.Debug("подробности")
	}
}

/*
Что важно запомнить:
  • slog vs log: log — простой текст без полей и уровней (ок для CLI/прототипа). slog — поля + уровни + JSON,
    стандарт для серверов. На уровне production логируют через slog.
  • Аргументы после msg идут ПАРАМИ ключ-значение: slog.Info("msg", "key", value, "key2", value2).
    Нечётное число аргументов — частая ошибка (повиснет "!BADKEY").
  • Handler решает ФОРМАТ: NewTextHandler (читать глазами, локально) vs NewJSONHandler (прод, машины).
  • logger.With(...) — приклеить постоянные поля (request_id, user_id) ко всем записям: основа трассировки запроса.
  • LogAttrs + slog.Int/String/... — типизированный и более быстрый путь (без аллокаций на пары interface{}).

Задача:
  1) Сделай JSON-logger, добавь через With поле service="api" и залогируй два разных события.
*/
