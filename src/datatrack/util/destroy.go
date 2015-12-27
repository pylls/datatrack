package util

import (
	"crypto/rand"
	"io"
)

func DestroyByteSlice(b []byte) {
	randomData := make([]byte, len(b))
	n, err := io.ReadFull(rand.Reader, randomData)
	if n != len(b) || err != nil {
		panic("failed to read random data")
	}
	for i := 0; i < len(b); i++ {
		b[i] = b[i] ^ randomData[i]
	}
}

func DestroyByteSliceShallow(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = b[i] | 0xFF // set all bits to I
		b[i] = b[i] & 0x00 // set all bits to 0
	}
}
