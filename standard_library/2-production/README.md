# Go stdlib — Уровень 2. Production-ready

> **Цель уровня:** перейти от «работает» к «можно катить в прод». Тесты, конкурентность,
> структурное логирование, безопасность, конфиги, шаблоны, бинарные данные.
> Это объём stdlib, достаточный для собеседования уровня **middle**.

## Как пользоваться этим уровнем

- **Только новое поверх [Базы](../1-base/README.md).** Если пакет уже был в Базе — здесь перечислены
  лишь функции, которых там не было (раздел «Углубление»). Новые пакеты — в разделе «Новые пакеты».
- **Любую функцию можно кликнуть** — откроется пример (`main.go`) с комментариями и результатом (`// => ...`).
- Запустить пример: зайди в папку пакета и выполни `go run main.go`.
  Для `testing` — `go test -v` (там лежат `main.go` + `main_test.go`).
- **Порядок — от самого нужного к нишевому.** Значок ⚠️ — нишевое/глубже: бери по необходимости.

> **Модель уровней:** `1-base → 2-production → 3-senior → 4-advanced`. Каждый уровень добавляет только новое.

---

## Часть A. Углубление базовых пакетов

### [strings](strings/main.go) — дополнительные операции
- [`Cut(s, sep)`](strings/main.go) — разрезать на «до» и «после» по первому разделителю (идиоматично).
- [`CutPrefix(s, p)`](strings/main.go) / [`CutSuffix(s, p)`](strings/main.go) — срезать и сообщить, было ли что срезать.
- [`Trim(s, set)`](strings/main.go) / [`TrimLeft`](strings/main.go) / [`TrimRight`](strings/main.go) / [`TrimFunc(s, f)`](strings/main.go) — гибкая обрезка.
- [`SplitN(s, sep, n)`](strings/main.go) — разбить максимум на `n` частей.
- [`Replace(s, old, new, n)`](strings/main.go) — заменить только `n` вхождений.
- [`LastIndex(s, sub)`](strings/main.go) — индекс последнего вхождения.
- [`Repeat(s, n)`](strings/main.go) — повторить строку.
- [`Map(f, s)`](strings/main.go) — преобразовать каждую руну.
- [`NewReplacer(pairs...)`](strings/main.go) — много замен за один проход.
- [`Builder.WriteByte`](strings/main.go) / [`WriteRune`](strings/main.go) / [`Grow(n)`](strings/main.go) — точная сборка.

### [slices](slices/main.go) — generic-операции (полный набор)
- [`ContainsFunc(s, f)`](slices/main.go) / [`IndexFunc(s, f)`](slices/main.go) — поиск по условию.
- [`EqualFunc(a, b, eq)`](slices/main.go) — сравнение через функцию.
- [`Compare(a, b)`](slices/main.go) / [`CompareFunc(...)`](slices/main.go) — лексикографическое сравнение.
- [`Compact(s)`](slices/main.go) / [`CompactFunc(s, eq)`](slices/main.go) — удалить соседние дубликаты.
- [`DeleteFunc(s, f)`](slices/main.go) — удалить по условию.
- [`Replace(s, i, j, v...)`](slices/main.go) — заменить диапазон.
- [`Grow(s, n)`](slices/main.go) / [`Clip(s)`](slices/main.go) — управление capacity.
- [`Reverse(s)`](slices/main.go) — развернуть.
- [`Concat(slices...)`](slices/main.go) — объединить несколько срезов.
- [`SortStableFunc(...)`](slices/main.go) / [`IsSortedFunc(...)`](slices/main.go) / [`BinarySearchFunc(...)`](slices/main.go) — варианты с компаратором.
- [`MaxFunc(s, cmp)`](slices/main.go) / [`MinFunc(s, cmp)`](slices/main.go) — экстремумы по компаратору.
- [`Sorted(seq)`](slices/main.go) / [`SortedFunc(...)`](slices/main.go) / [`Collect(seq)`](slices/main.go) / [`Values(s)`](slices/main.go) / [`All(s)`](slices/main.go) — мост к итераторам (Go 1.23).

### [maps](maps/main.go) — generic-операции
- [`EqualFunc(a, b, eq)`](maps/main.go) — сравнение значений через функцию.
- [`All(m)`](maps/main.go) — итератор пар ключ-значение.
- [`Insert(m, seq)`](maps/main.go) — влить пары из итератора.
- [`Collect(seq)`](maps/main.go) — собрать словарь из итератора.

