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
		fmt.Fprintf(sb, "usage: %s [options] command\n", os.Args[0])
		sb.WriteString("\nAvailable commands\n")
		tw := tabwriter.NewWriter(sb, 0, 4, 1, ' ', 0)
		for _, actionName := range actionString[1:] {
			fmt.Fprintf(tw, "  %s\t\t%s\n", actionName, commands[getAction(actionName)].Description)
		}
		tw.Flush()
		sb.WriteString("\nAvailable options:\n")
		fmt.Fprint(flag.CommandLine.Output(), sb.String())
		flag.PrintDefaults()
	}
	flag.BoolVar(&isGlobal, "g", false, "Run the command globally (can only be used with the 'use', 'edit', and 'delete' commands).")
	flag.Parse()
}
