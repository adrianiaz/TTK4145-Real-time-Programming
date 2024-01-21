package main

import (
	"fmt"
)

func main() {
	fmt.Println("booting server")
	go server()
	go client()
	select {}
}
