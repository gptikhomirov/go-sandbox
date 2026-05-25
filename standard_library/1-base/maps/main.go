/*
Пакет maps — generic-операции над map (Go 1.21+).

Зачем:   копировать/сливать map, сравнивать, доставать ключи в предсказуемом порядке.
Когда:   слияние «дефолты + оверрайды» в конфигах, сравнение состояний, вывод ключей по порядку.
Грабли:  порядок обхода map в Go СЛУЧАЕН by design — нельзя полагаться на него; чтобы вывести
         стабильно, соберите ключи и отсортируйте. map передаётся по ссылке — нужен Clone для копии.
*/
package main

import (
	"fmt"
	"slices"

	"maps"
)

func main() {
	defaults := map[string]int{"timeout": 30, "retries": 3}

	// Clone — независимая копия перед изменением (map, как и срез, передаётся по ссылке).
	cfg := maps.Clone(defaults)
	cfg["timeout"] = 60
	fmt.Println(defaults["timeout"], cfg["timeout"]) // => 30 60  (исходник не тронут)

	// Copy — слияние: значения из src перезаписывают dst. Паттерн «дефолты + оверрайды».
	maps.Copy(cfg, map[string]int{"retries": 5, "workers": 8})
	fmt.Println(cfg) // => map[retries:5 timeout:60 workers:8]

	// Equal — сравнить две map целиком (ключи + значения).
	fmt.Println(maps.Equal(defaults, map[string]int{"timeout": 30, "retries": 3})) // => true

	// Keys возвращает ИТЕРАТОР. Чтобы получить стабильный вывод — Collect + Sort.
	// Это решает проблему случайного порядка обхода.
	keys := slices.Collect(maps.Keys(cfg))
	slices.Sort(keys)
	fmt.Println(keys) // => [retries timeout workers]

	// DeleteFunc — удалить пары по условию (например, обнулённые).
	maps.DeleteFunc(cfg, func(k string, v int) bool { return v == 0 })
	fmt.Println(len(cfg)) // => 3
}

/*
Что запомнить (что чаще и почему):
  • Порядок обхода map случаен — НИКОГДА не выводите map напрямую, если важен порядок.
    Стандартный приём: keys := slices.Collect(maps.Keys(m)); slices.Sort(keys); затем по ключам.
  • maps.Clone vs ручное копирование циклом: Clone короче и быстрее; используют почти всегда.
    ВАЖНО: Clone делает SHALLOW-копию — вложенные срезы/мапы остаются общими.
  • maps.Copy(dst, src) — это слияние с перезаписью, а не замена. Частый паттерн конфигов.
  • maps.Keys/Values возвращают iter.Seq (итераторы 1.23), а не срезы — поэтому Collect.
    Если нужен просто срез ключей: slices.Collect(maps.Keys(m)).

Типичные сценарии:
  1) Конфиг:        cfg := maps.Clone(defaults); maps.Copy(cfg, userOverrides)
  2) Стабильный вывод: for _, k := range sortedKeys { fmt.Println(k, m[k]) }
  3) Сравнение:     if !maps.Equal(want, got) { t.Errorf(...) }
*/
