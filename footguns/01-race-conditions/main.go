package main

import (
	"fmt"
	"sync"
)

// This is an example of what *not* to do so that you can understand what a
// race condition is and how subtle they can be.
func raceConditions() int {
	// We declare this variable within this functions' scope.
	finalValue := 0

	// Set up our WaitGroups.
	wg := sync.WaitGroup{}

	for i := 0; i <= 10; i++ {
		i := i
		wg.Add(1)
		// Spawn all of our Goroutines
		go func() {
			defer wg.Done()
			// See the bug here? Basically, all of the Goroutines we've spawned try
			// to access this variable at the same time. So the final value can
			// change from run to run.
			finalValue += i
		}()
	}

	wg.Wait()

	return finalValue
}

// This is an example of how to avoid race conditions using a concurrency
// concept called mutxes. A mutex is short for mutual exclusion. There a couple
// different types of mutexes. However, those are beyond the scope of this
// lesson.
func mutexes() int {
	// We declare this variable within this functions' scope.
	finalValue := 0

	// Set up our WaitGroups.
	wg := sync.WaitGroup{}

	// Create our mutex object. If you pass this into another function, be sure
	// to pass a pointer to this otherwise it will pass a copy of this object
	// which is definitely not what you want to do.
	mux := sync.Mutex{}

	for i := 0; i <= 10; i++ {
		i := i
		wg.Add(1)
		// Spawn all of our Goroutines
		go func() {
			defer wg.Done()
			// Acquire our mutex lock. Only a single Goroutine can hold a lock at any
			// given time. The other Goroutines will block until it is their turn.
			mux.Lock()
			// Immediately defer unlocking the mutex when we're done.
			defer mux.Unlock()
			// Now we can safely modify the finalValue variable here.
			finalValue += i
		}()
	}

	wg.Wait()

	return finalValue
}

func runRaceConditionsAndMutexes() {
	resultMap := map[int]struct{}{}
	runs := 0
	for {
		if len(resultMap) == 10 {
			break
		}

		result := raceConditions()
		if _, ok := resultMap[result]; !ok {
			resultMap[result] = struct{}{}
		}

		runs++
	}

	results := []int{}
	for result := range resultMap {
		results = append(results, result)
	}

	fmt.Printf("No mutexes: %v (took %d runs)\n", results, runs)

	results = []int{}
	for i := 0; i <= 10; i++ {
		results = append(results, mutexes())
	}

	fmt.Printf("With mutexes: %v\n", results)
}

func main() {
	runRaceConditionsAndMutexes()
}
