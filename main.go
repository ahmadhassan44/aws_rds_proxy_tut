package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ahmadhassan44/aws_rds_proxy_tut/patterns"
)

func main() {
	ch := make(chan patterns.Result)
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	go func() { ch <- patterns.Web("golang") }()
	go func() { ch <- patterns.Image("golang") }()
	go func() { ch <- patterns.Video("golang") }()
	for i := 0; i < 3; i++ {
		fmt.Println(<-ch)
	}
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}
