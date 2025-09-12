package main

import (
	"fmt"

	"github.com/ahmadhassan44/aws_rds_proxy_tut/patterns"
)

func main() {

	ch := patterns.FanIn(patterns.RandomBoringStuff("Worker 1"), patterns.RandomBoringStuff("Worker 2"))
	fmt.Println("Im listening")
	for i := 0; ; i++ {
		fmt.Println(<-ch)
	}
	// fmt.Println("You are boring! I am leaving")
}
