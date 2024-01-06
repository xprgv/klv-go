package klv

import (
	"bytes"
	"fmt"
)

type Chunk struct {
	Key   []byte
	Value []byte

	length    uint64
	lengthRaw []byte

	raw []byte
}

func (c Chunk) String() string { return fmt.Sprintf("[%b] - [%b]", c.Key, c.Value) }

func (c Chunk) Raw() []byte { return c.raw }

func (c Chunk) BerLength() []byte { return c.lengthRaw }

type Chunks []Chunk

func (chunks Chunks) FindByKey(key []byte) Chunks {
	var result Chunks

	for _, chunk := range chunks {
		if bytes.Equal(chunk.Key, key) {
			result = append(result, chunk)
		}
	}

	return result
}

func (chunks Chunks) HasKey(key []byte) bool {
	for _, chunk := range chunks {
		if bytes.Equal(chunk.Key, key) {
			return true
		}
	}

	return false
}