### [context](context/main.go) — продвинутая отмена
- [`WithDeadline(parent, t)`](context/main.go) — отмена к конкретному моменту.
- [`WithCancelCause(parent)`](context/main.go) / [`Cause(ctx)`](context/main.go) — отмена с причиной.
- [`WithoutCancel(parent)`](context/main.go) — не наследовать отмену родителя.
- [`AfterFunc(ctx, f)`](context/main.go) — вызвать `f` при отмене.
- [`WithValue(parent, k, v)`](context/main.go) / [`ctx.Deadline()`](context/main.go) — значение запроса / дедлайн.

### [net/http](net-http/main.go) — production-настройки
- [`Server.Shutdown(ctx)`](net-http/main.go) — graceful shutdown.
- [`ListenAndServeTLS(...)`](net-http/main.go) — HTTPS-сервер.
- [`NewRequestWithContext(...)`](net-http/main.go) — запрос с context.
- [`StripPrefix(prefix, h)`](net-http/main.go) / [`FileServer(root)`](net-http/main.go) — отдача статики.
- [`SetCookie(w, c)`](net-http/main.go) / [`Request.Cookie(name)`](net-http/main.go) / [`Request.Cookies()`](net-http/main.go) — cookies.
- [`MaxBytesReader(w, body, n)`](net-http/main.go) — лимит размера тела.
- [`Request.BasicAuth()`](net-http/main.go) — basic-аутентификация.
- [`Request.ParseMultipartForm(n)`](net-http/main.go) / [`Request.FormFile(key)`](net-http/main.go) — загрузка файлов.
- [`Transport`](net-http/main.go) / [`Client.CloseIdleConnections()`](net-http/main.go) — настройка клиента.
- [`MethodGet`](net-http/main.go) / `MethodPost` / `MethodPut` / `MethodDelete` / `MethodPatch` — константы методов.

### [database/sql](database-sql/main.go) — пул и подготовленные запросы
- [`DB.SetMaxOpenConns(n)`](database-sql/main.go) / [`SetMaxIdleConns(n)`](database-sql/main.go) — настройка пула.
- [`DB.SetConnMaxLifetime(d)`](database-sql/main.go) / [`SetConnMaxIdleTime(d)`](database-sql/main.go) — lifetime соединений.
- [`DB.PrepareContext(...)`](database-sql/main.go) — prepared statement.
- [`DB.Stats()`](database-sql/main.go) — метрики пула.
- [`Named(name, value)`](database-sql/main.go) — именованный параметр.
- [`Rows.Columns()`](database-sql/main.go) / [`Rows.ColumnTypes()`](database-sql/main.go) — метаданные результата.

### [encoding/json](encoding-json/main.go) — дополнительное
- [`Valid(data)`](encoding-json/main.go) — проверить корректность JSON.
- [`Compact(dst, src)`](encoding-json/main.go) / [`Indent(...)`](encoding-json/main.go) — переформатирование.
- [`Encoder.SetIndent(...)`](encoding-json/main.go) / [`Encoder.SetEscapeHTML(b)`](encoding-json/main.go) — настройка вывода.
- [`Decoder.UseNumber()`](encoding-json/main.go) / [`Number`](encoding-json/main.go) — числа без потери точности.
- [`Decoder.More()`](encoding-json/main.go) / [`Decoder.Token()`](encoding-json/main.go) — потоковый разбор.

### [errors](errors/main.go) — продвинутая работа с ошибками
- [`Join(errs...)`](errors/main.go) — объединить несколько ошибок в одну (multi-error).

### [time](time/main.go) — расширенная работа со временем
- [`Date(...)`](time/main.go) / [`Unix(sec, nsec)`](time/main.go) — создать момент.
- [`ParseDuration(s)`](time/main.go) — разобрать `"1h30m"`.
- [`LoadLocation(name)`](time/main.go) / [`Time.In(loc)`](time/main.go) / [`UTC()`](time/main.go) / [`Local()`](time/main.go) — часовые пояса.
- [`After(d)`](time/main.go) / [`NewTimer(d)`](time/main.go) / [`NewTicker(d)`](time/main.go) — таймеры/тикеры.
- [`Time.Equal(t)`](time/main.go) — корректное сравнение моментов.
- [`Time.Truncate(d)`](time/main.go) / [`Round(d)`](time/main.go) — округление времени.
- [`Time.Unix()`](time/main.go) / [`UnixMilli()`](time/main.go) — в timestamp.
- [`Duration.String()`](time/main.go) / [`Hours()`](time/main.go) / [`Minutes()`](time/main.go) — представление длительности.

