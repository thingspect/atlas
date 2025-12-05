package radiobridge

import (
	"github.com/thingspect/atlas/pkg/decode"
)

const identLinkQuality = 0xfb

// linkQuality parses a Link Quality payload from a []byte according to the
// spec. Link Quality payloads are used to confirm server connectivity and
// provide device radio measurements.
func linkQuality(body []byte) ([]*decode.Point, error) {
	// Link Quality payload must be 5 bytes.
	if len(body) != 5 {
		return nil, decode.FormatErr("link quality", "bad length", body)
	}

	if body[1] != identLinkQuality {
		return nil, decode.FormatErr("link quality", "bad identifier", body)
	}

	// Parse protocol and count.
	proto := int32(body[0] >> 4)
	msgs := []*decode.Point{{Attr: "proto", Value: proto}}

	count := int32(body[0] & clearProto)
	msgs = append(msgs, &decode.Point{Attr: "count", Value: count})

	// Parse sub-band.
	msgs = append(msgs, &decode.Point{Attr: "sub_band", Value: int32(body[2])})

	// Parse device RSSI and SNR.
	msgs = append(msgs, &decode.Point{
		Attr: "device_rssi", Value: int32(int8(body[3])),
	})
	msgs = append(msgs, &decode.Point{
		Attr: "device_snr", Value: int32(int8(body[4])),
	})

	return msgs, nil
}
