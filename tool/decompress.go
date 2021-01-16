package tool

func Decompress(data []byte) ([]byte, error) {
	if len(data) < 1 {
		return nil, ErrInvalidHeader
	}

	algo := FindAlgoByHeaderByte(data[0])
	if algo == nil {
		return nil, ErrInvalidHeader
	}

	data = data[1:]

	res, err := algo.Backend.Decode(data)
	if err != nil {
		return nil, err
	}
	return res, nil
}
