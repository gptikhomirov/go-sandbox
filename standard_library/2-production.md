# Go stdlib — Уровень 2. Production-ready

> **Цель уровня:** перейти от «работает» к «можно катить в прод». Тесты, конкурентность,
> структурное логирование, безопасность, конфиги, шаблоны, работа с бинарными данными.
> Это объём stdlib, которого достаточно для собеседования уровня **middle**.
>
> **Только новое поверх [1-base](1-base.md).** Если пакет уже был в Базе — здесь перечислены
> лишь функции, которых там не было.

---

## Часть A. Углубление базовых пакетов

### fmt — расширенный вывод и парсинг
- `Print(...)` / `Sprint(...)` / `Sprintln(...)` — вывод/сборка без формата.
- `Fprint(w, ...)` / `Fprintln(w, ...)` — запись в writer без формата.
- `Scanf(format, ...)` — форматированное чтение из stdin.
- `Sscan(s, ...)` / `Sscanf(s, format, ...)` — парсинг значений прямо из строки.

### errors — продвинутая работа с ошибками
- `Join(errs...)` — объединяет несколько ошибок в одну (multi-error).

### log/slog — структурное логирование (стандарт для прода)
- `Info(msg, args...)` / `Warn(...)` / `Error(...)` / `Debug(...)` — лог с уровнем.
- `Log(ctx, level, msg, args...)` — лог с явным context/level.
- `New(handler)` — создаёт logger.
- `SetDefault(logger)` — задаёт logger по умолчанию.
- `Default()` — default logger.
- `With(args...)` — добавляет постоянные поля к logger.
- `WithGroup(name)` — группирует поля.
- `NewJSONHandler(w, opts)` — JSON handler (для прод-логов).
- `NewTextHandler(w, opts)` — text handler.
- `String/Int/Int64/Float64/Bool/Duration/Time/Any(key, value)` — типизированные атрибуты.
- `Group(key, args...)` — group attribute.
- `Logger.Info/Warn/Error/Debug(...)` — лог через конкретный logger.
- `Logger.With(...)` / `Logger.WithGroup(...)` — производный logger.
- `Logger.Enabled(ctx, level)` — проверяет, включён ли уровень.

### strings — дополнительные операции
- `ContainsAny(s, chars)` / `ContainsRune(s, r)` — поиск символов из набора.
- `Trim(s, cutset)` / `TrimLeft` / `TrimRight` / `TrimFunc(s, f)` — гибкая обрезка.
- `SplitN(s, sep, n)` — разбивает максимум на `n` частей.
- `Replace(s, old, new, n)` — заменяет `n` вхождений.
- `LastIndex(s, substr)` — индекс последнего вхождения.
- `Repeat(s, count)` — повторяет строку.
- `Cut(s, sep)` — разрезает на «до» и «после» по первому `sep` (идиоматично).
- `CutPrefix(s, prefix)` / `CutSuffix(s, suffix)` — срезает и сообщает успех.
- `Map(mapping, s)` — преобразует каждую rune.
- `NewReplacer(pairs...)` — многократная замена; `Replacer.Replace(s)`.
- `Builder.WriteByte(b)` / `WriteRune(r)` / `Grow(n)` — точная сборка.

### strconv — дополнительные конвертации
- `ParseUint(s, base, bitSize)` — string → uint64.
- `FormatUint(i, base)` / `FormatBool(b)` — uint64/bool → string.
- `Quote(s)` / `Unquote(s)` — экранирование и его снятие.
- `QuoteRune(r)` — экранированная rune-литерал.
- `AppendInt(dst, i, base)` / `AppendFloat(...)` — добавление в `[]byte` без аллокаций.

### math — расширенная математика
- `Pow10(n)` / `Mod(x, y)` — степень 10 / остаток float64.
- `Sin`/`Cos`/`Tan`/`Atan2(y, x)` — тригонометрия.
- `IsNaN(x)` / `IsInf(x, sign)` — проверки спецзначений.
- `Inf(sign)` / `NaN()` — создание спецзначений.
- `MaxFloat64` / `SmallestNonzeroFloat64` — границы float64.
- `MaxInt64` / `MinInt64` / `MaxInt32` — границы целых типов.

