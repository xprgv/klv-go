package klv

import (
	"errors"
	"io"
)

var (
	ErrPartialWrite = errors.New("partial write")
	ErrKeyTooLong   = errors.New("key is too long")
)

type Encoder interface {
	Encode(chunks Chunks) (err error)
}

type encoder struct {
	w         io.Writer
	keyLength uint
}

func NewEncoder(w io.Writer, keyLength uint) Encoder {
	return &encoder{
		w:         w,
		keyLength: keyLength,
	}
}

func (e *encoder) Encode(chunks Chunks) error {
	for _, chunk := range chunks {
		buf := []byte{}

		delta := int(e.keyLength - uint(len(chunk.Key)))
		if delta < 0 {
			return ErrKeyTooLong
		} else if delta > 0 {
			pad := make([]byte, delta)
			chunk.Key = append(chunk.Key, pad...)
		}

		buf = append(buf, chunk.Key...)
		buf = append(buf, BerEncodeChunk(len(chunk.Value))...)
		buf = append(buf, chunk.Value...)

		n, err := e.w.Write(buf)
		if err != nil {
			return err
		}
		if n != len(buf) {
			return ErrPartialWrite
		}
	}
	return nil
}
