package tool

func Compress(data []byte, algoID string) ([]byte, error) {
	algo := FindAlgoByID(algoID)
	if algo == nil {
		return nil, ErrAlgoNotFound
	}

	res := []byte{algo.HeaderByte}
	res = append(res, algo.Backend.Encode(data)...)

	return res, nil
}
