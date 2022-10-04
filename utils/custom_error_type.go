package utils

import "fmt"

// This is the struct that holds our custom error type
type CustomError struct {
	msg string
}

// This is a simple constructor function. It's not required but it can make
// things cleaner. Notice that it returns an error interface instead of a
// *CustomError
func NewCustomError(msg string) error {
	return &CustomError{
		msg: msg,
	}
}

// We must implement this method to satisfy the basic error interface.
// Note: For brevity, we only return the field that was set on our error type.
// In practice, you can do all kinds of additional text processing here.
func (c *CustomError) Error() string {
	return c.msg
}

// This method is specific to the CustomError type and cannot be used through
// the interface. We will explore how to do that.
func (c *CustomError) CustomFunc() string {
	return fmt.Sprintf("from custom error func: %s", c.Error())
}
