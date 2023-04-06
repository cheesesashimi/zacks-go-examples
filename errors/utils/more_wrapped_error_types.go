package utils

import (
	"errors"
	"fmt"
)

// A simple struct that holds an error associated with a given filename
type FileError struct {
	filename string
	err      error
}

// A helper function to create a new FileError instance
func NewFileError(filename string, err error) error {
	return &FileError{
		filename: filename,
		err:      err,
	}
}

// Implements a simple getter for the filename field since the field itself is private.
func (f *FileError) Filename() string {
	return f.filename
}

// Implements the error interface
func (f *FileError) Error() string {
	return fmt.Sprintf("an error occurred with file (%s): %s", f.filename, f.err)
}

// Implements the unwrap interface
func (f *FileError) Unwrap() error {
	return f.err
}

// This interrogates a given FileError or CustomWrappedError and prints information from it, if available.
func DebugFileAndCustomWrappedError(err error) {
	fmt.Println("original error text:", err)

	// Checks if we have a FileError
	var fErr *FileError
	if errors.As(err, &fErr) {
		// We have a FileError, so lets access fields on that struct.
		fmt.Println("we know the error occurred with this file:", fErr.filename)
	}

	// Checks if we have a CustomWrappedError
	var cErr *CustomWrappedError
	if errors.As(err, &cErr) {
		// We have a CustomWrappedError, so lets access fields on that struct.
		fmt.Println("we know we had the following message:", cErr.msg)
	}

	// If we can unwrap the given error to get the original error text, let's do
	// that here. This is somewhat flawed because we only get the first error in
	// the chain, assuming we can unwrap it at all.
	if unwrapped := errors.Unwrap(err); unwrapped != nil {
		fmt.Println("this is the unwrapped error we got:", unwrapped)
	}

	fmt.Println("===")
}

// Extracts a given error (using a provided matchFunc) that matches a given
// filename from an error chain.
func TraverseErrorChain(err error, matchFunc func(error) error) error {
	var unwrapped error = err

	// To find a specific within a given error chain, we can do this:
	//
	// Within a given error chain, one can only unwrap so far. Once we've
	// reached our limit of unwrapping, errors.Unwrap() will return nil. The
	// way errors.Unwrap() determines whether an error is unwrappable is
	// whether it implements the Unwrap() interface and whether calling that
	// interface returns a non-nil error.
	//
	// So what we do is we try to unwrap all the errors in a given error chain,
	// and then call errors.As() at each level to determine if it's the error
	// type we're interested in. It is worth noting that depending on the size of
	// the error chain as well as the various types contained therein, this can
	// be a computationally expensive operation.
	for unwrapped != nil {
		unwrapped = errors.Unwrap(unwrapped)

		if matched := matchFunc(unwrapped); matched != nil {
			return matched
		}
	}

	return nil
}
