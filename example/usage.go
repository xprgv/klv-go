package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/xprgv/klv-go"
)

func main() {
	var (
		chunks = klv.Chunks{
			{Key: []byte("hello"), Value: []byte("world")},
			{Key: []byte("user"), Value: []byte("xprgv")},
		}

		buffer = bytes.NewBuffer(nil)
	)

	for _, chunk := range chunks {
		fmt.Println("encoded:", string(chunk.Key), string(chunk.Value))
	}
	fmt.Println()

	if err := klv.NewEncoder(buffer, 5).Encode(chunks); err != nil {
		log.Fatal("Failed to encode: ", err)
	}

	chunks, err := klv.NewDecoder(buffer, 5).TakeAll()
	if err != nil {
		log.Fatal("Failed to decode: ", err)
	}

	for _, chunk := range chunks {
		fmt.Println("decoded:", string(chunk.Key), string(chunk.Value))
	}
}
