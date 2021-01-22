package bits

type Writer struct {
	Length int
	Data   []byte

	Current byte
	Written int
}

func (w *Writer) Write(bit byte) {
	w.Length++

	w.Current |= bit << w.Written

	w.Written++
	if w.Written >= 8 {
		w.Data = append(w.Data, w.Current)
		w.Current = 0
		w.Written = 0
	}
}

func (w *Writer) Close() {
	if w.Written > 0 {
		w.Data = append(w.Data, w.Current)
		w.Written = 0
		w.Current = 0
	}
}
