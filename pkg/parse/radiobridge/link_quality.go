package radiobridge

import (
	"github.com/thingspect/atlas/pkg/parse"
)

const identLinkQuality = 0xfb

// linkQuality parses a Link Quality payload from a []byte according to the
// spec. Link Quality payloads are used to confirm server connectivity and
// provide device radio measurements.
func linkQuality(body []byte) ([]*parse.Point, error) {
	// Link Quality payload must be 5 bytes.
	if len(body) != 5 {
		return nil, parse.ErrFormat("link quality", "bad length", body)
	}

	if body[1] != identLinkQuality {
		return nil, parse.ErrFormat("link quality", "bad identifier", body)
	}

	// Parse protocol and count.
	proto := int(body[0] >> 4)
	msgs := []*parse.Point{{Attr: "proto", Value: proto}}

	count := int(body[0] & clearProto)
	msgs = append(msgs, &parse.Point{Attr: "count", Value: count})

	// Parse sub-band.
	msgs = append(msgs, &parse.Point{Attr: "sub_band", Value: int(body[2])})

	// Parse device RSSI and SNR.
	msgs = append(msgs, &parse.Point{Attr: "dev_rssi",
		Value: int(int8(body[3]))})
	msgs = append(msgs, &parse.Point{Attr: "dev_snr",
		Value: int(int8(body[4]))})

	return msgs, nil
}
