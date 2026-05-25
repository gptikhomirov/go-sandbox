/*
Пакет strconv — конвертация строк в числа/bool и обратно.

Зачем:   данные извне приходят строками (HTTP, env, файлы) — их нужно превратить в типы.
Когда:   query-параметры (?page=2), переменные окружения, разбор CSV/конфигов.
Грабли:  Atoi падает на "2.0", " 2 " (пробелы) и "" — всегда проверяйте err и при ошибке
         отдавайте 400, а не игнорируйте. Результат ParseInt — всегда int64, приведите к нужному типу.
*/
package main

import (
	"fmt"
	"strconv"
)

func main() {
	// Atoi/Itoa — самый частый случай: строка <-> int. Например ?page=2 из URL.
	page, err := strconv.Atoi("2")
	if err != nil {
		fmt.Println("bad page:", err) // на практике → 400 Bad Request
	}
	fmt.Println(page + 1)        // => 3
	fmt.Println(strconv.Itoa(100)) // => 100

	// ParseInt с базой и разрядностью — для больших чисел и разных систем счисления.
	// base=0 определяет основание по префиксу: 0x → hex, 0b → bin, 0o → oct.
	n, _ := strconv.ParseInt("0xFF", 0, 64)
	fmt.Println(n) // => 255

	// ParseFloat — цены, координаты, метрики.
	price, _ := strconv.ParseFloat("19.99", 64)
	fmt.Printf("%.2f\n", price*2) // => 39.98

	// ParseBool — флаги из env/конфига: принимает 1/t/T/TRUE/true/0/f/false и т.п.
	debug, _ := strconv.ParseBool("true")
	fmt.Println(debug) // => true

	// Format* — обратное преобразование с контролем формата.
	fmt.Println(strconv.FormatInt(10, 2))          // => 1010  (двоичное)
	fmt.Println(strconv.FormatFloat(3.14159, 'f', 2, 64)) // => 3.14
}

/*
Что запомнить (что чаще и почему):
  • Atoi vs ParseInt: Atoi(s) — это ParseInt(s, 10, 0), короткая запись для обычного int в дес. системе.
    Берите ParseInt, только если нужна другая база (hex/bin) или явная разрядность (int64).
  • strconv vs fmt: strconv.Itoa(n) быстрее и понятнее, чем fmt.Sprintf("%d", n).
    fmt.Sprintf берут, когда в одной строке несколько значений сразу.
  • Всегда проверяйте err: «тихая» конвертация плохого ввода — источник багов и уязвимостей.
  • Парсинг из URL/формы: r.URL.Query().Get возвращает строку → её и скармливайте в Atoi/ParseInt.

Типичные сценарии:
  1) Пагинация:  page, err := strconv.Atoi(r.URL.Query().Get("page"))
  2) Конфиг:     maxConns, _ := strconv.Atoi(os.Getenv("MAX_CONNS"))
  3) Цена CSV:   price, err := strconv.ParseFloat(record[2], 64)
*/
