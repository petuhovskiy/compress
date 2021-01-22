package lossless

import (
	"encoding/binary"
	"fmt"
	"github.com/petuhovskiy/compress/lossless/bits"
)

type Huffman struct{}

func (Huffman) Encode(bytes []byte) []byte {
	hf := &huffmanCoder{}
	hf.rebuildTable()

	w := &bits.Writer{}
	for _, b := range bytes {
		hf.encode(int(b), w)
	}

	res := uvarint(uint64(len(bytes)))
	res = append(res, w.Data...)

	return res
}

func (Huffman) Decode(bytes []byte) (res []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	u, n := binary.Uvarint(bytes)
	bytesCount := int(u)
	bytes = bytes[n:]

	r := &bits.Reader{Data: bytes}

	hf := &huffmanCoder{}
	hf.rebuildTable()

	res = []byte{}
	for i := 0; i < bytesCount; i++ {
		res = append(res, byte(hf.decode(r)))
	}

	return res, nil
}
