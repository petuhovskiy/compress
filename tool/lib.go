package tool

import (
	"github.com/petuhovskiy/compress/lossless"
)

type Algo struct {
	ID         string
	HeaderByte byte
	Backend    lossless.Algorithm
}

var algos = []*Algo{
	//{
	//	ID:         "pmm",
	//	HeaderByte: 0x12,
	//	Backend:    lossless.PMM{},
	//},
	{
		ID:         "mock",
		HeaderByte: 0x13,
		Backend:    lossless.Mock{},
	},
	{
		ID:         "arithm_static",
		HeaderByte: 0x14,
		Backend:    lossless.Arithmetic{},
	},
	{
		ID:         "arithm_ppm",
		HeaderByte: 0x15,
		Backend:    lossless.ArithmeticPPM{},
	},
}

func FindAlgoByID(id string) *Algo {
	for _, algo := range algos {
		if algo.ID == id {
			return algo
		}
	}

	return nil
}

func FindAlgoByHeaderByte(b byte) *Algo {
	for _, algo := range algos {
		if algo.HeaderByte == b {
			return algo
		}
	}

	return nil
}
