package lossless

import "math/big"

var bigTable [65537]big.Int

type lr struct {
	l int // [
	r int // )

	bigL *big.Int
	bigR *big.Int
}

func newLR(l, r int) lr {
	return lr{
		l:    l,
		r:    r,
		bigL: &bigTable[l],
		bigR: &bigTable[r],
	}
}

func init() {
	for i := range bigTable {
		bigTable[i].SetInt64(int64(i))
	}
}
