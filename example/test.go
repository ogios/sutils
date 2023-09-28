package example

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
	length, _ := si.Next()
	fmt.Printf("length: %d\n", length)
	bs := make([]byte, length)
	length, _ = si.Read(bs)
	fmt.Println(string(bs))
}
