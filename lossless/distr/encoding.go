package distr

import "encoding/binary"

func EncodeFreq256(freqs []Freq) []byte {
	// 1kb must be enough
	res := make([]byte, 1024)

	ptr := 0

	for _, f := range freqs {
		ptr += binary.PutUvarint(res[ptr:], uint64(f.Count))
	}

	return res[:ptr]
}

func DecodeFreq256(bytes []byte) (freqs []Freq, count int, err error) {
	freqs = make([]Freq, 256, 256)

	ptr := 0

	for i := range freqs {
		freqs[i].Byte = byte(i)

		u, cnt := binary.Uvarint(bytes[ptr:])
		ptr += cnt
		freqs[i].Count = int(u)
	}

	return freqs, ptr, nil
}
