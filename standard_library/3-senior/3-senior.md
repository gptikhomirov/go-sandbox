# Go stdlib — Уровень 3. Senior / Middle+

> **Цель уровня:** выход на senior и подготовка к собеседованиям Middle+. Структуры данных
> и алгоритмы, внутренности рантайма, профилирование, отражение, бинарные форматы, хеши.
> В повседневной REST-разработке встречается редко, но именно это спрашивают «вглубь».
>
> **Только новое поверх [1-base](1-base/1-base.md) и [2-production](2-production.md).**

---

## Часть A. Структуры данных и алгоритмы (классика собесов)

### container/heap — приоритетная очередь (binary heap)
- `Init(h)` — инициализирует heap из произвольного среза.
- `Push(h, x)` — добавляет элемент, сохраняя инвариант.
- `Pop(h)` — извлекает min/root элемент.
- `Remove(h, i)` — удаляет элемент по индексу.
- `Fix(h, i)` — восстанавливает порядок после изменения элемента.
- `Interface` — реализуется поверх `sort.Interface` + `Push`/`Pop`.

### container/list — двусвязный список
- `New()` — создаёт список.
- `List.PushBack(v)` / `PushFront(v)` — добавление в конец/начало.
- `List.InsertAfter(v, mark)` / `InsertBefore(v, mark)` — вставка относительно элемента.
- `List.Remove(e)` — удаление.
- `List.MoveToBack(e)` / `MoveToFront(e)` — перемещение (основа LRU-кэша).
- `List.Front()` / `Back()` / `Len()` — доступ/размер.
- `Element.Next()` / `Prev()` / `Value` — обход и значение.

### math/bits — битовые операции (битовые трюки на собесах)
- `Len(x)` / `Len64(x)` / `Len32(x)` — позиция старшего бита.
- `OnesCount(x)` / `OnesCount64(x)` — popcount (число единичных битов).
- `LeadingZeros(x)` / `TrailingZeros(x)` — ведущие/хвостовые нули.
- `RotateLeft(x, k)` — циклический сдвиг.
- `Reverse(x)` / `ReverseBytes(x)` — разворот битов/байтов.
- `Add(x, y, carry)` / `Sub(x, y, borrow)` — арифметика с переносом.
- `Mul(x, y)` / `Div(hi, lo, y)` — полноразрядные умножение/деление.

---

## Часть B. Рантайм, профилирование, наблюдаемость

### reflect — runtime-интроспекция
- `TypeOf(v)` / `ValueOf(v)` — runtime type/value.
- `DeepEqual(a, b)` — глубокое сравнение (часто в тестах).
- `Zero(t)` / `New(t)` — zero value / pointer на новое значение.
- `Indirect(v)` — разыменование pointer/interface.
- `Type.Kind()` / `Name()` / `Elem()` / `Key()` — описание типа.
- `Type.NumField()` / `Field(i)` / `FieldByName(name)` — поля struct (теги, маршалинг).
- `Type.NumMethod()` / `Method(i)` — методы.
- `Value.Interface()` / `Kind()` / `Type()` — обратно к значению.
- `Value.CanSet()` / `Set(value)` / `SetString` / `SetInt` / `SetBool` — мутация.
- `Value.Field(i)` / `FieldByName(name)` / `Elem()` — навигация по значению.
- `Value.Call(args)` — динамический вызов функции/метода.
- `Value.IsNil()` / `IsZero()` / `Len()` — проверки.

### runtime — управление и информация о рантайме
- `NumCPU()` — число логических CPU.
- `GOMAXPROCS(n)` — число OS-потоков для Go-кода.
- `NumGoroutine()` — число goroutines.
- `GC()` — принудительный сбор мусора.
- `ReadMemStats(&m)` — статистика памяти; `MemStats`.
- `Caller(skip)` / `Callers(skip, pc)` — информация о стеке.
- `Goexit()` — завершение текущей goroutine.
- `KeepAlive(x)` — удержание объекта живым.

