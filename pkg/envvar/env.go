// Package envvar contains functions for type-safe handling and conversion of environment variables.
// Currently, the following types are supported:
//
//	int, float64, bool, string, []string, time.Time, time.Duration, net.IP, stage
//
// There are essentially two functions repeated for a variety of types..
// Get - look up the key as type T and return the fallback value if it cannot be found or parsed
//
//	Get<T>(key string, fallback T) (T, error)
//
// Examples:
//
//	retries := GetInt("RETRIES", 0) // the zero value is a good default here
//	clientTimeout := GetDuration("CLIENT_TIMEOUT", 30*time.Second)
//
// MustGet - look up the key as Type T and panic if it cannot be found or parsed
// Example:
//
//	authToken := MustGetString("AWS_AUTH_TOKEN")
package envvar

import (
	"encoding/json"
	"errors"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// lookupBool looks up and parses the specified environment variable, returning an error if it's missing or not a valid bool.
// See strconv.ParseBool for details on acceptable float strings.
func lookupBool(key string) (bool, error) {
	s, err := Lookup(key)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(s)
}

var ErrMissingKey = errors.New("missing environment variable key")
var ErrInvalidStage = errors.New(`invalid stage: expected "dev", "staging", "edge", or "prod"`)

// Lookup is as os.LookupEnv but returns an error specifying the missing key if it is not found.
func Lookup(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", ErrMissingKey
	}
	return val, nil
}

// MustGetString looks up the specified environment variable, panicking if it is missing (but NOT the empty string).
func MustGetString(key string) string {
	s, err := Lookup(key)
	if err != nil {
		logPanic(key).Caller(1).Err(err).Send()
	}
	return s
}

// GetString gets a string, falling back to the default if it is missing (but NOT if it is present and the empty string)
func GetString(key string, fallback string) string {
	s, err := Lookup(key)
	if err != nil {
		return fallback
	}
	return s
}

// GetBool gets a bool or return fallback if it doesn't exist.
// It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
func GetBool(key string, fallback bool) bool {
	b, err := lookupBool(key)
	if err != nil {
		return fallback
	}
	return b
}

// MustGetBool looks up and parses the specified environment variable as a bool, panicking if it is missing or invalid.
//
//	It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
func MustGetBool(key string) bool {
	b, err := lookupBool(key)
	if err != nil {
		logPanic(key).Err(err).Send()
	}
	return b
}

func lookupDuration(key string) (time.Duration, error) {
	s, err := Lookup(key)
	if err != nil {
		return 0, err
	}
	return time.ParseDuration(s)
}

// GetDuration looks up and parses the given environment variable as a time.Duration as though with time.ParseDuration,
// falling back to the default on missing or invalid value
func GetDuration(key string, fallback time.Duration) time.Duration {
	d, err := lookupDuration(key)
	if err != nil {
		return fallback
	}
	return d
}

// MustGetDuration looks up and parses the given environment variable as a time.Duration as though with time.ParseDuration,
// panicking on failure.
func MustGetDuration(key string) time.Duration {
	d, err := lookupDuration(key)
	if err != nil {
		logPanic(key).Caller(1).Err(err).Msg("missing or invalid environment variable: expected a time.Duration")
	}
	return d
}

// LookupFloat looks up and parses the specified environment variable, returning an error if it's missing or not a valid float.
// See strconv.ParseFloat(s, 64) for details on acceptable bool strings.
func lookupFloat(key string) (float64, error) {
	s, err := Lookup(key)
	if err != nil {
		return 0.0, err
	}
	return strconv.ParseFloat(s, 64)
}

// MustGetFloat looks up and parses the specified environment variable as a float64, panicking if it is missing or invalid.
func MustGetFloat(key string) float64 {
	x, err := lookupFloat(key)
	if err != nil {
		logPanic(key).Caller(1).Err(err).Send()
	}
	return x
}

// GetFloat gets the environment variable and  parses it as a float64, returning the fallback if it is missing or invalid.
func GetFloat(key string, fallback float64) float64 {
	x, err := lookupFloat(key)
	if err != nil {
		return fallback
	}
	return x
}

