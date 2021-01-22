package lossless

import (
	"encoding/binary"
	"fmt"
	"github.com/petuhovskiy/compress/lossless/bits"
)

type Huffman struct{}

func (Huffman) Encode(bytes []byte) []byte {
	var ppm [256]*huffmanCoder
	for i := range ppm {
		ppm[i] = &huffmanCoder{}
		ppm[i].rebuildTable()
	}

	prev := 0

	w := &bits.Writer{}
	for _, elem := range bytes {
		hf := ppm[prev]

		b := int(elem)
		hf.encode(b, w)

		hf.freq[b]++
		hf.rebuildTable()

		prev = b
	}
	w.Close()

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

	var ppm [256]*huffmanCoder
	for i := range ppm {
		ppm[i] = &huffmanCoder{}
		ppm[i].rebuildTable()
	}

	u, n := binary.Uvarint(bytes)
	bytesCount := int(u)
	bytes = bytes[n:]

	r := &bits.Reader{Data: bytes}

	prev := 0

	res = []byte{}
	for i := 0; i < bytesCount; i++ {
		hf := ppm[prev]

		b := hf.decode(r)

		hf.freq[b]++
		hf.rebuildTable()

		res = append(res, byte(b))

		prev = b
	}

	return res, nil
}
