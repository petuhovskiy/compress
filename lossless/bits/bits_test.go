package bits

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBits(t *testing.T) {
	f := 0

	for size := 0; size <= 256; size++ {
		var reference []byte

		w := &Writer{}
		for i := 0; i < size; i++ {
			b := byte(f & 0x1)
			reference = append(reference, b)

			w.Write(b)

			f++
			f %= 3
		}
		w.Close()

		r := &Reader{Data: w.Data}
		for i := 0; i < size; i++ {
			assert.Equal(t, reference[i], r.Read())
		}
	}
}
