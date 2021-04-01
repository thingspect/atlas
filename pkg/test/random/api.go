package random

import (
	"time"

	"github.com/google/uuid"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"google.golang.org/protobuf/types/known/timestamppb"
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
func Rule(prefix, orgID string) *common.Rule {
	return &common.Rule{
		Id:    uuid.NewString(),
		OrgId: orgID,
		Name:  prefix + "-" + String(10),
		Status: []common.Status{
			common.Status_ACTIVE,
			common.Status_DISABLED,
		}[Intn(2)],
		DeviceTag: prefix + "-" + String(10),
		Attr:      prefix + "-" + String(10),
		Expr:      []string{`true`, `false`}[Intn(2)],
	}
}

// Event generates a random event with prefixed identifiers.
func Event(prefix, orgID string) *api.Event {
	return &api.Event{
		OrgId:     orgID,
		UniqId:    prefix + "-" + String(16),
		RuleId:    uuid.NewString(),
		CreatedAt: timestamppb.New(time.Now().UTC().Truncate(time.Millisecond)),
		TraceId:   uuid.NewString(),
	}
}

// Alarm generates a random alarm with prefixed identifiers.
func Alarm(prefix, orgID, ruleID string) *api.Alarm {
	return &api.Alarm{
		Id:     uuid.NewString(),
		OrgId:  orgID,
		RuleId: ruleID,
		Name:   prefix + "-" + String(10),
		Status: []common.Status{
			common.Status_ACTIVE,
			common.Status_DISABLED,
		}[Intn(2)],
		Type:            api.AlarmType_APP,
		UserTags:        Tags(prefix, Intn(4)+1),
		SubjectTemplate: `rule name is: {{.rule.Name}}`,
		BodyTemplate:    `device status is: {{.device.Status}}`,
		RepeatInterval:  int32(Intn(99) + 1),
	}
}

// Alert generates a random alert with prefixed identifiers.
func Alert(prefix, orgID string) *api.Alert {
	return &api.Alert{
		OrgId:   orgID,
		UniqId:  prefix + "-" + String(16),
		AlarmId: uuid.NewString(),
		UserId:  uuid.NewString(),
		Status: []api.AlertStatus{
			api.AlertStatus_SENT,
			api.AlertStatus_ERROR,
		}[Intn(2)],
		Error:   []string{"", prefix + "-" + String(10)}[Intn(2)],
		TraceId: uuid.NewString(),
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
		Tags:   Tags(prefix, Intn(4)+1),
		AppKey: []string{"", String(30)}[Intn(2)],
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
