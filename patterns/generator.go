package patterns

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomBoringStuff(msg string) <-chan string {
	var words = []string{
		"apple", "banana", "cherry", "dog", "elephant",
		"forest", "guitar", "house", "island", "jungle",
		"king", "lemon", "mountain", "night", "ocean",
		"piano", "queen", "river", "sun", "tree",
		"umbrella", "violin", "wolf", "xylophone", "zebra",
	}
	ch := make(chan string)
	go func() {
		for i := 0; ; i++ {
			duration := time.Duration(rand.Intn(1e3)) * time.Millisecond
			time.Sleep(duration)
			ch <- fmt.Sprintf("Boring %s slept for: %v and said %s", msg, duration, words[rand.Intn(len(words))])
		}
	}()
	return ch
}
