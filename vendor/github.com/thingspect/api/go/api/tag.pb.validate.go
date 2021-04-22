// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/tag.proto

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
)

// Validate checks the field values on ListTagsRequest with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned. When asked to return all errors, validation continues
// after first violation, and the result is a list of violation errors wrapped
// in ListTagsRequestMultiError, or nil if none found. Otherwise, only the
// first error is returned, if any.
func (m *ListTagsRequest) Validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return ListTagsRequestMultiError(errors)
	}
	return nil
}

// ListTagsRequestMultiError is an error wrapping multiple validation errors
// returned by ListTagsRequest.Validate(true) if the designated constraints
// aren't met.
type ListTagsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListTagsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListTagsRequestMultiError) AllErrors() []error { return m }

// ListTagsRequestValidationError is the validation error returned by
// ListTagsRequest.Validate if the designated constraints aren't met.
type ListTagsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListTagsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListTagsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListTagsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListTagsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListTagsRequestValidationError) ErrorName() string { return "ListTagsRequestValidationError" }

// Error satisfies the builtin error interface
func (e ListTagsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListTagsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListTagsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListTagsRequestValidationError{}

// Validate checks the field values on ListTagsResponse with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned. When asked to return all errors, validation continues
// after first violation, and the result is a list of violation errors wrapped
// in ListTagsResponseMultiError, or nil if none found. Otherwise, only the
// first error is returned, if any.
func (m *ListTagsResponse) Validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return ListTagsResponseMultiError(errors)
	}
	return nil
}

// ListTagsResponseMultiError is an error wrapping multiple validation errors
// returned by ListTagsResponse.Validate(true) if the designated constraints
// aren't met.
type ListTagsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListTagsResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListTagsResponseMultiError) AllErrors() []error { return m }

// ListTagsResponseValidationError is the validation error returned by
// ListTagsResponse.Validate if the designated constraints aren't met.
type ListTagsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListTagsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListTagsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListTagsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListTagsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListTagsResponseValidationError) ErrorName() string { return "ListTagsResponseValidationError" }

// Error satisfies the builtin error interface
func (e ListTagsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListTagsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListTagsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListTagsResponseValidationError{}