// lookupInt looks up and parses the specified environment variable, returning an error if it's not a valid int.
// See strconv.Atoi(s) for details on acceptable int strings
func lookupInt(key string) (int, error) {
	s, err := Lookup(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(s)
}

// MustGetInt looks up and parses the specified environment variable as an int64, panicking if it is missing or invalid.
func MustGetInt(key string) int {
	n, err := lookupInt(key)
	if err != nil {
		logPanic(key).Caller(1).Err(err).Send()
	}
	return n
}

// GetInt looks up and parses the specified environment variable as an int64, returning the fallback if it is missing or invalid.
func GetInt(key string, fallback int) int {
	n, err := lookupInt(key)
	if err != nil {
		return fallback
	}
	return n
}

// LookupIP looks up and parses the given environment variable as though with(net.IP).UnmarshalText
func lookupIP(key string) (ip net.IP, err error) {
	s, err := Lookup(key)
	if err != nil {
		return ip, err
	}
	err = ip.UnmarshalText([]byte(s))
	return ip, err
}

// GetIP looks up and parses the given environment variable as though with (net.IP).UnmarshalText,
// falling back to the given default if it is missing or invalid.
func GetIP(key string, fallback net.IP) (ip net.IP) {
	ip, err := lookupIP(key)
	if err != nil {
		return fallback
	}
	return ip
}

// MustGetIP looks up and parses the given environment variable as a (net.IP) as though with IP.UnmarshalText
// panicking on failure.
func MustGetIP(key string) net.IP {
	ip, err := lookupIP(key)
	if err != nil {
		logPanic(key).Caller(1).Err(err).Msg("missing or invalid environment variable: expected a net.IP")
	}
	return ip
}

// lookupTime looks up and parses the given environment variable as though with time.Parse(time.RFC3339)
func lookupTime(key string) (t time.Time, err error) {
	s, err := Lookup(key)
	if err != nil {
		return t, err
	}
	err = t.UnmarshalText([]byte(s))
	return t, err
}

// GetTime looks up and parses the given environment variable as a time.RFC3339 datetime though with time.UnmarshalText,
// falling back to the given default if it is missing or invalid.
func GetTime(key string, fallback time.Time) (t time.Time) {
	t, err := lookupTime(key)
	if err != nil {
		return fallback
	}
	return t
}

// MustGetTime looks up and parses the given environment variable as a time.RFC3339 datetime though with time.UnmarshalText,
// panicking on failure.
func MustGetTime(key string) (t time.Time) {
	t, err := lookupTime(key)
	if err != nil {
		logPanic(key).Caller(1).Err(err).Msg("missing or invalid environment variable: expected a time in RFC3339")
	}
	return t
}

// GetStage looks up and parses specifically the 'ENV' or 'STAGE' environment variables (in that order). If neither
// are set, GetStage defaults to 'dev'.
func GetStage() string {

	val, _, _ := lookupStage()

	// this should be guaranteed by lookupStage().
	switch val {
	case "dev", "staging", "edge", "prod":
		return val
	default:
		panic(ErrInvalidStage)
	}
}

// MustGetStage looks up and parses specifically the 'ENV' or 'STAGE' environment variables (in that order). If neither
// are set, MustGetStage panics.
func MustGetStage() string {
	val, envErr, stageErr := lookupStage()
	if envErr != nil && stageErr != nil {
		logPanic("ENV/STAGE").
			Caller(1).
			AnErr("STAGE_ERR", stageErr).
			AnErr("ENV_ERR", envErr).
			Msg("missing or invalid environment variable")
	}
	// this should be guaranteed by lookupStage().
	switch val {
	case "dev", "staging", "edge", "prod":
		return val
	default:
		panic(ErrInvalidStage)
	}
}

func lookupStage() (val string, envErr, stageErr error) {
	if val, ok := os.LookupEnv("ENV"); !ok {
		envErr = ErrMissingKey
	} else if val = strings.ToLower(val); val != "dev" && val != "staging" && val != "prod" && val != "edge" {
		envErr = ErrInvalidStage
	} else {
		return val, nil, nil
	}

	if val, ok := os.LookupEnv("STAGE"); !ok {
		return "dev", envErr, ErrMissingKey
	} else if val = strings.ToLower(val); val != "dev" && val != "staging" && val != "prod" && val != "edge" {
		return "dev", envErr, ErrInvalidStage
	} else {
		return val, envErr, nil
	}
}

// GetStringArray gets the environment variable and parses it as a array of string([]string), returning the fallback if it is missing or invalid.
func GetStringArray(key string, fallback []string) (arr []string) {
	a, err := Lookup(key)
	if err != nil {
		return fallback
	}

	err = json.Unmarshal([]byte(a), &arr)
	if err != nil {
		return fallback
	}

	return arr
}

// GetStringArray gets the environment variable and parses it as a array of string([]string), panicking the fallback if it is missing or invalid.
func MustGetStringArray(key string, fallback []string) (arr []string) {
	a, err := Lookup(key)
	if err != nil {
		logPanic(key).Caller(1).Err(err).Send()
	}

	err = json.Unmarshal([]byte(a), &arr)
	if err != nil {
		logPanic(key).Caller(1).Err(err).Send()
	}

	return arr
}
