package bits

type Writer struct {
	Bytes []byte
	Exp   byte
}

func (w *Writer) Write(b bool) {
	if w.Exp == 0 {
		w.Exp = 1
		w.Bytes = append(w.Bytes, 0)
	}

	if b {
		w.Bytes[len(w.Bytes)-1] |= w.Exp
	}
	w.Exp <<= 1
}
