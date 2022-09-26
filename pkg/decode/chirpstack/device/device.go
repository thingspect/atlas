// Package device provides parse functions for ChirpStack device payloads.
package device

import (
	"fmt"

	"github.com/thingspect/atlas/pkg/decode"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Parse parses a device payload from a []byte according to the spec. Points,
// optional data []byte, and a timestamp are built from successful parse
// results. If a fatal error is encountered, it is returned along with any valid
// points.
func Parse(event string, body []byte) (
	[]*decode.Point, *timestamppb.Timestamp, []byte, error,
) {
	switch event {
	case "up":
		return deviceUp(body)
	case "join":
		msgs, err := deviceJoin(body)

		return msgs, timestamppb.Now(), nil, err
	case "ack":
		msgs, err := deviceAck(body)

		return msgs, timestamppb.Now(), nil, err
	case "log":
		msgs, err := deviceLog(body)

		return msgs, timestamppb.Now(), nil, err
	case "txack":
		msgs, err := deviceTXAck(body)

		return msgs, timestamppb.Now(), nil, err
	case "status":
		msgs, err := deviceStatus(body)

		return msgs, timestamppb.Now(), nil, err
	default:
		return nil, nil, nil, fmt.Errorf("%w: %s, %x", decode.ErrUnknownEvent,
			event, body)
	}
}
