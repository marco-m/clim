package clim

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"golang.org/x/exp/maps"
)

// ActionFn is the function type of the "action" returned by a parser.
type ActionFn[T any] func(ctx context.Context, uctx T) error

var (
	// User requested help.
	ErrHelp = errors.New("")
	// Either parse or validation error.
	ErrParse = errors.New("")
)

// ParseError creates an error that unwraps to [ErrParse].
// A user implementation of [Value] should use this function to return a parse
// error, so that the recommended mainInt function (see in directory examples)
// can stay generic.
func ParseError(format string, a ...any) error {
	return fmt.Errorf("%w%s", ErrParse, fmt.Sprintf(format, a...))
}

// helpError creates an error that unwraps to [ErrHelp].
// See in directory examples how to handle it.
func helpError(format string, a ...any) error {
	return fmt.Errorf("%w%s", ErrHelp, fmt.Sprintf(format, a...))
}

// CLI represents the top-level command, created with [New], and any
// subcommands, created with [CLI.AddCLI]. A [Flag] is added with [CLI.AddFlag].
type CLI[T any] struct {
	name        string
	desc        string
	long2flag   map[string]*Flag
	short2long  map[string]string
	longSeen    map[string]struct{} // Options seen on the command-line
	positionals []string
	//
	parent  string
	subCLIs []*CLI[T]
	action  ActionFn[T]
	groups  []cliGroup[T]
}

type cliGroup[T any] struct {
	name string
	clis []*CLI[T]
}

// New creates the top-level [CLI], representing the program itself,
// and sets action, to be returned by a successful parse.
// See [CLI.AddCLI] to add a sub CLI (subcommand).
func New[T any](name string, desc string, action ActionFn[T]) *CLI[T] {
	if name == "" {
		panic("clim.New: name cannot be empty")
	}
	return &CLI[T]{
		name:       name,
		desc:       desc,
		action:     action,
		long2flag:  make(map[string]*Flag),
		short2long: make(map[string]string),
		longSeen:   map[string]struct{}{},
	}
}

