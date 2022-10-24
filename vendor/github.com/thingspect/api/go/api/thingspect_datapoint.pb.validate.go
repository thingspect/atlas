// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/thingspect_datapoint.proto

package api

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
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
	_ = anypb.Any{}
	_ = sort.Sort
)

// define the regex for a UUID once up-front
var _thingspect_datapoint_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on PublishDataPointsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *PublishDataPointsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PublishDataPointsRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// PublishDataPointsRequestMultiError, or nil if none found.
func (m *PublishDataPointsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *PublishDataPointsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(m.GetPoints()) < 1 {
		err := PublishDataPointsRequestValidationError{
			field:  "Points",
			reason: "value must contain at least 1 item(s)",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetPoints() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, PublishDataPointsRequestValidationError{
						field:  fmt.Sprintf("Points[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, PublishDataPointsRequestValidationError{
						field:  fmt.Sprintf("Points[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PublishDataPointsRequestValidationError{
					field:  fmt.Sprintf("Points[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return PublishDataPointsRequestMultiError(errors)
	}

	return nil
}

// PublishDataPointsRequestMultiError is an error wrapping multiple validation
// errors returned by PublishDataPointsRequest.ValidateAll() if the designated
// constraints aren't met.
type PublishDataPointsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PublishDataPointsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PublishDataPointsRequestMultiError) AllErrors() []error { return m }

// PublishDataPointsRequestValidationError is the validation error returned by
// PublishDataPointsRequest.Validate if the designated constraints aren't met.
type PublishDataPointsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PublishDataPointsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PublishDataPointsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PublishDataPointsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PublishDataPointsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PublishDataPointsRequestValidationError) ErrorName() string {
	return "PublishDataPointsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e PublishDataPointsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPublishDataPointsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PublishDataPointsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PublishDataPointsRequestValidationError{}

// Validate checks the field values on ListDataPointsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListDataPointsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListDataPointsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListDataPointsRequestMultiError, or nil if none found.
func (m *ListDataPointsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ListDataPointsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetAttr()) > 40 {
		err := ListDataPointsRequestValidationError{
			field:  "Attr",
			reason: "value length must be at most 40 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetEndTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ListDataPointsRequestValidationError{
					field:  "EndTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ListDataPointsRequestValidationError{
					field:  "EndTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetEndTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ListDataPointsRequestValidationError{
				field:  "EndTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetStartTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ListDataPointsRequestValidationError{
					field:  "StartTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ListDataPointsRequestValidationError{
					field:  "StartTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetStartTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ListDataPointsRequestValidationError{
				field:  "StartTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	switch m.IdOneof.(type) {

	case *ListDataPointsRequest_UniqId:
		// no validation rules for UniqId

	case *ListDataPointsRequest_DeviceId:

		if m.GetDeviceId() != "" {

			if err := m._validateUuid(m.GetDeviceId()); err != nil {
				err = ListDataPointsRequestValidationError{
					field:  "DeviceId",
					reason: "value must be a valid UUID",
					cause:  err,
				}
				if !all {
					return err
				}
				errors = append(errors, err)
			}

		}

	default:
		err := ListDataPointsRequestValidationError{
			field:  "IdOneof",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)

	}

	if len(errors) > 0 {
		return ListDataPointsRequestMultiError(errors)
	}

	return nil
}

func (m *ListDataPointsRequest) _validateUuid(uuid string) error {
	if matched := _thingspect_datapoint_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// ListDataPointsRequestMultiError is an error wrapping multiple validation
// errors returned by ListDataPointsRequest.ValidateAll() if the designated
// constraints aren't met.
type ListDataPointsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListDataPointsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListDataPointsRequestMultiError) AllErrors() []error { return m }

// ListDataPointsRequestValidationError is the validation error returned by
// ListDataPointsRequest.Validate if the designated constraints aren't met.
type ListDataPointsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListDataPointsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListDataPointsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListDataPointsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListDataPointsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListDataPointsRequestValidationError) ErrorName() string {
	return "ListDataPointsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListDataPointsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListDataPointsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListDataPointsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListDataPointsRequestValidationError{}

// Validate checks the field values on ListDataPointsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListDataPointsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListDataPointsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListDataPointsResponseMultiError, or nil if none found.
func (m *ListDataPointsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ListDataPointsResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetPoints() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListDataPointsResponseValidationError{
						field:  fmt.Sprintf("Points[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListDataPointsResponseValidationError{
						field:  fmt.Sprintf("Points[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListDataPointsResponseValidationError{
					field:  fmt.Sprintf("Points[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ListDataPointsResponseMultiError(errors)
	}

	return nil
}

// ListDataPointsResponseMultiError is an error wrapping multiple validation
// errors returned by ListDataPointsResponse.ValidateAll() if the designated
// constraints aren't met.
type ListDataPointsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListDataPointsResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListDataPointsResponseMultiError) AllErrors() []error { return m }

// ListDataPointsResponseValidationError is the validation error returned by
// ListDataPointsResponse.Validate if the designated constraints aren't met.
type ListDataPointsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListDataPointsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListDataPointsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListDataPointsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListDataPointsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListDataPointsResponseValidationError) ErrorName() string {
	return "ListDataPointsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListDataPointsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListDataPointsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListDataPointsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListDataPointsResponseValidationError{}

// Validate checks the field values on LatestDataPointsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *LatestDataPointsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LatestDataPointsRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// LatestDataPointsRequestMultiError, or nil if none found.
func (m *LatestDataPointsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *LatestDataPointsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetStartTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, LatestDataPointsRequestValidationError{
					field:  "StartTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, LatestDataPointsRequestValidationError{
					field:  "StartTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetStartTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return LatestDataPointsRequestValidationError{
				field:  "StartTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	switch m.IdOneof.(type) {

	case *LatestDataPointsRequest_UniqId:
		// no validation rules for UniqId

	case *LatestDataPointsRequest_DeviceId:

		if m.GetDeviceId() != "" {

			if err := m._validateUuid(m.GetDeviceId()); err != nil {
				err = LatestDataPointsRequestValidationError{
					field:  "DeviceId",
					reason: "value must be a valid UUID",
					cause:  err,
				}
				if !all {
					return err
				}
				errors = append(errors, err)
			}

		}

	default:
		err := LatestDataPointsRequestValidationError{
			field:  "IdOneof",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)

	}

	if len(errors) > 0 {
		return LatestDataPointsRequestMultiError(errors)
	}

	return nil
}

func (m *LatestDataPointsRequest) _validateUuid(uuid string) error {
	if matched := _thingspect_datapoint_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// LatestDataPointsRequestMultiError is an error wrapping multiple validation
// errors returned by LatestDataPointsRequest.ValidateAll() if the designated
// constraints aren't met.
type LatestDataPointsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LatestDataPointsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LatestDataPointsRequestMultiError) AllErrors() []error { return m }

// LatestDataPointsRequestValidationError is the validation error returned by
// LatestDataPointsRequest.Validate if the designated constraints aren't met.
type LatestDataPointsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LatestDataPointsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LatestDataPointsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LatestDataPointsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LatestDataPointsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LatestDataPointsRequestValidationError) ErrorName() string {
	return "LatestDataPointsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e LatestDataPointsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLatestDataPointsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LatestDataPointsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LatestDataPointsRequestValidationError{}

// Validate checks the field values on LatestDataPointsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *LatestDataPointsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LatestDataPointsResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// LatestDataPointsResponseMultiError, or nil if none found.
func (m *LatestDataPointsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *LatestDataPointsResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetPoints() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, LatestDataPointsResponseValidationError{
						field:  fmt.Sprintf("Points[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, LatestDataPointsResponseValidationError{
						field:  fmt.Sprintf("Points[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return LatestDataPointsResponseValidationError{
					field:  fmt.Sprintf("Points[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return LatestDataPointsResponseMultiError(errors)
	}

	return nil
}

// LatestDataPointsResponseMultiError is an error wrapping multiple validation
// errors returned by LatestDataPointsResponse.ValidateAll() if the designated
// constraints aren't met.
type LatestDataPointsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LatestDataPointsResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LatestDataPointsResponseMultiError) AllErrors() []error { return m }

// LatestDataPointsResponseValidationError is the validation error returned by
// LatestDataPointsResponse.Validate if the designated constraints aren't met.
type LatestDataPointsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LatestDataPointsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LatestDataPointsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LatestDataPointsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LatestDataPointsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LatestDataPointsResponseValidationError) ErrorName() string {
	return "LatestDataPointsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e LatestDataPointsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLatestDataPointsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LatestDataPointsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LatestDataPointsResponseValidationError{}
