package models

// ValidationError represents model validation error.
type ValidationError struct {
	msg string
}

// Error is error interface implementation.
func (e *ValidationError) Error() string {
	return e.msg
}
