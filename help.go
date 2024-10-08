package clim

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

func (cli *CLI[T]) usage() error {
	var bld strings.Builder

	// Calculate the max width of the first column of commands.
	maxColWidth := 0
	for _, p := range cli.subCLIs {
		fmt.Fprintf(&bld, " %s", p.name)
		maxColWidth = max(maxColWidth, bld.Len())
		bld.Reset()
	}

	fmt.Fprintf(&bld, "%s -- %s\n\n", cli.rootToHere, cli.oneline)

	if cli.description != "" {
		for _, line := range strings.Split(cli.description, "\n") {
			fmt.Fprintf(&bld, " %s\n", line)
		}
		fmt.Fprintln(&bld)
	}

	fmt.Fprintf(&bld, "Usage: %s ", cli.rootToHere)
	if len(cli.subCLIs) > 0 {
		fmt.Fprintf(&bld, "<command> ")
	}
	fmt.Fprintf(&bld, "[options]")
	for _, pair := range cli.pairs {
		fmt.Fprintf(&bld, " %s", pair.Name)
	}
	fmt.Fprintf(&bld, "\n\n")

	if cli.examples != "" {
		fmt.Fprintf(&bld, "Examples:\n\n")
		for _, line := range strings.Split(cli.examples, "\n") {
			if line != "" {
				fmt.Fprintf(&bld, " %s\n", line)
			} else {
				fmt.Fprintln(&bld)
			}
		}
		fmt.Fprintln(&bld)
	}

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
			printSomeSubCommands(&bld, width, group.clis)
		}
	} else if len(cli.subCLIs) > 0 {
		printSomeSubCommands(&bld, width, cli.subCLIs)
	}

	cli.printOptions(&bld)

	cli.printPosArgs(&bld)

	if cli.footer != "" {
		fmt.Fprintln(&bld)
		for _, line := range strings.Split(cli.footer, "\n") {
			fmt.Fprintf(&bld, " %s\n", line)
		}
	}

	return newHelpError("%s", bld.String())
}

func (cli *CLI[T]) printOptions(bld *strings.Builder) {
	// First pass. Sort keys.
	longs := slices.Sorted(maps.Keys(cli.long2flag))

	// Second pass, calculate the max width of the first column.
	lines := make([]string, 0, len(longs)+1)
	var tmp strings.Builder
	maxColWidth := 0
	for _, long := range longs {
		flag := cli.long2flag[long]
		fmt.Fprintf(&tmp, " ")
		if flag.Short != "" {
			fmt.Fprintf(&tmp, "-%s, ", flag.Short)
		}
		fmt.Fprintf(&tmp, "--%s %s", flag.Long, flag.Label)
		lines = append(lines, tmp.String())
		maxColWidth = max(maxColWidth, tmp.Len())
		tmp.Reset()
	}
	// Same for -h
	fmt.Fprint(&tmp, " -h, --help")
	maxColWidth = max(maxColWidth, tmp.Len())
	tmp.Reset()

	// Third pass, add the second column.
	const gutter = 4
	fmt.Fprintf(bld, "Options:\n\n")
	for i, long := range longs {
		flag := cli.long2flag[long]
		fmt.Fprintf(bld, "%-*s%s", maxColWidth+gutter, lines[i], flag.Help)
		if flag.defValue != "" && !flag.Required {
			fmt.Fprintf(bld, " (default: %s)", flag.defValue)
		}
		if flag.Required {
			fmt.Fprintf(bld, " (required)")
		}
		fmt.Fprintf(bld, "\n")
	}
	if len(longs) > 0 {
		fmt.Fprintf(bld, "\n")
	}

	fmt.Fprintf(bld, "%-*s%s", maxColWidth+gutter,
		" -h, --help", "Print this help and exit\n")
}

func (cli *CLI[T]) printPosArgs(bld *strings.Builder) {
	if len(cli.pairs) == 0 {
		return
	}

	// First pass, calculate the max width of the first column.
	maxColWidth := 0
	for _, pair := range cli.pairs {
		maxColWidth = max(maxColWidth, len(pair.Name))
	}

	// Second pass, consider the second column.
	fmt.Fprintln(bld)
	const gutter = 6
	fmt.Fprintf(bld, "Positional arguments:\n\n")
	for _, pair := range cli.pairs {
		fmt.Fprintf(bld, " %-*s%s\n", maxColWidth+gutter, pair.Name, pair.Help)
	}
}

func printSomeSubCommands[T any](bld *strings.Builder, width int, subclis []*CLI[T]) {
	for _, cmd := range subclis {
		fmt.Fprintf(bld, " %-*s%s\n", width, cmd.name, cmd.oneline)
	}
	fmt.Fprintln(bld)
}
