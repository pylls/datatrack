package util

import (
	"testing"
)

func TestHash(t *testing.T) {

}

func BenchmarkHash1KiB(b *testing.B) {
	b.StopTimer()
	benchmarkHash(1, b)
}

func BenchmarkHash10KiB(b *testing.B) {
	b.StopTimer()
	benchmarkHash(10, b)
}

func BenchmarkHash100KiB(b *testing.B) {
	b.StopTimer()
	benchmarkHash(100, b)
}

func BenchmarkHash1MiB(b *testing.B) {
	b.StopTimer()
	benchmarkHash(1024, b)
}

func BenchmarkHash10MiB(b *testing.B) {
	b.StopTimer()
	benchmarkHash(10*1024, b)
}

func benchmarkHash(length int, b *testing.B) {
	b.SetBytes(int64(length * 1024))
	data := make([]byte, length*1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Hash(data)
	}
}
