package tool

import "fmt"

func Compress(data []byte, algoID string) ([]byte, error) {
	algo := FindAlgoByID(algoID)
	if algo == nil {
		return nil, ErrAlgoNotFound
	}

	res := []byte{algo.HeaderByte}
	res = append(res, algo.Backend.Encode(data)...)

	return res, nil
}

func CompressBest(data []byte) ([]byte, error) {
	beforeSize := len(data)

	var res []byte

	for i, algo := range Supported {
		dataCopy := make([]byte, len(data), len(data))
		copy(dataCopy, data)

		enc := algo.Backend.Encode(dataCopy)

		afterSize := len(enc)
		if beforeSize > 0 {
			factor := float64(afterSize) / float64(beforeSize) * 100
			fmt.Printf("Algo %s has a result of %.2f%%\n", algo.ID, 100-factor)
		}

		if i == 0 || len(enc) < len(res) {
			res = enc
		}
	}

	return res, nil
}
