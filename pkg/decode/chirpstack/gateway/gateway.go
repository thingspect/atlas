// Package gateway provides parse functions for ChirpStack gateway payloads.
package gateway

import (
	"fmt"

	"github.com/thingspect/atlas/pkg/decode"
)

// Parse parses a gateway payload from a []byte according to the spec. Points
// are built from successful parse results. If a fatal error is encountered, it
// is returned along with any valid points.
func Parse(event string, body []byte) ([]*decode.Point, error) {
	switch event {
	case "up":
		return gatewayUp(body)
	case "stats":
		return gatewayStats(body)
	case "ack":
		return gatewayAck(body)
	case "exec":
		return gatewayExec(body)
	default:
		return nil, fmt.Errorf("%w: %s, %x", decode.ErrUnknownEvent, event, body)
	}
}