### time — расширенная работа со временем
- `Date(...)` — создаёт `time.Time` вручную.
- `Unix(sec, nsec)` — время из Unix timestamp.
- `ParseDuration(s)` — парсит `"300ms"`, `"1h30m"` и т.п.
- `LoadLocation(name)` — загружает timezone.
- `After(duration)` / `NewTimer(duration)` / `NewTicker(duration)` — таймеры/тикеры.
- `time.Time.Equal(other)` — корректное сравнение моментов.
- `time.Time.UTC()` / `Local()` / `In(loc)` — смена timezone.
- `time.Time.Truncate(d)` / `Round(d)` — округление времени.
- `time.Time.Unix()` / `UnixMilli()` — в timestamp.
- `time.Duration.String()` / `Hours()` / `Minutes()` — представление duration.

### sort — добавочное
- `Sort(data)` / `Stable(data)` — сортировка через `sort.Interface`.
- `Reverse(data)` — разворачивает порядок сортировки.
- `Interface` — интерфейс `Len/Less/Swap` для своих типов.
- `SearchFloat64s(a, x)` — бинарный поиск в `[]float64`.

### slices — generic-операции (полный набор)
- `ContainsFunc(s, f)` / `IndexFunc(s, f)` — поиск по условию.
- `EqualFunc(a, b, eq)` — сравнение через функцию.
- `Compare(a, b)` / `CompareFunc(a, b, cmp)` — лексикографическое сравнение.
- `Compact(s)` / `CompactFunc(s, eq)` — удаляет соседние дубликаты.
- `DeleteFunc(s, del)` — удаляет по условию.
- `Replace(s, i, j, values...)` — заменяет диапазон.
- `Grow(s, n)` / `Clip(s)` — управление capacity.
- `Reverse(s)` — разворачивает slice.
- `Concat(slices...)` — объединяет несколько slice.
- `SortStableFunc(s, cmp)` / `IsSortedFunc(s, cmp)` / `BinarySearchFunc(...)` — варианты с comparator.
- `MaxFunc(s, cmp)` / `MinFunc(s, cmp)` — экстремумы по comparator.
- `Sorted(seq)` / `SortedFunc(seq, cmp)` / `Collect(seq)` / `Values(s)` / `All(s)` — мост к итераторам (Go 1.23).

### maps — generic-операции
- `EqualFunc(a, b, eq)` — сравнение value через функцию.
- `All(m)` — iterator пар key-value.
- `Insert(m, seq)` — вставляет пары из iterator.
- `Collect(seq)` — собирает map из iterator.

### context — продвинутая отмена
- `WithDeadline(parent, time)` — отмена к конкретному моменту.
- `WithValue(parent, key, value)` — request-scoped значение.
- `WithCancelCause(parent)` / `Cause(ctx)` — отмена с причиной.
- `WithoutCancel(parent)` — не наследует отмену родителя.
- `AfterFunc(ctx, fn)` — вызывает `fn` после отмены context.
- `ctx.Deadline()` — возвращает deadline, если есть.

### os — расширенная работа с ФС
- `Stat(name)` / `Lstat(name)` — инфо о файле (без/с переходом по symlink).
- `IsNotExist(err)` / `IsExist(err)` — классификация ошибок.
- `Mkdir` / `MkdirAll` / `Remove` / `RemoveAll` / `Rename` — операции с ФС.
- `ReadDir(name)` — список содержимого директории.
- `CreateTemp(dir, pattern)` / `MkdirTemp(dir, pattern)` — временные файл/директория.
- `Setenv` / `Unsetenv` / `Environ()` — управление окружением.
- `Getwd()` / `Chdir(dir)` — рабочая директория.
- `OpenFile(name, flag, perm)` — открытие с флагами (`O_APPEND`, `O_CREATE`, ...).
- `File.Seek` / `File.Sync` / `File.WriteString` — операции над открытым файлом.

### io — расширенные абстракции
- `CopyN(dst, src, n)` — копирует ровно `n` bytes.
- `ReadFull(r, buf)` / `ReadAtLeast(r, buf, min)` — гарантированное чтение.
- `LimitReader(r, n)` — ограничивает чтение.
- `MultiReader(...)` / `MultiWriter(...)` — объединение readers/writers.
- `TeeReader(r, w)` — читает и параллельно пишет.
- `Pipe()` — связанная пара reader/writer.
- `NopCloser(r)` — reader → ReadCloser с пустым Close.
- `Discard` — writer, который всё отбрасывает.
- `Seeker` / `ReaderAt` / `WriterAt` / `ReadWriter` — дополнительные интерфейсы.

