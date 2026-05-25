# Go stdlib — Уровень 1. База

> **Цель уровня:** беглый Go + SQL. Уметь решать алгоритмические задачи easy/medium и
> самостоятельно написать жизнеспособный REST API: хендлеры, ручки, миграции, базовый CRUD.
>
> **Модель уровней (кумулятивная, без пересечений):** `1-base → 2-production → 3-senior → 4-advanced`.
> Каждый следующий файл добавляет **только новые** пакеты и функции. Здесь — фундамент.
>
> **Формат:** `пакет` → `функция / тип / метод` — что делает.

---

## Часть A. Ядро языка (нужно всегда)

### fmt — вывод, форматирование, ошибки
- `Println(...)` — выводит значения с переносом строки.
- `Printf(format, ...)` — выводит форматированную строку.
- `Sprintf(format, ...)` — собирает форматированную строку и возвращает `string`.
- `Fprintf(w, format, ...)` — пишет форматированный текст в `io.Writer`.
- `Errorf(format, ...)` — создаёт ошибку; с `%w` заворачивает другую ошибку.
- `Scan(...)` — читает значения из stdin (для алгоритмических задач).
- `Fscan(r, ...)` — читает значения из `io.Reader`; удобно с `bufio.Reader`.

### errors — базовая работа с ошибками
- `New(text)` — создаёт простую ошибку.
- `Is(err, target)` — проверяет, есть ли конкретная ошибка внутри wrapped-chain.
- `As(err, &target)` — достаёт ошибку конкретного типа из chain.
- `Unwrap(err)` — достаёт внутреннюю ошибку.

### log — простое логирование
- `Println(...)` — пишет лог-строку с переносом.
- `Printf(format, ...)` — форматированная лог-строка.
- `Print(...)` — лог без форматирования.
- `Fatal(...)` — логирует и завершает процесс через `os.Exit(1)`.
- `Fatalf(format, ...)` — форматированный `Fatal`.
- `SetFlags(flags)` — задаёт формат метаданных (дата, время, файл).
- `SetOutput(w)` — перенаправляет вывод в `io.Writer`.

### strings — работа со строками
- `Contains(s, substr)` — проверяет наличие подстроки.
- `HasPrefix(s, prefix)` / `HasSuffix(s, suffix)` — проверяет prefix/suffix.
- `TrimSpace(s)` — убирает whitespace по краям.
- `TrimPrefix(s, prefix)` / `TrimSuffix(s, suffix)` — убирает prefix/suffix.
- `Split(s, sep)` — разбивает строку по разделителю.
- `Join(parts, sep)` — собирает `[]string` в строку.
- `Fields(s)` — разбивает строку по whitespace.
- `ReplaceAll(s, old, new)` — заменяет все вхождения.
- `ToLower(s)` / `ToUpper(s)` — меняет регистр.
- `EqualFold(a, b)` — сравнивает строки без учёта регистра.
- `Count(s, substr)` — считает количество вхождений.
- `Index(s, substr)` — индекс первого вхождения.
- `NewReader(s)` — превращает string в `io.Reader`.
- `Builder` — эффективная сборка строк.
- `Builder.WriteString(s)` — добавляет строку в builder.
- `Builder.String()` — возвращает итоговую строку.
- `Builder.Reset()` — очищает builder.

### strconv — конвертация строк, чисел, bool
- `Atoi(s)` / `Itoa(i)` — string ↔ int.
- `ParseInt(s, base, bitSize)` — string → int64.
- `ParseFloat(s, bitSize)` — string → float64.
- `ParseBool(s)` — string → bool.
- `FormatInt(i, base)` — int64 → string.
- `FormatFloat(f, fmt, prec, bitSize)` — float → string.

---

## Часть B. Алгоритмические задачи (easy/medium)

### math — базовая математика
- `Abs(x)` — модуль float64.
- `Max(x, y)` / `Min(x, y)` — максимум/минимум float64.
- `Ceil(x)` / `Floor(x)` / `Round(x)` / `Trunc(x)` — округления.
- `Sqrt(x)` — квадратный корень.
- `Pow(x, y)` — степень.
- `Log(x)` / `Log2(x)` / `Log10(x)` — логарифмы.
- `MaxInt` / `MinInt` — границы `int` (частая граница в алгоритмах).
- `Pi` / `E` — константы.

