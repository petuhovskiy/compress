package lossless

import (
	"github.com/petuhovskiy/compress/lossless/distr"
)

// Arithmetic is a plain algorithm, which is slow and ineffective (possibly square complexity).
// Frequency table is static and encoded in the header.
type Arithmetic struct{}

func (Arithmetic) buildEncodeTable(freq []distr.Freq) ([]lr, error) {
	const n = 65536

	encodeTable := make([]lr, 256, 256)
	l := 0
	for _, f := range freq {
		r := l + f.Count
		encodeTable[f.Byte] = newLR(l, r)

		l = r
	}

	if l != n {
		return nil, distr.ErrInvalidFrequency
	}

	return encodeTable, nil
}

func (Arithmetic) buildDecodeTable(freq []distr.Freq) ([]byte, error) {
	const n = 65536

	bytes := make([]byte, n, n)
	i := 0
	for _, f := range freq {
		for j := i; j < i+f.Count; j++ {
			bytes[j] = f.Byte
		}
		i += f.Count
	}

	return bytes, nil
}

func (Arithmetic) Encode(bytes []byte) []byte {
	fullFreq := distr.Frequency(bytes)
	freq := distr.Normalize65k(fullFreq)

	encodeTable, err := Arithmetic{}.buildEncodeTable(freq)
	if err != nil {
		panic(err)
	}

	enc := newArithmeticEncoder()

	for _, b := range bytes {
		lib := encodeTable[b]

		enc.encode(lib)
	}

	res := distr.EncodeFreq256(fullFreq)
	res = append(res, enc.bytes()...)

	return res
}

func (Arithmetic) Decode(bytes []byte) ([]byte, error) {
	fullFreq, count, err := distr.DecodeFreq256(bytes)
	if err != nil {
		return nil, err
	}

	freq := distr.Normalize65k(fullFreq)

	sum := 0
	for _, f := range fullFreq {
		sum += f.Count
	}

	decode, err := Arithmetic{}.buildDecodeTable(freq)
	if err != nil {
		return nil, err
	}

	encodeTable, err := Arithmetic{}.buildEncodeTable(freq)
	if err != nil {
		return nil, err
	}

	bytes = bytes[count:]

	decoder := newArithmeticDecoder(bytes, sum)

	result := []byte{}

	for i := sum - 1; i >= 0; i-- {
		rem := decoder.decodeRemainder()

		// decoded byte
		b := decode[rem]

		// info about encode
		node := encodeTable[b]
		decoder.next(node)

		result = append(result, b)
	}

	return result, nil
}