### bufio — тонкая настройка
- `NewReaderSize` / `NewWriterSize` — буферы заданного размера.
- `Reader.ReadBytes(delim)` / `ReadByte` / `ReadRune` / `Peek(n)` — точное чтение.
- `Writer.Write(p)` / `WriteByte` / `WriteRune` — запись.
- `Scanner.Bytes()` / `Scanner.Err()` / `Scanner.Buffer(buf, max)` — увеличение лимита токена.
- `ScanLines` / `ScanWords` / `ScanRunes` — готовые split-функции.

### path/filepath — дополнительное
- `Clean(path)` / `Split(path)` / `Rel(base, targ)` / `IsAbs(path)` — манипуляции с путём.
- `Glob(pattern)` / `Match(pattern, name)` — поиск/сопоставление.
- `ToSlash(path)` / `FromSlash(path)` — кроссплатформенные разделители.
- `Walk(root, fn)` — обход дерева (legacy-вариант `WalkDir`).

### encoding/json — дополнительное
- `Valid(data)` — проверяет валидность JSON.
- `Compact(dst, src)` / `Indent(dst, src, prefix, indent)` — переформатирование.
- `Encoder.SetIndent(...)` / `Encoder.SetEscapeHTML(on)` — настройка вывода.
- `Decoder.UseNumber()` — числа как `json.Number`.
- `Decoder.More()` / `Decoder.Token()` — потоковый разбор.
- `Number` — число без немедленной конвертации.

### net/http — production-настройки
- `ListenAndServeTLS(addr, cert, key, handler)` — HTTPS-сервер.
- `Server.Shutdown(ctx)` — graceful shutdown.
- `NewRequestWithContext(ctx, method, url, body)` — запрос с context.
- `StripPrefix(prefix, h)` — для статики/middleware.
- `FileServer(root)` — отдача статики.
- `SetCookie(w, cookie)` / `Request.Cookie(name)` / `Request.Cookies()` — cookies.
- `MaxBytesReader(w, body, n)` — лимит размера тела.
- `Request.BasicAuth()` — basic-аутентификация.
- `Request.ParseMultipartForm(maxMemory)` / `Request.FormFile(key)` — загрузка файлов.
- `Transport` — настройка пула соединений клиента.
- `Client.CloseIdleConnections()` — закрыть idle-соединения.
- `MethodGet` / `MethodPost` / `MethodPut` / `MethodDelete` / `MethodPatch` — константы методов.

### database/sql — пул и подготовленные запросы
- `DB.PrepareContext(ctx, query)` — prepared statement.
- `DB.SetMaxOpenConns(n)` / `SetMaxIdleConns(n)` — настройка пула.
- `DB.SetConnMaxLifetime(d)` / `SetConnMaxIdleTime(d)` — lifetime соединений.
- `DB.Stats()` — статистика пула (для метрик).
- `Named(name, value)` — named-параметр.
- `Rows.Columns()` / `Rows.ColumnTypes()` — метаданные результата.

---

## Часть B. Новые пакеты прода

### testing — тесты, бенчмарки, фаззинг
- `t.Run(name, fn)` — subtest (табличные тесты).
- `t.Helper()` — пометка helper-функции.
- `t.Fatal(...)` / `t.Fatalf(...)` — провал с остановкой.
- `t.Error(...)` / `t.Errorf(...)` — провал без остановки.
- `t.Log(...)` / `t.Logf(...)` — лог теста.
- `t.Skip(...)` / `t.Skipf(...)` — пропуск.
- `t.Cleanup(fn)` — cleanup после теста.
- `t.TempDir()` — временная директория.
- `t.Setenv(key, value)` — env на время теста.
- `t.Parallel()` — параллельный запуск.
- `b.N` / `b.ResetTimer()` / `b.StopTimer()` / `b.StartTimer()` / `b.ReportAllocs()` — бенчмарки.
- `f.Add(...)` / `f.Fuzz(fn)` — fuzz-тесты.
- `T` / `B` / `F` / `M` — типы тестов и `TestMain`.

