package main

func test([]byte) {
}

func main() {
	a := new([8]byte)
	b := make([]any, 1)
	b[0] = a

	test(b[0].([8]byte))
}
