package clim

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"golang.org/x/exp/maps"
)

var (
	ErrHelp  = errors.New("")
	ErrParse = errors.New("")
)

// parseError returns an error that unwraps to ErrParse.
func parseError(format string, a ...any) error {
	return fmt.Errorf("%w%s", ErrParse, fmt.Sprintf(format, a...))
}

// helpError returns an error that unwraps to ErrHelp.
func helpError(format string, a ...any) error {
	return fmt.Errorf("%w%s", ErrHelp, fmt.Sprintf(format, a...))
}

type Command struct {
	name        string
	desc        string
	long2flag   map[string]*Flag
	short2long  map[string]string
	positionals []string
}

func New(name string, desc string) *Command {
	if name == "" {
		panic("clim.New: name cannot be empty")
	}
	return &Command{
		name:       name,
		desc:       desc,
		long2flag:  make(map[string]*Flag),
		short2long: make(map[string]string),
	}
}

// AddFlag defines a flag with the specified short and long names and description.
// The type and value of the flag are represented by the first argument, of
// type [Value], which typically holds a user-defined implementation of [Value].
// For instance, the caller could create a flag that turns a comma-separated
// string into a slice of strings by giving the slice the methods of [Value]; in
// particular, [Set] would decompose the comma-separated string into the slice.
//
// Taken from std/flag and adapted.
func (cmd *Command) AddFlag(value Value, short, long, label, desc string) {
	//
	// Validate the short flag.
	//
	if short != "" {
		if strings.HasPrefix(short, "-") {
			panic("short flag name must not begin with '-'")
		}
		if strings.Contains(short, "=") {
			panic("short flag name must not contain '='")
		}

		// short can be empty.

		if len(short) > 1 {
			panic(fmt.Sprintf("short flag name %q must be exactly 1 character", short))
		}
		if short == "h" {
			panic(`cannot override short flag name "h"`)
		}
		if _, found := cmd.short2long[short]; found {
			panic(fmt.Sprintf("%s: short flag name %q already defined", cmd.name, short))
		}
	}

	//
	// Validate the long flag.
	//
	if strings.HasPrefix(long, "-") {
		panic(fmt.Sprintf("long flag name %q must not begin with '-'", long))
	}
	if strings.Contains(long, "=") {
		panic(fmt.Sprintf("long flag name %q must not contain '='", long))
	}
	if long == "" {
		panic("long flag name cannot be empty")
	}
	if len(long) < 3 {
		panic(fmt.Sprintf("long flag name %q must be at least 3 character", long))
	}
	if long == "help" {
		panic(`cannot override long flag name "help"`)
	}
	if _, found := cmd.long2flag[long]; found {
		panic(fmt.Sprintf("%s: long flag name %q already defined", cmd.name, long))
	}

	// A variable can be bound to only one flag.
	for k, flag := range cmd.long2flag {
		if flag.Value == value {
			panic(fmt.Sprintf("long flag name %q: variable already bound to flag %q",
				long, k))
		}
	}

	flag := &Flag{
		Short:    short,
		Long:     long,
		Label:    label,
		DefValue: value.String(),
		Desc:     desc,
		Value:    value,
	}

	if short != "" {
		cmd.short2long[short] = long
	}
	cmd.long2flag[long] = flag
}

func (cmd *Command) IntVar(dst *int, short, long, label string,
	defval int, desc string) {
	cmd.AddFlag(newIntValue(defval, dst), short, long, label, desc)
}

func (cmd *Command) StringVar(dst *string, short, long, label string,
	defval string, desc string) {
	cmd.AddFlag(newStringValue(defval, dst), short, long, label, desc)
}

func (cmd *Command) BoolVar(dst *bool, short, long string,
	defval bool, desc string) {
	cmd.AddFlag(newBoolValue(defval, dst), short, long, "", desc)
}

// Positionals returns the positional arguments, if any.
// Must be called after Parse.
func (cmd *Command) Positionals() []string {
	return cmd.positionals
}

func (cmd *Command) Parse(args []string) error {
	index := 0
	for {
		offset, err := cmd.parseOne(args[index:])
		if err != nil {
			return err
		}
		if offset == 0 {
			// Arrived at the end of args or at the end of the flags.
			cmd.positionals = args[index:]
			return nil
		}
		index += offset
	}
}

// _                           0  1              2            34  5
var flagRE = regexp.MustCompile(`^(?P<hyphens>-*)(?P<name>.*?)((=)(?P<value>.+))?$`)

func (cmd *Command) parseOne(args []string) (int, error) {
	if len(args) == 0 {
		return 0, nil
	}
	token := args[0]
	matches := flagRE.FindStringSubmatch(token)
	if len(matches) != 6 {
		return 0,
			parseError("clim internal error (regex); token: %q, matches: %q",
				token, matches)
	}
	hyphens := matches[1]
	name := matches[2]
	value := matches[5]

	// Any token with no hyphen suffix or with hyphen suffix of more than two is
	// a positional argument.
	if len(hyphens) == 0 || len(hyphens) > 2 {
		return 0, nil
	}

	// Special case: help
	if name == "h" || name == "help" {
		return 0, cmd.usage()
	}

	// Now we expect either a flag (short or long) or a parse error.

	long := name
	if len(name) == 1 {
		long = cmd.short2long[name]
		if long == "" {
			return 0, parseError("unrecognized flag %q", token)
		}
	}
	flag := cmd.long2flag[long]
	if flag == nil {
		return 0, parseError("unrecognized flag %q", token)
	}

	// Was the value provided in the same token, with "=" ?
	if len(value) > 0 {
		if err := flag.Value.Set(value); err != nil {
			return 0, parseError("setting %q: %s", token, err)
		}
		return 1, nil
	}

	if IsBoolValue(flag.Value) {
		if err := flag.Value.Set("true"); err != nil {
			return 0, parseError("setting %q: %s", token, err)
		}
		return 1, nil
	}

	if len(args) == 1 {
		return 0, parseError("flag %q requires a value", token)
	}
	nextValue := args[1]
	if err := flag.Value.Set(nextValue); err != nil {
		return 0, parseError("setting %q %q: %s", token, nextValue, err)
	}
	return 2, nil
}

func (cmd *Command) usage() error {
	// First pass. Sort keys.
	longs := maps.Keys(cmd.long2flag)
	slices.Sort(longs)

	// Second pass, calculate the max width of the first column.
	lines := make([]string, 0, len(longs)+1)
	var bld strings.Builder
	maxColWidth := 0
	for _, long := range longs {
		flag := cmd.long2flag[long]
		if flag.Short != "" {
			fmt.Fprintf(&bld, "-%s, ", flag.Short)
		}
		fmt.Fprintf(&bld, "--%s %s", flag.Long, flag.Label)
		lines = append(lines, bld.String())
		maxColWidth = max(maxColWidth, bld.Len())
		bld.Reset()
	}

	// Third pass, add the second column.
	const gutter = 4
	fmt.Fprintf(&bld, "%s -- %s\n\n", cmd.name, cmd.desc)
	fmt.Fprintf(&bld, "Usage: %s [options]\n\n", cmd.name)
	fmt.Fprintf(&bld, "Options:\n")
	for i, long := range longs {
		flag := cmd.long2flag[long]
		fmt.Fprintf(&bld, "%-*s%s (default: %s)\n", maxColWidth+gutter,
			lines[i], flag.Desc, flag.DefValue)
	}
	fmt.Fprintf(&bld, "\n%-*s%s", maxColWidth+gutter,
		"-h, --help", "Print this help and exit")
	return helpError("%s", bld.String())
}
