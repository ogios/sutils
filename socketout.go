package sutils

import (
	"bufio"
	"fmt"
	"io"
	"sync"

	"golang.org/x/exp/slog"
)

var OUT_TYPE_BYTES uint8 = 1
var OUT_TYPE_READER uint8 = 2
var SLICE_CAPACITY uint8 = 6
var ContentLengthMax uint8 = 8

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
		Raw:   make([]any, 0),
		Types: make([]uint8, 0),
		p:     0,
		pool: sync.Pool{
			New: func() any {
				buf := make([]byte, 1024)
				return &buf
			},
		},
	}
}

func (so *SBodyOUT) add(raw any, t uint8, lenb []byte) {
	so.Cond.L.Lock()
	defer so.Cond.L.Unlock()
	if so.p+1 >= SLICE_CAPACITY {
		so.Raw = append(make([]any, len(so.Raw)+int(SLICE_CAPACITY)), so.Raw...)
		so.Types = append(make([]uint8, len(so.Types)+int(SLICE_CAPACITY)), so.Types...)
		so.p = 0
	}
	so.Raw = append(so.Raw, lenb, raw)
	so.Types = append(so.Types, t)
	so.p += 1
}

func getlength(length int) ([]byte, error) {
	content_length := make([]byte, ContentLengthMax)
	var index uint8 = 0
	for length >= 255 {
		if index >= ContentLengthMax {
			return nil, fmt.Errorf("message length too long, max size is 255**%d", ContentLengthMax)
		}
		content_length[index] = byte(length % 255)
		length /= 255
		index++
	}
	content_length[index] = byte(length)
	return content_length, nil
}

func (so *SBodyOUT) AddBytes(raw []byte) error {
	length := len(raw)
	content_length, err := getlength(length)
	if err != nil {
		return err
	}
	so.add(raw, OUT_TYPE_BYTES, content_length)
	return nil
}

func (so *SBodyOUT) AddReader(raw io.Reader, length int) error {
	content_length, err := getlength(length)
	if err != nil {
		return err
	}
	so.add(raw, OUT_TYPE_READER, content_length)
	return nil
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
		slog.Debug(fmt.Sprintln("index", index))
		slog.Debug(fmt.Sprintln("input", input))
		if input == nil {
			continue
		}
		if index%2 == 0 {
			slog.Debug("Writing length bytes")
			writer.Write(input.([]byte))
		} else {
			t := so.Types[(index-1)/2]
			slog.Debug(fmt.Sprintln("Type:", t))
			switch t {
			case OUT_TYPE_BYTES:
				slog.Debug("Writing raw bytes")
				writer.Write(input.([]byte))
			case OUT_TYPE_READER:
				slog.Debug("Writing reader bytes")
				reader := bufio.NewReader(input.(io.Reader))
				temp := so.pool.Get().(*[]byte)
				for {
					read, err := reader.Read(*temp)
					if err == io.EOF {
						break
					}
					writer.Write((*temp)[:read])
				}
				so.pool.Put(temp)
			default:
				slog.Error(fmt.Sprintln("Unknow type:", t))
			}
		}
	}
	writer.Flush()
	return nil
}
