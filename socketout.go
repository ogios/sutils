package sutils

import (
	"io"
	"sync"
)

type SBodyOUT struct {
	Raw   []any
	Types []uint8
	Cond  sync.Cond
}

func NewSBodyOUT() SBodyOUT {
	return SBodyOUT{
		Cond:  *sync.NewCond(&sync.Mutex{}),
		Raw:   make([]any, 0),
		Types: make([]uint8, 0),
	}
}

func (so *SBodyOUT) AddBytes(raw []byte) {
	so.Cond.L.Lock()
	defer so.Cond.L.Unlock()
	so.Raw = append(so.Raw, &raw)
	so.Types = append(so.Types, 1)
}

func (so *SBodyOUT) AddReader(raw io.Reader) {
	so.Cond.L.Lock()
	defer so.Cond.L.Unlock()
	so.Raw = append(so.Raw, raw)
	so.Types = append(so.Types, 2)
}
