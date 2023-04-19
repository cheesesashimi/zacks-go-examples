package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cheesesashimi/zacks-go-examples/utils"
)

// So far, we've learned about channels as a way to communicate between
// separate Goroutines. Now it's time to introduce contexts.
//
// From the official Go docs: Package context defines the Context type, which
// carries deadlines, cancellation signals, and other request-scoped values
// across API boundaries and between processes.
//
// It is my opinion that the Go docs don't adequately communicate *why* this is
// useful. So I'll try to give a couple of scenarios where one might want to
// use contexts and how to use them.

// This function will start a Goroutine that runs infinitely until a
// cancellation signal is sent. It's worth mentioning that the Go convention is
// that contexts should be explicitly passed into functions (and should never
// be embedded in structs) and should always be the first argument to the
// function. I think there is a bit of leeway around functions which are
// methods on a struct, but I don't want to get into a flamewar right now :).
func startLongRunningProcess(ctx context.Context, name string) chan struct{} {
	// We return a channel which indicates that we've finished. This is because
	// there is a bit of a delay between the time the context is cancelled and
	// we've completed our shutdown. In this case, shutdown involves computing
	// and outputting the elapsed time of our function. This is not strictly
	// necessary, but we do it in this case because otherwise our output looks
	// funny.
	doneChan := make(chan struct{})

	go func() {
		utils.TimeIt(name, func() {
			fmt.Println(name, "started")
			for {
				select {
				// ctx.Done() returns a channel that we can read from to determine if the
				// context has been canceled. Any signal sent to this channel is
				// propagated to all calls of ctx.Done(), meaning that cancellations can
				// be easily synchronized across multiple Goroutines.
				case <-ctx.Done():
					return
				}
			}
		})

		// Just closing a channel will also cause a signal to be sent via the
		// channel. We can block until we receive that signal. Once again, this is
		// not strictly necessary.
		close(doneChan)
	}()

	return doneChan
}

func simpleCancellation() {
	// In Go, contexts are descended from a global Context object. If you've used
	// context.TODO(), this returns the global background context. Unfortunately,
	// the global background context doesn't have any way to be canceled.
	// Instead, what we must do is make a child context from the main background
	// context and attach a cancelation function to it. This is pretty simple to
	// do:
	ctx, cancel := context.WithCancel(context.Background())

	// It is considered good practice to defer cancelation of a context to make
	// sure that it actually does get cancelled if something unexpectedly goes
	// awry and the context is not otherwise canceled.
	defer cancel()

	doneChan := startLongRunningProcess(ctx, "simple-cancellation-1")
	delay := time.Millisecond * 100

	fmt.Println("Sending cancellation after", delay)
	time.Sleep(delay)
	cancel()

	// We can block until the context channel is closed.
	<-ctx.Done()

	// We also want to block until our long running process has finished shutting
	// down. As mentioned before, this is not strictly necessary.
	<-doneChan

	fmt.Println("Context cancelled!")
}

func simpleTimeout() {
	// We can wire up our timeout directly to the context. By doing this, we can
	// explicitly say that if the function does not return in t time, we should
	// automatically cancel its execution.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	doneChan := startLongRunningProcess(ctx, "timeout-cancellation-1")
	<-ctx.Done()

	// We also want to block until our long running process has finished shutting
	// down.
	<-doneChan

	fmt.Println("Context timeout reached!", ctx.Err())
}

// Deadlines are another way of representing the timeout paradigm. In this
// case, however, we specify a time in the future that we should have exited
// our function by.
func simpleDeadline() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Millisecond))
	defer cancel()

	doneChan := startLongRunningProcess(ctx, "deadline-cancellation-1")
	<-ctx.Done()
	<-doneChan
	fmt.Println("Context deadline reached!", ctx.Err())
}

// A single context may be shared across multiple Goroutines to synchronize
// their cancellation and shutdown. In this case, we share a single context
// across each Goroutine and then wait for them to be shut down.
func sharedContexts() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Because we want to wait for each of our Goroutines to shut down, we use a
	// WaitGroup. As mentioned before, this is not strictly necessary.
	wg := sync.WaitGroup{}

	for i := 1; i <= 10; i++ {
		wg.Add(1)

		i := i
		go func() {
			defer wg.Done()
			// We use this syntax to indicate that we want to block this Goroutine
			// until startLongRunningProcess closes the channel it returns.
			<-startLongRunningProcess(ctx, fmt.Sprintf("shared-context-%d", i))
		}()
	}

	<-ctx.Done()
	wg.Wait()
}

// Contexts can also have child contexts which inherit the parents cancellation
// signal. However, child contexts will not propagate their cancellation signal
// up to the parent context. This can be a useful pattern for when you want to
// have a shorter timeout for certain things.
func childContexts(parentTimeout time.Duration) {
	fmt.Printf("Parent context has timeout %s\n", parentTimeout)
	parentCtx, parentCancel := context.WithTimeout(context.Background(), parentTimeout)
	defer parentCancel()

	wg := sync.WaitGroup{}

	for i := 1; i <= 10; i++ {
		wg.Add(1)

		i := i
		go func() {
			defer wg.Done()
			// Create a child context from our parent context with an incrementally-increasing timeout.
			childTimeout := time.Millisecond * time.Duration(i*5)
			childCtx, childCancel := context.WithTimeout(parentCtx, childTimeout)
			defer childCancel()

			<-startLongRunningProcess(childCtx, fmt.Sprintf("child-context-%d", i))
		}()
	}

	<-parentCtx.Done()
	wg.Wait()
}

func runChildContexts() {
	// In this case, the children will reach their own individual timeout before
	// the parent timeout is reached.
	childContexts(time.Millisecond * 100)

	// In this case however, the parent timeout will be reached before the
	// children. All of the children will be shut down in response.
	childContexts(time.Millisecond)
}

func main() {
	simpleCancellation()
	simpleTimeout()
	simpleDeadline()
	sharedContexts()
	runChildContexts()
}
