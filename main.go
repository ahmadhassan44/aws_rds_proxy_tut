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
	web := patterns.Web("golang")
	image := patterns.Image("golang")
	video := patterns.Video("golang")
	elapsed := time.Since(start)
	fmt.Println(web)
	fmt.Println(image)
	fmt.Println(video)
	fmt.Println(elapsed)
}
