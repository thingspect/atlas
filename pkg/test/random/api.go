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
func Device(prefix, orgID string) *common.Device {
	return &common.Device{
		Id:     uuid.NewString(),
		OrgId:  orgID,
		UniqId: prefix + "-" + String(16),
		Name:   prefix + "-" + String(10),
		Status: []common.Status{
			common.Status_ACTIVE,
			common.Status_DISABLED,
		}[Intn(2)],
		Token: uuid.NewString(),
		Decoder: []common.Decoder{
			common.Decoder_RAW,
			common.Decoder_GATEWAY,
			common.Decoder_RADIO_BRIDGE_DOOR_V1,
			common.Decoder_RADIO_BRIDGE_DOOR_V2,
		}[Intn(4)],
		Tags: Tags(prefix, Intn(4)+1),
	}
}

// Rule generates a random rule with prefixed identifiers.
func Rule(prefix, orgID string) *api.Rule {
	return &api.Rule{
		Id:    uuid.NewString(),
		OrgId: orgID,
		Name:  prefix + "-" + String(10),
		Status: []common.Status{
			common.Status_ACTIVE,
			common.Status_DISABLED,
		}[Intn(2)],
		Tag:  prefix + "-" + String(10),
		Attr: prefix + "-" + String(10),
		Expr: []string{`true`, `false`}[Intn(2)],
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
		Status: []common.Status{
			common.Status_ACTIVE,
			common.Status_DISABLED,
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
