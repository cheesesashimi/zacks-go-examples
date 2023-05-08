package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/cheesesashimi/zacks-go-examples/utils"
)

// A Goroutine is a lightweight thread managed by the Go runtime. The Go
// runtime basically schedules Goroutines across multiple OS threads, reusing
// threads wherever possible.
//
// However, to better understand contexts, one must first understand Goroutines
// and channels.

func spawning() {
	// Starting a Goroutine is very simple:
	go utils.NamedSleep("spawning", 500*time.Millisecond)

	// You can also start them as a closure. Although care must be taken with
	// respect to variable scoping to ensure that multiple Goroutines don't try
	// to read or write the same value simultaneously.
	go func() {
		fmt.Println("hello from another Goroutine! ID:", utils.GetGoroutineID())
	}()

	// If you compile and run this program as-is, neither of the above print
	// statements will appear to execute. The reason is because the main
	// Goroutine does not wait for either of these Goroutines to finish
	// executing. Furthermore, these Goroutines have no way to block the main
	// Goroutine until they're finished. This leads to something called a "race
	// condition" that we'll cover later.

	// It would be great if we only blocked the main Goroutine for as long as it
	// takes the child Goroutines to execute. While this is certainly possible
	// using a more advanced technique that we'll cover later, we'll use a bit of
	// a hack to achieve that goal for now:

	// Uncommenting this line will cause the main Goroutine to sleep, providing
	// enough time for the above Goroutines to print their output:
	// time.Sleep(time.Millisecond * 600)
}

// Another technique that one can do to wait on a Goroutine is to use the
// WaitGroup API.
func waitingWithAWaitgroup(name string) {
	wg := sync.WaitGroup{}

	for i := 0; i <= 10; i++ {
		// In Go, the go loop value is actually a pointer to a singular memory
		// address. Under the hood, the Go runtime changes the value of this
		// pointer without changing the actual memory address itself. In
		// single-threaded Go programs, this is usually not a problem (unless
		// you're trying to build a slice of pointers, that is). However, since
		// we're trying to spawn multiple Goroutines in this case, it is highly
		// probable that multiple Goroutines will try to read from the same memory
		// location at once leading to unintended behavior.
		//
		// What this syntax does is it makes a local copy of the value to make a
		// local copy of the value to use. Try commenting out this line and
		// uncommenting the corresponding fmt.Printf statements to see what happens
		// :).
		// fmt.Printf("%p\n", &i)
		i := i
		// fmt.Printf("%p\n", &i)

		// For each Goroutine you want to wait on, increment the WaitGroup.
		wg.Add(1)

		go func() {
			// Upon completion, the spawned Goroutine should mark that it is done.
			// This is best done with a defer statement. Defers are run in LIFO
			// (last-in, first-out) order before a function returns.
			defer wg.Done()
			utils.NamedSleep(fmt.Sprintf("%s-%d", name, i), 500*time.Millisecond)
		}()
	}

	// Once all of the Goroutines are started, we need to wait for them to finish
	// executing. So we wait!
	wg.Wait()
	fmt.Println("done waiting")
}

// A Goroutine can have child Goroutines. However, it's considered good
// practice for the parent Goroutine to ensure that its children have shut
// down. Here, we'll just call the waitingWithAWaitgroup function in 10
// separate Goroutines with our own Waitgroup.
func childGoroutines() {
	wg := sync.WaitGroup{}

	for i := 0; i <= 10; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			name := fmt.Sprintf("child-%d", i)
			utils.TimeIt(name, func() {
				waitingWithAWaitgroup(name)
			})
		}()
	}

	wg.Wait()
}

func main() {
	// These are ordered inversely for a reason.
	waitingWithAWaitgroup("waiting-with-waitgroup")
	childGoroutines()
	spawning()
}
