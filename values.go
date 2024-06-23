// This file contains all the types that clim knows how to parse.
// The majority of the code is taken from the std/flag package and adapted.

package clim

import (
	"fmt"
	"strconv"
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
	Short    string // The short flag, optional.
	Long     string // The long flag, mandatory.
	Label    string // Placeholder in usage message.
	DefValue string // Default value, always as text; for usage message.
	Desc     string // Description.
	Value    Value  // Final value, once parsed.
}

//
// int Value
//

type intValue int

func IntVal(dst *int, defval int) *intValue {
	*dst = defval
	return (*intValue)(dst)
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		err = fmt.Errorf("could not parse %q as int (%s)", s, err)
	}
	*i = intValue(v)
	return err
}

func (i *intValue) String() string { return strconv.Itoa(int(*i)) }

//
// string Value
//

type stringValue string

func StringVal(dst *string, defval string) *stringValue {
	*dst = defval
	return (*stringValue)(dst)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) String() string { return string(*s) }

//
// bool Value
//

type boolValue bool

func BoolVal(dst *bool, defval bool) *boolValue {
	*dst = defval
	return (*boolValue)(dst)
}

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		err = fmt.Errorf("could not parse %q as bool (%s)", s, err)
	}
	*b = boolValue(v)
	return err
}

func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }

func (b *boolValue) IsBoolFlag() bool { return true }
