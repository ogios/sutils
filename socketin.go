package sutils

import (
	"bufio"
	"errors"
	"fmt"
	"math"
)

type SBodyIN struct {
	raw           *bufio.Reader
	readed        int
	CurrSecLength int
	BufferSize    int
	buffer        [16]byte
}

func NewSBodyIn(reader *bufio.Reader) SBodyIN {
	return SBodyIN{
		raw:           reader,
		readed:        0,
		BufferSize:    1024,
		CurrSecLength: 0,
		buffer:        [16]byte{},
	}
}

func (si *SBodyIN) Next() (int, error) {
	if si.CurrSecLength < si.readed {
		return 0, errors.New("please read all of current section")
	}
	readlen, err := si.raw.ReadBytes(0)
	if err != nil {
		return 0, err
	}
	if len(readlen) < 1 {
		return 0, errors.New("no length provided")
	}
	rawlen := readlen[:len(readlen)-1]
	total := 0
	for index, b := range rawlen {
		pow := float64(len(rawlen) - 1 - index)
		feat := int(math.Pow(255, pow))
		total += int(b) * feat
	}
	si.CurrSecLength = total
	si.readed = 0
	return total, nil
}

func (si *SBodyIN) GetSec() ([]byte, error) {
	if si.readed < si.CurrSecLength {
		bs := make([]byte, si.CurrSecLength-si.readed)
		readed, err := si.raw.Read(bs)
		if err != nil {
			return nil, err
		}
		si.readed += readed
		return bs, nil
	} else {
		length, err := si.Next()
		if err != nil {
			return nil, err
		}
		temp := make([]byte, length)
		readlength, err := si.raw.Read(temp)
		if err != nil {
			return nil, err
		}
		if readlength != length {
			fmt.Printf("Wrong length: read:%d - next:%d\n", readlength, length)
		}
		return temp, nil
	}
}

func (si *SBodyIN) Read(buf []byte) (int, error) {
	if si.CurrSecLength == si.readed {
		return 0, errors.New("no more bytes for current section")
	}
	if len(buf) <= si.CurrSecLength-si.readed {
		return si.raw.Read(buf)
	} else {
		temp := make([]byte, si.CurrSecLength-si.readed)
		length, err := si.raw.Read(temp)
		if err != nil {
			return 0, err
		}
		copy(buf[:length], temp)
		return length, nil
	}
}
