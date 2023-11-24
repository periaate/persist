package persist

import (
	"fmt"
	"testing"
	"time"
)

func TestHash(t *testing.T) {
	hfn := NewHashU64()
	c := uint64(100_000)
	ti := StartTime()
	for key := uint64(0); key < c; key++ {
		_ = hfn(key)
	}
	ti.Count("", int(c))
}
func TestXxHash(t *testing.T) {
	c := uint64(100_000)
	ti := StartTime()
	for key := uint64(0); key < c; key++ {
		_ = HashString("HELLO, WORLD")
	}
	ti.Count("", int(c))
}

type Timer struct {
	TotalDuration time.Duration
	LastTime      time.Time
}

func (t *Timer) Count(s string, i int) {
	currentTime := time.Now()

	if t.LastTime.IsZero() {
		t.LastTime = currentTime
	}

	elapsedTime := currentTime.Sub(t.LastTime)
	t.TotalDuration += elapsedTime
	avgDuration := t.TotalDuration / time.Duration(i+1)

	fmt.Printf("%s — Total: %v — avg/item: %v\n", s, t.TotalDuration, avgDuration)

	t.LastTime = currentTime
}

func (t *Timer) Start() {
	t.LastTime = time.Now()
	t.TotalDuration = 0
}

func StartTime() *Timer {
	t := Timer{}
	t.Start()
	return &t
}
