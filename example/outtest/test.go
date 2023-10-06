package main

import (
	"fmt"
	"net"
	"os"

	"github.com/ogios/sutils"
)

var PATH string = "/home/ogios/work/go/sutils/example/outtest/read.txt"

func test() {
	f, err := os.OpenFile(PATH, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	stat, err := f.Stat()
	if err != nil {
		panic(err)
	}

	fmt.Println("readed")

	out := sutils.NewSBodyOUT()
	length := fmt.Sprintf("File length: %d", stat.Size())
	out.AddBytes([]byte(length))
	out.AddReader(f, int(stat.Size()))

	fmt.Println("data added")
	fmt.Println(out.Raw...)
	fmt.Println(out.Types)

	conn, err := net.Dial("tcp", "127.0.0.1:15002")
	if err != nil {
		panic(err)
	}
	err = out.WriteTo(conn)
	if err != nil {
		panic(err)
	}

	fmt.Println("write done")
	f.Close()
	err = conn.Close()
	if err != nil {
		panic(err)
	}
}
