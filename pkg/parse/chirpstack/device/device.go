// Package device provides parse functions for ChirpStack device payloads.
package device

import (
	"fmt"
	"time"

	"github.com/thingspect/atlas/pkg/parse"
)

// Device parses a device payload from a []byte according to the spec. Points
// and (optional) data []byte and times are built from successful parse results.
// If a fatal error is encountered, it is returned along with any valid points.
func Device(event string, body []byte) ([]*parse.Point, time.Time, []byte,
	error) {
	switch event {
	case "up":
		return deviceUp(body)
	case "join":
		msgs, ts, err := deviceJoin(body)
		return msgs, ts, nil, err
	case "ack":
		msgs, err := deviceAck(body)
		return msgs, time.Now(), nil, err
	case "error":
		msgs, err := deviceError(body)
		return msgs, time.Now(), nil, err
	case "txack":
		msgs, err := deviceTxAck(body)
		return msgs, time.Now(), nil, err
	default:
		return nil, time.Time{}, nil, fmt.Errorf("%w: %s, %x",
			parse.ErrUnknownEvent, event, body)
	}
}
