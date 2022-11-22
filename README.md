# Go Disk Buffer

Package `buffer` helps to work with huge amount of data, which cannot be stored in RAM. Instead of keeping all data in RAM `buffer.Buffer` can store the data on a disk in a temporary file.

**Features:**

- `buffer.Buffer` is compatible with `io.Reader` and `io.Writer` interfaces
- `buffer.Buffer` can replace `bytes.Buffer` (except some methods – check [Unavailable methods](#unavailable-methods))
- You can encrypt data on a disk. Just use `Buffer.EnableEncryption` method

**Notes:**

- It is **not** recommended to use zero value of `buffer.Buffer`. Use `buffer.NewBuffer()` or `buffer.NewBufferWithMaxMemorySize()` instead
- `buffer.Buffer` is **not** thread-safe!
- `buffer.Buffer` uses a directory returned by `os.TempDir()` to store temp files. You can change the directory with `Buffer.ChangeTempDir` method

##

- [Example](#example)
- [Benchmark](#benchmark)
- [Available methods](#available-methods)
  - [Read](#read)
  - [Write](#write)
  - [Other](#other)
- [Unavailable methods](#unavailable-methods)
  - [Can be added](#can-be-added)
  - [Won't be adeed](#wont-be-adeed)

## Example

With `bytes.Buffer`

```go
package main

import (
    "bytes"
    "fmt"
)

func main() {
    b := bytes.Buffer{}
    // b := bytes.NewBuffer(nil)
    // b := bytes.NewBufferString("")

    b.Write([]byte("Hello,"))
    b.WriteByte(' ')
    b.Write([]byte("World!"))

    data := b.Next(13)
    fmt.Println(string(data)) // "Hello, World!"
}
```

With `github.com/ShoshinNikita/go-disk-buffer.Buffer`

```go
package main

import (
    "fmt"

    buffer "github.com/ShoshinNikita/go-disk-buffer"
)

func main() {
    b := buffer.NewBufferWithMaxMemorySize(7) // store only 7 bytes in RAM
    // b := buffer.NewBuffer(nil)
    // b := buffer.NewBufferString("")

    b.Write([]byte("Hello,"))
    b.WriteByte(' ')
    b.Write([]byte("World!")) // will be stored on a disk in a temporary file

    data := b.Next(13)
    fmt.Println(string(data)) // "Hello, World!"
}
```

## Benchmark

**CPU:** Intel Core i7-3630QM  
**RAM:** 8 GB  
**Disk:** HDD, 5400 rpm

```
Buffer_size_is_greater_than_data/bytes.Buffer-8     1000       1591091 ns/op      10043209 B/op     36 allocs/op
Buffer_size_is_greater_than_data/utils.Buffer-8     1000       1346077 ns/op       6901679 B/op     26 allocs/op

Buffer_size_is_equal_to_data/bytes.Buffer-8         1000       1760100 ns/op      10043195 B/op     36 allocs/op
Buffer_size_is_equal_to_data/utils.Buffer-8         2000       1357077 ns/op       7434159 B/op     27 allocs/op

Buffer_size_is_less_than_data/bytes.Buffer-8          50      36522090 ns/op     177848123 B/op     53 allocs/op
Buffer_size_is_less_than_data/utils.Buffer-8          10     110406320 ns/op     112327659 B/op     62 allocs/op
```

## Available methods

### Read

- `Read(p []byte) (n int, err error)`
- `ReadAt(b []byte, off int64) (n int, err error)`
- `ReadByte() (byte, error)`
- `Next(n int) []byte`
- `WriteTo(w io.Writer) (n int64, err error)`

### Write

- `Write(p []byte) (n int, err error)`
- `WriteByte(c byte) error`
- `WriteRune(r rune) (n int, err error)`
- `WriteString(s string) (n int, err error)`
- `ReadFrom(r io.Reader) (n int64, err error)`

### Other

- `Len() int`
- `Cap() int` – equal to `Len()` method
- `Reset()`

## Unavailable methods

### Can be added

- `ReadBytes(delim byte) (line []byte, err error)`
- `ReadString(delim byte) (line string, err error)`
- `ReadRune() (r rune, size int, err error)` – **help wanted** (check `Buffer.readRune()` method)

### Won't be adeed

- `Bytes() []byte`

  **Reason:** `go-disk-buffer` was created to store a huge amount of data. If your data can fit in RAM, you should use `bytes.Buffer`

- `String() string`

  **Reason:** see the previous reason

- `Grow(n int)`

  **Reason:** we can allocate the memory only in RAM. It doesn't make sense to allocate space on a disk

- `Truncate(n int)`
- `UnreadByte() error`
- `UnreadRune() error`