### [os](os/main.go) — расширенная работа с ФС
- [`Stat(name)`](os/main.go) / [`Lstat(name)`](os/main.go) — информация о файле.
- [`IsNotExist(err)`](os/main.go) / [`IsExist(err)`](os/main.go) — классификация ошибок.
- [`Mkdir`](os/main.go) / [`MkdirAll`](os/main.go) / [`RemoveAll`](os/main.go) / [`Rename`](os/main.go) — операции с ФС.
- [`ReadDir(name)`](os/main.go) — содержимое директории.
- [`CreateTemp(...)`](os/main.go) / [`MkdirTemp(...)`](os/main.go) — временные файл/папка.
- [`Setenv`](os/main.go) / [`Unsetenv`](os/main.go) / [`Environ()`](os/main.go) — окружение.
- [`Getwd()`](os/main.go) / [`Chdir(dir)`](os/main.go) — рабочая директория.
- [`OpenFile(name, flag, perm)`](os/main.go) — открытие с флагами (`O_APPEND`, `O_CREATE`).
- [`File.Seek`](os/main.go) / [`File.Sync`](os/main.go) / [`File.WriteString`](os/main.go) — операции над файлом.

### [io](io/main.go) — расширенные абстракции
- [`LimitReader(r, n)`](io/main.go) — ограничить чтение `n` байтами (безопасность).
- [`CopyN(dst, src, n)`](io/main.go) — скопировать ровно `n` байт.
- [`ReadFull(r, buf)`](io/main.go) / [`ReadAtLeast(...)`](io/main.go) — гарантированное чтение.
- [`MultiReader(...)`](io/main.go) / [`MultiWriter(...)`](io/main.go) — объединение потоков.
- [`TeeReader(r, w)`](io/main.go) — читать и попутно копировать (хеш на лету).
- [`Pipe()`](io/main.go) / [`NopCloser(r)`](io/main.go) / [`Discard`](io/main.go) — труба / пустой Close / «мусорка».

### [strconv](strconv/main.go) — дополнительные конвертации
- [`ParseUint(...)`](strconv/main.go) / [`FormatUint(...)`](strconv/main.go) — беззнаковые числа.
- [`FormatBool(b)`](strconv/main.go) — bool → строка.
- [`Quote(s)`](strconv/main.go) / [`Unquote(s)`](strconv/main.go) / [`QuoteRune(r)`](strconv/main.go) — экранирование в кавычки.
- [`AppendInt(...)`](strconv/main.go) / [`AppendFloat(...)`](strconv/main.go) — добавление в `[]byte` без аллокаций.

### [bufio](bufio/main.go) — тонкая настройка
- [`Reader.Peek(n)`](bufio/main.go) — заглянуть вперёд без сдвига.
- [`Reader.ReadBytes(d)`](bufio/main.go) / [`ReadByte`](bufio/main.go) / [`ReadRune`](bufio/main.go) — точное чтение.
- [`Scanner.Buffer(buf, max)`](bufio/main.go) — снять лимит длины строки (~64 КБ).
- [`Scanner.Bytes()`](bufio/main.go) / [`Scanner.Err()`](bufio/main.go) — токен как `[]byte` / ошибка.
- [`ScanLines`](bufio/main.go) / [`ScanWords`](bufio/main.go) / [`ScanRunes`](bufio/main.go) — готовые split-функции.
- [`NewReaderSize`](bufio/main.go) / [`NewWriterSize`](bufio/main.go) / [`Writer.WriteByte`](bufio/main.go) / [`WriteRune`](bufio/main.go) — буферы и точная запись.

### [sort](sort/main.go) — добавочное
- [`Sort(data)`](sort/main.go) / [`Stable(data)`](sort/main.go) — сортировка через `sort.Interface`.
- [`Interface`](sort/main.go) — интерфейс `Len/Less/Swap` для своих типов.
- [`Reverse(data)`](sort/main.go) — развернуть порядок сортировки.
- [`SearchFloat64s(a, x)`](sort/main.go) — бинарный поиск в `[]float64`.

