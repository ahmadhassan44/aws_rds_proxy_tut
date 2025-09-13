package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ahmadhassan44/aws_rds_proxy_tut/patterns"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := patterns.Web("golang")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}
