# Go stdlib — Уровень 1. База

> **Цель уровня:** свободно писать на Go и работать с SQL. Уметь решать алгоритмические
> задачи (easy/medium) и самостоятельно сделать рабочий REST API: маршруты, ручки, чтение/запись
> данных, базовый CRUD.

## Как пользоваться этим уровнем

- Каждый **пакет** — это набор готовых инструментов. Под каждым перечислены его функции.
- **Любую функцию можно кликнуть** — откроется живой пример (`main.go`), где она показана
  с пошаговыми комментариями и результатом (`// => ...`).
- Запустить пример: зайди в папку пакета и выполни `go run main.go`.
- **Порядок — сверху вниз, от простого к сложному.** Иди по списку, не перепрыгивай.
- Значок ⚠️ — тема посложнее: требует базовых знаний или отдельной настройки. Если только
  начал — пройди всё до неё, а к ⚠️ вернись позже.

> **Модель уровней:** `1-base → 2-production → 3-senior → 4-advanced`. Каждый следующий уровень
> добавляет только новые пакеты и функции, без повторов.

---

## Часть A. Основы языка (с них начинают)

### [fmt](fmt/main.go) — вывод на экран и сборка текста
- [`Println(...)`](fmt/main.go) — печатает значения через пробел и переводит строку.
- [`Printf(format, ...)`](fmt/main.go) — печатает по шаблону (`%s`, `%d`, `%v`, `%.2f`).
- [`Sprintf(format, ...)`](fmt/main.go) — НЕ печатает, а возвращает готовую строку.
- [`Fprintf(w, format, ...)`](fmt/main.go) — пишет по шаблону в приёмник (`io.Writer`).
- [`Errorf(format, ...)`](fmt/main.go) — создаёт ошибку; `%w` оборачивает другую ошибку.
- [`Fscan(r, ...)`](fmt/main.go) — читает значения из источника (`io.Reader`).
- [`Scan(...)`](fmt/main.go) — читает значения с клавиатуры (stdin).

### [strings](strings/main.go) — действия со строками (текстом)
- [`Contains(s, substr)`](strings/main.go) — есть ли подстрока внутри.
- [`HasPrefix(s, p)`](strings/main.go) / [`HasSuffix(s, p)`](strings/main.go) — начинается/заканчивается ли.
- [`TrimSpace(s)`](strings/main.go) — убирает пробелы по краям.
- [`TrimPrefix(s, p)`](strings/main.go) / [`TrimSuffix(s, p)`](strings/main.go) — убирает начало/конец.
- [`Split(s, sep)`](strings/main.go) — режет строку по разделителю в список.
- [`Fields(s)`](strings/main.go) — режет по пробелам на слова.
- [`Join(parts, sep)`](strings/main.go) — склеивает список строк в одну.
- [`ReplaceAll(s, old, new)`](strings/main.go) — заменяет все вхождения.
- [`ToLower(s)`](strings/main.go) / [`ToUpper(s)`](strings/main.go) — меняет регистр.
- [`EqualFold(a, b)`](strings/main.go) — сравнивает без учёта регистра.
- [`Count(s, substr)`](strings/main.go) — сколько раз встретилось.
- [`Index(s, substr)`](strings/main.go) — позиция первого вхождения.
- [`NewReader(s)`](strings/main.go) — превращает строку в источник (`io.Reader`).
- [`Builder`](strings/main.go) — эффективная сборка длинной строки по кусочкам.

### [strconv](strconv/main.go) — строки ↔ числа и логические значения
- [`Atoi(s)`](strconv/main.go) — строка → int.
- [`Itoa(i)`](strconv/main.go) — int → строка.
- [`ParseInt(s, base, bit)`](strconv/main.go) — строка → int64 (можно указать систему счисления).
- [`ParseFloat(s, bit)`](strconv/main.go) — строка → дробное число.
- [`ParseBool(s)`](strconv/main.go) — строка → true/false.
- [`FormatInt(i, base)`](strconv/main.go) — int64 → строка в нужной системе счисления.
- [`FormatFloat(f, fmt, prec, bit)`](strconv/main.go) — дробное → строка с нужной точностью.

### [errors](errors/main.go) — работа с ошибками
- [`New(text)`](errors/main.go) — создаёт ошибку с текстом.
- [`Is(err, target)`](errors/main.go) — это именно та ошибка? (проверка сквозь обёртки).
- [`As(err, &target)`](errors/main.go) — достать ошибку нужного типа (когда у неё есть доп. данные).
- [`Unwrap(err)`](errors/main.go) — снять один слой обёртки.

