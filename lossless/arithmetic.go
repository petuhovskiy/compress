package lossless

import (
	"encoding/binary"
	"github.com/petuhovskiy/compress/lossless/distr"
	"math/big"
)

type lr struct {
	l int // [
	r int // )

	bigL *big.Int
	bigR *big.Int
}

func buildEncodeTable(freq []distr.Freq) ([]lr, error) {
	const n = 65536

	encodeTable := make([]lr, 256, 256)
	l := 0
	for _, f := range freq {
		r := l + f.Count
		encodeTable[f.Byte] = lr{
			l:    l,
			r:    r,
			bigL: big.NewInt(int64(l)),
			bigR: big.NewInt(int64(r)),
		}

		l = r
	}

	if l != n {
		return nil, distr.ErrInvalidFrequency
	}

	return encodeTable, nil
}

func buildDecodeTable(freq []distr.Freq) ([]byte, error) {
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

func uvarint(u uint64) []byte {
	b := make([]byte, 6, 6)
	n := binary.PutUvarint(b, u)
	return b[:n]
}

type Arithmetic struct{}

func (Arithmetic) Encode(bytes []byte) []byte {
	fullFreq := distr.Frequency(bytes)
	freq := distr.Normalize65k(fullFreq)

	const m = 16

	encodeTable, err := buildEncodeTable(freq)
	if err != nil {
		panic(err)
	}

	l := big.NewInt(0)
	diff := big.NewInt(1)

	tmp := &big.Int{}

	for _, b := range bytes {
		lib := encodeTable[b]

		l.Lsh(l, m)

		l.Add(l, tmp.Mul(diff, lib.bigL))
		diff.Mul(diff, tmp.Sub(lib.bigR, lib.bigL))
	}

	r := diff.Add(diff, l)
	mid := r.Sub(r, big.NewInt(1))
	trailingBits := 0
	for i := l.BitLen() - 1; i >= 0; i-- {
		if r.Bit(i) != l.Bit(i) {
			trailingBits = i
			break
		}
	}

	mid.Rsh(mid, uint(trailingBits))

	enc := distr.EncodeFreq256(fullFreq)
	enc = append(enc, uvarint(uint64(trailingBits))...)
	enc = append(enc, mid.Bytes()...)

	return enc
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

	decode, err := buildDecodeTable(freq)
	if err != nil {
		return nil, err
	}

	encodeTable, err := buildEncodeTable(freq)
	if err != nil {
		return nil, err
	}

	bytes = bytes[count:]

	trailingBits, count := binary.Uvarint(bytes)
	bytes = bytes[count:]

	num := (&big.Int{}).SetBytes(bytes)

	num.Lsh(num, uint(trailingBits))

	const m = 16

	tmp := &big.Int{}

	result := []byte{}

	for i := sum - 1; i >= 0; i-- {
		rem := tmp.Rsh(num, uint(i*m)).Int64()

		// decoded byte
		b := decode[rem]

		// info about encode
		node := encodeTable[b]

		tmp.Lsh(node.bigL, uint(i*m))
		num.Sub(num, tmp)

		tmp.Sub(node.bigR, node.bigL)
		num.Div(num, tmp)

		result = append(result, b)
	}

	return result, nil
}
