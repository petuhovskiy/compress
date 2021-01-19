package lossless

import "encoding/binary"

func uvarint(u uint64) []byte {
	b := make([]byte, binary.MaxVarintLen64, binary.MaxVarintLen64)
	n := binary.PutUvarint(b, u)
	return b[:n]
}
