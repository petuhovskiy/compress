package lossless

import (
	"encoding/binary"
	"github.com/petuhovskiy/compress/lossless/distr"
)

// ArithmeticPPM is an adaptive algorithm, which counts frequencies online to encode frequent bytes efficiently.
// Data is divided in blocks, and every block is encoded almost independently.
type ArithmeticPPM struct{}

func (ArithmeticPPM) allBytes() []byte {
	var res []byte
	for i := 0; i <= 0xff; i++ {
		res = append(res, byte(i))
	}
	return res
}

func (ArithmeticPPM) encodeBlock(block []byte) []byte {
	freq := distr.Frequency(ArithmeticPPM{}.allBytes())

	encoded := uvarint(uint64(len(block)))
	enc := newArithmeticEncoder()

	for _, b := range block {
		normFreq := distr.Normalize65k(freq)
		encodeTable, _ := Arithmetic{}.buildEncodeTable(normFreq)

		lib := encodeTable[b]
		enc.encode(lib)

		freq[b].Count++
	}

	encoded = append(encoded, enc.bytes()...)
	return encoded
}

func (ArithmeticPPM) decodeBlock(encoded []byte) ([]byte, error) {
	freq := distr.Frequency(ArithmeticPPM{}.allBytes())

	u, offset := binary.Uvarint(encoded)
	encoded = encoded[offset:]
	bytesCount := int(u)

	decoder := newArithmeticDecoder(encoded, bytesCount)

	var result []byte
	for i := 0; i < bytesCount; i++ {
		rem := decoder.decodeRemainder()

		normFreq := distr.Normalize65k(freq)
		decodeTable, _ := Arithmetic{}.buildDecodeTable(normFreq)
		encodeTable, _ := Arithmetic{}.buildEncodeTable(normFreq)

		// decoded byte
		b := decodeTable[rem]
		result = append(result, b)
		freq[b].Count++

		// info about encode
		node := encodeTable[b]
		decoder.next(node)
	}

	return result, nil
}

func (ArithmeticPPM) Encode(bytes []byte) []byte {
	const blockSize = 16 * 1024

	res := []byte{}

	for offset := 0; offset < len(bytes); offset += blockSize {
		var block []byte

		if nextOffset := offset + blockSize; nextOffset >= len(bytes) {
			block = bytes[offset:]
		} else {
			block = bytes[offset:nextOffset]
		}

		encoded := ArithmeticPPM{}.encodeBlock(block)

		res = append(res, uvarint(uint64(len(encoded)))...)
		res = append(res, encoded...)
	}

	return res
}

func (ArithmeticPPM) Decode(bytes []byte) ([]byte, error) {
	result := []byte{}

	for offset := 0; offset < len(bytes); {
		u, count := binary.Uvarint(bytes[offset:])
		offset += count

		length := int(u)
		block := bytes[offset : offset+length]
		offset += length

		decoded, err := ArithmeticPPM{}.decodeBlock(block)
		if err != nil {
			return nil, err
		}

		result = append(result, decoded...)
	}

	return result, nil
}
