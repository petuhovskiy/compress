package distr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testEnc(t *testing.T, b []byte) {
	freq := Frequency(b)

	encoded := EncodeFreq256(freq)
	decoded, count, err := DecodeFreq256(encoded)
	assert.Nil(t, err)

	assert.Equal(t, len(encoded), count)
	assert.EqualValues(t, freq, decoded)
}

func TestEncoding(t *testing.T) {
	testEnc(t, []byte{})
	testEnc(t, []byte{0x0, 0x1, 0x2, 0x0, 0x1, 0x2})
	testEnc(t, []byte{0x0, 0x0, 0x0, 0x1})
}
