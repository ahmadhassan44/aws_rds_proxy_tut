package main

import (
	"fmt"
)

func main() {

	ch := randomBoringStuff("Worker 1")
	fmt.Println("Im listening")
	for i := 0; ; i++ {
		fmt.Println(<-ch)
	}
	fmt.Println("You are boring! I am leaving")
}
