package main

import (
	"errors"
	"fmt"
)

// One of the most common interfaces that Go developers use is the built-in
// error interface. The error interface is also fairly simple as well:
//
// type error interface {
//		Error() string
//    Unwrap() error
// }
//
// You'll observe that all of our custom error types implement both an Error()
// method and an Unwrap() method. This is because there is another interface
// which specifies the Unwrap() method. A custom error type does not require an
// Unwrap() method. However, if this custom error type will wrap other errors,
// it is advised that one be added.

// This is a very basic customError type that contains an info field for
// additional info as well as an err field for an underlying error.
type customError struct {
	info string
	err  error
}

// This is a simple constructor. Notice what type it returns. In this case, it
// returns an error type and not the underlying concrete type (customError).
// This works because customError satisfies the error interface.
func newCustomError(info string, err error) error {
	return &customError{
		info: info,
		err:  err,
	}
}

// This returns the original underlying error.
func (c *customError) Unwrap() error {
	return c.err
}

// This returns a string representation of the customError.
func (c *customError) Error() string {
	return fmt.Sprintf("customError (%s): %s", c.info, c.err)
}

// This is an additional method that is specific to the customError type. It
// does not implement interfaces that we're aware of.
func (c *customError) Info() string {
	return c.info
}

// This is a simple error type that does not have an Unwrap() method.
type anotherErrType string

func (a anotherErrType) Error() string {
	return string(a)
}

// Some packages define exported error values like this that one can compare to
// understand if they've reached a certain condition. More on this later.
var sentinalErr error = fmt.Errorf("sentinal error")

func returnSentinalErr() error {
	return sentinalErr
}

func main() {
	// Why are custom error types useful?
	//
	// In short, it allows one to define multiple distinct error types
	// representing a wide range of situations. It then allows one to wrap and
	// unwrap errors and make decisions based upon an errors underlying type as
	// opposed to its contents.
	origErr := fmt.Errorf("an error")
	//anotherErr := fmt.Errorf("an error")
	anotherErr := anotherErrType("an error")

	// Why would one do this?
	//
	// Even though these two errors have the same text contained within, they are
	// completely different types.
	fmt.Println("origErr == anotherErr?", origErr == anotherErr)
	fmt.Println("origErr.Error() == anotherErr.Error()?", origErr.Error() == anotherErr.Error())

	// So instead of inferring whether the error is the same based upon testing
	// equality, one should either use the sentinal error pattern or custom error
	// types. This is an example of the sentinal error pattern where we compare
	// the error against a known error.
	fmt.Println("returnSentinalErr() == sentinalErr?", returnSentinalErr() == sentinalErr)

	// But what happens if the sentinalError is wrapped in another error?
	fmt.Println("newCustomError(\"from sentinal\", sentinalErr) == sentinalErr?", newCustomError("from sentinal", sentinalErr) == sentinalErr)

	// Instead, one should use errors.Is() since it unwraps the errors to inform whether it is a match.
	fmt.Println("errors.Is(newCustomError(\"from sentinal\", sentinalErr), sentinalErr)?", errors.Is(newCustomError("from sentinal", sentinalErr), sentinalErr))

	// Sometimes, you may want to base this off of the errors' underlying type instead. This is where the errors.As() function comes into play.
	custErr := newCustomError("custom info", origErr)
	var ce *customError
	fmt.Println("errors.As(custErr, &ce)?", errors.As(custErr, &ce))

	fmt.Println("errors.As(origErr, &ce)?", errors.As(origErr, &ce))

	// In this case, our error is very deeply wrapped. What happens in this case?
	deeplyWrappedError := fmt.Errorf("error: %w", fmt.Errorf("something went wrong: %w", fmt.Errorf("something else: %w", custErr)))
	fmt.Println("deeplyWrappedError matches custErr?", errors.As(deeplyWrappedError, &ce))
	fmt.Println("deeplyWrappedError contains origErr?", errors.Is(deeplyWrappedError, origErr))

	errs := []error{
		sentinalErr,
		origErr,
		anotherErr,
		custErr,
		deeplyWrappedError,
		newCustomError("other info", origErr),
		newCustomError("custom info", sentinalErr),
	}

	for _, err := range errs {
		fmt.Printf("original error (%T): %s\n", err, err)

		// Using errors.As() allows us to figure out if a given error matches an
		// underlying error type. This works even in deeply-nested cases. It will
		// perform a type assertion and get the original concrete type if it matches.
		var ce *customError
		if errors.As(err, &ce) {
			if ce.Info() == "custom info" {
				fmt.Println("we've got a matching error:", ce)
			} else {
				fmt.Println("we've matched our type but not our info")
			}
			if errors.Is(ce, origErr) {
				fmt.Println("we matched our original error, even after being wrapped")
			} else if errors.Is(ce, sentinalErr) {
				fmt.Println("we found a wrapped sentinal error")
			}
		}
	}
}
