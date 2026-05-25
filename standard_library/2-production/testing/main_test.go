/*
Тесты к main.go. Демонстрируют основные возможности пакета testing.
Запуск:  go test -v   (бенчмарки: go test -bench .   фаззинг: go test -fuzz FuzzMax)
*/
package main

import "testing"

// TestMain (тип M) — точка входа ВСЕХ тестов пакета: общая настройка до и очистка после.
// Если её нет — тесты запускаются напрямую. m.Run() запускает все TestXxx.
func TestMain(m *testing.M) {
	// здесь могла бы быть общая подготовка (поднять тестовую БД и т.п.)
	m.Run()
	// здесь — общая очистка
}

// Табличный тест: список случаев + t.Run на каждый (подтест со своим именем).
func TestMax(t *testing.T) {
	cases := []struct {
		name    string
		a, b, w int
	}{
		{"a больше", 5, 2, 5},
		{"b больше", 2, 5, 5},
		{"равны", 3, 3, 3},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) { // подтест — видно в выводе как TestMax/a_больше
			assertEqual(t, Max(c.a, c.b), c.w)
		})
	}
}

// Helper — вспомогательная функция проверки. t.Helper() говорит: «при ошибке показывай номер
// строки ВЫЗОВА, а не этой функции» — иначе все ошибки указывали бы сюда.
func assertEqual(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		// Error/Errorf — пометить тест проваленным, но ПРОДОЛЖИТЬ. Fatal/Fatalf — провалить и ОСТАНОВИТЬ.
		t.Errorf("получили %d, хотим %d", got, want)
	}
}

func TestHelpers(t *testing.T) {
	t.Log("просто запись в лог теста (видна при -v)") // Log/Logf — диагностика

	dir := t.TempDir()       // временная папка, авто-удаляется после теста
	t.Setenv("MODE", "test") // переменная окружения только на время теста (вернётся обратно)
	t.Cleanup(func() {       // отложенная очистка (выполнится в конце теста)
		_ = dir
	})

	if dir == "" {
		t.Fatal("TempDir пустой") // Fatal остановил бы тест немедленно
	}
	// ВАЖНО: t.Setenv нельзя сочетать с t.Parallel в одном тесте (Go это запрещает).
}

// t.Parallel — разрешить параллельный запуск с другими Parallel-тестами (без Setenv в этом же тесте!).
func TestParallel(t *testing.T) {
	t.Parallel()
	if Max(1, 2) != 2 {
		t.Error("ожидали 2")
	}
}

func TestSkip(t *testing.T) {
	if testing.Short() {
		t.Skip("пропускаем в режиме -short") // Skip/Skipf — пропустить тест с причиной
	}
}

// Бенчмарк (тип B): измеряет скорость. b.N подбирается автоматически.
func BenchmarkMax(b *testing.B) {
	b.ReportAllocs() // показать аллокации
	b.ResetTimer()   // не учитывать подготовку выше
	for i := 0; i < b.N; i++ {
		Max(i, i+1)
	}
	b.StopTimer()  // приостановить таймер (например, для дорогой проверки)
	_ = 1 + 1
	b.StartTimer() // и снова включить
}

// Фаззинг (тип F): подаёт случайные данные и ищет вход, который ломает функцию.
// Запуск: go test -fuzz FuzzMax
func FuzzMax(f *testing.F) {
	f.Add(1, 2) // seed — стартовые примеры
	f.Fuzz(func(t *testing.T, a, b int) {
		m := Max(a, b)
		if m < a || m < b { // свойство, которое ВСЕГДА должно выполняться
			t.Fatalf("Max(%d,%d)=%d меньше аргумента", a, b, m)
		}
	})
}

/*
Что важно запомнить:
  • Файлы тестов — *_test.go; функции TestXxx(t *testing.T). Запуск: go test (не go run).
  • Табличные тесты + t.Run — главный паттерн: список случаев, на каждый свой подтест с именем.
  • Error/Errorf vs Fatal/Fatalf: Error помечает провал и ИДЁТ дальше; Fatal останавливает тест.
    В табличных тестах внутри t.Run обычно Error (чтобы проверить все случаи).
  • t.Helper() в функциях-проверках — чтобы номер строки ошибки указывал на ВЫЗОВ, а не на хелпер.
  • t.TempDir()/t.Setenv()/t.Cleanup() — авто-очистка: не надо вручную удалять и возвращать окружение.
  • t.Parallel() — параллельный прогон; b.N/ResetTimer/ReportAllocs — бенчмарки; f.Add/f.Fuzz — фаззинг.
  • Типы: T (тест), B (бенчмарк), F (фаззинг), M (TestMain — общая обвязка пакета).

Задача:
  1) Напиши функцию Min и табличный тест к ней через t.Run с тремя случаями.
*/
