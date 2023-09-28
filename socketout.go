package sutils

import (
	"bufio"
	"fmt"
	"io"
	"sync"
)

var OUT_TYPE_BYTES uint8 = 1
var OUT_TYPE_READER uint8 = 2
var SLICE_CAPACITY uint8 = 3

type SBodyOUT struct {
	Raw   []any
	Types []uint8
	Cond  sync.Cond
	pool  sync.Pool
	p     uint8
}

func NewSBodyOUT() SBodyOUT {
	return SBodyOUT{
		Cond:  *sync.NewCond(&sync.Mutex{}),
		Raw:   make([]any, SLICE_CAPACITY),
		Types: make([]uint8, SLICE_CAPACITY),
		p:     0,
		pool: sync.Pool{
			New: func() any {
				buf := make([]byte, 1024)
				return &buf
			},
		},
	}
}

func (so *SBodyOUT) add(raw any, t uint8) {
	so.Cond.L.Lock()
	defer so.Cond.L.Unlock()
	if so.p+1 >= SLICE_CAPACITY {
		so.Raw = append(make([]any, len(so.Raw)+int(SLICE_CAPACITY)), so.Raw...)
		so.Types = append(make([]uint8, len(so.Types)+int(SLICE_CAPACITY)), so.Types...)
		so.p = 0
	}
	so.Raw = append(so.Raw, raw)
	so.Types = append(so.Types, t)
	so.p += 1
}

func (so *SBodyOUT) AddBytes(raw []byte) {
	so.add(raw, OUT_TYPE_BYTES)
}

func (so *SBodyOUT) AddReader(raw io.Reader) {
	so.add(raw, OUT_TYPE_READER)
}

func (so *SBodyOUT) WriteTo(output io.Writer) (write_err error) {
	defer func() {
		if err := recover(); err != nil {
			write_err = fmt.Errorf("Unexpected write error: %w", err)
		} else {
			write_err = nil
		}
	}()
	so.Cond.L.Lock()
	defer so.Cond.L.Unlock()
	writer := bufio.NewWriter(output)
	for index, input := range so.Raw {
		t := so.Types[index]
		switch t {
		case OUT_TYPE_BYTES:
			writer.Write(input.([]byte))
		case OUT_TYPE_READER:
			reader := bufio.NewReader(input.(io.Reader))
			temp := so.pool.Get().(*[]byte)
			for {
				read, err := reader.Read(*temp)
				if err == io.EOF {
					break
				}
				writer.Write((*temp)[:read])
			}
			writer.Flush()
			so.pool.Put(temp)
		}
	}
	return nil
}
