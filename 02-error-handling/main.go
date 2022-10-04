package main

import "fmt"

func iReturnAnError() error {
	return fmt.Errorf("i'm the returned error")
}

func iAlsoReturnAnError() error {
	return fmt.Errorf("i'm also a returned error")
}

func iReturnANilError() error {
	return nil
}

func checkError(err error) {
	// Generally speaking, if an error is nil, no error occurred.
	if err == nil {
		fmt.Println("no error returned!")
	} else {
		fmt.Printf("an error was returned: %s: %T\n", err, err)
	}
}

func iDoStuff() error {
	// Sometimes, checking these errors can get repetitive. However, it
	// explicitly identifies what should happen whenever we encounter an error.
	// In this specific example, we return immediately upon discovering an error,
	// returning the error value back to the caller unmodified so that the caller
	// can determine what to do with it.
	//
	// In some cases, we may want to do something different; this will be covered
	// later.
	if err := iReturnAnError(); err != nil {
		return err
	}

	if err := iAlsoReturnAnError(); err != nil {
		return err
	}

	if err := iReturnAnError(); err != nil {
		return err
	}

	return nil
}

func main() {
	checkError(iReturnANilError())

	checkError(iReturnAnError())

	checkError(iDoStuff())
}
