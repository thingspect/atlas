// Package config provides functions to parse fields of configuration structs.
package config

import (
	"encoding/base64"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// String retrieves the environment variable by its key. If the variable is
// present, it is returned as a string, otherwise the default value is returned.
func String(key string, defVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return defVal
}

// StringSlice retrieves the environment variable by its comma-delimited key. If
// the variable is present, it is returned as a string slice, otherwise the
// default value is returned.
func StringSlice(key string, defVal []string) []string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.Split(val, ",")
	}

	return defVal
}

// Int retrieves the environment variable by its key. If the variable is
// present, it is returned as an int, otherwise the default value is returned.
func Int(key string, defVal int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defVal
	}

	valInt, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("Int strconv.Atoi key, err: %v, %v", key, err)
	}

	return valInt
}

// Bool retrieves the environment variable by its key. If the variable is
// present, it is returned as a bool, otherwise the default value is returned.
func Bool(key string, defVal bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defVal
	}

	valBool, err := strconv.ParseBool(val)
	if err != nil {
		log.Fatalf("Bool strconv.ParseBool key, err: %v, %v", key, err)
	}

	return valBool
}

// Duration retrieves the environment variable by its key. If the variable is
// present, it is returned as a duration, otherwise the default value is
// returned.
func Duration(key string, defVal time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defVal
	}

	valDur, err := time.ParseDuration(val)
	if err != nil {
		log.Fatalf("Duration time.ParseDuration key, err: %v, %v", key, err)
	}

	return valDur
}

// ByteSlice retrieves the environment variable by its key, using standard
// (padded) base64 encoding. If the variable is present, it is returned as a
// byte slice. Due to use by encryption keys, no default argument is offered,
// and an empty slice (which will be rejected by length checks) is returned if
// the key is not present.
func ByteSlice(key string) []byte {
	val, ok := os.LookupEnv(key)
	if !ok {
		return []byte{}
	}

	valByte, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		log.Fatalf("ByteSlice base64.StdEncoding.DecodeString key, err: %v, %v",
			key, err)
	}

	return valByte
}
