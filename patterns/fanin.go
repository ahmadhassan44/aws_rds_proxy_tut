package patterns

func FanIn(ch1, ch2 <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			c <- <-ch1
		}
	}()
	go func() {
		for {
			c <- <-ch2
		}
	}()
	return c
}