### [math](math/main.go) — расширенная математика
- [`Pow10(n)`](math/main.go) / [`Mod(x, y)`](math/main.go) — степень 10 / остаток float64.
- [`Sin`](math/main.go) / [`Cos`](math/main.go) / [`Tan`](math/main.go) / [`Atan2(y, x)`](math/main.go) — тригонометрия.
- [`IsNaN(x)`](math/main.go) / [`IsInf(x, s)`](math/main.go) / [`Inf(s)`](math/main.go) / [`NaN()`](math/main.go) — спецзначения.
- [`MaxFloat64`](math/main.go) / [`MaxInt64`](math/main.go) / [`MinInt64`](math/main.go) / [`MaxInt32`](math/main.go) — границы типов.

### [fmt](fmt/main.go) — расширенный вывод и парсинг
- [`Print(...)`](fmt/main.go) / [`Sprint(...)`](fmt/main.go) / [`Sprintln(...)`](fmt/main.go) — вывод/сборка без формата.
- [`Fprint(w, ...)`](fmt/main.go) / [`Fprintln(w, ...)`](fmt/main.go) — запись в writer без формата.
- [`Sscan(s, ...)`](fmt/main.go) / [`Sscanf(s, fmt, ...)`](fmt/main.go) — разбор значений из строки.
- [`Scanf(fmt, ...)`](fmt/main.go) — форматированное чтение из stdin.

### [path/filepath](path-filepath/main.go) — дополнительное
- [`Clean(path)`](path-filepath/main.go) — нормализовать путь (защита от `..`).
- [`Split(path)`](path-filepath/main.go) / [`Rel(base, targ)`](path-filepath/main.go) / [`IsAbs(path)`](path-filepath/main.go) — манипуляции.
- [`Glob(pattern)`](path-filepath/main.go) / [`Match(pattern, name)`](path-filepath/main.go) — поиск/сопоставление.
- [`ToSlash(p)`](path-filepath/main.go) / [`FromSlash(p)`](path-filepath/main.go) — кроссплатформенные разделители.

### [log/slog](log-slog/main.go) — структурное логирование (стандарт для прода)
- [`Info`](log-slog/main.go) / [`Warn`](log-slog/main.go) / [`Error`](log-slog/main.go) / [`Debug(msg, args...)`](log-slog/main.go) — лог с уровнем.
- [`New(handler)`](log-slog/main.go) / [`SetDefault(l)`](log-slog/main.go) / [`Default()`](log-slog/main.go) — создать/назначить logger.
- [`NewJSONHandler(...)`](log-slog/main.go) / [`NewTextHandler(...)`](log-slog/main.go) — формат вывода.
- [`Logger.With(args...)`](log-slog/main.go) / [`WithGroup(name)`](log-slog/main.go) — постоянные поля / группа.
- [`String`](log-slog/main.go) / [`Int`](log-slog/main.go) / [`Any(k, v)`](log-slog/main.go) / [`Group(...)`](log-slog/main.go) — типизированные атрибуты.
- [`Logger.Enabled(ctx, level)`](log-slog/main.go) — включён ли уровень.

---

## Часть B. Новые пакеты прода

### [testing](testing/main.go) — тесты, бенчмарки, фаззинг
- [`t.Run(name, fn)`](testing/main_test.go) — подтест (табличные тесты).
- [`t.Helper()`](testing/main_test.go) — пометка helper-функции.
- [`t.Fatal`](testing/main_test.go) / [`Fatalf`](testing/main_test.go) — провал с остановкой.
- [`t.Error`](testing/main_test.go) / [`Errorf`](testing/main_test.go) — провал без остановки.
- [`t.Log`](testing/main_test.go) / [`Logf`](testing/main_test.go) / [`Skip`](testing/main_test.go) / [`Skipf`](testing/main_test.go) — лог / пропуск.
- [`t.Cleanup(fn)`](testing/main_test.go) / [`t.TempDir()`](testing/main_test.go) / [`t.Setenv(k, v)`](testing/main_test.go) — авто-очистка.
- [`t.Parallel()`](testing/main_test.go) — параллельный запуск.
- [`b.N`](testing/main_test.go) / [`ResetTimer`](testing/main_test.go) / [`StopTimer`](testing/main_test.go) / [`StartTimer`](testing/main_test.go) / [`ReportAllocs`](testing/main_test.go) — бенчмарки.
- [`f.Add(...)`](testing/main_test.go) / [`f.Fuzz(fn)`](testing/main_test.go) — фаззинг.
- [`T`](testing/main_test.go) / [`B`](testing/main_test.go) / [`F`](testing/main_test.go) / [`M`](testing/main_test.go) — типы тестов и TestMain.

