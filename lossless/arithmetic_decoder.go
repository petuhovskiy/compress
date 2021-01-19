package lossless

import (
	"encoding/binary"
	"math/big"
)

type arithmeticDecoder struct {
	num *big.Int
	tmp *big.Int

	i int
}

func newArithmeticDecoder(bytes []byte, bytesSize int) *arithmeticDecoder {
	trailingBits, count := binary.Uvarint(bytes)
	bytes = bytes[count:]

	num := (&big.Int{}).SetBytes(bytes)

	num.Lsh(num, uint(trailingBits))

	return &arithmeticDecoder{
		num: num,
		tmp: &big.Int{},
		i:   bytesSize - 1,
	}
}

func (d *arithmeticDecoder) decodeRemainder() int64 {
	const m = 16

	rem := d.tmp.Rsh(d.num, uint(d.i*m)).Int64()
	return rem
}

func (d *arithmeticDecoder) next(step lr) {
	const m = 16

	i := d.i
	d.i--

	d.tmp.Lsh(step.bigL, uint(i*m))
	d.num.Sub(d.num, d.tmp)

	d.tmp.Sub(step.bigR, step.bigL)
	d.num.Div(d.num, d.tmp)
}
