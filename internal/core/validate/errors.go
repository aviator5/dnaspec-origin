package validate

import (
	"fmt"
)

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

// ValidationErrors represents a collection of validation errors
type ValidationErrors []ValidationError

// Error implements the error interface for multiple errors
func (errs ValidationErrors) Error() string {
	if len(errs) == 0 {
		return "no errors"
	}
	if len(errs) == 1 {
		return errs[0].Error()
	}
	return fmt.Sprintf("%d validation errors", len(errs))
}

// Add appends a new validation error
func (errs *ValidationErrors) Add(field, message string) {
	*errs = append(*errs, ValidationError{Field: field, Message: message})
}

// IsEmpty returns true if there are no errors
func (errs ValidationErrors) IsEmpty() bool {
	return len(errs) == 0
}
