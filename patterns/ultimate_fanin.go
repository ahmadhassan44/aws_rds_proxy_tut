package patterns

func UltimateFanIn(ch1, ch2 <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			select {
			case msg := <-ch1:
				c <- msg
			case msg := <-ch2:
				c <- msg
			}
		}
	}()
	return c
}