### net/http/httptest — тестирование HTTP
- `NewRequest(method, target, body)` — тестовый запрос.
- `NewRecorder()` — fake `ResponseWriter`.
- `NewServer(handler)` / `NewTLSServer(handler)` — тестовый сервер.
- `ResponseRecorder.Code` / `.Body` / `.Header()` / `.Result()` — проверка ответа.

### sync — синхронизация goroutines
- `Mutex.Lock()` / `Unlock()` / `TryLock()` — взаимное исключение.
- `RWMutex.Lock/Unlock/RLock/RUnlock/TryLock/TryRLock` — RW-блокировка.
- `WaitGroup.Add(delta)` / `Done()` / `Wait()` — ожидание группы goroutines.
- `Once.Do(fn)` — однократное выполнение.
- `Cond.Wait()` / `Signal()` / `Broadcast()` — условная переменная.
- `Pool.Get()` / `Pool.Put(x)` — пул переиспользуемых объектов.
- `Map.Load/Store/Delete/Range` — concurrent map.
- `Map.LoadOrStore(key, value)` / `LoadAndDelete(key)` — атомарные комбо.

### sync/atomic — атомарные значения
- `atomic.Int32` / `Int64` / `Uint32` / `Uint64` / `Bool` / `Pointer[T]` — атомарные типы.
- `Load()` / `Store(v)` / `Add(delta)` / `Swap(new)` — операции.
- `CompareAndSwap(old, new)` — CAS.

### bytes — операции с `[]byte` (аналог strings)
- `NewBuffer(buf)` / `NewBufferString(s)` / `NewReader(b)` — конструкторы.
- `Contains` / `HasPrefix` / `HasSuffix` / `Equal` / `Compare` — проверки.
- `Split` / `Join` / `TrimSpace` / `Trim` / `ReplaceAll` / `ToLower` / `ToUpper` / `Repeat` — манипуляции.
- `Cut(s, sep)` — разрез по разделителю.
- `Buffer.Write(p)` / `WriteString(s)` / `Read(p)` — потоковая запись/чтение.
- `Buffer.Bytes()` / `String()` / `Reset()` / `Grow(n)` — управление буфером.

### regexp — регулярные выражения
- `MustCompile(expr)` / `Compile(expr)` — компиляция (Must — или panic).
- `MatchString(pattern, s)` / `Match(pattern, b)` — быстрая проверка.
- `QuoteMeta(s)` — экранирование regexp-символов.
- `Regexp.MatchString(s)` / `Match(b)` — проверка совпадения.
- `Regexp.FindString(s)` / `FindAllString(s, n)` — поиск совпадений.
- `Regexp.FindStringSubmatch(s)` — совпадение + capture groups.
- `Regexp.ReplaceAllString(src, repl)` / `ReplaceAllFunc(...)` — замена.
- `Regexp.Split(s, n)` — разбиение по regexp.
- `Regexp.SubexpNames()` — имена групп.

### unicode — категории символов
- `IsLetter(r)` / `IsDigit(r)` / `IsNumber(r)` / `IsSpace(r)` — классификация rune.
- `IsUpper(r)` / `IsLower(r)` — регистр.
- `ToUpper(r)` / `ToLower(r)` — смена регистра rune.

### unicode/utf8 — UTF-8
- `RuneCountInString(s)` / `RuneCount(p)` — число rune (≠ длине в байтах!).
- `ValidString(s)` / `Valid(p)` — проверка валидности UTF-8.
- `DecodeRuneInString(s)` / `DecodeLastRuneInString(s)` — декодирование.
- `EncodeRune(p, r)` / `RuneLen(r)` — кодирование/длина.

### flag — CLI-флаги и конфиги
- `String` / `Int` / `Bool` / `Duration` / `Float64(name, value, usage)` — объявление флагов.
- `Var(value, name, usage)` — custom-флаг.
- `Parse()` — парсит флаги.
- `Args()` / `Arg(i)` / `NArg()` — позиционные аргументы.
- `NewFlagSet(name, errorHandling)` — отдельный набор флагов (для subcommands).

### os/exec — запуск внешних команд
- `Command(name, args...)` / `CommandContext(ctx, name, args...)` — создание команды.
- `LookPath(file)` — поиск executable в PATH.
- `Cmd.Run()` / `Start()` / `Wait()` — запуск.
- `Cmd.Output()` / `CombinedOutput()` — захват вывода.
- `Cmd.Stdin/Stdout/Stderr/Env/Dir` — настройка.

