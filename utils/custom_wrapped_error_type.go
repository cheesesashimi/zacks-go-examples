package utils

import "fmt"

// This is the struct that holds our custom error type
type CustomWrappedError struct {
	msg string
	err error
}

// This is a simple constructor function. It's not required but it can make
// things cleaner. Notice that it returns an error interface instead of a
// *CustomWrappedError
func NewCustomWrappedError(msg string, err error) error {
	return &CustomWrappedError{
		msg: msg,
		err: err,
	}
}

// We must implement this method to satisfy the basic error interface.
// Note: For brevity, we only return the field that was set on our error type.
// In practice, you can do all kinds of additional text processing here.
func (c *CustomWrappedError) Error() string {
	return fmt.Sprintf("%s: %s", c.msg, c.err)
}

// We must implement this method to satisfy the error unwrap interface.
func (c *CustomWrappedError) Unwrap() error {
	return c.err
}

// This method is specific to the CustomError type and cannot be used through
// the interface. We will explore how to do that.
func (c *CustomWrappedError) CustomFunc() string {
	return fmt.Sprintf("from custom wrapped error func: %s", c.Error())
}