### [sync](sync/main.go) — синхронизация горутин
- [`Mutex.Lock`](sync/main.go) / [`Unlock`](sync/main.go) / [`TryLock`](sync/main.go) — взаимное исключение.
- [`RWMutex`](sync/main.go) — раздельные блокировки чтения/записи.
- [`WaitGroup.Add`](sync/main.go) / [`Done`](sync/main.go) / [`Wait`](sync/main.go) — дождаться группы горутин.
- [`Once.Do(fn)`](sync/main.go) — выполнить ровно один раз.
- [`Cond`](sync/main.go) — условная переменная (Wait/Signal/Broadcast).
- [`Pool.Get`](sync/main.go) / [`Put`](sync/main.go) — пул переиспользуемых объектов.
- [`Map.Load`](sync/main.go) / [`Store`](sync/main.go) / [`Delete`](sync/main.go) / [`Range`](sync/main.go) / [`LoadOrStore`](sync/main.go) — потокобезопасный словарь.

### [sync/atomic](sync-atomic/main.go) — атомарные значения
- [`atomic.Int64`](sync-atomic/main.go) / [`Int32`](sync-atomic/main.go) / [`Bool`](sync-atomic/main.go) / [`Pointer[T]`](sync-atomic/main.go) — атомарные типы.
- [`Load()`](sync-atomic/main.go) / [`Store(v)`](sync-atomic/main.go) / [`Add(d)`](sync-atomic/main.go) / [`Swap(new)`](sync-atomic/main.go) — операции.
- [`CompareAndSwap(old, new)`](sync-atomic/main.go) — CAS (основа lock-free).

### [bytes](bytes/main.go) — операции с `[]byte` (аналог strings)
- [`NewBuffer`](bytes/main.go) / [`NewReader(b)`](bytes/main.go) — конструкторы.
- [`Contains`](bytes/main.go) / [`HasPrefix`](bytes/main.go) / [`Equal`](bytes/main.go) / [`Split`](bytes/main.go) / [`Cut`](bytes/main.go) — как в strings, но для байтов.
- [`Buffer.Write`](bytes/main.go) / [`WriteString`](bytes/main.go) / [`ReadString`](bytes/main.go) / [`Bytes()`](bytes/main.go) / [`String()`](bytes/main.go) / [`Reset()`](bytes/main.go) — растущий буфер.

### [regexp](regexp/main.go) — регулярные выражения
- [`MustCompile(expr)`](regexp/main.go) / [`Compile(expr)`](regexp/main.go) — компиляция.
- [`MatchString(p, s)`](regexp/main.go) — быстрая проверка.
- [`FindString(s)`](regexp/main.go) / [`FindAllString(s, n)`](regexp/main.go) — поиск совпадений.
- [`FindStringSubmatch(s)`](regexp/main.go) — совпадение + группы.
- [`ReplaceAllString(...)`](regexp/main.go) / [`ReplaceAllStringFunc(...)`](regexp/main.go) — замена.
- [`Split(s, n)`](regexp/main.go) / [`QuoteMeta(s)`](regexp/main.go) — разбиение / экранирование.

### [unicode](unicode/main.go) — категории символов
- [`IsLetter`](unicode/main.go) / [`IsDigit`](unicode/main.go) / [`IsSpace`](unicode/main.go) / [`IsNumber(r)`](unicode/main.go) — классификация руны.
- [`IsUpper`](unicode/main.go) / [`IsLower`](unicode/main.go) / [`ToUpper`](unicode/main.go) / [`ToLower(r)`](unicode/main.go) — регистр руны.

