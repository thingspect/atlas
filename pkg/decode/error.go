package decode

import (
	"fmt"

	"github.com/thingspect/atlas/pkg/consterr"
)

// ErrFormat is returned when a payload format is malformed.
const ErrFormat consterr.Error = "format"

// FormatErr returns an error due to a malformed payload.
func FormatErr(function string, reason string, body []byte) error {
	return fmt.Errorf("%s %w %s: %x", function, ErrFormat, reason, body)
}
