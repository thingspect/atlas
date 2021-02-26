package random

import (
	"github.com/google/uuid"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
)

// Org generates a random org with prefixed identifiers.
func Org(prefix string) *api.Org {
	return &api.Org{
		Id:   uuid.NewString(),
		Name: prefix + "-" + String(10),
	}
}

// Device generates a random device with prefixed identifiers.
func Device(prefix, orgID string) *api.Device {
	return &api.Device{
		Id:     uuid.NewString(),
		OrgId:  orgID,
		UniqId: prefix + "-" + String(16),
		Name:   prefix + "-" + String(10),
		Status: []api.Status{
			api.Status_ACTIVE,
			api.Status_DISABLED,
		}[Intn(2)],
		Token: uuid.NewString(),
		Decoder: []api.Decoder{
			api.Decoder_RAW,
			api.Decoder_GATEWAY,
			api.Decoder_RADIO_BRIDGE_DOOR_V1,
			api.Decoder_RADIO_BRIDGE_DOOR_V2,
		}[Intn(4)],
		Tags: Tags(prefix, Intn(4)+1),
	}
}

// User generates a random user with prefixed identifiers.
func User(prefix, orgID string) *api.User {
	return &api.User{
		Id:    uuid.NewString(),
		OrgId: orgID,
		Email: prefix + "-" + Email(),
		Role: []common.Role{
			common.Role_CONTACT,
			common.Role_VIEWER,
			common.Role_BUILDER,
			common.Role_ADMIN,
			common.Role_SYS_ADMIN,
		}[Intn(5)],
		Status: []api.Status{
			api.Status_ACTIVE,
			api.Status_DISABLED,
		}[Intn(2)],
		Tags: Tags(prefix, Intn(4)+1),
	}
}

// Tags generates n random tags with prefixed identifiers.
func Tags(prefix string, n int) []string {
	tags := []string{}
	for i := 0; i < n; i++ {
		tags = append(tags, prefix+"-"+String(10))
	}

	return tags
}