### [unicode/utf8](unicode-utf8/main.go) — UTF-8
- [`RuneCountInString(s)`](unicode-utf8/main.go) — число символов (≠ длине в байтах!).
- [`ValidString(s)`](unicode-utf8/main.go) / [`Valid(p)`](unicode-utf8/main.go) — корректность UTF-8.
- [`DecodeRuneInString(s)`](unicode-utf8/main.go) / [`RuneLen(r)`](unicode-utf8/main.go) — декодирование / длина руны.

### [crypto/rand](crypto-rand/main.go) — криптостойкая случайность
- [`Read(b)`](crypto-rand/main.go) — заполнить байты криптослучайными данными (токены).
- [`Reader`](crypto-rand/main.go) — глобальный secure-источник.
- [`Int(rand, max)`](crypto-rand/main.go) — случайное число в `[0, max)`.

### [crypto/sha256](crypto-sha256/main.go) — хеш SHA-256
- [`Sum256(data)`](crypto-sha256/main.go) — хеш одним вызовом.
- [`New()`](crypto-sha256/main.go) — потоковый хеш (io.Writer) для больших данных.

### [crypto/hmac](crypto-hmac/main.go) — подпись секретным ключом
- [`New(hash, key)`](crypto-hmac/main.go) — создать HMAC.
- [`Equal(a, b)`](crypto-hmac/main.go) — безопасное сравнение подписей.

### [crypto/subtle](crypto-subtle/main.go) — сравнение секретов constant-time
- [`ConstantTimeCompare(x, y)`](crypto-subtle/main.go) — сравнение без timing-утечки.

### [encoding/base64](encoding-base64/main.go) — Base64
- [`StdEncoding.EncodeToString`](encoding-base64/main.go) / [`DecodeString`](encoding-base64/main.go) — кодирование/декодирование.
- [`URLEncoding`](encoding-base64/main.go) / [`RawURLEncoding`](encoding-base64/main.go) — URL-safe варианты (токены).

### [encoding/hex](encoding-hex/main.go) — hex
- [`EncodeToString(src)`](encoding-hex/main.go) / [`DecodeString(s)`](encoding-hex/main.go) — байты ↔ hex-строка.

### [encoding/csv](encoding-csv/main.go) — CSV
- [`NewReader(r)`](encoding-csv/main.go) / [`Read()`](encoding-csv/main.go) / [`ReadAll()`](encoding-csv/main.go) — чтение записей.
- [`Reader.Comma`](encoding-csv/main.go) / [`FieldsPerRecord`](encoding-csv/main.go) — настройка разбора.
- [`NewWriter(w)`](encoding-csv/main.go) / [`Write(rec)`](encoding-csv/main.go) / [`WriteAll(recs)`](encoding-csv/main.go) / [`Flush()`](encoding-csv/main.go) — запись.

### [flag](flag/main.go) — CLI-флаги и конфиги
- [`String`](flag/main.go) / [`Int`](flag/main.go) / [`Bool`](flag/main.go) / [`Duration`](flag/main.go) / [`Float64(name, def, usage)`](flag/main.go) — объявление флагов.
- [`Parse()`](flag/main.go) — разобрать флаги.
- [`Args()`](flag/main.go) / [`Arg(i)`](flag/main.go) / [`NArg()`](flag/main.go) — позиционные аргументы.
- [`NewFlagSet(name, mode)`](flag/main.go) / [`Var(...)`](flag/main.go) — отдельный набор / свой тип флага.

### [os/exec](os-exec/main.go) — запуск внешних команд
- [`Command(name, args...)`](os-exec/main.go) / [`CommandContext(...)`](os-exec/main.go) — создать команду.
- [`LookPath(file)`](os-exec/main.go) — найти программу в PATH.
- [`Cmd.Output()`](os-exec/main.go) / [`CombinedOutput()`](os-exec/main.go) / [`Run()`](os-exec/main.go) / [`Start()`](os-exec/main.go) / [`Wait()`](os-exec/main.go) — запуск.
- [`Cmd.Dir`](os-exec/main.go) / [`Env`](os-exec/main.go) / [`Stdin`](os-exec/main.go) / [`Stdout`](os-exec/main.go) / [`Stderr`](os-exec/main.go) — настройка.

