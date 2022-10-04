package utils

import "fmt"

func PrintErrorContentAndType(err error) {
	fmt.Printf("%s: %T\n", err, err)
}