### runtime/debug — отладка и тюнинг
- `Stack()` / `PrintStack()` — stack trace.
- `ReadBuildInfo()` — build info бинарника (версия, зависимости).
- `SetGCPercent(n)` — частота GC.
- `SetMemoryLimit(limit)` — soft memory limit.
- `FreeOSMemory()` — вернуть память ОС.

### runtime/pprof — профилирование
- `StartCPUProfile(w)` / `StopCPUProfile()` — CPU-профиль.
- `WriteHeapProfile(w)` — heap-профиль.
- `Lookup(name)` / `Profiles()` — доступ к профилям.
- `Do(ctx, labels, fn)` / `SetGoroutineLabels(ctx)` — pprof-метки.

### net/http/pprof — pprof через HTTP
- `Index` / `Cmdline` / `Profile` / `Symbol` / `Trace` — HTTP-хендлеры pprof.
- `Handler(name)` — хендлер конкретного профиля.

### expvar — экспорт метрик через HTTP
- `NewInt(name)` / `NewFloat(name)` / `NewString(name)` / `NewMap(name)` — публичные переменные.
- `Publish(name, v)` — публикация custom-переменной.
- `Handler()` — HTTP-хендлер `/debug/vars`.

---

## Часть C. Бинарные форматы и хеши

### encoding/binary — бинарное кодирование чисел
- `Read(r, order, data)` / `Write(w, order, data)` — чтение/запись бинарных данных.
- `Size(v)` — размер бинарного представления.
- `BigEndian` / `LittleEndian` / `NativeEndian` — порядок байтов.
- `PutUvarint(buf, x)` / `Uvarint(buf)` — varint-кодирование.

### encoding/gob — бинарная сериализация Go-значений
- `NewEncoder(w)` / `NewDecoder(r)` — encoder/decoder.
- `Encoder.Encode(v)` / `Decoder.Decode(&v)` — (де)сериализация.
- `Register(value)` — регистрация конкретного типа за интерфейсом.

### hash — интерфейсы хешей
- `Hash` / `Hash32` / `Hash64` — интерфейсы хеш-функций.
- `Hash.Write(p)` / `Sum(b)` / `Reset()` / `Size()` — методы.

### hash/crc32 — CRC-32 (контрольные суммы)
- `ChecksumIEEE(data)` — CRC-32 IEEE одним вызовом.
- `New(tab)` / `NewIEEE()` — streaming.
- `MakeTable(poly)` — таблица полинома.

### hash/fnv — быстрый non-crypto хеш (шардирование, мапы)
- `New32a()` / `New64a()` — FNV-1a 32/64-bit.
- `New32()` / `New64()` / `New128()` / `New128a()` — остальные варианты.

---

## Часть D. Файловые абстракции и сетевые утилиты

### io/fs — абстракция файловой системы (тестируемость)
- `FS` / `File` / `DirEntry` / `FileInfo` — интерфейсы.
- `ReadFile(fsys, name)` / `ReadDir(fsys, name)` / `Stat(fsys, name)` — операции поверх FS.
- `Glob(fsys, pattern)` / `WalkDir(fsys, root, fn)` — поиск/обход.
- `Sub(fsys, dir)` — sub-filesystem.
- `ValidPath(name)` — валидация пути FS.

### net/http/httputil — отладка HTTP и reverse proxy
- `DumpRequest(req, body)` / `DumpRequestOut(req, body)` / `DumpResponse(resp, body)` — дамп для логов/отладки.
- `NewSingleHostReverseProxy(target)` — простой reverse proxy.
- `ReverseProxy` — настраиваемый reverse proxy.

### net/http/httptrace — трассировка клиентских запросов
- `WithClientTrace(ctx, trace)` — добавляет trace в context.
- `ClientTrace` — колбэки на этапах DNS/connect/TLS/headers.

### testing/quick — property-based тестирование
- `Check(f, config)` — проверяет свойство на случайных входах.
- `CheckEqual(f, g, config)` — проверяет эквивалентность двух функций.
