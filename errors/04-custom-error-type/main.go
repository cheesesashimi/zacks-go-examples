package main

import (
	"errors"

	"github.com/cheesesashimi/zacks-go-examples/errors/utils"
)

func main() {
	// In Golang, it is possible to create a new type based upon a primitive type
	// (such as a string, integer, etc.). This allows one to lean on the compiler
	// to ensure that you only pass the correct type. Errors are no different.
	//
	// The way we do this in Golang is to take advantage of the fact that errors
	// are an interface. For brevity, an interface is requiring that a given type
	// implement certain methods. In the case of errors, they must implement the
	// error interface, which is a method called Error() which returns a string.
	//
	// Errors can also optionally define the unwrap interface, which allows one
	// to nest another error inside.

	// Here, we create a new instance of our custom error:
	err1 := utils.NewCustomError("i'm a custom error")

	// Notice that we have a new error type. This will become more useful to us later.
	utils.PrintErrorContentAndType(err1)

	// We can also make our custom error types be able to unwrap errors. To do that, we must implement the Unwrap interface.
	err2 := utils.NewCustomWrappedError("i'm a custom wrapped error", err1)

	utils.PrintErrorContentAndType(err2)

	// What happens when we unwrap our wrapped error?
	utils.PrintErrorContentAndType(errors.Unwrap(err2))

	// What happens when we unwrap our unwrapped error?
	utils.PrintErrorContentAndType(errors.Unwrap(err1))
}
