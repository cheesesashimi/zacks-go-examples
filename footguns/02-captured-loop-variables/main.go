package main

import "fmt"

// In Go, the go loop value is actually a pointer to a singular memory address.
// Under the hood, the Go runtime changes the value of this pointer without
// changing the actual memory address itself. This can cause subtle problems
// when copying values to a new slice or passing into another function.

// In this example, we do not "capture" the loop variable.
func notCapturingTheLoopVariable(items []string) {
	// We have an array of pointers to strings.
	copiedItems := []*string{}

	for i, item := range items {
		// For each string, output what the string value is and what its memory address is.
		fmt.Printf("index: %d\tvalue: %s\titerator memory location: %p\tlist memory location: %p\n", i, item, &item, &items[i])

		// For each string, copy the pointer into our copied items array.
		copiedItems = append(copiedItems, &item)
	}

	// Notice how all of the copiedItems slice is just a pointer to the same
	// memory location
	printCopiedItems(copiedItems)

	// What happens if we mutate one of the items in our copied items list?
	*copiedItems[0] = "hello"

	// Notice how all of the items in our copiedItems slice changed?
	printCopiedItems(copiedItems)
}

// In this example, we capture the loop variable within a local scope.
func capturingTheLoopVariable(items []string) {
	copiedItems := []*string{}

	for i, item := range items {
		// We make a local copy of item here, which allows us to capture the loop
		// variable. This allocates additional memory and copies the value of the
		// iterator to our local variable. We are free to do anything we want with
		// our copy of the value since we will only mutate our copy of it.
		//
		// For simplicity and to be explicit about what's happening, I called the
		// captured loop variable "capturedItem". In practice though, you can do
		// something like item := item because of Go variable scoping rules.
		capturedItem := item

		// For each string, output what the string value is and what its memory address is.
		fmt.Printf("index: %d\tvalue: %s\titerator memory location: %p\tlist memory location: %p\n", i, capturedItem, &capturedItem, &items[i])

		// For each string, copy the pointer into our copied items array.
		copiedItems = append(copiedItems, &capturedItem)
	}

	printCopiedItems(copiedItems)
}

// Sometimes, one may want to look an item in a slice up via its index instead
// of using the Go iteration pattern:
func explicitSliceAccess(items []string) {
	copiedItems := []*string{}

	for i := range items {
		// Instead of copying the value from the slice iteration, we can look up
		// the value by using the slice index. While this works, it can have
		// unintended consequences as we'll see shortly.
		fmt.Printf("index: %d\tvalue: %s\tlist memory location: %p\n", i, items[i], &items[i])

		copiedItems = append(copiedItems, &items[i])
	}

	printCopiedItems(copiedItems)

	modifyCopiedItems(copiedItems, items)

	// Uh-oh! We've mutated our original list! How!?! Why!?! This happens because
	// in Go, the underlying array for slices is passed by reference, not value
	// when it is passed into another function. With this in mind, it is possible
	// to accidentally mutate the underlying slice or map, as we did here.
}

func explicitSliceAccessWithCopy(items []string) {
	copiedItems := []*string{}

	for i := range items {
		// We manually look up the value of items.
		item := items[i]

		fmt.Printf("index: %d\tvalue: %s\tlist memory location: %p\n", i, item, &item)

		copiedItems = append(copiedItems, &item)
	}

	modifyCopiedItems(copiedItems, items)

	// In this case, we did not mutate the underlying slice because we made a
	// copy of our value and then assigned the memory location of that to our
	// copiedItems list. When we modified the copiedItems slice, the change only
	// occurred to our copy of the slice; not to the backing array.
}

func mutatingTheSlice(items []string) {
	for i, item := range items {
		// What happens when we modify the iterator value?
		item += "-hello"

		fmt.Printf("index: %d\tvalue: %s\titerator memory location: %p\tlist memory location: %p\n", i, item, &item, &items[i])
	}

	// Nothing! Because the value that the iterator value points to is a copy of what is in our slice.
	isMutated(items)

	// But what happens when we do something like this?
	for i := range items {
		items[i] += "-hello"
		item := items[i]

		fmt.Printf("index: %d\tvalue: %s\titerator memory location: %p\tlist memory location: %p\n", i, item, &item, &items[i])
	}

	// The underlying slice is mutated in this particular case.
	isMutated(items)
}

func isMutated(items []string) bool {
	mutated := items[0] != "one"
	fmt.Printf("Is mutated? %v\n", mutated)
	return mutated
}

// Modifies the copied items and prints value from the original items slice to see if it changed.
func modifyCopiedItems(copiedItems []*string, originalItems []string) {
	i := 0

	// What happens when I change the value in our copied items?
	*copiedItems[i] = "hello"

	// Lets print it out here.
	printCopiedItems(copiedItems)

	// But what about our source slice?
	fmt.Printf("value: %s\tlist memory location: %p\tunderlying slice changed? %v\n", originalItems[i], &originalItems[i], originalItems[i] == *copiedItems[i])
}

// Prints the copied items slice.
func printCopiedItems(copiedItems []*string) {
	fmt.Printf("Copied items: %v\n", copiedItems)

	for i, item := range copiedItems {
		// Note: The list memory location ends up as a pointer to a pointer.
		fmt.Printf("index: %d\tvalue: %s\titerator memory location: %p\tlist memory location: %p\t\n", i, *item, item, &copiedItems[i])
	}
}

// Resets the provided slice back to its defaults.
func resetItemsValues(items []string) {
	items[0] = "one"
	items[1] = "two"
	items[2] = "three"
}

func main() {
	items := []string{"one", "two", "three"}

	fmt.Println("Not capturing the loop variable:")
	notCapturingTheLoopVariable(items)
	fmt.Println("")

	fmt.Println("Capturing the loop variable:")
	capturingTheLoopVariable(items)
	fmt.Println("")

	fmt.Println("Explicit slice access:")
	explicitSliceAccess(items)
	// Did explicitSliceAccess() unintentionally modify our items slice?
	// If so, set it back to the original value.
	isMutated(items)
	// Spoiler: It is mutated! So we set it back to the expected value.
	resetItemsValues(items)
	fmt.Println("")

	fmt.Println("Explicit slice access with copy:")
	explicitSliceAccessWithCopy(items)
	// Did explicitSliceAccessWithCopy() unintentionally modify our items slice?
	// If so, set it back to the original value.
	isMutated(items)
	// Spoiler: It is not mutated!
	fmt.Println("")

	fmt.Println("Mutating the slice:")
	mutatingTheSlice(items)
	// We call this again here to show that the original slice is indeed mutated.
	isMutated(items)
	resetItemsValues(items)
	fmt.Println("")
}
