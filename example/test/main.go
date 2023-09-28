package main

import (
	"bufio"
	"fmt"
	"net"

	st "github.com/ogios/sutils"
)

func main() {
	ln, _ := net.Listen("tcp", ":15002")
	conn, _ := ln.Accept()

	si := st.NewSBodyIn(bufio.NewReader(conn))

	// next
	length, err := si.Next()
	if err != nil {
		panic(err)
	}
	fmt.Printf("length: %d\n", length)

	// get half of it
	bs := make([]byte, int(length/2))
	length, err = si.Read(bs)
	if err != nil {
		panic(err)
	}

	// get the rest
	app, err := si.GetSec()
	if err != nil {
		panic(err)
	}

	// print
	final := make([]byte, length)
	final = append(final, bs...)
	final = append(final, app...)
	fmt.Println(string(final))
}
