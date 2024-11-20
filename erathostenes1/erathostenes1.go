package main

import "fmt"

const N = 100

func odds(out chan<- int) {
	fmt.Println(2) // print first prime
	for i := 3; i < 10*N; i += 2 {
		out<- i
	}
	close(out)
}

func sieve(in <-chan  int, out chan<- int) {
    prime, ok := <-in
	if !ok {
		close(out)
		return
	}
	fmt.Println(prime) // print prime
	for num := range in {
		if num % prime != 0 {
			out<- num
		}
	}
	close(out)
}

func main() {
    // Declare channels
	channels := make([]chan int, N)

	// Initialize channels 
	for i := range channels {
		channels[i] = make(chan int)
	}

    fmt.Println("The first", N, "prime numbers are:");

    // Connect/start goroutines
	go odds(channels[0])
	
	for i := 0; i < N-1; i++ {
		go sieve(channels[i], channels[i+1])
	}

	for range channels[N-1] {
	} // Do nothing
    // Await termination

    fmt.Println("Done!")
}