### sort — сортировка и бинарный поиск
- `Ints(x)` / `Strings(x)` / `Float64s(x)` — сортирует базовые slice.
- `Slice(x, less)` — сортирует slice по custom-функции.
- `SliceStable(x, less)` — стабильная сортировка.
- `Search(n, f)` — бинарный поиск по условию.
- `SearchInts(a, x)` — бинарный поиск в отсортированном `[]int`.

### slices — современные операции над slice
- `Contains(s, v)` / `Index(s, v)` — поиск элемента.
- `Equal(a, b)` — сравнивает два slice.
- `Clone(s)` — shallow copy.
- `Delete(s, i, j)` — удаляет диапазон.
- `Insert(s, i, values...)` — вставляет элементы.
- `Sort(s)` / `SortFunc(s, cmp)` — сортировка.
- `IsSorted(s)` — проверка сортировки.
- `BinarySearch(s, target)` — бинарный поиск.
- `Max(s)` / `Min(s)` — экстремумы.

### maps — операции над map
- `Clone(m)` — shallow copy.
- `Copy(dst, src)` — копирует пары.
- `Equal(a, b)` — сравнивает две map.
- `DeleteFunc(m, del)` — удаляет пары по условию.
- `Keys(m)` / `Values(m)` — iterator ключей/значений.

### cmp — сравнение значений (для `SortFunc` и дженериков)
- `Compare(a, b)` — возвращает `-1`, `0` или `1`.
- `Less(a, b)` — проверяет `a < b`.
- `Or(vals...)` — первое non-zero значение.
- `Ordered` — constraint для ordered-типов.

### bufio — быстрый ввод/вывод (важно для алгоритмов)
- `NewReader(r)` / `NewWriter(w)` — buffered reader/writer.
- `NewScanner(r)` — scanner для построчного чтения.
- `Reader.ReadString(delim)` — читает строку до delimiter.
- `Writer.WriteString(s)` — пишет string в buffer.
- `Writer.Flush()` — сбрасывает buffer (не забывать!).
- `Scanner.Scan()` — переходит к следующему токену.
- `Scanner.Text()` — текущий токен как string.
- `Scanner.Split(splitFunc)` — задаёт способ разбиения (например, `bufio.ScanWords`).

---

## Часть C. Базовый REST API + SQL

### time — время и duration
- `Now()` — текущее время.
- `Since(t)` / `Until(t)` — прошло/осталось времени.
- `Sleep(duration)` — пауза goroutine.
- `Parse(layout, value)` — парсит строку во время.
- `time.Time.Format(layout)` — форматирует время в строку.
- `time.Time.Add(duration)` / `time.Time.Sub(other)` — арифметика времени.
- `time.Time.Before(other)` / `time.Time.After(other)` — сравнение.
- `time.Duration.Seconds()` / `Milliseconds()` — duration в числах.

### context — отмена и timeout запросов
- `Background()` — корневой context для `main`.
- `TODO()` — временный context.
- `WithCancel(parent)` — context с ручной отменой.
- `WithTimeout(parent, duration)` — отмена по timeout.
- `ctx.Done()` — channel, закрывается при отмене.
- `ctx.Err()` — причина отмены.
- `ctx.Value(key)` — значение из context.

### os — файлы, env, аргументы
- `ReadFile(name)` / `WriteFile(name, data, perm)` — файл целиком.
- `Open(name)` / `Create(name)` — открыть/создать файл.
- `Getenv(key)` — env-переменная (часто для конфигов сервиса).
- `LookupEnv(key)` — env-переменная + bool «задана ли».
- `Args` — аргументы запуска.
- `Stdin` / `Stdout` / `Stderr` — стандартные потоки.
- `Exit(code)` — завершить процесс с кодом.

### io — абстракции чтения/записи
- `ReadAll(r)` — читает весь reader в `[]byte` (например, тело запроса).
- `Copy(dst, src)` — копирует stream.
- `EOF` — ошибка конца stream.
- `Reader` / `Writer` — ключевые интерфейсы стандартной библиотеки.
- `Closer` / `ReadCloser` — интерфейсы с `Close()`.

