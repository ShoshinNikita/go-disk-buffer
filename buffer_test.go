package buffer

import (
	"bytes"
	"io"
	"math/rand"
	"os"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

func TestBuffer_CheckBufferAndFileSize(t *testing.T) {
	tests := []struct {
		maxSize int
		//
		data []byte
		//
		bufferSize int
		fileSize   int
	}{
		{
			maxSize:    15,
			data:       make([]byte, 10),
			bufferSize: 10,
			fileSize:   0,
		},
		{
			maxSize:    15,
			data:       make([]byte, 15),
			bufferSize: 15,
			fileSize:   0,
		},
		{
			maxSize:    15,
			data:       make([]byte, 16),
			bufferSize: 15,
			fileSize:   1,
		},
		{
			maxSize:    15,
			data:       make([]byte, 20),
			bufferSize: 15,
			fileSize:   5,
		},
		{
			maxSize:    20,
			data:       make([]byte, 1<<20),
			bufferSize: 20,
			fileSize:   1<<20 - 20,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run("", func(t *testing.T) {
			t.Parallel()

			require := require.New(t)

			b := NewBufferWithMaxMemorySize(tt.maxSize)
			defer b.Reset()

			n, err := b.Write(tt.data)
			require.Nil(err, "error during Write()")

			// Checks
			require.Equal(len(tt.data), n, "not all data written")

			require.Equal(len(tt.data), b.Len(), "Len() method returned wrong value")

			require.Equal(tt.bufferSize, b.buff.Len(), "buffer contains wrong amount of bytes")

			if len(tt.data) <= tt.maxSize {
				require.Equal("", b.filename, "buffer created excess file")

				// Must skip file checks
				return
			}

			f, err := os.Open(b.filename)
			require.Nilf(err, "can't open file %s", b.filename)
			defer f.Close()

			fileSize := func() int {
				info, err := f.Stat()
				if err != nil {
					return 0
				}

				return int(info.Size())
			}()

			require.Equal(tt.fileSize, fileSize, "buffer contains wrong amount of bytes")
		})

	}
}

func TestBuffer_WriteAndRead(t *testing.T) {
	tests := []struct {
		maxSize       int
		readSliceSize int
		//
		data [][]byte
		//
		res []byte
	}{
		{
			maxSize:       20,
			readSliceSize: 256,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:       1,
			readSliceSize: 256,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:       5,
			readSliceSize: 10,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:       5,
			readSliceSize: 20,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:       5,
			readSliceSize: 10,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:       5,
			readSliceSize: 5,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:       0,
			readSliceSize: 5,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:       0,
			readSliceSize: 5,
			data:          [][]byte{},
			res:           nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run("", func(t *testing.T) {
			t.Parallel()

			require := require.New(t)

			b := NewBufferWithMaxMemorySize(tt.maxSize)
			defer b.Reset()

			var dataWritten int
			for _, d := range tt.data {
				n, err := b.Write(d)
				dataWritten += len(d)

				require.Nil(err, "error during Write()")
				require.Equal(len(d), n, "not all data written")
				require.Equal(dataWritten, b.Len(), "Len() method returned wrong value")
			}

			res, err := readByChunks(require, b, tt.readSliceSize)
			require.Nil(err, "error during Read()")
			require.Equalf(tt.res, res, "wrong content was read")

			require.Equal(0, b.Len(), "Buffer must be empty")
		})
	}
}

func TestBuffer_ReadByte(t *testing.T) {
	require := require.New(t)

	data := []byte("1234")

	b := NewBufferWithMaxMemorySize(len(data) / 2)
	b.Write([]byte(data))

	for i := 0; i < len(data); i++ {
		c, err := b.ReadByte()
		require.Nil(err)
		require.Equal(data[i], c)
	}
}

// TODO
func TestBuffer_ReadRune(t *testing.T) {
	t.Skip("skip the test because \"Buffer.ReadRune()\" method is not finished")

	require := require.New(t)

	data := []byte("Hello | ✓ | 123456 | Привет!")

	b := NewBufferWithMaxMemorySize(len(data) / 2)
	b.Write([]byte(data))

	for _, rn := range string(data) {
		r, size, err := b.readRune()
		require.Nil(err)
		require.Equal(utf8.RuneLen(rn), size)
		require.Equal(rn, r)

	}
}

func TestBuffer_Next(t *testing.T) {
	tests := []struct {
		originalData []byte

		readChunk    int
		receivedData []byte
	}{
		{
			originalData: []byte("Hello, world!"),
			readChunk:    0,
			receivedData: []byte{},
		},
		{
			originalData: []byte("Hello, world!"),
			readChunk:    5,
			receivedData: []byte("Hello"),
		},
		{
			originalData: []byte("Hello, world!"),
			readChunk:    13,
			receivedData: []byte("Hello, world!"),
		},
		{
			originalData: []byte("Hello, world!"),
			readChunk:    20,
			receivedData: []byte("Hello, world!"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run("", func(t *testing.T) {
			require := require.New(t)

			b := NewBuffer(tt.originalData)
			defer b.Reset()

			data := b.Next(tt.readChunk)
			require.Equal(tt.receivedData, data)
		})
	}
}

func TestBuffer_ReadFrom(t *testing.T) {
	tests := []struct {
		before []byte
		data   []byte
		after  []byte
	}{
		{
			before: []byte("hello"),
			data:   []byte(""),
			after:  []byte("world"),
		},
		{
			before: []byte(""),
			data:   []byte("test"),
			after:  []byte(""),
		},
		{
			before: []byte("hello"),
			data:   []byte(": some test message: "),
			after:  []byte("world"),
		},
		{
			before: []byte("test"),
			data:   []byte(generateRandomString(1000)),
			after:  []byte("!!!"),
		},
		{
			before: []byte(""),
			data:   []byte(generateRandomString(2047)),
			after:  []byte(""),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run("", func(t *testing.T) {
			t.Parallel()

			fullMsg := append(append(tt.before, tt.data...), tt.after...)

			require := require.New(t)

			b := NewBuffer(nil)
			defer b.Reset()

			// Write before
			n, err := b.Write([]byte(tt.before))
			require.Nil(err)
			require.Equal(len(tt.before), n)

			// Write the data
			buffer := bytes.NewBuffer(nil)
			buffer.Write(tt.data)

			n1, err := b.ReadFrom(buffer)
			require.Nil(err)
			require.Equal(len(tt.data), int(n1))

			// Write after
			n, err = b.Write([]byte(tt.after))
			require.Nil(err)
			require.Equal(len(tt.after), n)

			// Check

			buffData, err := readByChunks(require, b, 32)
			require.Nil(err)
			require.Equal(fullMsg, buffData)
		})
	}

}

func TestBuffer_WriteSmth(t *testing.T) {
	tests := []struct {
		desc  string
		value interface{} // string, byte or rune
		size  int
	}{
		{desc: "write byte", value: byte('t'), size: 1},
		{desc: "write rune (cyrillic)", value: rune('П'), size: 2},
		{desc: "write rune (symbol)", value: rune('✓'), size: 3},
		{desc: "write string", value: "hello", size: 5},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.desc, func(t *testing.T) {
			require := require.New(t)

			b := NewBuffer(nil)
			defer b.Reset()

			switch v := tt.value.(type) {
			case byte:
				err := b.WriteByte(v)
				require.Nil(err)
			case rune:
				n, err := b.WriteRune(v)
				require.Nil(err)
				require.Equal(tt.size, n)
			case string:
				n, err := b.WriteString(v)
				require.Nil(err)
				require.Equal(tt.size, n)
			}
		})
	}
}

func TestBuffer_WriteTo(t *testing.T) {
	tests := []struct {
		data []byte
	}{
		{data: []byte(generateRandomString(1))},
		{data: []byte(generateRandomString(61))},
		{data: []byte(generateRandomString(513))},
		{data: []byte(generateRandomString(2056))},
	}

	for _, tt := range tests {
		tt := tt

		t.Run("", func(t *testing.T) {
			t.Parallel()

			require := require.New(t)

			b := NewBuffer(nil)
			defer b.Reset()

			// Write

			n, err := b.Write(tt.data)
			require.Nil(err)
			require.Equal(len(tt.data), n)

			// WriteTo
			buffer := bytes.NewBuffer(nil)
			n1, err := b.WriteTo(buffer)
			require.Nil(err)
			require.Equal(int64(len(tt.data)), n1)
			require.Equal(tt.data, buffer.Bytes())
		})
	}
}

func TestBuffer_FuzzTest(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 100; i++ {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			require := require.New(t)

			var (
				sliceSize      = rand.Intn(1<<10) + 1
				bufferSize     = rand.Intn(sliceSize * 2) // can be zero
				writeChunkSize = rand.Intn(sliceSize) + 1
				readChunkSize  = rand.Intn(sliceSize) + 1
			)

			defer func() {
				// Log only when test is failed
				if t.Failed() {
					t.Logf("sliceSize: %d; bufferSize: %d; writeChunkSize: %d; readChunkSize: %d\n",
						sliceSize, bufferSize, writeChunkSize, readChunkSize)
				}
			}()

			slice := make([]byte, sliceSize)
			for i := range slice {
				slice[i] = byte(rand.Intn(128))
			}

			b := NewBufferWithMaxMemorySize(bufferSize)
			defer b.Reset()
			// Write slice by chunks
			writeByChunks(require, b, slice, writeChunkSize)

			res, err := readByChunks(require, b, readChunkSize)
			require.Nil(err, "error during Read()")
			require.Equal(slice, res, "wrong content was read")
		})
	}
}

func writeByChunks(require *require.Assertions, b *Buffer, source []byte, chunk int) {
	// Write slice by chunks
	for i := 0; i < len(source); i += chunk {
		bound := i + chunk
		if bound > len(source) {
			bound = len(source)
		}

		_, err := b.Write(source[i:bound])
		require.Nil(err)
		require.Equal(bound, b.Len())
	}
}

func readByChunks(require *require.Assertions, b *Buffer, chunk int) ([]byte, error) {
	var (
		res      []byte
		dataRead int
		bufSize  = b.Len()
	)

	data := make([]byte, chunk)
	for {
		n, err := b.Read(data)
		dataRead += n
		data = data[:n]
		res = append(res, data...)
		data = data[:cap(data)]

		require.Equal(dataRead, bufSize-b.Len(), "Len() method returned wrong value")
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}
	}

	return res, nil
}

func BenchmarkBuffer(b *testing.B) {
	benchs := []struct {
		description    string
		dataSize       int
		maxBufferSize  int
		writeChunkSize int
		readChunkSize  int
	}{
		{
			description:    "Buffer size is greater than data",
			dataSize:       1 << 20, // 1MB
			maxBufferSize:  2 << 20, // 2MB
			writeChunkSize: 1024,
			readChunkSize:  2048,
		},
		{
			description:    "Buffer size is equal to data",
			dataSize:       1 << 20, // 1MB
			maxBufferSize:  1 << 20, // 1MB
			writeChunkSize: 1024,
			readChunkSize:  2048,
		},
		{
			description:    "Buffer size is less than data",
			dataSize:       20 << 20, // 20MB
			maxBufferSize:  1 << 20,  // 1MB
			writeChunkSize: 1024,
			readChunkSize:  2048,
		},
	}

	for _, bench := range benchs {
		b.Run(bench.description, func(b *testing.B) {
			slice := make([]byte, bench.dataSize)
			for i := range slice {
				slice[i] = byte(rand.Intn(128))
			}

			b.ResetTimer()

			b.Run("bytes.Buffer", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					buff := bytes.NewBuffer(nil)

					err := writeByChunksBenchmark(buff, slice, bench.writeChunkSize)
					if err != nil {
						b.Fatalf("error during Write(): %s", err)
					}

					_, err = readByChunksBenchmark(buff, bench.readChunkSize)
					if err != nil {
						b.Fatalf("error during Read(): %s", err)
					}
				}
			})

			b.Run("utils.Buffer", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					buff := NewBufferWithMaxMemorySize(bench.maxBufferSize)

					err := writeByChunksBenchmark(buff, slice, bench.writeChunkSize)
					if err != nil {
						b.Fatalf("error during Write(): %s", err)
					}

					_, err = readByChunksBenchmark(buff, bench.readChunkSize)
					if err != nil {
						b.Fatalf("error during Read(): %s", err)
					}
				}
			})
		})
	}

}

func writeByChunksBenchmark(w io.Writer, source []byte, chunk int) error {
	// Write slice by chunks
	for i := 0; i < len(source); i += chunk {
		bound := i + chunk
		if bound > len(source) {
			bound = len(source)
		}

		_, err := w.Write(source[i:bound])
		if err != nil {
			return err
		}
	}

	return nil
}

func readByChunksBenchmark(r io.Reader, chunk int) ([]byte, error) {
	var res []byte

	data := make([]byte, chunk)
	for {
		n, err := r.Read(data)
		data = data[:n]
		res = append(res, data...)
		data = data[:cap(data)]

		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}
	}

	return res, nil
}
