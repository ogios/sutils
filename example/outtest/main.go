package main

import (
	"fmt"

	"github.com/ogios/sutils"
)

func gen() []byte {
	a := make([]byte, 0)
	for i := 0; i < cap(a); i++ {
		a[i] = byte(i)
	}
	return a
}

func main() {
	// test()

	out := sutils.NewSBodyOUT()
	out.AddBytes(gen())
	fmt.Println(out.Raw[0])
}
