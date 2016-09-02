package util

import (
	"crypto/subtle"
	"encoding/binary"
	"errors"
	"math"
)

func Itob(i int) (r []byte) {
	r = make([]byte, 8)
	binary.PutVarint(r, int64(i))
	return
}

func Btoi(i []byte) (r int) {
	t, _ := binary.Varint(i)
	return int(t)
}

func Equal(x, y []byte) bool {
	if len(x) != len(y) {
		return false
	}
	// ConstantTimeCompare returns 1 iff the two equal length slices, x and y, have equal contents.
	// The time taken is a function of the length of the slices and is independent of the contents.
	return subtle.ConstantTimeCompare(x, y) == 1
}

func ToByteArray32(s []byte) (*[32]byte, error) {
	if len(s) != 32 {
		return nil, errors.New("slice has to be 32 bytes long")
	}
	var result [32]byte
	for i := 0; i < 32; i++ {
		result[i] = s[i]
	}
	return &result, nil
}

func ToByteArray64(s []byte) (*[64]byte, error) {
	if len(s) != 64 {
		return nil, errors.New("slice has to be 64 bytes long")
	}
	var result [64]byte
	for i := 0; i < 64; i++ {
		result[i] = s[i]
	}
	return &result, nil
}

func Pow(a, b int) int {
	return int(math.Pow(float64(a), float64(b)))
}
