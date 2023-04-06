package main

import (
	"errors"
	"fmt"

	"github.com/cheesesashimi/zacks-go-examples/errors/utils"
)

func main() {
	// At their core, Golang errors are a built-in type used to convey
	// information about the outcome of an operation.

	// One can construct an error in two main ways:
	err := errors.New("hello, i'm an error")
	utils.PrintErrorContentAndType(err)

	// -or-

	err = fmt.Errorf("hello, i'm also an error")
	utils.PrintErrorContentAndType(err)

	// Functions can return an error either as the sole return value.
	err = iReturnAnError()
	utils.PrintErrorContentAndType(err)

	// -or-

	// In conjunction with another value. It is Golang convention that the error
	// be the last return value, as shown here.
	_, err = iAlsoReturnAnError()
	utils.PrintErrorContentAndType(err)
}

func iReturnAnError() error {
	return fmt.Errorf("i'm a returned error")
}

func iAlsoReturnAnError() (bool, error) {
	return false, fmt.Errorf("i'm also a returned error")
}