### [log](log/main.go) — служебные сообщения (логи)
- [`Println(...)`](log/main.go) / [`Printf(...)`](log/main.go) / [`Print(...)`](log/main.go) — лог с датой/временем.
- [`Fatal(...)`](log/main.go) / [`Fatalf(...)`](log/main.go) — залогировать и завершить программу.
- [`SetFlags(flags)`](log/main.go) — настроить пометки (дата, время, файл:строка).
- [`SetOutput(w)`](log/main.go) — перенаправить лог в другой приёмник.

---

## Часть B. Данные и алгоритмы

### [math](math/main.go) — математика для дробных чисел (float64)
- [`Abs(x)`](math/main.go) — модуль (убирает минус).
- [`Ceil(x)`](math/main.go) / [`Floor(x)`](math/main.go) / [`Round(x)`](math/main.go) / [`Trunc(x)`](math/main.go) — округления.
- [`Sqrt(x)`](math/main.go) — квадратный корень.
- [`Pow(x, y)`](math/main.go) — степень.
- [`Log(x)`](math/main.go) / [`Log2(x)`](math/main.go) / [`Log10(x)`](math/main.go) — логарифмы.
- [`MaxInt`](math/main.go) / [`MinInt`](math/main.go) — границы int (удобно как «бесконечность» в алгоритмах).
- [`Pi`](math/main.go) / [`E`](math/main.go) — константы.
- _Для целых чисел бери встроенные `min()`/`max()` — пакет math для них не нужен._

### [slices](slices/main.go) — действия над срезами (списками)
- [`Contains(s, v)`](slices/main.go) / [`Index(s, v)`](slices/main.go) — есть ли элемент и где.
- [`Sort(s)`](slices/main.go) — сортировка чисел/строк по возрастанию.
- [`SortFunc(s, cmp)`](slices/main.go) — сортировка по своему правилу.
- [`Max(s)`](slices/main.go) / [`Min(s)`](slices/main.go) — наибольший/наименьший элемент.
- [`Clone(s)`](slices/main.go) — независимая копия (срезы делят данные!).
- [`Equal(a, b)`](slices/main.go) — поэлементное сравнение двух срезов.
- [`Contains(s, v)`](slices/main.go) — проверка наличия без ручного цикла.
- [`BinarySearch(s, t)`](slices/main.go) — быстрый поиск в отсортированном срезе.
- [`Insert(s, i, v...)`](slices/main.go) / [`Delete(s, i, j)`](slices/main.go) — вставка/удаление по индексу.
- [`IsSorted(s)`](slices/main.go) — проверка, отсортирован ли срез.

### [maps](maps/main.go) — действия над словарями (map)
- [`Clone(m)`](maps/main.go) — независимая копия словаря.
- [`Copy(dst, src)`](maps/main.go) — перенести/наложить пары (с заменой совпадающих ключей).
- [`Equal(a, b)`](maps/main.go) — сравнить два словаря.
- [`Keys(m)`](maps/main.go) / [`Values(m)`](maps/main.go) — перебор ключей/значений.
- [`DeleteFunc(m, f)`](maps/main.go) — удалить пары по условию.

### [sort](sort/main.go) — сортировка (классический способ)
- [`Ints(x)`](sort/main.go) / [`Strings(x)`](sort/main.go) / [`Float64s(x)`](sort/main.go) — сортировка базовых списков.
- [`Slice(x, less)`](sort/main.go) — сортировка по своему правилу (через функцию «i раньше j?»).
- [`SliceStable(x, less)`](sort/main.go) — то же, но сохраняет порядок равных.
- [`SearchInts(a, x)`](sort/main.go) — бинарный поиск в отсортированном `[]int`.
- [`Search(n, f)`](sort/main.go) — бинарный поиск по условию.
- _Для простых чисел/строк сейчас короче `slices.Sort` — но `sort.Slice` встречается часто._

### [cmp](cmp/main.go) — сравнение значений
- [`Compare(a, b)`](cmp/main.go) — возвращает -1 / 0 / 1.
- [`Or(vals...)`](cmp/main.go) — первое непустое значение (и для сортировки по нескольким полям).
- [`Less(a, b)`](cmp/main.go) — проверяет `a < b`.

