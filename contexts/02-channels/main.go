package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/cheesesashimi/zacks-go-examples/utils"
)

// Channels are a typed conduit through which you can send and receive values
// with the channel operator, <-. With channels, it's possible to communicate
// between Goroutines instead of sharing memory, without using traditional
// concurrency primitives such as mutexes, etc.

// If you'll recall from the Goroutines example, we'd like to get some kind of
// signal that the child Goroutine is completed. Channels are perfect for this
// because we can block a Goroutine waiting on a value to be sent through a
// channel -or- for the channel be closed.
func waitingWithAChannel() {
	// First, we need to create a channel using the make() function. There is no
	// literal syntax for creating a channel.
	doneChan := make(chan struct{})

	// We spawn our Goroutine as a closure so we can reference the channel. An
	// alternative is that one can pass in the channel as we'll se later.
	go func() {
		utils.NamedSleep("waiting-with-a-channel", 500*time.Millisecond)
		// We close the channel when we're done sleeping.
		close(doneChan)
	}()

	// This syntax does two things:
	// 1. Blocks the current Goroutine until something is sent over the channel.
	// 2. When paired with an assignment operator (e.g., := or =), it will assign
	// the value sent through the channel to the awaiting variable.
	<-doneChan
}

// In this case, we want to do more than just block the main Goroutine while we
// wait for it to finish.
func sendValueOverChannel() {
	sumChan := make(chan int)
	defer close(sumChan)

	go func() {
		// In this case, we execute our sum
		sumChan <- utils.Sum(utils.GenerateRandomNumbers(0, 100, 100))
	}()

	// This syntax will block the current Goroutine until a value is published
	// over the channel.
	sum := <-sumChan
	fmt.Println("Sum:", sum)
}

// There is a lot happening within this function, so lets take it step by step:
func iteratingOverChannels() {
	// Create a channel to report the summed values over. This channel serves
	// two purposes:
	// 1. Allows us to block the main Goroutine until the summation is complete.
	// 2. Allows us to get the summed value without having to share memory
	// between the main Goroutine and the Goroutines we start.
	sumChan := make(chan int)

	// This channel is where we send our random numbers to be summed.
	numChan := make(chan int)

	// Start a Goroutine that generates random numbers and sends them iteratively
	// to the number channel. It will close the channel when it is finished.
	go func() {
		n := 100

		for i := 0; i <= n; i++ {
			numChan <- utils.GenerateRandomNumber(0, 100)
		}

		// Close our channel when we've generated all of our numbers.
		close(numChan)
	}()

	// Start a Goroutine that consumes numbers from the number channel and sums
	// them up. Upon completion, it will send them over the sum channel.
	go func() {
		sum := 0

		// We range over the channel until the channel is closed. The close signal
		// becomes the loop exit condition.
		for num := range numChan {
			sum += num
		}

		sumChan <- sum
		close(sumChan)
	}()

	// We block until the sumChan has a value written to it and it is closed.
	// This approach has two distinct advantages:
	// 1. We've done this without mutexes (well, ones that we have to manage,
	// anyway).
	// 2. We don't need to allocate memory for a full list of random numbers.
	// It's very similar to using generators (yield keyword) in Python with the
	// major difference being that the producer and consumer might not be running
	// in the same thread.
	fmt.Println(<-sumChan)
}

// It is possible for multiple Goroutines to read from a single channel. It
// should be mentioned that the first Goroutine available gets the value from
// the channel. The value is *not* broadcast to all Goroutines listening on
// that channel. It is possible to write code which does that, but that's out
// of scope for this lesson.
func multipleGoroutinesReadingAndWritingToTheSameChannel() {
	// This is the channel that all of the Goroutines listen on.
	numChan := make(chan int)

	// This producer function will produce 100 random numbers and send them over
	// the common channel.
	producerFunc := func() {
		for i := 0; i < 10; i++ {
			num := utils.GenerateRandomNumber(0, 100)
			fmt.Printf("sent %d from producer Goroutine %d\n", num, utils.GetGoroutineID())
			numChan <- num
		}
	}

	// This consumer function is executed within each Goroutine we start. It
	// consumes numbers from the common channel and sums them up.
	consumerFunc := func() {
		id := utils.GetGoroutineID()
		sum := 0
		for num := range numChan {
			fmt.Printf("received %d in consumer Goroutine %d\n", num, id)
			sum += num
		}
		fmt.Printf("Goroutine %d finished with sum: %d\n", id, sum)
	}

	// Start our producer Goroutines which generate random numbers.
	// While it is possible to use channels to determine when our Goroutines are
	// finished, it can get pretty complicated. Instead, we'll use a WaitGroup
	// here for simplicity.
	producerWaitGroup := sync.WaitGroup{}
	for i := 1; i <= 5; i++ {
		producerWaitGroup.Add(1)
		go func() {
			defer producerWaitGroup.Done()
			producerFunc()
		}()
	}

	// Start our consumer Goroutines that consume the random numbers.
	// We use a separate WaitGroup for our consumer Goroutines since they'll shut
	// down when the number channel is closed.
	consumerWaitGroup := sync.WaitGroup{}
	for i := 1; i <= 5; i++ {
		consumerWaitGroup.Add(1)
		go func() {
			defer consumerWaitGroup.Done()
			consumerFunc()
		}()
	}

	// Wait for all of the producer functions to complete executing.
	producerWaitGroup.Wait()

	// Close our channel. This will cause the consumer Goroutines to shut down.
	close(numChan)

	// Wait for our consumer Goroutines to finish.
	consumerWaitGroup.Wait()
}

