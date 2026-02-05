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

	// Non-interactive mode flags
	noTUI          bool
	profileFlag    string
	nameFlag       string
	emailFlag      string
	signingKeyFlag string
	keyFormatFlag  string
	gpgProgramFlag string
	yesFlag        bool
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

	// Existing flags
	flag.BoolVar(&isGlobal, "g", false, "Run the command globally (can only be used with the 'use', 'edit', and 'delete' commands).")

	// Non-interactive mode flags
	flag.BoolVar(&noTUI, "no-tui", false, "Disable TUI prompts for automated/scripted usage.")
	flag.StringVar(&profileFlag, "profile", "", "Profile name (for create/use/delete in --no-tui mode).")
	flag.StringVar(&nameFlag, "name", "", "Git user name (for create in --no-tui mode).")
	flag.StringVar(&emailFlag, "email", "", "Git user email (for create in --no-tui mode).")
	flag.StringVar(&signingKeyFlag, "signing-key", "", "Signing key (GPG key ID, SSH key path, or X.509 certificate ID).")
	flag.StringVar(&keyFormatFlag, "key-format", "", "Signing key format: 'openpgp', 'ssh', or 'x509' (default: openpgp if --signing-key is set).")
	flag.StringVar(&gpgProgramFlag, "gpg-program", "", "GPG program to use (default: gpg). Only applicable for openpgp format.")
	flag.BoolVar(&yesFlag, "yes", false, "Confirm destructive operations without prompting (for --no-tui mode).")

	flag.Parse()
}