### path/filepath — пути файловой системы
- `Join(elem...)` — собирает путь с учётом ОС.
- `Base(path)` / `Dir(path)` / `Ext(path)` — части пути.
- `Abs(path)` — абсолютный путь.
- `WalkDir(root, fn)` — рекурсивный обход (например, для миграций).

### net/url — URL и query-параметры
- `Parse(rawURL)` — парсит URL.
- `QueryEscape(s)` / `QueryUnescape(s)` — экранирование query value.
- `URL.Query()` — query params как `url.Values`.
- `URL.String()` — собирает URL в строку.
- `Values.Get(key)` — первое значение query-параметра.
- `Values.Set(key, value)` / `Values.Add(key, value)` — задать/добавить.
- `Values.Encode()` — собирает query string.

### encoding/json — JSON для тел запросов/ответов
- `Marshal(v)` — Go value → JSON `[]byte`.
- `MarshalIndent(v, prefix, indent)` — красивый JSON.
- `Unmarshal(data, &v)` — JSON → Go value.
- `NewEncoder(w)` / `NewDecoder(r)` — streaming encoder/decoder (идеально для HTTP).
- `Encoder.Encode(v)` — пишет JSON в `ResponseWriter`.
- `Decoder.Decode(&v)` — читает JSON из тела запроса.
- `Decoder.DisallowUnknownFields()` — строгая валидация полей.
- `RawMessage` — отложенный decode.
- `Marshaler` / `Unmarshaler` — интерфейсы для custom (де)сериализации.

### net/http — HTTP server и client
- `ListenAndServe(addr, handler)` — запускает HTTP-сервер.
- `HandleFunc(pattern, fn)` — регистрирует функцию-handler.
- `Handle(pattern, handler)` — регистрирует `http.Handler`.
- `Redirect(w, r, url, code)` — отправляет redirect.
- `Error(w, msg, code)` — отправляет error response.
- `NotFound(w, r)` — отправляет 404.
- `ServeMux` — стандартный router (с Go 1.22 умеет `GET /path/{id}`).
- `Handler` — интерфейс `ServeHTTP(w, r)`.
- `HandlerFunc` — adapter из функции в Handler.
- `Server` — настраиваемый сервер (timeouts, addr).
- `Client` — HTTP-клиент.
- `NewRequest(method, url, body)` — создаёт запрос.
- `Get(url)` / `Post(url, contentType, body)` — простые запросы.
- `Client.Do(req)` — выполняет запрос.
- `Request.Context()` — context запроса.
- `Request.PathValue(name)` — path-параметр (`{id}` из маршрута).
- `Request.FormValue(key)` — значение из query/form.
- `Request.ParseForm()` — парсит query/form.
- `ResponseWriter.Header()` — заголовки ответа.
- `ResponseWriter.WriteHeader(code)` — статус-код.
- `ResponseWriter.Write(data)` — тело ответа.
- Константы статусов: `StatusOK` (200), `StatusCreated` (201), `StatusNoContent` (204), `StatusBadRequest` (400), `StatusUnauthorized` (401), `StatusForbidden` (403), `StatusNotFound` (404), `StatusConflict` (409), `StatusInternalServerError` (500).

### database/sql — работа с БД (ядро REST API)
- `Open(driverName, dataSourceName)` — создаёт DB handle (pool).
- `DB.PingContext(ctx)` — проверяет соединение.
- `DB.QueryContext(ctx, query, args...)` — SELECT с множеством строк.
- `DB.QueryRowContext(ctx, query, args...)` — SELECT одной строки.
- `DB.ExecContext(ctx, query, args...)` — INSERT/UPDATE/DELETE.
- `DB.BeginTx(ctx, opts)` — начинает транзакцию.
- `DB.Close()` — закрывает handle.
- `Row.Scan(dest...)` — читает одну строку в переменные.
- `Rows.Next()` — переход к следующей строке.
- `Rows.Scan(dest...)` — читает текущую строку.
- `Rows.Err()` — ошибка после итерации.
- `Rows.Close()` — закрывает rows.
- `Tx.ExecContext(...)` / `Tx.QueryContext(...)` — запросы в транзакции.
- `Tx.Commit()` / `Tx.Rollback()` — завершение транзакции.
- `NullString` / `NullInt64` / `NullBool` / `NullTime` — nullable-типы для колонок.
- `ErrNoRows` — строка не найдена.
