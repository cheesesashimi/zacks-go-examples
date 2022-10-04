package main

import (
	"fmt"

	"github.com/cheesesashimi/golang-errors/utils"
)

type resolvableError struct {
	knowledgebaseURL string
	err              error
}

func (r *resolvableError) Error() string {
	return fmt.Sprintf("resolvable error: '%s', see: %s", r.err, r.knowledgebaseURL)
}

func main() {
	// It is possible to add additional fields to an error struct. In fact, one
	// can have any number of arbitrary fields attached to an error struct. In
	// this specific example, we can attach an arbitrary URL whenever we
	// encounter a specific error.
	errs := []*resolvableError{
		{
			knowledgebaseURL: "https://link.to.kb.article/123",
			err:              fmt.Errorf("file not found"),
		},
		{
			knowledgebaseURL: "https://link.to.kb.article/456",
			err:              fmt.Errorf("user not authorized"),
		},
		{
			knowledgebaseURL: "https://link.to.kb.article/789",
			err:              fmt.Errorf("objects do not match"),
		},
	}

	for _, err := range errs {
		utils.PrintErrorContentAndType(err)
	}
}
