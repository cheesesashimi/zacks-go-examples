package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
)

// In this section, we will interrogate a given error to determine if there is
// additional data that can we can use and we can decide a different path
// forward, depending on the data contained therein.

// This function accepts a path to a file. It attempts to do the following:
// 1. Open the file and read its contents into memory.
// 2. Decode the JSON contents of the file into an untyped map, then print the map.
func readAJSONFile(path string) error {
	fmt.Println("reading from:", path)

	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		// Add some additional context to the error returned by ioutil.ReadFile
		return fmt.Errorf("readAJSONFile file error: %w", err)
	}

	dst := map[string]interface{}{}
	if err := json.Unmarshal(fileBytes, &dst); err != nil {
		// Add some additional context to the error returned by json.Unmarshal
		return fmt.Errorf("readAJSONFile JSON error: %w", err)
	}

	fmt.Println("this is our data:", dst)

	return nil
}

// This function will try to fallback to additional JSON files based upon
// information about the error we've found.
func readJSONFileAndFallback(path string) error {
	if err := readAJSONFile(path); err != nil {
		// Just print out our error.
		fmt.Println("uh-oh:", err)

		// Does our error chain contain an fs.PathError?
		var pErr *fs.PathError
		// This checks if we've got an fs.PathError someplace in our error chain.
		// As discussed elsewehere, errors.As() traverses all the errors in a given
		// error chain, performing a type assertion at each level to determine if
		// the error matches the given type.
		// See: https://cs.opensource.google/go/go/+/refs/tags/go1.19.2:src/io/fs/fs.go;l=244-248
		if errors.As(err, &pErr) {
			// As discussed elsewhere, errors.Is() traverses all the errors in a
			// given error chain, determining if a given error value matches. This
			// checks if we have os.ErrNotExist someplace in our error chain. This is
			// a different operation than errors.As() because os.ErrNotExist is a
			// sentinel error instead of a specific error type.
			//
			// See: https://cs.opensource.google/go/go/+/refs/tags/go1.19.2:src/internal/oserror/errors.go;drc=fb4f7fdb26da9ed0fee6beab280c84b399edaa42;l=16
			if errors.Is(pErr, os.ErrNotExist) {
				fmt.Println("we couldn't read from", pErr.Path, "falling back to malformed.json")
				return readJSONFileAndFallback("malformed.json")
			}

			return fmt.Errorf("an unknown fs.PathError occurred: %w", pErr)
		}

		// Does our error chain contain a JSON error?
		// See: https://cs.opensource.google/go/go/+/master:src/encoding/json/scanner.go;l=47-50?q=json.SyntaxError&ss=go%2Fgo
		var jsonErr *json.SyntaxError
		if errors.As(err, &jsonErr) {
			fmt.Println("we couldn't parse JSON due to a syntax error, falling back to good.json")
			return readJSONFileAndFallback("good.json")
		}
	}

	return nil
}

func main() {
	if err := readJSONFileAndFallback("/file/does/not/exist/go/away"); err != nil {
		panic(err)
	}
}
