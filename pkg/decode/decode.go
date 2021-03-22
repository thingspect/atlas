// Package decode provides helper structs and functions for use by decoder
// function subpackages.
package decode

import (
	"time"

	"github.com/thingspect/atlas/pkg/consterr"
)

// ValidWindow is the window that a payload's timestamp is considered valid.
// This is based on the expected battery life of a gateway when it can be
// queueing payloads. If a timestamp is outside that window, the gateway likely
// has bogus time.
const ValidWindow = 4 * time.Hour

const ErrUnknownEvent consterr.Error = "unknown event type"

// Point represents an attribute-value pair as produced by a decoder function.
// Values should conform to common.DataPoint types, specifically int32, float64,
// string, bool, and []byte.
type Point struct {
	Attr  string
	Value interface{}
}
