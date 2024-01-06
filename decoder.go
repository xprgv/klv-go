package klv

import (
	"errors"
	"io"
)

var (
	ErrKeyLength   = errors.New("incorrect key length")
	ErrValueLength = errors.New("incorrect value length")
)

type Decoder interface {
	Take() (Chunk, error)
	TakeAll() (Chunks, error)
}

type decoder struct {
	r         io.Reader
	keyLength int
}

func NewDecoder(r io.Reader, keyLength int) Decoder {
	return &decoder{
		r:         r,
		keyLength: keyLength,
	}
}

func (d *decoder) Take() (Chunk, error) {
	var (
		rawBuf = make([]byte, 0)
		keyBuf = make([]byte, d.keyLength)
	)

	n, err := d.r.Read(keyBuf)
	if err != nil {
		return Chunk{}, err
	}

	if n != d.keyLength {
		return Chunk{}, ErrKeyLength
	}

	rawBuf = append(rawBuf, keyBuf...)

	length, lengthRaw, err := BerDecode(d.r)
	if err != nil {
		return Chunk{}, err
	}

	rawBuf = append(rawBuf, lengthRaw...)

	valueBuf := make([]byte, length)

	n, err = d.r.Read(valueBuf)
	if err != nil {
		return Chunk{}, err
	}

	if length != uint64(n) {
		return Chunk{}, ErrValueLength
	}

	rawBuf = append(rawBuf, valueBuf...)

	return Chunk{
		raw:       rawBuf,
		Key:       keyBuf,
		length:    length,
		lengthRaw: lengthRaw,
		Value:     valueBuf,
	}, nil
}

func (d *decoder) TakeAll() (Chunks, error) {
	chunks := Chunks{}
	for {
		chunk, err := d.Take()
		if err != nil {
			if err == io.EOF {
				return chunks, nil
			}
			return chunks, err
		}

		chunks = append(chunks, chunk)
	}
}
