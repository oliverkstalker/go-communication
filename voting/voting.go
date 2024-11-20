package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Votes struct {
	a, b int
}

// Simulate a polling station sending random votes
func station(out chan<- Votes) {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
		aVotes := rand.Intn(100)
		out <- Votes{aVotes, 100 - aVotes}
	}
	close(out)
}

// Collector node aggregates votes from two input channels and sends to an output channel
func collector(in1, in2 <-chan Votes, out chan<- Votes) {
	for {
		select {
		case v1, ok1 := <-in1:
			if ok1 {
				out <- v1
			}
		case v2, ok2 := <-in2:
			if ok2 {
				out <- v2
			}
		}

		// If both input channels are closed, close the output channel and exit
		if in1 == nil && in2 == nil {
			close(out)
			return
		}
	}
}


func main() {
	rand.Seed(time.Now().UnixNano())

	// Number of polling stations
	numStations := 8

	// Channels for communication
	votesChannels := make([]chan Votes, numStations+7) // 8 stations + 7 collectors
	for i := range votesChannels {
		votesChannels[i] = make(chan Votes)
	}

	// Launch polling stations
	for i := 0; i < numStations; i++ {
		go station(votesChannels[i])
	}

	// Launch collector nodes
	for i := 0; i < 7; i++ {
		go collector(votesChannels[2*i], votesChannels[2*i+1], votesChannels[numStations+i])
	}

	// Final tally at the root collector
	tally := Votes{}
	go func() {
		for v := range votesChannels[len(votesChannels)-1] {
			tally.a += v.a
			tally.b += v.b
			fmt.Printf("Current tally: A=%d, B=%d\n", tally.a, tally.b)
		}
	}()

	// Wait for the final result
	time.Sleep(5 * time.Second)
	tot := tally.a + tally.b

	fmt.Printf("All votes counted. Total votes: %d\n", tot)
	var winner string
	switch {
	case tally.a > tally.b:
		winner = "A"
	case tally.a < tally.b:
		winner = "B"
	default:
		winner = "undetermined"
	}
	fmt.Printf("The winner is: %s\n", winner)
	if winner == "B" {
		fmt.Println("A: This must be FRAUD!!!")
	}
}