### [embed](embed/main.go) — встраивание файлов в бинарник
- [`//go:embed file`](embed/main.go) — директива встраивания (в string/[]byte/FS).
- [`embed.FS`](embed/main.go) / [`FS.ReadFile`](embed/main.go) / [`FS.Open`](embed/main.go) / [`FS.ReadDir`](embed/main.go) — доступ к встроенному.

### [html/template](html-template/main.go) — HTML-шаблоны с авто-escaping
- [`New(name)`](html-template/main.go) / [`Parse(text)`](html-template/main.go) / [`Must(t, err)`](html-template/main.go) — создание.
- [`ParseFiles(...)`](html-template/main.go) / [`ParseGlob(p)`](html-template/main.go) — загрузка из файлов.
- [`Template.Execute(w, data)`](html-template/main.go) / [`ExecuteTemplate(...)`](html-template/main.go) — рендеринг.
- [`Template.Funcs(map)`](html-template/main.go) — пользовательские функции.

### [text/template](text-template/main.go) — текстовые шаблоны (конфиги, кодогенерация)
- [`New(name)`](text-template/main.go) / [`Parse(text)`](text-template/main.go) — создание.
- [`Template.Execute(w, data)`](text-template/main.go) / [`ExecuteTemplate(...)`](text-template/main.go) — рендеринг.
- [`Template.Funcs(map)`](text-template/main.go) — пользовательские функции.

### [mime](mime/main.go) — MIME-типы
- [`TypeByExtension(ext)`](mime/main.go) — MIME по расширению.
- [`ParseMediaType(v)`](mime/main.go) — разбор `Content-Type` на тип + параметры.

### [mime/multipart](mime-multipart/main.go) — формы с файлами
- [`NewWriter(w)`](mime-multipart/main.go) / [`CreateFormFile(...)`](mime-multipart/main.go) / [`WriteField(...)`](mime-multipart/main.go) / [`FormDataContentType()`](mime-multipart/main.go) — сборка формы.
- [`NewReader(r, boundary)`](mime-multipart/main.go) / [`NextPart()`](mime-multipart/main.go) / [`Part.FormName()`](mime-multipart/main.go) / [`FileName()`](mime-multipart/main.go) — разбор формы.

### [math/rand/v2](math-rand-v2/main.go) — псевдослучайные числа (НЕ для крипто)
- [`IntN(n)`](math-rand-v2/main.go) / [`N(n)`](math-rand-v2/main.go) / [`Float64()`](math-rand-v2/main.go) — случайное в диапазоне / дробное.
- [`Perm(n)`](math-rand-v2/main.go) / [`Shuffle(n, swap)`](math-rand-v2/main.go) — перестановка / перемешивание.
- [`New(src)`](math-rand-v2/main.go) / [`NewPCG(...)`](math-rand-v2/main.go) — свой генератор (воспроизводимость).

### ⚠️ [net/http/httptest](net-http-httptest/main.go) — тестирование HTTP
- [`NewRequest(method, target, body)`](net-http-httptest/main.go) — тестовый запрос.
- [`NewRecorder()`](net-http-httptest/main.go) — фейковый ResponseWriter (`.Code`, `.Body`, `.Result()`).
- [`NewServer(handler)`](net-http-httptest/main.go) — реальный локальный сервер (`.URL`, `.Close()`).

### ⚠️ [crypto/tls](crypto-tls/main.go) — TLS (тонкая настройка шифрования)
- [`Config`](crypto-tls/main.go) — конфигурация (MinVersion, ServerName, Certificates).
- [`LoadX509KeyPair(cert, key)`](crypto-tls/main.go) — загрузить сертификат и ключ.
- [`Dial(network, addr, cfg)`](crypto-tls/main.go) — TLS-соединение поверх TCP.

### ⚠️ [net](net/main.go) — TCP/UDP/DNS (ниже HTTP)
- [`Dial(net, addr)`](net/main.go) / [`Listen(net, addr)`](net/main.go) — клиент / сервер на TCP.
- [`SplitHostPort(hp)`](net/main.go) / [`JoinHostPort(h, p)`](net/main.go) / [`ParseIP(s)`](net/main.go) — разбор адресов.
- [`LookupHost(host)`](net/main.go) / [`LookupIP(host)`](net/main.go) — DNS-резолвинг.
- [`Conn`](net/main.go) / [`Listener`](net/main.go) — ключевые интерфейсы (основа net/http).
