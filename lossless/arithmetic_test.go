package lossless

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArithmetic_Encode(t *testing.T) {
	text := ""
	for i := 0; i < 256; i++ {
		for c := 'a'; c <= 'z'; c++ {
			text += string(c)
		}
	}

	before := []byte(text)
	after := Arithmetic{}.Encode(before)

	fmt.Printf("Compression factor (less is better): %.5f\n", float64(len(after))/float64(len(before)))
}

func BenchmarkArithmetic_Full(b *testing.B) {
	const size = 32 * 1024

	b.SetBytes(size)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		bytes := make([]byte, size, size)
		for i := range bytes {
			bytes[i] = byte(i % 32)
		}

		encoded := Arithmetic{}.Encode(bytes)

		_, err := Arithmetic{}.Decode(encoded)
		assert.Nil(b, err)
	}
}
