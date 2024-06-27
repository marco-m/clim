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
	String() string
	Set(string) error
}

// boolFlag is an interface to be implemented by boolean types (in addition to
// the [Value] interface), to indicate that the flag can be supplied without
// "=value" text.
// No need to type assert on this; instead, use function [IsBoolValue].
// Should be called boolValue to match the Value interface, but we use
// boolValue for the concrete implementation... To be reconsidered.
type boolFlag interface {
	IsBoolFlag() bool
}

// IsBoolValue returns true if 'value' implements the [boolFlag] interface.
func IsBoolValue(value Value) bool {
	if x, ok := value.(boolFlag); ok {
		return x.IsBoolFlag()
	}
	return false
}

// A Flag represents the state of a flag.
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

func Int(dst *int, defval int) *intValue {
	*dst = defval
	return (*intValue)(dst)
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return fmt.Errorf("could not parse %q as int (%s)", s, err)
	}
	*i = intValue(v)
	return nil
}

func (i *intValue) String() string { return strconv.Itoa(int(*i)) }

//
// int slice Value
//

type intSliceValue []int

// parse a comma-separated list of integers
func IntSlice(dst *[]int, defval []int) *intSliceValue {
	*dst = defval
	return (*intSliceValue)(dst)
}

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

func String(dst *string, defval string) *stringValue {
	*dst = defval
	return (*stringValue)(dst)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) String() string { return string(*s) }

//
// string slice value
//

type stringSliceValue []string

// parse a comma-separated list of strings
func StringSlice(dst *[]string, defval []string) *stringSliceValue {
	*dst = defval
	return (*stringSliceValue)(dst)
}

func (s *stringSliceValue) Set(val string) error {
	*s = strings.Split(val, ",")
	return nil
}

func (s *stringSliceValue) String() string { return strings.Join(*s, ",") }

//
// bool Value
//

type boolValue bool

func Bool(dst *bool, defval bool) *boolValue {
	*dst = defval
	return (*boolValue)(dst)
}

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return fmt.Errorf("could not parse %q as bool (%s)", s, err)
	}
	*b = boolValue(v)
	return nil
}

func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }

func (b *boolValue) IsBoolFlag() bool { return true }
