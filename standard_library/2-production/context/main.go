/*
Пакет context (уровень 2) — только новое поверх base.
В base уже: Background, TODO, WithTimeout, WithCancel, Done, Err, Value.
Здесь: WithDeadline, WithValue (подробнее), WithCancelCause/Cause (отмена с причиной),
WithoutCancel, AfterFunc, ctx.Deadline.

Зачем: отмена к конкретному МОМЕНТУ (а не «через N»), отмена с понятной ПРИЧИНОЙ,
запуск действия при отмене (AfterFunc) и «отвязка» от родительской отмены (WithoutCancel).

Как запустить:  go run main.go
*/
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	// WithDeadline — отмена к КОНКРЕТНОМУ времени (а WithTimeout — через интервал).
	deadline := time.Now().Add(20 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	if dl, ok := ctx.Deadline(); ok {
		fmt.Println("дедлайн установлен:", dl.After(time.Now())) // => дедлайн установлен: true
	}

	// WithCancelCause — отмена с ПРИЧИНОЙ. Cause(ctx) вернёт переданную ошибку (понятнее, чем просто "canceled").
	ctx2, cancel2 := context.WithCancelCause(context.Background())
	cancel2(errors.New("пользователь отменил"))
	fmt.Println(ctx2.Err())            // => context canceled (общая причина)
	fmt.Println(context.Cause(ctx2))   // => пользователь отменил (конкретная причина)

	// AfterFunc запускает функцию, КОГДА контекст отменится/истечёт. Удобно для очистки.
	done := make(chan struct{})
	ctx3, cancel3 := context.WithCancel(context.Background())
	context.AfterFunc(ctx3, func() {
		close(done) // сработает после cancel3()
	})
	cancel3()
	<-done
	fmt.Println("AfterFunc сработал после отмены") // => AfterFunc сработал после отмены

	// WithoutCancel — производный контекст, который НЕ отменится вместе с родителем.
	// Нужно, когда фоновую задачу нельзя обрывать вместе с запросом (например, дописать аудит-лог).
	parent, pcancel := context.WithCancel(context.Background())
	detached := context.WithoutCancel(parent)
	pcancel()                          // отменяем родителя
	fmt.Println(parent.Err() != nil)   // => true (родитель отменён)
	fmt.Println(detached.Err() == nil) // => true (а этот продолжает жить)
}

/*
Что важно запомнить:
  • WithTimeout vs WithDeadline: «через N» против «к моменту T». WithTimeout = WithDeadline(now+N).
  • WithCancelCause + Cause(ctx) — отмена с ПРИЧИНОЙ: ctx.Err() остаётся общим (canceled/deadline),
    а Cause возвращает твою конкретную ошибку. Лучше для диагностики, чем «голый» Canceled.
  • AfterFunc — «сделай это при отмене» без ручного select по ctx.Done() в отдельной горутине.
  • WithoutCancel — отвязать фоновую работу от отмены запроса (аудит, дозапись), чтобы её не оборвало.

Задача:
  1) Сделай WithTimeout на 50мс, через WithCancelCause отмени его раньше с причиной errors.New("rate limited"),
     и выведи context.Cause.
*/