// AddFlag adds a [Flag] to cli.
// The type and value of the flag are represented by the field [Flag.Value],
// which holds either one of the implementation of [Value] from the clim package
// (e.g. [Int], [IntSlice], [Bool], ...) or a user-defined one.
//
// Taken from std/flag and adapted.
func (cli *CLI[T]) AddFlag(flag *Flag) {
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
		if _, found := cli.short2long[flag.Short]; found {
			panic(fmt.Sprintf("%s: short flag name %q already defined", cli.name, flag.Short))
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
	if _, found := cli.long2flag[flag.Long]; found {
		panic(fmt.Sprintf("%s: long flag name %q already defined", cli.name, flag.Long))
	}

	// A variable can be bound to only one flag.
	for k, fl := range cli.long2flag {
		if fl.Value == flag.Value {
			panic(fmt.Sprintf("long flag name %q: variable already bound to flag %q",
				flag.Long, k))
		}
	}

	flag.defValue = flag.Value.String()
	if flag.Label == "" && !isBoolValue(flag.Value) {
		flag.Label = strings.ToUpper(flag.Long)
	}

	if flag.Short != "" {
		cli.short2long[flag.Short] = flag.Long
	}
	cli.long2flag[flag.Long] = flag
}

// AddCLI adds a sub CLI, that is, a subcommand, with name and desc,
// and sets action, to be returned by a successful parse.
func (cli *CLI[T]) AddCLI(name string, desc string, action ActionFn[T]) *CLI[T] {
	subCLI := New[T](name, desc, action)
	subCLI.parent = cli.name
	cli.subCLIs = append(cli.subCLIs, subCLI)
	return subCLI
}

// AddGroup adds the subclis to the group name.
func (cli *CLI[T]) AddGroup(name string, clis ...*CLI[T]) {
	cli.groups = append(cli.groups, cliGroup[T]{name, clis})
}

// Args returns the positional arguments, if any.
// Must be called after Parse.
// WARNING will probably disappear, replaced by support for positional
// arguments parsing.
func (cmd *CLI[T]) Args() []string {
	return cmd.positionals
}

// Parse recursively processes args, calling the needed subCLI, and returns
// the associated action.
func (cli *CLI[T]) Parse(args []string) (ActionFn[T], error) {
	index := 0
	for {
		long, offset, err := cli.parseOne(args[index:])
		if err != nil {
			return nil, err
		}
		cli.longSeen[long] = struct{}{}
		if offset == 0 {
			// Arrived at the end of args or at the end of the flags.
			break
		}
		index += offset
	}

	// Are we missing any required options?
	var missing []string
	for name, flag := range cli.long2flag {
		if !flag.Required {
			continue
		}
		if _, found := cli.longSeen[name]; !found {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		slices.Sort(missing)
		return nil, ParseError("missing required options: %s",
			strings.Join(missing, ", "))
	}

	cli.positionals = args[index:]

	if len(cli.subCLIs) == 0 {
		return cli.run, nil
	}

	// If we are here, we have subcommands.

	if len(cli.positionals) == 0 {
		return nil, ParseError("expected a command")
	}
	command := cli.positionals[0]
	for _, p := range cli.subCLIs {
		if p.name == command {
			return p.Parse(cli.positionals[1:])
		}
	}

	return nil, ParseError("unrecognized command %q", command)
}

// _                           0  1              2            34  5
var flagRE = regexp.MustCompile(`^(?P<hyphens>-*)(?P<name>.*?)((=)(?P<value>.+))?$`)

// parseOne parses the first option in 'args', handling all possible cases:
// --foo bar              consumes two items in 'args'
// --foo=bar              consumes one item in 'args'
// --zoo     (boolean)    consumes one item in 'args'
// end of options, beginning of positional arguments
//
// it returns the tuple (long, number_of_items_consumed (0, 1 or 2), error).
// long is used by the caller to enforce required options.
func (cli *CLI[T]) parseOne(args []string) (string, int, error) {
	if len(args) == 0 {
		return "", 0, nil
	}
	token := args[0]
	matches := flagRE.FindStringSubmatch(token)
	if len(matches) != 6 {
		return "", 0,
			ParseError("clim internal error (regex); token: %q, matches: %q",
				token, matches)
	}
	hyphens := matches[1]
	name := matches[2]
	value := matches[5]

	// Any token with no hyphen suffix or with hyphen suffix of more than two is
	// a positional argument.
	if len(hyphens) == 0 || len(hyphens) > 2 {
		return "", 0, nil
	}

	// Special case: help
	if name == "h" || name == "help" {
		return "", 0, cli.usage()
	}

	// Now we expect either a flag (short or long) or a parse error.

	long := name
	if len(name) == 1 {
		long = cli.short2long[name]
		if long == "" {
			return "", 0, ParseError("unrecognized flag %q", token)
		}
	}
	flag := cli.long2flag[long]
	if flag == nil {
		return "", 0, ParseError("unrecognized flag %q", token)
	}

	// Was the value provided in the same token, with "=" ?
	if len(value) > 0 {
		if err := flag.Value.Set(value); err != nil {
			return "", 0, ParseError("setting %q: %s", token, err)
		}
		return long, 1, nil
	}

	if isBoolValue(flag.Value) {
		if err := flag.Value.Set("true"); err != nil {
			return "", 0, ParseError("setting %q: %s", token, err)
		}
		return long, 1, nil
	}

	if len(args) == 1 {
		return "", 0, ParseError("flag %q requires a value", token)
	}
	nextValue := args[1]
	if err := flag.Value.Set(nextValue); err != nil {
		return "", 0, ParseError("setting %q %q: %s", token, nextValue, err)
	}
	return long, 2, nil
}

func (cli *CLI[T]) usage() error {
	var bld strings.Builder

	// Calculate the max width of the first column of commands.
	maxColWidth := 0
	for _, p := range cli.subCLIs {
		fmt.Fprintf(&bld, " %s", p.name)
		maxColWidth = max(maxColWidth, bld.Len())
		bld.Reset()
	}

	parentAndMe := cli.name
	if cli.parent != "" {
		parentAndMe = cli.parent + " " + cli.name
	}
	fmt.Fprintf(&bld, "%s -- %s\n\n", parentAndMe, cli.desc)
	fmt.Fprintf(&bld, "Usage: %s ", parentAndMe)
	if len(cli.subCLIs) > 0 {
		fmt.Fprintf(&bld, "<command> ")
	}
	fmt.Fprintf(&bld, "[options]\n\n")
	if len(cli.groups) > 0 {
		fmt.Fprintf(&bld, "available commands:\n\n")
	} else if len(cli.subCLIs) > 0 {
		fmt.Fprintf(&bld, "Commands:\n\n")
	}

	const gutter = 4
	width := maxColWidth + gutter

	// Render the commands, per group.
	if len(cli.groups) > 0 {
		for _, group := range cli.groups {
			fmt.Fprintf(&bld, "%s:\n\n", group.name)
			for _, cmd := range group.clis {
				fmt.Fprintf(&bld, " %-*s%s\n", width, cmd.name, cmd.desc)
			}
			fmt.Fprintln(&bld)
		}
	} else if len(cli.subCLIs) > 0 {
		for _, cmd := range cli.subCLIs {
			fmt.Fprintf(&bld, " %-*s%s\n", width, cmd.name, cmd.desc)
		}
		fmt.Fprintln(&bld)
	}

	return helpError("%s", bld.String()+cli.usageOptions())
}

func (cli *CLI[T]) usageOptions() string {
	// First pass. Sort keys.
	longs := maps.Keys(cli.long2flag)
	slices.Sort(longs)

	// Second pass, calculate the max width of the first column.
	lines := make([]string, 0, len(longs)+1)
	var bld strings.Builder
	maxColWidth := 0
	for _, long := range longs {
		flag := cli.long2flag[long]
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
		flag := cli.long2flag[long]
		fmt.Fprintf(&bld, "%-*s%s", maxColWidth+gutter, lines[i], flag.Desc)
		if flag.defValue != "" && !flag.Required {
			fmt.Fprintf(&bld, " (default: %s)", flag.defValue)
		}
		if flag.Required {
			fmt.Fprintf(&bld, " (required)")
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

func (cli *CLI[T]) run(ctx context.Context, uctx T) error {
	if cli.action == nil {
		return ParseError("command '%s': no action registered", cli.name)
	}
	return cli.action(ctx, uctx)
}