### [bufio](bufio/main.go) — чтение текста по строкам/словам
- [`NewScanner(r)`](bufio/main.go) — сканер для построчного чтения.
- [`Scanner.Scan()`](bufio/main.go) — перейти к следующей строке/слову.
- [`Scanner.Text()`](bufio/main.go) — текущая строка/слово.
- [`Scanner.Split(fn)`](bufio/main.go) — режим разбиения (`bufio.ScanWords` — по словам).
- [`NewReader(r)`](bufio/main.go) / [`Reader.ReadString(d)`](bufio/main.go) — чтение до разделителя.
- [`NewWriter(w)`](bufio/main.go) / [`Writer.WriteString(s)`](bufio/main.go) / [`Writer.Flush()`](bufio/main.go) — буферизованная запись.

---

## Часть C. Файлы, время, веб и база данных

### [time](time/main.go) — время и даты
- [`Now()`](time/main.go) — текущий момент.
- [`Since(t)`](time/main.go) / [`Until(t)`](time/main.go) — сколько прошло / осталось.
- [`Sleep(d)`](time/main.go) — пауза.
- [`Parse(layout, s)`](time/main.go) — текст → время.
- [`Time.Format(layout)`](time/main.go) — время → текст (шаблон-образец `2006-01-02 15:04:05`).
- [`Time.Add(d)`](time/main.go) / [`Time.Sub(t)`](time/main.go) — прибавить / разница.
- [`Time.Before(t)`](time/main.go) / [`Time.After(t)`](time/main.go) — сравнение моментов.
- [`Duration.Seconds()`](time/main.go) / [`Milliseconds()`](time/main.go) — длительность в числах.

### [os](os/main.go) — файлы, окружение, аргументы
- [`ReadFile(name)`](os/main.go) / [`WriteFile(name, data, perm)`](os/main.go) — прочитать/записать файл целиком.
- [`Open(name)`](os/main.go) / [`Create(name)`](os/main.go) — открыть/создать файл для потоковой работы.
- [`Getenv(key)`](os/main.go) — переменная окружения.
- [`LookupEnv(key)`](os/main.go) — переменная окружения + «задана ли».
- [`Remove(name)`](os/main.go) — удалить файл.
- [`Args`](os/main.go) — аргументы запуска программы.
- [`Stdin`](os/main.go) / [`Stdout`](os/main.go) / [`Stderr`](os/main.go) — стандартные потоки.
- [`Exit(code)`](os/main.go) — завершить программу с кодом.

### [io](io/main.go) — потоки данных (Reader и Writer)
- [`Reader`](io/main.go) / [`Writer`](io/main.go) — главные интерфейсы: «откуда читать» / «куда писать».
- [`ReadAll(r)`](io/main.go) — прочитать весь источник в `[]byte`.
- [`Copy(dst, src)`](io/main.go) — перелить данные из источника в приёмник.
- [`EOF`](io/main.go) — сигнал «данные закончились» (не ошибка!).
- [`Closer`](io/main.go) / [`ReadCloser`](io/main.go) — интерфейсы с методом `Close()`.

### [path/filepath](path-filepath/main.go) — пути к файлам и папкам
- [`Join(elem...)`](path-filepath/main.go) — собрать путь с правильным разделителем.
- [`Base(p)`](path-filepath/main.go) / [`Dir(p)`](path-filepath/main.go) / [`Ext(p)`](path-filepath/main.go) — имя / папка / расширение.
- [`Abs(p)`](path-filepath/main.go) — абсолютный путь.
- [`WalkDir(root, fn)`](path-filepath/main.go) — рекурсивный обход папки.

### [net/url](net-url/main.go) — веб-адреса и параметры запроса
- [`Parse(rawURL)`](net-url/main.go) — разобрать адрес на части.
- [`URL.Query()`](net-url/main.go) — параметры запроса как `Values`.
- [`Values.Get(key)`](net-url/main.go) — значение параметра по имени.
- [`Values.Set(k, v)`](net-url/main.go) / [`Values.Add(k, v)`](net-url/main.go) — задать/добавить параметр.
- [`Values.Encode()`](net-url/main.go) — собрать строку параметров (с экранированием).
- [`QueryEscape(s)`](net-url/main.go) / [`QueryUnescape(s)`](net-url/main.go) — экранировать/раскодировать значение.

