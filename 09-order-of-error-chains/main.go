package main

import (
	"errors"
	"fmt"

	"github.com/cheesesashimi/golang-errors/utils"
)

func main() {
	filename := "/a/nonexistant/file"
	parseMessage := "contents malformed"
	parseError := fmt.Errorf("contents in an unknown format")

	// When you have multiple error types in a given error chain, it generally does
	// not matter where in the error chain a given error occurs unless you're
	// looking for multiple errors of the same kind within the same chain (more on this later).
	//
	// In general, you should nest your errors in a way that makes sense to an
	// end-user when the details are printed by the Error() function. This is
	// also a good reminder that the output of the Error() function is
	// intended only for consumption by humans (except potentially in a unit-test
	// context).
	//
	// A Go program should interrogate a given error type using the fields and
	// methods that are available to it. In this way, it doesn't matter what
	// order you nest your errors on your error chain.
	errorsToDebug := []error{
		// This first error has a FileError on the outside with a CustomWrappedError contained within it.
		utils.NewFileError(filename, utils.NewCustomWrappedError(parseMessage, parseError)),
		// This second error has a CustomWrappedError as the outer layer, with a FileError contained within it.
		utils.NewCustomWrappedError(parseMessage, utils.NewFileError(filename, parseError)),

		// These two are here just to show what happens with non-special error types.
		parseError,
		fmt.Errorf("outer error: %w", parseError),
	}

	// Iterate through our errors to debug and figure out what we can about them.
	for _, errToDebug := range errorsToDebug {
		utils.DebugFileAndCustomWrappedError(errToDebug)
	}

	// Order *does* make a difference when you have multiple errors of the same
	// concrete type within the same error chain. This is because errors.As()
	// will only return the first example it finds.
	//
	// For this next example, we want to interrogate the error contained within
	// the innermost FileError (the one associated with the filename
	// /yet/another/nonexistant/file).
	nestedFileErrors := utils.NewFileError("/a/nonexistant/file",
		utils.NewFileError("/another/nonexistant/file",
			utils.NewFileError("/yet/another/nonexistant/file", fmt.Errorf("innermost error"))))

	// Let's first try our debug function:
	utils.DebugFileAndCustomWrappedError(nestedFileErrors)

	// The output only matched /a/nonexistant/file.
	// So how do we interrogate the innermost FileError, which is associated with /yet/another/nonexistant/file?
	// Here's how:
	ourErr := utils.TraverseErrorChain(nestedFileErrors, func(err error) error {
		// This function defines what we're looking for and how it matches what we're
		// looking for. We pass this into utils.TraverseErrorChain, which traverses the
		// complete error chain for a given error until we either match what we're
		// looking for or until we've gone as far down the error chain as possible.

		// At each unwrapping level, we do the following:
		// 1. Call errors.As() at the current level to determine if we have a
		// FileError.
		// 2. Then, we call the Filename() method on the FileError and compare it
		// to our provided filename.
		var fErr *utils.FileError
		if errors.As(err, &fErr) && fErr.Filename() == "/yet/another/nonexistant/file" {
			// We have a FileError and the filename matches the one we're looking for,
			// so we're done!
			return fErr
		}

		// We haven't found anything, so we return nil here.
		return nil
	})

	utils.DebugFileAndCustomWrappedError(ourErr)
}
