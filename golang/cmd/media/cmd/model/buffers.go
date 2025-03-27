package model

import (
	"bytes"
	"io"
	"strings"
)

type IBuffer interface {
	GetData() []byte
	Reset()
	Append(sb *strings.Builder)
}

type M3u8Buffer struct {
	Data     []byte
	capacity int
}

// GetData implements IBuffer.
func (b *M3u8Buffer) GetData() []byte {
	return b.Data
}

func NewM3u8Buffer(size, overflowBuffer int) *M3u8Buffer {
	return &M3u8Buffer{
		capacity: size + overflowBuffer, //extra for overflow buffer
		Data:     make([]byte, 0, size+overflowBuffer),
	}
}

func (b *M3u8Buffer) Reset() {
	b.Data = make([]byte, 0, cap(b.Data))
}

func (b *M3u8Buffer) Append(sb *strings.Builder) {
	b.Data = append(b.Data, sb.String()...)
}

func (b *M3u8Buffer) ReadFrom(r io.Reader) (int64, error) {
	b.Data = b.Data[:0]
	var total int64
	for {
		if len(b.Data) == cap(b.Data) {
			newData := make([]byte, len(b.Data), cap(b.Data)*2)
			copy(newData, b.Data)
			b.Data = newData
		}
		n, err := r.Read(b.Data[len(b.Data):cap(b.Data)])
		b.Data = b.Data[:len(b.Data)+n]
		total += int64(n)
		if err == io.EOF {
			return total, nil
		}
		if err != nil {
			return total, err
		}
	}
}

func (b *M3u8Buffer) ReadLine() (string, error) {
	i := bytes.IndexByte(b.Data, '\n')
	if i >= 0 {
		line := string(b.Data[:i+1])
		b.Data = b.Data[i+1:]
		return line, nil
	}
	if len(b.Data) > 0 {
		line := string(b.Data)
		b.Data = b.Data[:0]
		err := io.EOF
		return line, err
	}
	return "", io.EOF
}
