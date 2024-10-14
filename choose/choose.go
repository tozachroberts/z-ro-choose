package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type State struct {
	choice    string
	committed bool
}

type Event struct {
	kind string
}

func emitEvents(e chan Event) {
	fmt.Println("ztest 1")
	for _, v := range [6]string{"left", "right", "click", "left", "click", "right"} {
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(3)
		log.Printf("event %s sleeping %d\n", v, n)
		time.Sleep(time.Duration(n) * time.Second)
		event := Event{v}
		e <- event
	}
	close(e)
}

func onEvent(e chan Event) {
	var s = State{"", false}
	for event := range e {
		switch kind := event.kind; kind {
		case "left":
			s.choice = "left"
			s.committed = false
		case "right":
			s.choice = "right"
			s.committed = false
		case "click":
			s.committed = true
		}
	}
	log.Printf("state: %+v\n", s)
}

func main() {
	log.SetPrefix("choose: ")
	log.SetFlags(0)

	e := make(chan Event, 1)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		emitEvents(e)
	}()
	go func() {
		defer wg.Done()
		onEvent(e)
	}()
	wg.Wait()
}
