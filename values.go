// This file contains all the types that clim knows how to parse.
// The majority of the code is taken from the std/flag package and adapted.

package clim

import (
	"fmt"
	"strconv"
	"strings"
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

// A Flag represents the state of a flag.
// See also [CLI.AddFlag].
type Flag struct {
	Value    Value  // Final value, once parsed, mandatory.
	Short    string // Short flag, optional.
	Long     string // Long flag, mandatory.
	Label    string // Placeholder in usage message, optional.
	Desc     string // Description, optional.
	Required bool   // Optional, default false.
	//
	defValue string // Default value, for usage message. Taken from Value.
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