// Up until now, our channel reads block the current Goroutine until a value is
// available on them. This may be undesirable in certain situations. Instead,
// one can use the select {} syntax to read values from multiple channels
// simultaneously.
func nonBlockingChannelReads() {
	cumulativeChan := make(chan int)

	chan1 := make(chan int)
	chan2 := make(chan int)
	chan3 := make(chan int)

	// Start Goroutines to sum numbers and send the results to the provided
	// channel upon completion.
	go func() {
		chan1 <- utils.Sum(utils.GenerateRandomNumbers(0, 100, 100))
		fmt.Println("Goroutine 1 finished")
	}()

	go func() {
		chan2 <- utils.Sum(utils.GenerateRandomNumbers(0, 100, 100))
		fmt.Println("Goroutine 2 finished")
	}()

	go func() {
		chan3 <- utils.Sum(utils.GenerateRandomNumbers(0, 100, 100))
		fmt.Println("Goroutine 3 finished")
	}()

	go func() {
		// Because our channels will only return a single value, we must keep track
		// of which channel has been read so we know whether we can terminate our
		// loop. If we don't do this, we will wait indefinitely.
		chan1Finished := false
		chan2Finished := false
		chan3Finished := false

		// We need a loop so that we can keep reading from all of our channels.
		for {
			// Using the select syntax, we can listen on any number of channels. This
			// is very similar to the switch {} syntax with a few exceptions.
			select {
			// Using this syntax, we can figure out if we should expect more values
			// to be reported over the channel.
			case chan1Value, isFinished := <-chan1:
				fmt.Println("chan 1 value:", chan1Value)
				chan1Finished = isFinished
				cumulativeChan <- chan1Value
			case chan2Value, isFinished := <-chan2:
				fmt.Println("chan 2 value:", chan2Value)
				chan2Finished = isFinished
				cumulativeChan <- chan2Value
			case chan3Value, isFinished := <-chan3:
				fmt.Println("chan 3 value:", chan3Value)
				chan3Finished = isFinished
				cumulativeChan <- chan3Value
			default:
				if chan1Finished && chan2Finished && chan3Finished {
					close(cumulativeChan)
					fmt.Println("cumulative Goroutine finished")
					return
				}
			}
		}
	}()

	sum := 0
	for num := range cumulativeChan {
		sum += num
	}

	fmt.Println("Cumulative Value:", sum)
}

// In general, something to keep in mind when you start a Goroutine is how to
// stop it. Up until now, This becomes critical when dealing with error
// handling and multiple running Goroutines, etc. as you don't want to leak
// Goroutines or have deadlocks, etc.
func shuttingDownAGoroutine() {
	shutdownChan := make(chan struct{})

	go func() {
		delay := time.Millisecond * 10
		for {
			select {
			case <-shutdownChan:
				fmt.Println("Received shutdown signal")
				return
			default:
				time.Sleep(delay)
				fmt.Printf("Slept for %s\n", delay)
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)
	shutdownChan <- struct{}{}
	fmt.Println("Goroutine is now shut down")
}

func main() {
	rand.Seed(time.Now().Unix())

	waitingWithAChannel()
	sendValueOverChannel()
	iteratingOverChannels()
	multipleGoroutinesReadingAndWritingToTheSameChannel()
	nonBlockingChannelReads()
	shuttingDownAGoroutine()
}
