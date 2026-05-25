/*
Пакет time — время, длительности, измерение интервалов.

Зачем:   таймстампы записей, форматирование дат в API, измерение длительности операций.
Когда:   created_at/updated_at, разбор дат из запросов, TTL токенов/кэша, метрики.
Грабли:  layout форматирования — это НЕ strftime, а конкретная дата "2006-01-02 15:04:05"
         (1 2 3 4 5 6 7 по порядку). Сравнивать время через == нельзя (учитывает монотонные
         часы и локацию) — используйте Equal/Before/After.
*/
package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()

	// Format использует «магическую» референсную дату: Mon Jan 2 15:04:05 2006.
	// RFC3339 — стандарт для API и логов.
	fmt.Println(now.Format("2006-01-02 15:04")) // => 2026-05-25 15:49 (пример)

	// Parse — обратная операция: строка → time.Time (дата из query-параметра).
	t, err := time.Parse("2006-01-02", "2024-03-15")
	if err != nil {
		fmt.Println("parse error:", err)
	}
	fmt.Println(t.Year(), t.Month(), t.Day()) // => 2024 March 15

	// Арифметика: Add сдвигает, Sub даёт разницу (Duration).
	expires := now.Add(24 * time.Hour) // TTL токена/кэша
	fmt.Println(t.Before(now))         // => true  (2024 раньше «сейчас»)
	_ = expires

	// Измерение длительности — классический паттерн Since(start) для метрик/профайла.
	start := time.Now()
	time.Sleep(10 * time.Millisecond) // имитация работы
	elapsed := time.Since(start)
	fmt.Println(elapsed.Milliseconds() >= 10) // => true
}

/*
Что запомнить (что чаще и почему):
  • Длительность задают умножением на константу: 5*time.Second, 24*time.Hour. НЕ передавайте
    «голое» число — time.Sleep(5) это 5 наносекунд, частая ошибка.
  • time.Since(start) vs time.Now().Sub(start): Since — короткая запись того же самого, её и берут.
  • Format/Parse используют ОДИН и тот же layout-шаблон. Запомните RFC3339 (time.RFC3339) —
    им сериализуют время в JSON/API по умолчанию.
  • Сравнение моментов: Before/After/Equal, а не </>/==. UTC() приводит к единой зоне перед сравнением/хранением.
  • В БД таймстампы храните в UTC (now.UTC()), форматируйте в локальную зону только при выводе.

Типичные сценарии:
  1) Запись в БД:   user.CreatedAt = time.Now().UTC()
  2) TTL токена:    exp := time.Now().Add(15 * time.Minute)
  3) Метрика:       defer func(s time.Time){ log.Println(time.Since(s)) }(time.Now())
*/
