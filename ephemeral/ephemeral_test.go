package ephemeral

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"testing"
)

func TestEncDec(t *testing.T) {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		t.Fatalf("failed to generate key: %s", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("failed to create AES block cipher: %s", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatalf("failed to create GCM cipher-mode: %s", err)
	}
	Cipher = aead

	msg := []byte("hello world")
	ct := Encrypt(msg)
	pt := Decrypt(ct)
	if !bytes.Equal(msg, pt) {
		t.Fatal("invalid ciphertext")
	}
}
