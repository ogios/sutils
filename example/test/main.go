package main

import "fmt"

func a() (text any) {
	defer func() {
		fmt.Println(text)
		text = 2
	}()
	return 1
}

func main() {
	b := a()
	fmt.Println(b)
}
