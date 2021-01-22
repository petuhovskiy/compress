package bits

type Reader struct {
	Data     []byte
	BitsRead int
}

func (r *Reader) Read() byte {
	b := r.Data[0] & 1
	r.Data[0] >>= 1

	r.BitsRead++
	if r.BitsRead % 8 == 0 {
		r.Data = r.Data[1:]
	}

	return b
}
