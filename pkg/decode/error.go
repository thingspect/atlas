package decode

import (
	"fmt"

	"github.com/thingspect/atlas/pkg/consterr"
)

const errFormat consterr.Error = "format"

// ErrFormat returns an error due to a malformed payload.
func ErrFormat(function string, reason string, body []byte) error {
	return fmt.Errorf("%s %w %s: %x", function, errFormat, reason, body)
}
