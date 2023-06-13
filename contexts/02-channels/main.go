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

	// This function accepts a channel as an argument for where to send values
	// to. It also accepts an ID, which is solely for identification reasons.
	sendRandomNumbersToChannel := func(destChan chan int, id int) {
		for _, num := range utils.GenerateRandomNumbers(0, 100, 100) {
			// time.Sleep(time.Millisecond * time.Duration(utils.GenerateRandomNumber(0, 100)))
			destChan <- num
		}
		close(destChan) // Comment this line out and see what happens :).
		fmt.Printf("Goroutine %d finished, channel %d closed\n", id, id)
	}

	// Start Goroutines to sum numbers and send the results to the provided
	// channel upon completion.
	go sendRandomNumbersToChannel(chan1, 1)
	go sendRandomNumbersToChannel(chan2, 2)
	go sendRandomNumbersToChannel(chan3, 3)

	go func() {
		// We must keep track of which channel has been read so we know whether we
		// can terminate our loop. If we don't do this, we will wait indefinitely.
		// Alternatively, if we don't close our channels, we'll deadlock. Try
		// commenting out the "close(destChan)" line in
		// sendRandomNumbersToChannel() and see what happens :).
		chan1Finished := false
		chan2Finished := false
		chan3Finished := false

		// Because we want to do the same thing for each channel, we call this
		// function with our channel ID, the value received, and whether it was OK.
		// We infer from the the ok value whether the channel has closed or not.
		handleChannelEvent := func(channelID, value int, ok bool) bool {
			finished := !ok
			fmt.Printf("chan %d value: %d, ok? %v, finished? %v\n", channelID, value, ok, finished)
			cumulativeChan <- value
			return finished
		}

		printChannelState := func() {
			fmt.Printf("chan 1 finished? %v, chan 2 finished? %v, chan 3 finished? %v\n", chan1Finished, chan2Finished, chan3Finished)
		}

		fmt.Printf("start: ")
		printChannelState()

		// We need a loop so that we can keep reading from all of our channels.
		for {
			// Using the select syntax, we can listen on any number of channels. This
			// is very similar to the switch {} syntax with a few exceptions.
			select {
			// Using this syntax, we can figure out if we should expect more values
			// to be reported over the channel. The ok values indicate whether the
			// channel is closed, which we use to infer whether we should break out
			// of our loop or not. See: https://go.dev/ref/spec#Receive_operator for
			// details.
			case val, ok := <-chan1:
				chan1Finished = handleChannelEvent(1, val, ok)
			case val, ok := <-chan2:
				chan2Finished = handleChannelEvent(2, val, ok)
			case val, ok := <-chan3:
				chan3Finished = handleChannelEvent(3, val, ok)
			default:
				fmt.Printf("default case: ")
				printChannelState()
			}

			// After we've processed our channel value, we need to determine whether
			// we've reached the required exit condition for our loop. If we've
			// reached it, we need to close our cumulative channel and return from
			// here. We need to guarantee that the exit condition is evaluated on
			// each loop iteration. Because of that requirement, we cannot use a
			// default case on the above select statement. Here's why:
			//
			// While we cannot read from a closed channel and will get a runtime
			// panic if we try to do so, the use of the receive operator means that
			// each statement in the select will still be evaluated even if the
			// channel is closed. The assignment will result in a nil or zero value
			// for the given type as well as a false boolean value. However, for the
			// purposes of the select case, that case will still evaluate to "true",
			// executing the code for that specific case and not falling through to
			// the default case.
			//
			// To see this in action, you can do the following:
			// $ go build -o channels && ./channels | grep "default case"
			if chan1Finished && chan2Finished && chan3Finished {
				close(cumulativeChan)
				fmt.Printf("cumulative Goroutine finished: ")
				printChannelState()
				return
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
