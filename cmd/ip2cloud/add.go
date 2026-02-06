package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/devanshbatham/ip2cloud/internal/store"
)

func runAdd(args []string) {
	addUsage := func() {
		fmt.Fprintf(os.Stderr, "Usage: ip2cloud add <provider> [-f file] [cidrs...]\n\n")
		fmt.Fprintf(os.Stderr, "Add CIDR ranges for a cloud provider.\n")
		fmt.Fprintf(os.Stderr, "CIDRs can be passed as arguments, from a file (-f), or piped via stdin (-f -).\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		fmt.Fprintf(os.Stderr, "  -f string    Read CIDRs from a file (use '-' for stdin)\n")
		fmt.Fprintf(os.Stderr, "  -build       Rebuild binary trie after adding (default: true)\n")
	}

	if len(args) < 1 || args[0] == "-h" || args[0] == "--help" {
		addUsage()
		if len(args) >= 1 {
			os.Exit(0)
		}
		os.Exit(1)
	}

	provider := args[0]

	fs := flag.NewFlagSet("add", flag.ExitOnError)
	file := fs.String("f", "", "Read CIDRs from a file (use '-' for stdin)")
	rebuild := fs.Bool("build", true, "Rebuild binary trie after adding")
	fs.Usage = addUsage
	fs.Parse(args[1:])

	cidrs := fs.Args()

	if *file != "" {
		var r *os.File
		if *file == "-" {
			r = os.Stdin
		} else {
			var err error
			r, err = os.Open(*file)
			if err != nil {
				fatal("opening file: %v", err)
			}
			defer r.Close()
		}
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line != "" && !strings.HasPrefix(line, "#") {
				cidrs = append(cidrs, line)
			}
		}
		if err := sc.Err(); err != nil {
			fatal("reading input: %v", err)
		}
	}

	if len(cidrs) == 0 {
		fatal("no CIDRs provided. Use arguments, -f file, or -f - for stdin.")
	}

	s, err := store.DefaultStore()
	if err != nil {
		fatal("%v", err)
	}

	if s.ProviderExists(provider) {
		fmt.Printf("Provider '%s' already exists. Overwrite? [y/N]: ", provider)
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			fmt.Println("Aborted.")
			return
		}
		if err := s.OverwriteRanges(provider, cidrs); err != nil {
			fatal("overwriting ranges: %v", err)
		}
		fmt.Printf("Overwrote %s with %d ranges\n", provider, len(cidrs))
	} else {
		if err := s.AddRanges(provider, cidrs); err != nil {
			fatal("adding ranges: %v", err)
		}
		fmt.Printf("Added %d ranges to %s\n", len(cidrs), provider)
	}

	if *rebuild {
		if _, err := s.Build(); err != nil {
			fatal("rebuild: %v", err)
		}
		fmt.Println("Rebuilt binary trie")
	}
}
