// Package consterr provides types and functions to create constant errors.
package consterr

// Error supports constant errors and implements the error interface.
type Error string

// Verify Error implements error.
var _ error = new(Error)

// Error returns the error.
func (e Error) Error() string { return string(e) }
