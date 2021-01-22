package lossless

import (
	"github.com/petuhovskiy/compress/lossless/bits"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHuffman(t *testing.T) {
	w := &bits.Writer{}
	hf := newHuffmanCoder()
	hf.rebuildTable()

	for i := 0; i < bytesCount; i++ {
		hf.encode(i, w)
	}

	r := &bits.Reader{Data: w.Data}
	for i := 0; i < bytesCount; i++ {
		assert.Equal(t, i, hf.decode(r))
	}
}
