package patterns

import "time"

func IDontWaitLongForBoringStuff(ch1, ch2 <-chan string, timeout time.Duration) <-chan string {
	ch := make(chan string)
	go func() {
		for {
			select {
			case msg := <-ch1:
				ch <- msg
			case msg := <-ch2:
				ch <- msg
			case <-time.After(timeout):
				ch <- "Both of them are too boring"
			}
		}
	}()
	return ch
}