### [encoding/json](encoding-json/main.go) — формат JSON
- [`Marshal(v)`](encoding-json/main.go) — данные Go → JSON.
- [`MarshalIndent(v, ...)`](encoding-json/main.go) — «красивый» JSON с отступами.
- [`Unmarshal(data, &v)`](encoding-json/main.go) — JSON → данные Go.
- [`NewEncoder(w)`](encoding-json/main.go) / [`Encoder.Encode(v)`](encoding-json/main.go) — писать JSON в приёмник (поток).
- [`NewDecoder(r)`](encoding-json/main.go) / [`Decoder.Decode(&v)`](encoding-json/main.go) — читать JSON из источника (поток).
- [`Decoder.DisallowUnknownFields()`](encoding-json/main.go) — строгая проверка лишних полей.

### ⚠️ [context](context/main.go) — отмена и ограничение по времени
- [`Background()`](context/main.go) — начальный (пустой) контекст.
- [`TODO()`](context/main.go) — заглушка, «контекст сюда ещё не прокинут».
- [`WithTimeout(parent, d)`](context/main.go) — контекст, который сам отменится через время.
- [`WithCancel(parent)`](context/main.go) — контекст с ручной отменой.
- [`ctx.Done()`](context/main.go) — сигнал «пора останавливаться».
- [`ctx.Err()`](context/main.go) — причина отмены.
- [`ctx.Value(key)`](context/main.go) — достать значение, привязанное к запросу.

### ⚠️ [net/http](net-http/main.go) — веб-сервер и веб-клиент
- [`ListenAndServe(addr, h)`](net-http/main.go) — запустить HTTP-сервер.
- [`NewServeMux()`](net-http/main.go) / [`ServeMux`](net-http/main.go) — роутер (с Go 1.22 умеет `GET /path/{id}`).
- [`HandleFunc(pattern, fn)`](net-http/main.go) — привязать функцию-обработчик к адресу.
- [`Handler`](net-http/main.go) / [`HandlerFunc`](net-http/main.go) — интерфейс обработчика и переходник из функции.
- [`Server`](net-http/main.go) — настраиваемый сервер (адрес, таймауты).
- [`Request.PathValue(name)`](net-http/main.go) — часть пути (`{id}` из маршрута).
- [`Request.FormValue(key)`](net-http/main.go) — значение из query/формы.
- [`ResponseWriter.Header()`](net-http/main.go) / [`WriteHeader(code)`](net-http/main.go) / [`Write(data)`](net-http/main.go) — формирование ответа.
- [`Get(url)`](net-http/main.go) — простой GET-запрос (клиент).
- [`Client`](net-http/main.go) / [`Client.Do(req)`](net-http/main.go) — настраиваемый клиент и выполнение запроса.
- [`Redirect`](net-http/main.go) / [`Error`](net-http/main.go) / [`NotFound`](net-http/main.go) — типовые ответы.
- Константы статусов: [`StatusOK`](net-http/main.go) (200), `StatusCreated` (201), `StatusNoContent` (204), `StatusBadRequest` (400), `StatusUnauthorized` (401), `StatusForbidden` (403), `StatusNotFound` (404), `StatusConflict` (409), `StatusInternalServerError` (500).

### ⚠️ [database/sql](database-sql/main.go) — работа с базой данных (SQL)
- [`Open(driver, dsn)`](database-sql/main.go) — подготовить подключение к базе.
- [`DB.PingContext(ctx)`](database-sql/main.go) — проверить, что база доступна.
- [`DB.QueryRowContext(...)`](database-sql/main.go) — получить одну строку.
- [`DB.QueryContext(...)`](database-sql/main.go) — получить много строк.
- [`DB.ExecContext(...)`](database-sql/main.go) — изменить данные (INSERT/UPDATE/DELETE).
- [`DB.BeginTx(ctx, ...)`](database-sql/main.go) — начать транзакцию.
- [`DB.Close()`](database-sql/main.go) — закрыть подключение.
- [`Row.Scan(...)`](database-sql/main.go) — прочитать одну строку в переменные.
- [`Rows.Next()`](database-sql/main.go) / [`Rows.Scan(...)`](database-sql/main.go) / [`Rows.Err()`](database-sql/main.go) / [`Rows.Close()`](database-sql/main.go) — перебор многих строк.
- [`Tx.Commit()`](database-sql/main.go) / [`Tx.Rollback()`](database-sql/main.go) — подтвердить/откатить транзакцию.
- [`ErrNoRows`](database-sql/main.go) — строка не найдена.
- [`NullString`](database-sql/main.go) / `NullInt64` / `NullBool` / `NullTime` — типы для колонок, которые могут быть NULL.
