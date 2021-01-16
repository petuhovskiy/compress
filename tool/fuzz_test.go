package tool

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func testAllAlgos(t *testing.T, src []byte) {
	for _, algo := range algos {
		data1, err := Compress(src, algo.ID)
		assert.Nil(t, err)

		data2, err := Decompress(data1)
		assert.Nil(t, err)

		assert.Equal(t, src, data2)
	}
}

func TestAllSamples(t *testing.T) {
	testAllAlgos(t, []byte{0x00, 0x00, 0x00, 0x01})
	testAllAlgos(t, []byte{})
	testAllAlgos(t, []byte{0x0})
	testAllAlgos(t, []byte{0x01})
	testAllAlgos(t, []byte{0xff})
	testAllAlgos(t, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestAllFuzz(t *testing.T) {
	rnd := rand.New(rand.NewSource(123))

	for i := 0; i < 1; i++ {
		sz := rnd.Intn( 4)
		data := make([]byte, sz, sz)
		_, _= rnd.Read(data)

		testAllAlgos(t, data)
	}
}
