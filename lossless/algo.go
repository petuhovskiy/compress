package lossless

type Algorithm interface {
	Encode([]byte) []byte
	Decode([]byte) ([]byte, error)
}
