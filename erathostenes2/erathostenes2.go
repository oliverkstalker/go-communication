package main

import "fmt"

const N = 100

func odds(out chan<- int, feedback <-chan int, done chan<- bool) {
	fmt.Println(2) // Print the first prime, 2
	for i := 3; ; i += 2 {
		select {
		case out <- i: // Send the next odd number
		case lastPrime := <-feedback: // Terminate when the last sieve element sends its prime
			fmt.Println(lastPrime) // Print the last prime number
			close(out)             // Close the odd number stream
			done <- true           // Notify the main program
			return
		}
	}
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
	feedback := make(chan int)
	done := make(chan bool)

    fmt.Println("The first", N, "prime numbers are:");

    // Connect/start goroutines
	go odds(channels[0], feedback, done)
	
	for i := 0; i < N-1; i++ {
		go sieve(channels[i], channels[i+1])
	}

    // Await termination
	go func () {
		for lastPrime := range channels[N-1] {
			feedback <- lastPrime
		}
		close(feedback)
	}()

	// wait for termination signal
	<-done

    fmt.Println("Done!")
}