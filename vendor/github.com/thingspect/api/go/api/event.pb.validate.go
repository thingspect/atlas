// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/event.proto

package api

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/ptypes"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = ptypes.DynamicAny{}
)

// Validate checks the field values on Event with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *Event) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for OrgId

	// no validation rules for UniqId

	// no validation rules for RuleId

	if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return EventValidationError{
				field:  "CreatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for TraceId

	return nil
}

// EventValidationError is the validation error returned by Event.Validate if
// the designated constraints aren't met.
type EventValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e EventValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e EventValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e EventValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e EventValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e EventValidationError) ErrorName() string { return "EventValidationError" }

// Error satisfies the builtin error interface
func (e EventValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sEvent.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = EventValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = EventValidationError{}

// Validate checks the field values on ListEventsRequest with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *ListEventsRequest) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for RuleId

	if v, ok := interface{}(m.GetEndTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ListEventsRequestValidationError{
				field:  "EndTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetStartTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ListEventsRequestValidationError{
				field:  "StartTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	switch m.IdOneof.(type) {

	case *ListEventsRequest_UniqId:
		// no validation rules for UniqId

	case *ListEventsRequest_DeviceId:
		// no validation rules for DeviceId

	default:
		return ListEventsRequestValidationError{
			field:  "IdOneof",
			reason: "value is required",
		}

	}

	return nil
}

// ListEventsRequestValidationError is the validation error returned by
// ListEventsRequest.Validate if the designated constraints aren't met.
type ListEventsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListEventsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListEventsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListEventsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListEventsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListEventsRequestValidationError) ErrorName() string {
	return "ListEventsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListEventsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListEventsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListEventsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListEventsRequestValidationError{}

// Validate checks the field values on ListEventsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ListEventsResponse) Validate() error {
	if m == nil {
		return nil
	}

	for idx, item := range m.GetEvents() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListEventsResponseValidationError{
					field:  fmt.Sprintf("Events[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// ListEventsResponseValidationError is the validation error returned by
// ListEventsResponse.Validate if the designated constraints aren't met.
type ListEventsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListEventsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListEventsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListEventsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListEventsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListEventsResponseValidationError) ErrorName() string {
	return "ListEventsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListEventsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListEventsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListEventsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListEventsResponseValidationError{}

// Validate checks the field values on LatestEventsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *LatestEventsRequest) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for RuleId

	return nil
}

// LatestEventsRequestValidationError is the validation error returned by
// LatestEventsRequest.Validate if the designated constraints aren't met.
type LatestEventsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LatestEventsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LatestEventsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LatestEventsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LatestEventsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LatestEventsRequestValidationError) ErrorName() string {
	return "LatestEventsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e LatestEventsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLatestEventsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LatestEventsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LatestEventsRequestValidationError{}

// Validate checks the field values on LatestEventsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *LatestEventsResponse) Validate() error {
	if m == nil {
		return nil
	}

	for idx, item := range m.GetEvents() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return LatestEventsResponseValidationError{
					field:  fmt.Sprintf("Events[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// LatestEventsResponseValidationError is the validation error returned by
// LatestEventsResponse.Validate if the designated constraints aren't met.
type LatestEventsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LatestEventsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LatestEventsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LatestEventsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LatestEventsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LatestEventsResponseValidationError) ErrorName() string {
	return "LatestEventsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e LatestEventsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLatestEventsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LatestEventsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LatestEventsResponseValidationError{}
