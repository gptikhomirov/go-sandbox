# Go stdlib — Уровень 4. Прочее (знать, что существует)

> **Цель уровня:** узкоспециализированные пакеты. В повседневной работе и на собесах почти не
> встречаются, но полезно знать, что они есть, чтобы не изобретать велосипед, когда понадобятся:
> низкоуровневая память и syscalls, сжатие и архивы, парсинг Go-кода, инспекция бинарников, почта.
>
> **Только новое поверх [1-base](1-base/1-base.md), [2-production](2-production.md), [3-senior](3-senior.md).**

---

## Часть A. Низкий уровень

### unsafe — небезопасные операции с памятью
- `Pointer` — небезопасный указатель (мост между типами).
- `Sizeof(x)` / `Alignof(x)` / `Offsetof(field)` — размер/выравнивание/смещение поля.

### syscall — прямые системные вызовы (legacy; см. `golang.org/x/sys`)
- `Syscall(...)` / `Syscall6(...)` / `RawSyscall(...)` — raw-вызовы.
- `Getpid()` / `Getppid()` — PID процессов.
- `Kill(pid, sig)` — отправка сигнала.
- `Open` / `Close` / `Read` / `Write` — низкоуровневые файловые операции.

---

## Часть B. Сжатие и архивы

### compress/gzip — gzip
- `NewReader(r)` / `NewWriter(w)` / `NewWriterLevel(w, level)` — reader/writer.
- `Writer.Write(p)` / `Close()` / `Flush()` — запись потока.

### compress/zlib — zlib
- `NewReader(r)` / `NewWriter(w)` / `NewWriterLevel(w, level)` — reader/writer.

### compress/flate — DEFLATE
- `NewReader(r)` / `NewWriter(w, level)` — низкоуровневый DEFLATE.

### archive/zip — ZIP-архивы
- `OpenReader(name)` / `NewReader(r, size)` / `NewWriter(w)` — чтение/запись.
- `File.Open()` — файл внутри архива.
- `Writer.Create(name)` / `CreateHeader(header)` / `Close()` — добавление файлов.

### archive/tar — TAR-архивы
- `NewReader(r)` / `NewWriter(w)` — reader/writer.
- `Reader.Next()` — следующий entry.
- `Writer.WriteHeader(h)` / `Write(p)` / `Close()` — запись.
- `Header` — метаданные entry.

---

## Часть C. Дополнительные форматы данных

### encoding/xml — XML
- `Marshal(v)` / `MarshalIndent(v, prefix, indent)` — Go value → XML.
- `Unmarshal(data, &v)` — XML → Go value.
- `NewEncoder(w)` / `NewDecoder(r)` — streaming.
- `Decoder.Token()` — потоковый разбор по токенам.

### path — slash-пути (URL/embed, не для файловой системы)
- `Join(elem...)` / `Clean(path)` — манипуляции.
- `Base(path)` / `Dir(path)` / `Ext(path)` / `Split(path)` — части пути.
- `Match(pattern, name)` — сопоставление.

### text/scanner — лексический сканер
- `Scanner.Init(src)` — инициализация.
- `Scanner.Scan()` / `TokenText()` — следующий токен и его текст.
- `Scanner.Peek()` / `Next()` — посмотреть/прочитать rune.

### html — HTML-escaping (без шаблонов)
- `EscapeString(s)` / `UnescapeString(s)` — экранирование и его снятие.

---

## Часть D. Инструментарий Go (анализ кода)

### go/ast — синтаксическое дерево Go
- `Inspect(node, fn)` / `Walk(visitor, node)` — обход AST.
- `File` / `FuncDecl` / `GenDecl` / `TypeSpec` / `Ident` — узлы дерева.

### go/parser — парсинг Go-кода
- `ParseFile(fset, filename, src, mode)` — парсит файл.
- `ParseDir(fset, path, filter, mode)` — парсит директорию.
- `ParseExpr(x)` — парсит выражение.

### go/token — позиции и токены
- `NewFileSet()` — создаёт file set.
- `FileSet.Position(pos)` — позиция в исходнике.
- `IsKeyword(name)` — проверка ключевого слова.

### go/format — форматирование исходников
- `Source(src)` — форматирует Go-код (как `gofmt`).
- `Node(dst, fset, node)` — форматирует узел AST.

### go/types — проверка типов
- `Config.Check(path, fset, files, info)` — type-check пакета.
- `Identical(x, y)` / `AssignableTo(V, T)` / `ConvertibleTo(V, T)` — отношения типов.

### go/doc — извлечение документации
- `New(pkg, importPath, mode)` / `NewFromFiles(...)` — построение doc-пакета.

### go/build — контекст сборки
- `Default` — контекст по умолчанию.
- `Context.Import(path, srcDir, mode)` — метаданные пакета.

---

## Часть E. Инспекция бинарников

### debug/elf — ELF (Linux)
- `Open(name)` / `NewFile(r)` — открыть файл.
- `File.Section(name)` / `Symbols()` / `Close()` — секции и символы.

### debug/macho — Mach-O (macOS)
- `Open(name)` / `NewFile(r)` — открыть файл.
- `File.Section(name)` / `Symtab` / `Close()` — секции и таблица символов.

### debug/pe — PE (Windows)
- `Open(name)` / `NewFile(r)` — открыть файл.
- `File.Section(name)` / `ImportedSymbols()` / `Close()` — секции и импорты.

### debug/dwarf — отладочная информация DWARF
- `Data.Reader()` — reader по DWARF.
- `Reader.Next()` / `Seek(off)` / `SkipChildren()` — навигация.
- `Entry.Val(attr)` — значение атрибута.

### plugin — динамическая загрузка плагинов
- `Open(path)` — открыть `.so`-плагин.
- `Plugin.Lookup(symName)` — найти символ.

---

## Часть F. Почта и текстовые протоколы

### net/mail — разбор email-адресов и сообщений
- `ParseAddress(address)` / `ParseAddressList(list)` — парсинг адресов.
- `ReadMessage(r)` — чтение сообщения.

### net/smtp — SMTP-клиент
- `SendMail(addr, auth, from, to, msg)` — отправка письма.
- `PlainAuth(...)` / `CRAMMD5Auth(...)` — аутентификация.

### net/textproto — текстовые протоколы (основа HTTP/SMTP)
- `NewReader(r)` / `NewWriter(w)` — reader/writer.
- `CanonicalMIMEHeaderKey(s)` — каноничный вид ключа заголовка.
- `Reader.ReadLine()` / `ReadMIMEHeader()` — чтение строк/заголовков.

---

## Часть G. Тестовые утилиты

### testing/fstest — тестовая файловая система
- `MapFS` — in-memory FS для тестов.
- `TestFS(fsys, expected...)` — проверка корректности реализации FS.

### testing/iotest — тестирование reader/writer
- `ErrReader(err)` / `OneByteReader(r)` / `HalfReader(r)` / `TimeoutReader(r)` — «вредные» readers для проверки устойчивости.
- `TruncateWriter(w, n)` — writer с ограничением.
