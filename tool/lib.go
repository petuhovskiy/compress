package tool

import (
	"github.com/petuhovskiy/compress/lossless"
)

type Algo struct {
	ID          string
	HeaderByte  byte
	Backend     lossless.Algorithm
	Description string
}

var Supported = []*Algo{
	{
		ID:          "mock",
		HeaderByte:  0x13,
		Backend:     lossless.Mock{},
		Description: "This algorithm does nothing.",
	},
	{
		ID:          "arithm_static",
		HeaderByte:  0x14,
		Backend:     lossless.Arithmetic{},
		Description: "Arithmetic algo, which writes header with frequencies, and works very slowly on data >32k bytes.",
	},
	{
		ID:          "ppm",
		HeaderByte:  0x15,
		Backend:     lossless.ArithmeticPPM{},
		Description: "Arithmetic algo, which groups data in blocks and calculates frequencies on the fly. Works a bit slow.",
	},
	{
		ID:          "huffman",
		HeaderByte:  0x16,
		Backend:     lossless.Huffman{},
		Description: "Huffman algorithm, with ppm enabled. Context is previous byte in the sequence.",
	},
}

func FindAlgoByID(id string) *Algo {
	for _, algo := range Supported {
		if algo.ID == id {
			return algo
		}
	}

	return nil
}

func FindAlgoByHeaderByte(b byte) *Algo {
	for _, algo := range Supported {
		if algo.HeaderByte == b {
			return algo
		}
	}

	return nil
}
