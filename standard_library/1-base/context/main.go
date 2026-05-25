/*
Пакет context — отмена, таймауты и передача request-scoped значений.

Зачем:   ограничить время операции, отменить работу при разрыве соединения, протащить request_id.
Когда:   запросы к БД/внешним API (всегда с таймаутом), graceful-остановка воркеров.
Грабли:  ctx — ВСЕГДА первый аргумент функции. WithCancel/WithTimeout возвращают cancel —
         его нужно вызвать (defer cancel()), иначе утечёт горутина/таймер. WithValue — только
         для метаданных запроса, НЕ для передачи обязательных параметров.
*/
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// slowOperation имитирует долгую работу, которая УВАЖАЕТ отмену: select между
// результатом и ctx.Done(). Так и пишут обращения к внешним сервисам.
func slowOperation(ctx context.Context) error {
	select {
	case <-time.After(200 * time.Millisecond): // «работа» заняла 200мс
		return nil
	case <-ctx.Done(): // контекст отменён или истёк таймаут
		return ctx.Err() // DeadlineExceeded или Canceled
	}
}

func main() {
	// WithTimeout — ограничить операцию по времени. Главный приём для внешних вызовов.
	// defer cancel() обязателен — освобождает ресурсы, даже если уложились в срок.
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := slowOperation(ctx) // 200мс работы при лимите 50мс → таймаут
	fmt.Println(errors.Is(err, context.DeadlineExceeded)) // => true  → отдать 504

	// WithCancel — ручная отмена (остановить воркеров по сигналу/событию).
	ctx2, cancel2 := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel2() // отменяем извне
	}()
	<-ctx2.Done()
	fmt.Println(ctx2.Err()) // => context canceled

	// WithValue — request-scoped данные (request_id, user). Только метаданные запроса.
	ctx3 := context.WithValue(context.Background(), "requestID", "abc-123")
	fmt.Println(ctx3.Value("requestID")) // => abc-123
}

/*
Что запомнить (что чаще и почему):
  • WithTimeout vs WithDeadline vs WithCancel:
      WithTimeout  — «отмени через N» — чаще всего (запросы к БД/HTTP). Внутри это WithDeadline(now+N).
      WithDeadline — «отмени к моменту T» — когда дедлайн известен абсолютно (конец окна обработки).
      WithCancel   — отмена «по событию», без времени (остановка пайплайна, graceful shutdown).
  • defer cancel() ставьте СРАЗУ после создания контекста — забытый cancel это утечка ресурсов
    (go vet это даже подсказывает).
  • Background() vs TODO(): Background — корневой контекст в main/тестах; TODO — заглушка,
    «контекст сюда ещё не прокинут». В реальном коде контекст приходит сверху (r.Context() в хендлере).
  • WithValue — для request_id, трассировки, юзера из middleware. НЕ кладите туда то, что должно
    быть явным аргументом функции.

Типичные сценарии:
  1) Запрос к БД:   ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second); defer cancel()
  2) Отмена в хендлере: всё, что приняли r.Context(), отменится при разрыве соединения клиентом.
  3) Воркер-пул:    <-ctx.Done() // выйти, когда контекст отменён
*/
