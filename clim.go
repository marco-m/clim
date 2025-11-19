package clim

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

var (
	// ErrHelp is returned when the user asked for help.
	// Check for this case with errors.Is(err, clim.ErrHelp)
	ErrHelp = errors.New("")
	// ErrParse is returned in case of parse (by clim) or validation error
	// (by the user program).
	// Check for this case with errors.Is(err, clim.ErrParse).
	// See also NewParseError.
	ErrParse = errors.New("")
)

// NewParseError creates an error that unwraps to [ErrParse].
// A user implementation of [Value] should use this function to return a parse
// error, so that the recommended mainInt function (see in directory examples)
// can stay generic.
func NewParseError(format string, a ...any) error {
	return fmt.Errorf("%w%s", ErrParse, fmt.Sprintf(format, a...))
}

// newHelpError creates an error that unwraps to [ErrHelp].
// See in directory examples how to handle it.
func newHelpError(format string, a ...any) error {
	return fmt.Errorf("%w%s", ErrHelp, fmt.Sprintf(format, a...))
}

// CLI represents the top-level command, created with [New], and any
// subcommands, created with [CLI.addCLI]. A [Flag] is added with [CLI.AddFlag].
type CLI struct {
	name         string
	oneline      string
	description  string
	examples     string
	footer       string
	orderedFlags []string // The flags in the order the user added them
	long2flag    map[string]*Flag
	short2long   map[string]string
	posargs      *[]string
	pairs        []Pair
	longSeen     map[string]struct{} // Options seen on the command-line
	positionals  []string
	//
	parent     *CLI
	rootToHere string
	subCLIs    []*CLI
	groups     []cliGroup
}

type cliGroup struct {
	name string
	clis []*CLI
}

// NewTop creates the top-level [CLI], representing the program itself.
// Parameter 'name' is the program name; parameter 'oneline' is the one-line
// description.
// To add a subcommand, see [NewSub].
func NewTop(name string, oneline string) (*CLI, error) {
	if name == "" {
		return nil, NewParseError("cli name cannot be empty")
	}
	return newCli(nil, name, oneline), nil
}

// NewSub creates a subcommand and adds it to the 'parent' node, which itself
// could be the top command (created by [NewTop]) or an intermediate subcommand.
// Parameter 'name' is the name of the subcommand; parameter 'oneline' is the
// one-line description.
func NewSub(parent *CLI, name string, oneline string) (*CLI, error) {
	if name == "" {
		return nil, NewParseError("cli name cannot be empty")
	}
	if parent == nil {
		return nil, NewParseError("parent cli cannot be nil")
	}
	child := newCli(parent, name, oneline)
	if parent.posargs != nil {
		return nil,
			NewParseError(
				"%s: already have pos args; cannot have also subcommand %q",
				parent.name, child.name)
	}
	for _, sc := range parent.subCLIs {
		if child.name == sc.name {
			// `banana: long flag name "count" already defined`
			return nil,
				NewParseError("%s: subcommand %q already defined",
					parent.rootToHere, child.name)
		}
	}
	parent.subCLIs = append(parent.subCLIs, child)
	return child, nil
}

func newCli(parent *CLI, name string, oneline string) *CLI {
	child := &CLI{
		parent:     parent,
		name:       name,
		oneline:    oneline,
		long2flag:  make(map[string]*Flag),
		short2long: make(map[string]string),
		// name2posarg: make(map[string]*PosArg),
		longSeen: map[string]struct{}{},
	}
	child.rootToHere = strings.Join(pathRootToNode(child), " ")
	return child
}

func (cli *CLI) Command() string {
	return cli.rootToHere
}

func (cli *CLI) SetDescription(desc string) {
	cli.description = strings.TrimSpace(desc)
}

func (cli *CLI) SetExamples(examples string) {
	cli.examples = strings.TrimSpace(examples)
}

func (cli *CLI) SetFooter(footer string) {
	cli.footer = strings.TrimSpace(footer)
}

