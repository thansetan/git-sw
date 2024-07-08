package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

var (
	isGlobal      bool
	allowedGlobal = map[Action]struct{}{
		USE:    {},
		EDIT:   {},
		DELETE: {},
	}
)

func parseFlag() {
	flag.Usage = func() {
		sb := new(strings.Builder)
		sb.WriteString(fmt.Sprintf("usage: %s [options] command\n", os.Args[0]))
		sb.WriteString("Available commands: ")
		sb.WriteByte('\n')
		tw := tabwriter.NewWriter(sb, 0, 4, 1, ' ', 0)
		for _, actionName := range actionString[1:] {
			fmt.Fprintf(tw, "  %s\t\t%s\n", actionName, commands[getAction(actionName)].Description)
		}
		tw.Flush()
		sb.WriteByte('\n')
		sb.WriteString("Available options: ")
		sb.WriteByte('\n')
		fmt.Fprint(flag.CommandLine.Output(), sb.String())
		flag.PrintDefaults()
	}
	flag.BoolVar(&isGlobal, "g", false, "whether to run the command globally (can only be used with the 'use', 'edit', and 'delete' commands).")
	flag.Parse()
}
