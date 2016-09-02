package model

import "crypto/sha512"

// HashOutputLen is the output length (in bytes) of Hash(...).
const HashOutputLen = 32

// Hash hashes the provided data with SHA-512 (first HashOutputLen-byte output size).
func hash(data ...[]byte) []byte {
	hasher := sha512.New()

	for i := 0; i < len(data); i++ {
		hasher.Write(data[i])
	}

	return hasher.Sum(nil)[:HashOutputLen]
}
