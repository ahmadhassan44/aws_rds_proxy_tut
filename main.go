package main

import (
	"fmt"
	"time"

	"github.com/ahmadhassan44/aws_rds_proxy_tut/patterns"
)

func main() {

	ch := patterns.IDontWaitLongForBoringStuff(patterns.RandomBoringStuff("Worker 1"), patterns.RandomBoringStuff("Worker 2"), 600*time.Millisecond)
	fmt.Println("Im listening")
	for i := 0; ; i++ {
		msg := <-ch
		if msg == "Both of them are too boring" {
			fmt.Println("Both of them are too boring")
			break
		}
		fmt.Println(msg)
	}
	// fmt.Println("You are boring! I am leaving")
}
