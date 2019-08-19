# Disk Buffer

## Available methods

### Read

- `Read(p []byte) (n int, err error)`
- `ReadByte() (byte, error)`
- `Next(n int) []byte`

### Write

- `Write(p []byte) (n int, err error)`
- `WriteByte(c byte) error`
- `WriteRune(r rune) (n int, err error)`
- `WriteString(s string) (n int, err error)`

### Other

- `Len() int`
- `Reset()`
