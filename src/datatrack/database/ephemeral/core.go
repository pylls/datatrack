package ephemeral

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

// Cipher is used to encrypt and decrypt.
var Cipher cipher.AEAD

// Setup creates a new ephemeral key and sets up the Cipher using
// AES-GCM mode for Encrypt and Decrypt.
func Setup() {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	Cipher = aead
}

// Encrypt encrypts data and returns the ciphertext. Panics on any failure.
func Encrypt(data []byte) []byte {
	nonce := make([]byte, Cipher.NonceSize(), Cipher.NonceSize()+len(data)+Cipher.Overhead())
	_, err := rand.Read(nonce)
	if err != nil {
		panic(err)
	}

	return append(nonce, Cipher.Seal(nil, nonce, data, nil)...)
}

// Decrypt attempts to decrypt the ciphertext. Panics on any failure.
func Decrypt(ct []byte) []byte {
	// special case
	if ct == nil {
		return nil
	}

	pt, err := Cipher.Open(nil, ct[:Cipher.NonceSize()], ct[Cipher.NonceSize():], nil)
	if err != nil {
		panic(err)
	}

	return pt
}
