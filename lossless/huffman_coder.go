package lossless

import (
	"container/heap"
	"github.com/petuhovskiy/compress/lossless/bits"
)

const (
	bytesCount = 1 << 8
)

type huffmanNode struct {
	// contains number for jumping by single bit.
	// if number if positive, this is the byte.
	// otherwise, this is the offset to jump backwards.
	next [2]int
}

type huffmanTemp struct {
	sum    int
	isByte bool
	value  int
	height int
}

func (h huffmanTemp) getNext(pos int) int {
	if h.isByte {
		return h.value
	}
	return h.value - pos
}

type huffmanHeap []huffmanTemp

func (h huffmanHeap) Len() int {
	return len(h)
}

func (h huffmanHeap) Less(i, j int) bool {
	if h[i].sum != h[j].sum {
		return h[i].sum < h[j].sum
	}

	return h[i].height < h[j].height
}

func (h huffmanHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *huffmanHeap) Push(x interface{}) {
	*h = append(*h, x.(huffmanTemp))
}

func (h *huffmanHeap) Pop() interface{} {
	n := len(*h)
	x := (*h)[n-1]
	*h = (*h)[0 : n-1]
	return x
}

type huffmanCoder struct {
	freq  [bytesCount]int
	table [bytesCount - 1]huffmanNode

	temp     huffmanHeap
	tempBits []int
}

func newHuffmanCoder() *huffmanCoder {
	return &huffmanCoder{
		temp:     make([]huffmanTemp, bytesCount),
		tempBits: make([]int, bytesCount),
	}
}

func (h *huffmanCoder) rebuildTable() {
	temp := h.temp[:0]

	for i, f := range h.freq {
		temp = append(temp, huffmanTemp{
			sum:    f,
			isByte: true,
			value:  i,
		})
	}

	heap.Init(&temp)

	pos := 0
	for size := bytesCount; size > 1; size-- {
		a := heap.Pop(&temp).(huffmanTemp)
		b := heap.Pop(&temp).(huffmanTemp)

		h.table[pos] = huffmanNode{
			next: [2]int{a.getNext(pos), b.getNext(pos)},
		}

		heap.Push(&temp, huffmanTemp{
			sum:    a.sum + b.sum,
			isByte: false,
			value:  pos,
			height: a.height + b.height + 1,
		})
		pos++
	}
}

func (h *huffmanCoder) encode(b int, w *bits.Writer) {
	it := 0
	temp := h.tempBits[:0]

	prev := -1
	for i, node := range h.table {
		if (prev == -1 && node.next[0] == b) || (prev != -1 && node.next[0]+i == prev) {
			prev = i
			temp = append(temp, 0)
			it++
		} else if (prev == -1 && node.next[1] == b) || (prev != -1 && node.next[1]+i == prev) {
			prev = i
			temp = append(temp, 1)
			it++
		}
	}

	for it > 0 {
		it--
		w.Write(byte(temp[it] & 1))
	}
}

func (h *huffmanCoder) decode(r *bits.Reader) int {
	node := len(h.table) - 1

	for {
		b := r.Read()
		next := h.table[node].next[b]
		if next >= 0 {
			return next
		}

		node += next
	}
}
