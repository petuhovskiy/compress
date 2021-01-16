package tool

import "errors"

var (
	ErrAlgoNotFound  = errors.New("algo not found, look at supported algos")
	ErrInvalidHeader = errors.New("invalid data in header byte")
)
