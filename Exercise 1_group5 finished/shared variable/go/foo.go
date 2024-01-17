// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
	//"time"
)

func server(c1 chan int, q1 chan int, q2 chan int) {
	// setup
	var i = 0
	//var temp int;
    var quitCounter int = 0
	for {
		select {

		case temp := <-c1:
			if temp == 1 {
				i++
			}
			if temp == -1 {
				i--
			}
		case <-q1:
			quitCounter++
		case <-q2:
			quitCounter++
		}
		if quitCounter == 2 {
			Println("The magic number is:", i)
			return
		}
	}

}

func incrementing(c chan int, q chan int) {
	//TODO: increment i 1000000 times
	for i := 0; i < 1000000; i++ {
		c <- 1
	}
	q <- 0
    return
}

func decrementing(c chan int, q chan int) {
	//TODO: decrement i 1000000 times

	for i := 0; i < 1000001; i++ {
		c <- -1
	}
	q <- 0
    return
}

func main() {
	// What does GOMAXPROCS do? What happens if you set it to 1? Only one thread will be active.
	runtime.GOMAXPROCS(3)
	ch1 := make(chan int)
	quit1 := make(chan int)
	quit2 := make(chan int)

	// TODO: Spawn both functions as goroutines
	go incrementing(ch1, quit1)
	go decrementing(ch1, quit2)
	go server(ch1, quit1, quit2)
	// We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
	// We will do it properly with channels soon. For now: Sleep.
	select {}
}
