package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/devanshbatham/ip2cloud/internal/store"
)

func runRemove(args []string) {
	removeUsage := func() {
		fmt.Fprintf(os.Stderr, "Usage: ip2cloud remove <provider> [-build]\n\n")
		fmt.Fprintf(os.Stderr, "Remove a cloud provider and its CIDR ranges.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		fmt.Fprintf(os.Stderr, "  -build    Rebuild binary trie after removing (default: true)\n")
	}

	if len(args) < 1 || args[0] == "-h" || args[0] == "--help" {
		removeUsage()
		if len(args) >= 1 {
			os.Exit(0)
		}
		os.Exit(1)
	}

	provider := args[0]

	fs := flag.NewFlagSet("remove", flag.ExitOnError)
	rebuild := fs.Bool("build", true, "Rebuild binary trie after removing")
	fs.Usage = removeUsage
	fs.Parse(args[1:])

	s, err := store.DefaultStore()
	if err != nil {
		fatal("%v", err)
	}

	if err := s.RemoveProvider(provider); err != nil {
		fatal("%v", err)
	}
	fmt.Printf("Removed provider '%s'\n", provider)

	if *rebuild {
		if _, err := s.Build(); err != nil {
			fatal("rebuild: %v", err)
		}
		fmt.Println("Rebuilt binary trie")
	}
}
