// This file contains all the types that clim knows how to parse.
// The majority of the code is taken from the std/flag package and adapted.

package clim

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

// Value is the interface for a type to be parsable by clim.
// (Idea taken from std/flag).
type Value interface {
	// String is called by help to print the default value.
	String() string
	// Set is called by [CLI.Parse].
	Set(string) error
}

// boolFlag is an interface to be implemented by boolean types (in addition to
// the [Value] interface), to indicate that the flag can be supplied without
// "=value" text.
// No need to type assert on this; instead, use function [isBoolValue].
// Should be called boolValue to match the Value interface, but we use
// boolValue for the concrete implementation... To be reconsidered.
type boolFlag interface {
	IsBoolFlag() bool
}

// isBoolValue returns true if 'value' implements the [boolFlag] interface.
func isBoolValue(value Value) bool {
	if x, ok := value.(boolFlag); ok {
		return x.IsBoolFlag()
	}
	return false
}

//
// int Value
//

type intValue int

// Int creates a [Value] that parses an integer into dst.
// See also [Flag] and [CLI.AddFlag].
func Int(dst *int, defval int) *intValue {
	*dst = defval
	return (*intValue)(dst)
}

// Set will be called by the parsing machinery.
func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return fmt.Errorf("could not parse %q as int (%s)", s, err)
	}
	*i = intValue(v)
	return nil
}

// String is called by help to print the default value.
func (i *intValue) String() string { return strconv.Itoa(int(*i)) }

//
// int slice Value
//

type intSliceValue []int

// IntSlice creates a [Value] that parses an comma-separated list of integers
// into dst.
// See also [Flag] and [CLI.AddFlag].
func IntSlice(dst *[]int, defval []int) *intSliceValue {
	*dst = defval
	return (*intSliceValue)(dst)
}

// Set is called by [CLI.Parse].
func (is *intSliceValue) Set(val string) error {
	// Reset any default values
	*is = make(intSliceValue, 0)
	for _, s := range strings.Split(val, ",") {
		v, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("could not parse %q as int (%s)", s, err)
		}
		*is = append(*is, v)
	}
	return nil
}

// String is called by help to print the default value.
func (is *intSliceValue) String() string {
	vals := make([]string, 0, len(*is))
	for _, i := range *is {
		vals = append(vals, strconv.Itoa(i))
	}
	return strings.Join(vals, ",")
}

//
// float64 Value
//

type float64Value float64

// Float64 creates a [Value] that parses a float into dst.
// See also [Flag] and [CLI.AddFlag].
func Float64(dst *float64, defval float64) *float64Value {
	*dst = defval
	return (*float64Value)(dst)
}

// Set will be called by the parsing machinery.
func (f *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("could not parse %q as float (%s)", s, err)
	}
	*f = float64Value(v)
	return nil
}

// String is called by help to print the default value.
func (f *float64Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 64) }

//
// string Value
//

type stringValue string

// String creates a [Value] that parses a string into dst.
// See also [Flag] and [CLI.AddFlag].
func String(dst *string, defval string) *stringValue {
	*dst = defval
	return (*stringValue)(dst)
}

// Set is called by [CLI.Parse].
func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

// String is called by help to print the default value.
func (s *stringValue) String() string { return string(*s) }

//
// string slice value
//

type stringSliceValue []string

// StringSlice creates a [Value] that parses a comma-separated list of strings
// into dst.
// See also [Flag] and [CLI.AddFlag].
func StringSlice(dst *[]string, defval []string) *stringSliceValue {
	*dst = defval
	return (*stringSliceValue)(dst)
}

// Set is called by [CLI.Parse].
func (s *stringSliceValue) Set(val string) error {
	*s = strings.Split(val, ",")
	return nil
}

// String is called by help to print the default value.
func (s *stringSliceValue) String() string { return strings.Join(*s, ",") }

//
// bool Value
//

type boolValue bool

// Bool creates a [Value] that parses a boolean into dst.
// See also [Flag] and [CLI.AddFlag].
func Bool(dst *bool, defval bool) *boolValue {
	*dst = defval
	return (*boolValue)(dst)
}

// Set is called by [CLI.Parse].
func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return fmt.Errorf("could not parse %q as bool (%s)", s, err)
	}
	*b = boolValue(v)
	return nil
}

// String is called by help to print the default value.
func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }

func (b *boolValue) IsBoolFlag() bool { return true }

//
// time.Duration Value
//

type durationValue time.Duration

// Duration creates a [Value] that parses a time.Duration into dst.
// See also [Flag] and [CLI.AddFlag].
func Duration(dst *time.Duration, defval time.Duration) *durationValue {
	*dst = defval
	return (*durationValue)(dst)
}

// Set will be called by the parsing machinery.
func (ll *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*ll = durationValue(v)
	return nil
}

// String is called by help to print the default value.
func (ll *durationValue) String() string { return time.Duration(*ll).String() }

//
// slog.Level Value
//

type logLevelValue slog.Level

// LogLevel creates a [Value] that parses a slog.Level into dst.
// See also [Flag] and [CLI.AddFlag].
func LogLevel(dst *slog.Level, defval slog.Level) *logLevelValue {
	*dst = defval
	return (*logLevelValue)(dst)
}

// Set will be called by the parsing machinery.
func (ll *logLevelValue) Set(s string) error {
	var logLevel slog.Level
	err := logLevel.UnmarshalText([]byte(s))
	if err != nil {
		return fmt.Errorf("could not parse %q as slog.Level (%s)", s, err)
	}
	*ll = logLevelValue(logLevel)
	return nil
}

// String is called by help to print the default value.
func (ll *logLevelValue) String() string { return slog.Level(*ll).String() }
