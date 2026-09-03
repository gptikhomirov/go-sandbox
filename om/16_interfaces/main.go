package main

import "time"

type Film struct {
	Name     string
	Duration time.Duration
}

func (f *Film) Stop() {
	f.Duration = 0
}

type Stopable interface {
	Stop()
}

func stop(p Stopable) {
	p.Stop()
}

func main() {
	timer := time.NewTicker(1 * time.Second)

	stop(timer)
}
