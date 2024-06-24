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
	// User requested help.
	ErrHelp = errors.New("")
	// Either parse or validation error.
	ErrParse = errors.New("")
)

// parseError returns an error that unwraps to [ErrParse].
func parseError(format string, a ...any) error {
	return fmt.Errorf("%w%s", ErrParse, fmt.Sprintf(format, a...))
}

// helpError returns an error that unwraps to [ErrHelp].
func helpError(format string, a ...any) error {
	return fmt.Errorf("%w%s", ErrHelp, fmt.Sprintf(format, a...))
}

type Command struct {
	name        string
	desc        string
	long2flag   map[string]*Flag
	short2long  map[string]string
	positionals []string
	//
	parent  string
	parsers []*Command
	action  func() error
	groups  []group
}

type group struct {
	name     string
	commands []*Command
}

// New creates the top-level command, representing the program itself.
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
func (cmd *Command) AddFlag(flag *Flag) {
	//
	// Validate the short flag.
	//
	if flag.Short != "" {
		if strings.HasPrefix(flag.Short, "-") {
			panic("short flag name must not begin with '-'")
		}
		if strings.Contains(flag.Short, "=") {
			panic("short flag name must not contain '='")
		}

		// short can be empty.

		if len(flag.Short) > 1 {
			panic(fmt.Sprintf("short flag name %q must be exactly 1 character", flag.Short))
		}
		if flag.Short == "h" {
			panic(`cannot override short flag name "h"`)
		}
		if _, found := cmd.short2long[flag.Short]; found {
			panic(fmt.Sprintf("%s: short flag name %q already defined", cmd.name, flag.Short))
		}
	}

	//
	// Validate the long flag.
	//
	if strings.HasPrefix(flag.Long, "-") {
		panic(fmt.Sprintf("long flag name %q must not begin with '-'", flag.Long))
	}
	if strings.Contains(flag.Long, "=") {
		panic(fmt.Sprintf("long flag name %q must not contain '='", flag.Long))
	}
	if flag.Long == "" {
		panic("long flag name cannot be empty")
	}
	if len(flag.Long) < 2 {
		panic(fmt.Sprintf("long flag name %q must be at least 2 character", flag.Long))
	}
	if flag.Long == "help" {
		panic(`cannot override long flag name "help"`)
	}
	if _, found := cmd.long2flag[flag.Long]; found {
		panic(fmt.Sprintf("%s: long flag name %q already defined", cmd.name, flag.Long))
	}

	// A variable can be bound to only one flag.
	for k, fl := range cmd.long2flag {
		if fl.Value == flag.Value {
			panic(fmt.Sprintf("long flag name %q: variable already bound to flag %q",
				flag.Long, k))
		}
	}

	flag.defValue = flag.Value.String()
	if flag.Label == "" && !IsBoolValue(flag.Value) {
		flag.Label = strings.ToUpper(flag.Long)
	}

	if flag.Short != "" {
		cmd.short2long[flag.Short] = flag.Long
	}
	cmd.long2flag[flag.Long] = flag
}

// AddCommand adds subcommand 'name'.
func (cmd *Command) AddParser(name string, desc string) *Command {
	parser := New(name, desc)
	parser.parent = cmd.name
	cmd.parsers = append(cmd.parsers, parser)
	return parser
}

func (cmd *Command) Group(name string, commands ...*Command) {
	cmd.groups = append(cmd.groups, group{name, commands})
}

// Args returns the positional arguments, if any.
// Must be called after Parse.
// WARNING will probably disappear, replaced by support for positional
// arguments parsing.
func (cmd *Command) Args() []string {
	return cmd.positionals
}

func (cmd *Command) Parse(args []string) (func() error, error) {
	index := 0
	for {
		offset, err := cmd.parseOne(args[index:])
		if err != nil {
			return nil, err
		}
		if offset == 0 {
			// Arrived at the end of args or at the end of the flags.
			break
		}
		index += offset
	}

	cmd.positionals = args[index:]

	if len(cmd.parsers) == 0 {
		return cmd.run, nil
	}
	// If we are here, we have subcommands.

	if len(cmd.positionals) == 0 {
		return nil, parseError("expected a command")
	}
	command := cmd.positionals[0]
	for _, p := range cmd.parsers {
		if p.name == command {
			return p.Parse(cmd.positionals[1:])
		}
	}

	return nil, parseError("unrecognized command %q", command)
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
	var bld strings.Builder

	// Calculate the max width of the first column of commands.
	maxColWidth := 0
	for _, p := range cmd.parsers {
		fmt.Fprintf(&bld, " %s", p.name)
		maxColWidth = max(maxColWidth, bld.Len())
		bld.Reset()
	}

	parentAndMe := cmd.name
	if cmd.parent != "" {
		parentAndMe = cmd.parent + " " + cmd.name
	}
	fmt.Fprintf(&bld, "%s -- %s\n\n", parentAndMe, cmd.desc)
	fmt.Fprintf(&bld, "Usage: %s ", parentAndMe)
	if len(cmd.parsers) > 0 {
		fmt.Fprintf(&bld, "<command> ")
	}
	fmt.Fprintf(&bld, "[options]\n\n")
	if len(cmd.parsers) > 0 {
		fmt.Fprintf(&bld, "available commands:\n\n")
	}

	// Render the commands, per group.
	const gutter = 4
	width := maxColWidth + gutter
	for _, group := range cmd.groups {
		fmt.Fprintf(&bld, "%s:\n\n", group.name)
		for _, cmd := range group.commands {
			fmt.Fprintf(&bld, " %-*s%s\n", width, cmd.name, cmd.desc)
		}
		fmt.Fprintln(&bld)
	}

	return helpError("%s", bld.String()+cmd.usageOptions())
}

func (cmd *Command) usageOptions() string {
	// First pass. Sort keys.
	longs := maps.Keys(cmd.long2flag)
	slices.Sort(longs)

	// Second pass, calculate the max width of the first column.
	lines := make([]string, 0, len(longs)+1)
	var bld strings.Builder
	maxColWidth := 0
	for _, long := range longs {
		flag := cmd.long2flag[long]
		fmt.Fprintf(&bld, " ")
		if flag.Short != "" {
			fmt.Fprintf(&bld, "-%s, ", flag.Short)
		}
		fmt.Fprintf(&bld, "--%s %s", flag.Long, flag.Label)
		lines = append(lines, bld.String())
		maxColWidth = max(maxColWidth, bld.Len())
		bld.Reset()
	}
	// Same for -h
	fmt.Fprint(&bld, " -h, --help")
	maxColWidth = max(maxColWidth, bld.Len())
	bld.Reset()

	// Third pass, add the second column.
	const gutter = 4
	fmt.Fprintf(&bld, "Options:\n\n")
	for i, long := range longs {
		flag := cmd.long2flag[long]
		fmt.Fprintf(&bld, "%-*s%s", maxColWidth+gutter, lines[i], flag.Desc)
		if flag.defValue != "" {
			fmt.Fprintf(&bld, " (default: %s)", flag.defValue)
		}
		fmt.Fprintf(&bld, "\n")
	}
	if len(longs) > 0 {
		fmt.Fprintf(&bld, "\n")
	}

	fmt.Fprintf(&bld, "%-*s%s", maxColWidth+gutter,
		" -h, --help", "Print this help and exit")
	return bld.String()
}

func (cmd *Command) Action(fn func() error) {
	cmd.action = fn
}

func (cmd *Command) run() error {
	if cmd.action == nil {
		return parseError("command %q: no action registered", cmd.name)
	}
	return cmd.action()
}
