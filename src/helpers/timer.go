package helpers

import (
	"fmt"
	"time"
)

type Timer struct {
	start int64
	prev  int64
	id    string
}

func StartTimer(id string) *Timer {
	t := &Timer{}
	t.id = id
	t.start = time.Now().UnixNano()
	t.prev = time.Now().UnixNano()
	return t
}

func (t *Timer) Check(str string) {
	now := time.Now().UnixNano()
	fmt.Printf("[+%vns] [%v] %v\n", (now - t.prev), t.id, str)
	t.prev = time.Now().UnixNano()
}

func (t *Timer) Finish(str string) {
	now := time.Now().UnixNano()
	fmt.Printf("[+%vns] [%v] %v\n", (now - t.start), t.id, str)
}
