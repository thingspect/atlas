package random

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/thingspect/proto/go/api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Org generates a random org with prefixed identifiers.
func Org(prefix string) *api.Org {
	return &api.Org{
		Id:          uuid.NewString(),
		Name:        prefix + "-" + String(10),
		DisplayName: prefix + "-" + String(10),
		Email:       prefix + "-" + Email(),
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

// Rule generates a random rule with prefixed identifiers.
func Rule(prefix, orgID string) *api.Rule {
	return &api.Rule{
		Id:    uuid.NewString(),
		OrgId: orgID,
		Name:  prefix + "-" + String(10),
		Status: []api.Status{
			api.Status_ACTIVE,
			api.Status_DISABLED,
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
		Status: []api.Status{
			api.Status_ACTIVE,
			api.Status_DISABLED,
		}[Intn(2)],
		Type: []api.AlarmType{
			api.AlarmType_APP,
			api.AlarmType_SMS,
			api.AlarmType_EMAIL,
		}[Intn(3)],
		UserTags:        Tags(prefix, Intn(4)+1),
		SubjectTemplate: `rule name is: {{.rule.Name}}`,
		BodyTemplate:    `device status is: {{.device.Status}}`,
		//nolint:gosec // Safe conversion for limited values.
		RepeatInterval: int32(Intn(99) + 1),
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
		Name:  prefix + "-" + String(10),
		Email: prefix + "-" + Email(),
		Role: []api.Role{
			api.Role_CONTACT,
			api.Role_VIEWER,
			api.Role_PUBLISHER,
			api.Role_BUILDER,
			api.Role_ADMIN,
			api.Role_SYS_ADMIN,
		}[Intn(6)],
		Status: []api.Status{
			api.Status_ACTIVE,
			api.Status_DISABLED,
		}[Intn(2)],
		Tags: Tags(prefix, Intn(4)+1),
	}
}

// SMSUser generates a random SMS user with prefixed identifiers.
func SMSUser(prefix, orgID string) *api.User {
	return &api.User{
		Id:    uuid.NewString(),
		OrgId: orgID,
		Name:  prefix + "-" + String(10),
		Email: prefix + "-" + Email(),
		// https://en.wikipedia.org/wiki/555_(telephone_number)
		Phone: "+1" + strconv.Itoa(Intn(900)+100) + "5550" +
			strconv.Itoa(Intn(100)+100),
		Role: []api.Role{
			api.Role_CONTACT,
			api.Role_VIEWER,
			api.Role_PUBLISHER,
			api.Role_BUILDER,
			api.Role_ADMIN,
			api.Role_SYS_ADMIN,
		}[Intn(6)],
		Status: []api.Status{
			api.Status_ACTIVE,
			api.Status_DISABLED,
		}[Intn(2)],
		Tags: Tags(prefix, Intn(4)+1),
	}
}

// AppUser generates a random mobile application user with prefixed identifiers.
func AppUser(prefix, orgID string) *api.User {
	return &api.User{
		Id:    uuid.NewString(),
		OrgId: orgID,
		Name:  prefix + "-" + String(10),
		Email: prefix + "-" + Email(),
		Role: []api.Role{
			api.Role_CONTACT,
			api.Role_VIEWER,
			api.Role_PUBLISHER,
			api.Role_BUILDER,
			api.Role_ADMIN,
			api.Role_SYS_ADMIN,
		}[Intn(6)],
		Status: []api.Status{
			api.Status_ACTIVE,
			api.Status_DISABLED,
		}[Intn(2)],
		Tags:   Tags(prefix, Intn(4)+1),
		AppKey: String(30),
	}
}

// Key generates a random API key with prefixed identifiers.
func Key(prefix, orgID string) *api.Key {
	return &api.Key{
		Id:    uuid.NewString(),
		OrgId: orgID,
		Name:  prefix + "-" + String(10),
		Role: []api.Role{
			api.Role_CONTACT,
			api.Role_VIEWER,
			api.Role_BUILDER,
			api.Role_ADMIN,
			api.Role_SYS_ADMIN,
		}[Intn(5)],
	}
}

// Tags generates n random tags with prefixed identifiers.
func Tags(prefix string, n int) []string {
	tags := []string{}
	for range n {
		tags = append(tags, prefix+"-"+String(10))
	}

	return tags
}
