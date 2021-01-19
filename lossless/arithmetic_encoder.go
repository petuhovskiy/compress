package lossless

import (
	"math/big"
)

type arithmeticEncoder struct {
	l    *big.Int
	diff *big.Int
	tmp  *big.Int
}

func newArithmeticEncoder() *arithmeticEncoder {
	return &arithmeticEncoder{
		l:    big.NewInt(0),
		diff: big.NewInt(1),
		tmp:  &big.Int{},
	}
}

func (e *arithmeticEncoder) encode(step lr) {
	const m = 16

	e.l.Lsh(e.l, m)

	e.l.Add(e.l, e.tmp.Mul(e.diff, step.bigL))
	e.diff.Mul(e.diff, e.tmp.Sub(step.bigR, step.bigL))
}

func (e *arithmeticEncoder) bytes() []byte {
	l := e.l

	r := &big.Int{}
	r.Add(e.diff, e.l)

	mid := r.Sub(r, big.NewInt(1))
	trailingBits := 0
	for i := l.BitLen() - 1; i >= 0; i-- {
		if r.Bit(i) != l.Bit(i) {
			trailingBits = i
			break
		}
	}

	mid.Rsh(mid, uint(trailingBits))

	res := uvarint(uint64(trailingBits))
	res = append(res, mid.Bytes()...)

	return res
}
