/*
Пакет time (уровень 2) — только новое поверх base.
В base уже: Now, Since/Until, Sleep, Parse, Format, Add/Sub, Before/After, Seconds/Milliseconds.
Здесь: Date/Unix (создать момент), ParseDuration, таймеры/тикеры (After/NewTimer/NewTicker),
часовые пояса (UTC/Local/In/LoadLocation), Truncate/Round, Unix()/UnixMilli, Duration.String/Hours/Minutes.

Как запустить:  go run main.go
*/
package main

import (
	"fmt"
	"time"
)

func main() {
	// Date — собрать момент вручную. Unix — из секунд эпохи (timestamp).
	d := time.Date(2024, time.March, 15, 10, 30, 0, 0, time.UTC)
	fmt.Println(d.Format("2006-01-02 15:04")) // => 2024-03-15 10:30
	fmt.Println(d.Unix())                      // => 1710498600 (секунд с 1970)
	fmt.Println(time.Unix(0, 0).UTC().Year())  // => 1970

	// ParseDuration читает человекочитаемую длительность из строки (конфиги, флаги).
	dur, _ := time.ParseDuration("1h30m")
	fmt.Println(dur.Hours())    // => 1.5
	fmt.Println(dur.String())   // => 1h30m0s
	fmt.Println(dur.Minutes())  // => 90

	// Часовые пояса: UTC()/Local()/In(loc). LoadLocation грузит зону по имени.
	msk, err := time.LoadLocation("Europe/Moscow")
	if err == nil {
		fmt.Println(d.In(msk).Hour() - d.UTC().Hour()) // => 3 (Москва +3 к UTC)
		// Equal — правильное сравнение моментов (учитывает зоны), в отличие от ==.
		fmt.Println(d.Equal(d.In(msk))) // => true (один и тот же миг, разные зоны)
	}

	// Truncate/Round — округление времени (например, «начало часа» для группировки метрик).
	fmt.Println(d.Truncate(time.Hour).Format("15:04")) // => 10:00 (вниз до начала часа)
	fmt.Println(d.Round(time.Hour).Format("15:04"))    // => 11:00 (к ближайшему часу: 10:30 -> 11:00)

	// Тикер: After даёт канал, который «выстрелит» через время. NewTicker — повторяющийся.
	select {
	case <-time.After(10 * time.Millisecond):
		fmt.Println("сработал таймер") // => сработал таймер
	}
}

/*
Что важно запомнить:
  • Создать момент: Date(...) вручную или Unix(sec, nsec) из timestamp. Format/Parse — преобразование с текстом (см. base).
  • ParseDuration("1h30m") — стандарт для таймаутов в конфигах/флагах. Обратно — Duration.String().
  • Хранить и сравнивать время лучше в UTC; In(loc)/Local() — только для ОТОБРАЖЕНИЯ пользователю.
  • Equal vs ==: для моментов используй Equal (== учитывает монотонные часы и представление, легко ошибиться).
  • After/NewTimer — одноразовый сигнал; NewTicker — периодический (не забудь ticker.Stop(), иначе утечёт).
  • Truncate(time.Hour) — обрезать до начала часа/дня: частый приём для агрегации по времени.

Задача:
  1) Из строки "250ms" получи Duration и передай в time.Sleep. Выведи dur.Milliseconds().
*/
