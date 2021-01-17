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

	const n = 65536
	const m = 16
	bigN := big.NewInt(n)

	encodeTable, err := buildEncodeTable(freq)
	if err != nil {
		panic(err)
	}

	l := big.NewInt(0)
	r := big.NewInt(1)
	base := big.NewInt(1)

	for _, b := range bytes {
		lib := encodeTable[b]

		diff := (&big.Int{}).Sub(r, l)

		tmp := &big.Int{}
		base = base.Mul(base, bigN)

		l = l.Mul(l, bigN)
		nl := (&big.Int{}).Add(l, tmp.Mul(diff, lib.bigL))
		nr := (&big.Int{}).Add(l, tmp.Mul(diff, lib.bigR))

		l = nl
		r = nr

		//spew.Dump("lr", l, r)
	}

	enc := distr.EncodeFreq256(fullFreq)
	//enc = append(enc, bytes...)

	//spew.Dump(l)
	//spew.Dump(r)
	//
	//spew.Dump(l.BitLen())
	//spew.Dump(r.BitLen())
	//
	//spew.Dump((&big.Int{}).Sub(r, l))
	//spew.Dump((&big.Int{}).Sub(r, l).BitLen())

	mid := r.Sub(r, big.NewInt(1))
	trailingBits := 0
	for i := l.BitLen()-1; i >= 0; i-- {
		if r.Bit(i) != l.Bit(i) {
			//spew.Dump("BIT", i, r.Bit(i), l.Bit(i))
			trailingBits = i
			break
		}
	}

	//spew.Dump(l.BitLen())
	//spew.Dump(r.BitLen())
	//spew.Dump(base.BitLen())
	//
	//spew.Dump(l.TrailingZeroBits())
	//spew.Dump(r.TrailingZeroBits())

	//trailingBits := mid.TrailingZeroBits()
	mid.Rsh(mid, uint(trailingBits))

	//spew.Dump(mid)

	//mid.Lsh(mid, uint(trailingBits))
	//spew.Dump("after", mid)
	//mid.Rsh(mid, uint(trailingBits))

	// TODO: ?
	enc = append(enc, uvarint(uint64(trailingBits))...)
	enc = append(enc, mid.Bytes()...)

	//spew.Dump(l.Bits())
	//spew.Dump(r.Bits())
	//
	//spew.Dump(l.String())
	//spew.Dump(r.String())
	//
	//spew.Dump(base.Bits())

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

	//if len(bytes) != sum {
	//	return nil, distr.ErrInvalidFrequency
	//}

	num := (&big.Int{}).SetBytes(bytes)

	num.Lsh(num, uint(trailingBits))

	//spew.Dump(num)

	const n = 65536
	const m = 16
	//bigN := big.NewInt(n)

	result := []byte{}

	for i := sum - 1; i >= 0; i-- {
		cur := (&big.Int{}).Rsh(num, uint(i*m))
		//spew.Dump(cur)

		abc := cur.Int64()
		if abc == n {
			abc--
		}
		b := decode[abc]

		node := encodeTable[b]

		nl := (&big.Int{}).Lsh(node.bigL, uint(i*m))
		num = num.Sub(num, nl)
		//num = num.Mul(num, bigN)
		num = num.Div(num, nl.Sub(node.bigR, node.bigL))

		result = append(result, b)

		//spew.Dump("num", num)

		//spew.Dump(abc, b)
	}

	return result, nil
}