// A Flag represents the state of a flag.
// See also [CLI.AddFlag].
type Flag struct {
	Value    Value  // Final value, once parsed, mandatory.
	Short    string // Short flag, optional.
	Long     string // Long flag, mandatory.
	Label    string // Placeholder in usage message, optional.
	Help     string // Help text, optional.
	Required bool   // Optional, default false.
	//
	defValue string // Default value, for usage message. Taken from Value.
}

// AddFlags adds 'flags' to cli.
// The type and value of each flag are represented by the field [Flag.Value],
// which holds either one of the implementation of [Value] from the clim package
// (e.g. [Int], [IntSlice], [Bool], ...) or a user-defined one.
//
// Taken from std/flag and adapted.
func (cli *CLI) AddFlags(flags ...*Flag) error {
	for _, f := range flags {
		if err := cli.addFlag(f); err != nil {
			return err
		}
	}
	return nil
}

func (cli *CLI) addFlag(flag *Flag) error {
	//
	// Validate the short flag.
	//
	if flag.Short != "" {
		if strings.HasPrefix(flag.Short, "-") {
			return NewParseError("short flag name must not begin with '-'")
		}
		if strings.Contains(flag.Short, "=") {
			return NewParseError("short flag name must not contain '='")
		}

		// short can be empty.

		if len(flag.Short) > 1 {
			return NewParseError("short flag name %q must be exactly 1 character",
				flag.Short)
		}
		if flag.Short == "h" {
			return NewParseError(`cannot override short flag name "h"`)
		}
		if _, found := cli.short2long[flag.Short]; found {
			return NewParseError("%s: short flag name %q already defined",
				cli.name, flag.Short)
		}
	}

	//
	// Validate the long flag.
	//
	if strings.HasPrefix(flag.Long, "-") {
		return NewParseError("long flag name %q must not begin with '-'", flag.Long)
	}
	if strings.Contains(flag.Long, "=") {
		return NewParseError("long flag name %q must not contain '='", flag.Long)
	}
	if flag.Long == "" {
		return NewParseError("long flag name cannot be empty")
	}
	if len(flag.Long) < 2 {
		return NewParseError("long flag name %q must be at least 2 characters",
			flag.Long)
	}
	if flag.Long == "help" {
		return NewParseError(`cannot override long flag name "help"`)
	}
	if _, found := cli.long2flag[flag.Long]; found {
		return NewParseError("%s: long flag name %q already defined",
			cli.name, flag.Long)
	}

	// A variable can be bound to only one flag.
	for k, fl := range cli.long2flag {
		if fl.Value == flag.Value {
			return NewParseError(
				"long flag name %q: variable already bound to flag %q",
				flag.Long, k)
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
	cli.orderedFlags = append(cli.orderedFlags, flag.Long)

	return nil
}

// AddGroup adds the subclis to the group name.
func (cli *CLI) AddGroup(name string, clis ...*CLI) error {
	if len(clis) == 0 {
		return NewParseError("AddGroup %s: child list is empty", name)
	}
	for _, child := range clis {
		if !slices.Contains(cli.subCLIs, child) {
			return NewParseError("AddGroup %s: child %s is missing previous AddCLI",
				name, child.name)
		}
	}
	cli.groups = append(cli.groups, cliGroup{name, clis})
	return nil
}

// Parse processes args, following subcommands (if any).
// It returns the subcommand path.
func (cli *CLI) Parse(args []string) (command string, err error) {
	index := 0

	// Parse all the options. At the end of the loop, 'index' points to the
	// beginning (if any) of the positional arguments.
	for {
		long, offset, err := cli.parseOne(args[index:])
		if err != nil {
			return "", err
		}
		cli.longSeen[long] = struct{}{}
		if offset == 0 {
			// Arrived at the end of the options.
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
		return "", NewParseError("missing required options: %s",
			strings.Join(missing, ", "))
	}

	//
	// Process the remaining of args (if any).
	//

	cli.positionals = args[index:]

	if len(cli.subCLIs) > 0 && cli.posargs != nil {
		return "", fmt.Errorf(
			"clim: internal error: command %q has both subcommands and pos args",
			cli.rootToHere)
	}

	//
	// Subcommand.
	//
	if len(cli.subCLIs) > 0 {
		if len(cli.positionals) == 0 {
			return "", NewParseError("expected a command")
		}
		command := cli.positionals[0]
		for _, p := range cli.subCLIs {
			if p.name == command {
				return p.Parse(cli.positionals[1:])
			}
		}
		return "", NewParseError("unrecognized command %q", command)
	}

	//
	// Positional arguments.
	//
	if cli.posargs != nil {
		*cli.posargs = cli.positionals
	}

	return cli.rootToHere, nil
}

type Pair struct {
	Name string
	Help string
}

func (cli *CLI) AddPosArgs(values *[]string, pairs ...Pair) error {
	if len(cli.subCLIs) > 0 {
		// FIXME this is NOT a parse error!!!
		return NewParseError("%s: already have subcommands; cannot have also pos args",
			cli.name)
	}

	cli.posargs = values
	cli.pairs = pairs
	names := make(map[string]int, len(pairs))

	for idx, pair := range pairs {
		if idx2, found := names[pair.Name]; found {
			return NewParseError(
				"%s: pos arg at index %d (%q) was already defined at index %d",
				cli.name, idx, pair.Name, idx2)
		}
		if pair.Name == "" {
			return NewParseError("%s: pos arg at index %d (%q) cannot be empty",
				cli.name, idx, pair.Name)
		}

		// A variable can be bound to only one flag.
		// for k, fl := range cli.long2flag {
		// 	if fl.Value == flag.Value {
		// 		return NewParseError"long flag name %q: variable already bound to flag %q",
		// 			flag.Long, k))
		// 	}
		// }

		names[pair.Name] = idx
	}
	return nil
}

// CountTrue returns the number of args that are true.
func CountTrue(args ...bool) int {
	n := 0
	for _, v := range args {
		if v {
			n++
		}
	}
	return n
}

// ExitCode is an helper function to reduce the boilerplate in your program.
// It returns an appropriate exit code:
//   - 0  if no error, or help requested
//   - 1  if any other errors
//
// Usage:
//
//	func main() {
//	  os.Exit(MainInt(os.Args[1:]))
//	}
//
//	func MainInt(args []string) int {
//		if err := mainErr(args); err != nil {
//			fmt.Println(err)
//			return clim.ExitCode(err)
//		}
//		return 0
//	}
func ExitCode(err error) int {
	if err != nil {
		if errors.Is(err, ErrHelp) {
			return 0
		}
		return 1
	}
	return 0
}

// regex to match an option on the command-line.
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
func (cli *CLI) parseOne(args []string) (string, int, error) {
	if len(args) == 0 {
		return "", 0, nil
	}
	token := args[0]
	matches := flagRE.FindStringSubmatch(token)
	if len(matches) != 6 {
		return "", 0,
			NewParseError("clim internal error (regex); token: %q, matches: %q",
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
			return "", 0, NewParseError("unrecognized flag %q", token)
		}
	}
	flag := cli.long2flag[long]
	if flag == nil {
		return "", 0, NewParseError("unrecognized flag %q", token)
	}

	// Was the value provided in the same token, with "=" ?
	if len(value) > 0 {
		if err := flag.Value.Set(value); err != nil {
			return "", 0, NewParseError("setting %q: %s", token, err)
		}
		return long, 1, nil
	}

	if isBoolValue(flag.Value) {
		if err := flag.Value.Set("true"); err != nil {
			return "", 0, NewParseError("clim internal error: setting %q: %s",
				token, err)
		}
		return long, 1, nil
	}

	if len(args) == 1 {
		return "", 0, NewParseError("flag %q requires a value", token)
	}
	nextValue := args[1]
	if err := flag.Value.Set(nextValue); err != nil {
		return "", 0, NewParseError("setting %q %q: %s", token, nextValue, err)
	}
	return long, 2, nil
}

// pathRootToNode returns the CLI names in the tree path from the root to
// 'node'.
// TODO write test and add this to all errors?
func pathRootToNode(node *CLI) []string {
	var path []string
	cursor := node
	for {
		path = append(path, cursor.name)
		if cursor.parent == nil {
			break
		}
		cursor = cursor.parent
	}
	slices.Reverse(path)
	return path
}
