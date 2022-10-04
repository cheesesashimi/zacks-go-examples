package main

import (
	"compress/gzip"
	"errors"
	"fmt"

	"github.com/cheesesashimi/golang-errors/utils"
)

const errText string = "i'm an error"

func basicEquality() {
	// Two errors containing the same text can actually be completely different errors:
	err1 := errors.New(errText)
	err2 := fmt.Errorf(errText)

	// Cmoparing their contents shows that they contain the same text. However,
	// this practice is discouraged for reasons we'll discover later.
	fmt.Printf("%s == %s? %v\n", err1, err2, err1.Error() == err2.Error())

	// Instead, it is typically better to compare the errors directly
	fmt.Printf("%s == %s? %v\n", err1, err2, err1 == err2)

	// Both are the same type, so why are they not equal?
	// In short, they're pointers to different areas in memory.
	fmt.Printf("%T: %p\n%T: %p\n", err1, err1, err2, err2)
}

func sentinalErrors() {
	// The first is that certain packages export an error that we can test
	// against. (Note: Golang cannot have error constants)
	//
	// For example, the gzip package has a couple of exported errors contained
	// within it that one can check against:
	utils.PrintErrorContentAndType(gzip.ErrHeader)
	utils.PrintErrorContentAndType(gzip.ErrChecksum)
}

func customErrorTypes() {
	err1 := errors.New(errText)

	// Other packages sometimes define their own custom error types.
	customErr := utils.NewCustomError(errText)

	// Trying to compare this to the above errors will fail the equality test.
	// Because not only does this refer to a different memory location, but the
	// error type is also different.
	fmt.Printf("%s (%T) == %s (%T)? %v\n", customErr, customErr, err1, err1, customErr == err1)
}

func typeAssertions() {
	customErr := utils.NewCustomError(errText)

	// We can perform a type assertion to ensure that we a customError. All a
	// type assetion does is assert that a variable is the type we expect it to
	// be.
	assertedCustomErr, ok := customErr.(*utils.CustomError)
	fmt.Printf("%s: %T, ok? %v\n", customErr, customErr, ok)

	// This now means we can access the methods specific to the customError type:
	fmt.Println(assertedCustomErr.CustomFunc())

	// While this works, what happens when we wrap this error inside of another?
	wrapped := fmt.Errorf("outer: %w", customErr)
	utils.PrintErrorContentAndType(wrapped)

	// Our type assertion fails because the error whose type we're asserting is
	// not the same type as the one we're looking for.
	_, ok = wrapped.(*utils.CustomError)
	fmt.Printf("%s: %T, ok? %v\n", wrapped, wrapped, ok)

	// But if we unwrap it, the type assertion succeeds! What gives!?!
	unwrapped := errors.Unwrap(wrapped)
	assertedCustomErr, ok = unwrapped.(*utils.CustomError)
	fmt.Printf("%s: %T, ok? %v\n", unwrapped, unwrapped, ok)
	fmt.Println(assertedCustomErr.CustomFunc())
}

func errorsIs() {
	// First, we'll look at the errors.Is() function. This function compares two
	// errors to determine if they're the same error, unwrapping as necessary.
	err1 := errors.New(errText)
	err2 := fmt.Errorf("wrapped: %w", err1)
	customErr := utils.NewCustomError(errText)
	wrappedCustomErr := utils.NewCustomWrappedError("wrapped", customErr)

	fmt.Printf("%v\n", errors.Is(err2, err1))                  // This evaluates to True because we've wrapped one error in another.
	fmt.Printf("%v\n", errors.Is(err1, errors.Unwrap(err2)))   // This evaluates to True because we're comparing the unwrapped error to the original error.
	fmt.Printf("%v\n", errors.Is(wrappedCustomErr, customErr)) // This evaluates to True because we've wrapped our custom error type within another error.
	fmt.Printf("%v\n", errors.Is(customErr, err1))             // This evaluates to False because even though both errors have the same text, they are not the same type.
}

func errorsAs() {
	customErr := utils.NewCustomError(errText)
	wrappedCustomErr := fmt.Errorf("wrapped custom error: %w", customErr)
	doublyWrappedCustomErr := fmt.Errorf("double wrapped: %w", wrappedCustomErr)

	// For this example, our goal is to get the underlying customError, if one is available.
	errs := []error{
		customErr,
		fmt.Errorf(errText), // This is a generic error
		wrappedCustomErr,
		errors.New(errText), // This is also a generic error
		doublyWrappedCustomErr,
	}

	for _, err := range errs {
		// First, we declare an uninitialized pointer for the customError type.
		var cErr *utils.CustomError
		// Next, we pass in the error that we want to know more about, as well as a
		// pointer to the uninitialized pointer (I know it's weird...). errors.As()
		// will unwrap all errors in the chain and if it finds one that's the same
		// type as the one we're expecting, it will initialize the memory that cErr
		// points to.
		if errors.As(err, &cErr) {
			// We've found one that matches our customError type.
			fmt.Printf("customError found: calling CustomFunc() ")
			utils.PrintErrorContentAndType(err)
			fmt.Println(cErr.CustomFunc())
		} else {
			// We haven't found one that matches our customError type.
			fmt.Printf("not a customError: ")
			utils.PrintErrorContentAndType(err)
		}
	}
}

func main() {
	// Not all errors indicate that something is wrong. In fact, in certain
	// cases, an error with specific contents can be expected. For example, a
	// common pattern when reading a file is to check if any errors returned are
	// an end-of-file (EOF) error. Just like other types, one can check that an
	// error meets certain criteria. However, there are a few pitfalls to be
	// aware of.

	// For example, basic equality works differently for errors:
	basicEquality()

	// So how do we ensure that errors are what we expect them to be? The first
	// technique is relying upon a sentinal error value that is exported from a
	// given package.
	sentinalErrors()

	// Other packages define their own custom error types that can contain more information:
	customErrorTypes()

	// Using those custom error types and a technique called type assertions, we
	// can get more information. However, as we'll soon learn, using type
	// assertions in this manner is also full of pitfalls:
	typeAssertions()

	// So now that it seems like we've exhausted all possibilities for looking at
	// errors, how do we actually work with errors in a more useful way?
	//
	// Starting in Golang 1.13, the standard library gained the errors.Is and
	// errors.As functions. These functions can be used to work with errors much
	// more efficiently than repeatedly checking for certain error types,
	// unwrapping, etc. To that end, they perform the unwrapping and type
	// asserting for you.
	//
	// Let's first look at errors.Is():
	errorsIs()

	// While errors.Is() compares two errors and their types to determine
	// equality, sometimes, we need access to the underlying error type to access
	// either a specific field or a method. This is where errors.As() comes into
	// play:
	errorsAs()
}
