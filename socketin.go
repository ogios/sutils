package sutils

import (
	"bufio"
	"errors"
	"fmt"
	"math"

	"golang.org/x/exp/slog"
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
	// readlen := make([]byte, ContentLengthMax)
	index := -1
	t := 0

	var main func() error
	main = func() error {
		current, err := si.raw.ReadByte()
		if err != nil {
			return err
		}
		if current != 255 {
			index++
			err := main()
			if err != nil {
				return err
			}
			feat := int(math.Pow(255, float64(index)))
			t += int(current) * feat
			index--
		}
		return nil
	}
	err := main()
	if err != nil {
		return 0, err
	}
	return t, nil

	// readed, err := si.raw.Read(readlen)
	// slog.Debug(fmt.Sprintln("Length:", readlen))
	// if err != nil || readed != int(ContentLengthMax) {
	// 	return 0, err
	// }
	// if len(readlen) < 1 {
	// 	return 0, errors.New("no length provided")
	// }
	// rawlen := readlen[:len(readlen)-1]
	// total := 0
	// for index, b := range rawlen {
	// 	feat := int(math.Pow(255, float64(index)))
	// 	total += int(b) * feat
	// }
	// si.CurrSecLength = total
	// si.readed = 0
	// return total, nil
}

func (si *SBodyIN) GetSec() ([]byte, error) {
	if si.readed < si.CurrSecLength {
		bs := make([]byte, si.CurrSecLength-si.readed)
		slog.Debug(fmt.Sprintln("Get old:", len(bs)))
		readed, err := si.raw.Read(bs)
		if err != nil {
			return nil, err
		}
		si.readed += readed
		return bs, nil
	} else {
		slog.Debug("Get new")
		length, err := si.Next()
		if err != nil {
			return nil, err
		}
		temp := make([]byte, length)
		readlength, err := si.Read(temp)
		if err != nil {
			return nil, err
		}
		if readlength != length {
			slog.Error(fmt.Sprintf("Wrong length: read:%d - next:%d\n", readlength, length))
		}
		return temp, nil
	}
}

func (si *SBodyIN) Read(buf []byte) (int, error) {
	if si.CurrSecLength == si.readed {
		return 0, errors.New("no more bytes for current section")
	}
	if len(buf) <= si.CurrSecLength-si.readed {
		i, err := si.raw.Read(buf)
		if err != nil {
			return 0, err
		}
		si.readed += i
		return i, err
	} else {
		temp := make([]byte, si.CurrSecLength-si.readed)
		length, err := si.raw.Read(temp)
		if err != nil {
			return 0, err
		}
		si.readed += length
		copy(buf[:length], temp)
		return length, nil
	}
}
