package utils

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Accepts a name and a function and times how long it takes the function to run.
func TimeIt(name string, f func()) {
	start := time.Now()
	defer func() {
		fmt.Printf("%s finished running in Goroutine %d in %s\n", name, GetGoroutineID(), time.Since(start))
	}()

	f()
}

// Sleeps for a given time and outputs its name upon completion.
func NamedSleep(name string, d time.Duration) {
	TimeIt(name, func() {
		time.Sleep(d)
	})
}

// Simple function that sums a slice of ints.
func Sum(nums []int) int {
	sum := 0

	for _, num := range nums {
		sum += num
	}

	return sum
}

func GenerateRandomNumber(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func GenerateRandomNumbers(min, max, n int) []int {
	nums := make([]int, n)

	for i := 0; i < n; i++ {
		nums[i] = GenerateRandomNumber(min, max)
	}

	return nums
}

// Copied from https://gist.github.com/metafeather/3615b23097836bc36579100dac376906.
// Generally, this isn't a good idea. See: https://go.dev/doc/faq#no_goroutine_id.
// However, we'll do it here for purely educational purposes.
func GetGoroutineID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
