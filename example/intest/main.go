package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/ogios/sutils"
)

func main() {
	ln, _ := net.Listen("tcp", ":15002")
	conn, _ := ln.Accept()

	si := sutils.NewSBodyIn(bufio.NewReader(conn))

	// next
	length, err := si.Next()
	if err != nil {
		panic(err)
	}
	fmt.Printf("length: %d\n", length)
}