### crypto/rand — криптостойкая случайность
- `Read(b)` — заполняет bytes криптослучайными данными.
- `Reader` — глобальный secure random.
- `Int(rand, max)` — случайный `big.Int` < max.

### crypto/sha256 — SHA-256
- `Sum256(data)` — хеш одним вызовом.
- `New()` — streaming SHA-256.

### crypto/hmac — HMAC (подпись токенов, webhooks)
- `New(hash, key)` — создаёт HMAC.
- `Equal(a, b)` — безопасное сравнение MAC.

### crypto/subtle — constant-time сравнения
- `ConstantTimeCompare(x, y)` — сравнение без timing-leak.

### crypto/tls — TLS
- `Config` — конфигурация TLS.
- `LoadX509KeyPair(certFile, keyFile)` — загрузка cert/key.
- `Dial(network, addr, config)` — TLS-соединение.

### encoding/base64 — Base64
- `StdEncoding.EncodeToString(src)` / `DecodeString(s)` — кодирование/декодирование.
- `URLEncoding` / `RawURLEncoding` — URL-safe варианты (для токенов).

### encoding/hex — hex
- `EncodeToString(src)` / `DecodeString(s)` — bytes ↔ hex string.

### encoding/csv — CSV
- `NewReader(r)` / `NewWriter(w)` — reader/writer.
- `Reader.Read()` / `ReadAll()` — чтение записей.
- `Reader.Comma` / `Reader.FieldsPerRecord` — настройка.
- `Writer.Write(record)` / `WriteAll(records)` / `Flush()` — запись.

### mime — MIME-типы
- `TypeByExtension(ext)` — MIME по расширению.
- `ParseMediaType(v)` — разбор `Content-Type` с параметрами.

### mime/multipart — multipart-формы и загрузка файлов
- `NewReader(r, boundary)` / `NewWriter(w)` — reader/writer.
- `Reader.NextPart()` / `Reader.ReadForm(maxMemory)` — чтение частей.
- `Writer.CreateFormFile(field, filename)` / `CreateFormField(field)` — части формы.
- `Writer.WriteField(field, value)` / `FormDataContentType()` / `Close()` — запись формы.
- `Part.FormName()` / `FileName()` — метаданные части.

### html/template — HTML-шаблоны с авто-escaping
- `New(name)` / `ParseFiles(...)` / `ParseGlob(pattern)` — создание/парсинг.
- `Must(t, err)` — обёртка-or-panic.
- `Template.Execute(w, data)` / `ExecuteTemplate(w, name, data)` — рендеринг.
- `Template.Funcs(funcMap)` — пользовательские функции.

### text/template — текстовые шаблоны (конфиги, кодогенерация)
- `New(name)` / `Parse(text)` / `ParseFiles(...)` — создание/парсинг.
- `Template.Execute(w, data)` / `ExecuteTemplate(w, name, data)` — рендеринг.
- `Template.Funcs(funcMap)` — пользовательские функции.

### embed — встраивание файлов в бинарник
- `//go:embed file.txt` — встроить один файл.
- `//go:embed templates/*` — встроить набор.
- `FS` — embedded filesystem.
- `FS.ReadFile(name)` / `FS.Open(name)` / `FS.ReadDir(name)` — доступ к встроенному.

### math/rand/v2 — псевдослучайные числа (не для крипто)
- `IntN(n)` / `Int64N(n)` / `Int32N(n)` — случайное в `[0, n)`.
- `Float64()` / `Float32()` — случайное float.
- `Perm(n)` — случайная перестановка `[0, n)`.
- `Shuffle(n, swap)` — перемешивание.
- `New(source)` — отдельный генератор.

### net — TCP/UDP/DNS (когда нужно ниже HTTP)
- `Dial(network, address)` / `DialTimeout(...)` — подключение.
- `Listen(network, address)` — слушать адрес.
- `SplitHostPort(hostport)` / `JoinHostPort(host, port)` — разбор адресов.
- `LookupHost(host)` / `LookupIP(host)` — DNS-резолвинг.
- `ParseIP(s)` — парсинг IP.
- `Conn` / `Listener` — ключевые интерфейсы.
