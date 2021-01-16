package distr

import "sort"

type Freq struct {
	Byte  byte
	Count int
}

type FreqSort []Freq

func (f FreqSort) Len() int {
	return len(f)
}

func (f FreqSort) Less(i, j int) bool {
	return f[i].Count > f[j].Count
}

func (f FreqSort) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func Frequency(bytes []byte) []Freq {
	freq := make([]Freq, 256, 256)
	for i := range freq {
		freq[i].Byte = byte(i)
	}

	for _, b := range bytes {
		freq[b].Count++
	}

	sort.Sort(FreqSort(freq))
	return freq
}

// Normalize65k computes byte frequency distribution, normalized to 65536.
func Normalize65k(freqSlice []Freq) []Freq {
	var freq []Freq
	freq = append(freq, freqSlice...)

	sum := 0
	for _, f := range freq {
		sum += f.Count
	}

	n := 65536

	for j := range freq {
		i := len(freq) - 1 - j

		cnt := freq[i].Count
		if cnt == 0 {
			continue
		}

		norm := float64(cnt) / float64(sum) * float64(n)
		pts := int(norm)
		if pts < 1 {
			pts = 1
		}

		freq[i].Count = pts
		sum -= cnt
		n -= pts
	}

	freq[0].Count += n

	return freq
}
