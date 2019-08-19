# Go Disk Buffer

This package help to work with huge amount of data, which cannot be stored in RAM. Instead of keeping all data in RAM `go-disk-buffer.Buffer` stores the data on a disk in a temporary file.

**Notes:**

- `go-disk-buffer.Buffer` is compatible with `io.Reader` and `io.Writer` interfaces
- `go-disk-buffer.Buffer` can replace `bytes.Buffer` (except some methods – check [Unavailable methods](#unavailable-methods))
- `buffer.Buffer` is **not** thread-safe!

##

- [Example](#example)
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

## Available methods

### Read

- `Read(p []byte) (n int, err error)`
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
