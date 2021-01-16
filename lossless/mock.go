package lossless

type Mock struct{}

func (Mock) Encode(bytes []byte) []byte {
	return bytes
}

func (Mock) Decode(bytes []byte) ([]byte, error) {
	return bytes, nil
}
