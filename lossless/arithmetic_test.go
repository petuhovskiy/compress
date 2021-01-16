package lossless

import (
	"fmt"
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

	fmt.Printf("Compression factor (less is better): %.5f\n", float64(len(after)) / float64(len(before)))
}
