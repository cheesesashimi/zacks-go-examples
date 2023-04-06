package main

import (
	"errors"
	"fmt"

	"github.com/cheesesashimi/zacks-go-examples/errors/utils"
)

func basicWrappingAndUnwrapping() {
	// Sometimes, an error on its own may not tell the full story of what
	// happened, or you may wish to add some additional context. Beginning with
	// Golang v1.13, one can wrap an error inside another error, and unwrap it
	// later (if desired).
	err1 := fmt.Errorf("i tried doing something")
	utils.PrintErrorContentAndType(err1)

	// To wrap an error, one should use the fmt.Errorf function. The magic string
	// formatting function to use here is %w.
	err2 := fmt.Errorf("i can't do that dave: %w", err1)
	utils.PrintErrorContentAndType(err2)

	// When printed, it looks like we've only concatenated err1 onto err2. In
	// reality, err2 is completely embedded inside err1 and we are able to get
	// err1 back if we desire:
	unwrappedErr := errors.Unwrap(err2)
	utils.PrintErrorContentAndType(unwrappedErr)

	// Unwrapping this further, we end up with a nil object because there is
	// nothing else to unwrap. This occurs either when we've reached the
	// beginning of an error chain or we've found an unwrappable error (more on
	// this later)
	unwrappedErr = errors.Unwrap(unwrappedErr)
	utils.PrintErrorContentAndType(unwrappedErr)
}

func multipleWrappingAndUnwrapping() {
	// It is possible to have many errors wrapped within one another
	multipleWrappedErrs := fmt.Errorf("base error")

	messages := []string{
		"unauthorized user",
		"cannot open file",
		"access denied",
	}

	for _, message := range messages {
		// Just keep wrapping errors for each of our messages
		multipleWrappedErrs = fmt.Errorf("%s: %w", message, multipleWrappedErrs)
	}

	utils.PrintErrorContentAndType(multipleWrappedErrs)

	// Now, let's unwrap all of our wrapped errors!
	for errors.Unwrap(multipleWrappedErrs) != nil {
		multipleWrappedErrs = errors.Unwrap(multipleWrappedErrs)
		utils.PrintErrorContentAndType(multipleWrappedErrs)
	}
}

func main() {
	basicWrappingAndUnwrapping()
	multipleWrappingAndUnwrapping()
}
