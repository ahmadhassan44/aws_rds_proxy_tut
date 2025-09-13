package patterns

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	Web   = fakeSearch("web")
	Image = fakeSearch("image")
	Video = fakeSearch("video")
)

type Result struct {
	kind   string
	query  string
	result string
}
type Search func(query string) Result

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		return Result{
			kind:   kind,
			query:  query,
			result: fmt.Sprintf("%s result for %s", kind, query),
		}
	}
}
