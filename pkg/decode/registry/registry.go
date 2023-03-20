// Package registry provides data payload decoder function mappings.
package registry

import (
	"fmt"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/pkg/consterr"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/atlas/pkg/decode/globalsat"
	"github.com/thingspect/atlas/pkg/decode/radiobridge"
)

// ErrNotFound is returned when a decoder function name is not registered.
const ErrNotFound consterr.Error = "decoder function not found"

// Registry holds decoder-function mappings.
type Registry struct {
	funcs map[api.Decoder]func(body []byte) ([]*decode.Point, error)
}

// noOpDecoder passes through data payloads without decoding. This is for
// devices that do not use or do not yet have decoders.
func noOpDecoder(_ []byte) ([]*decode.Point, error) { return nil, nil }

// New returns a Registry with all decoder function mappings loaded.
func New() *Registry {
	return &Registry{
		funcs: map[api.Decoder]func(body []byte) ([]*decode.Point, error){
			api.Decoder_RAW:                  noOpDecoder,
			api.Decoder_GATEWAY:              noOpDecoder,
			api.Decoder_RADIO_BRIDGE_DOOR_V1: radiobridge.Door,
			api.Decoder_RADIO_BRIDGE_DOOR_V2: radiobridge.Door,
			api.Decoder_GLOBALSAT_CO2:        globalsat.CO2,
			api.Decoder_GLOBALSAT_CO:         globalsat.CO,
			api.Decoder_GLOBALSAT_PM25:       globalsat.PM25,
		},
	}
}

// Decode parses a device data payload by matching Registry function.
func (reg Registry) Decode(decoder api.Decoder, body []byte) (
	[]*decode.Point, error,
) {
	decFunc, ok := reg.funcs[decoder]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, decoder)
	}

	// Decode data payload.
	return decFunc(body)
}
