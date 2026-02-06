package main

import (
	"fmt"
	"os"
)

var version = "dev"

const usage = `ip2cloud - Check IPs against cloud provider ranges

Usage:
  ip2cloud [flags] [ip ...]     Lookup IPs from stdin or arguments
  ip2cloud build [flags]        Build binary trie from provider data
  ip2cloud add <provider> ...   Add CIDR ranges to a provider
  ip2cloud remove <provider>    Remove a provider and its ranges
  ip2cloud list                 List providers and range counts
  ip2cloud version              Print version

Lookup Flags:
  -p, -provider string   Only match specific providers (comma-separated, e.g., aws,azure)
  -j, -json              Print output in JSON format
  -w int                 Number of concurrent workers (default: NumCPU)

Build Flags:
  -seed string           Seed data from a directory of .txt files (default: embedded data)

Add Flags:
  -f string              Read CIDRs from a file (use '-' for stdin)
  -build                 Rebuild binary trie after adding (default: true)

Remove Flags:
  -build                 Rebuild binary trie after removing (default: true)

Examples:
  cat ips.txt | ip2cloud              Lookup IPs from stdin
  ip2cloud 8.8.8.8 3.5.1.1           Lookup specific IPs
  ip2cloud -p aws < ips.txt           Only show AWS matches
  ip2cloud -j < ips.txt               Output as JSON
  ip2cloud add mycloud 10.0.0.0/8     Add a CIDR range
  ip2cloud remove mycloud             Remove a provider
  ip2cloud list                       List all providers
  ip2cloud build                      Rebuild trie from embedded data

Run 'ip2cloud <command> -h' for command-specific help.
`

func main() {
	if len(os.Args) < 2 {
		runLookup(os.Args[1:])
		return
	}

	switch os.Args[1] {
	case "build":
		runBuild(os.Args[2:])
	case "add":
		runAdd(os.Args[2:])
	case "remove":
		runRemove(os.Args[2:])
	case "list":
		runList()
	case "-v", "--version", "version":
		fmt.Printf("ip2cloud version %s\n", version)
	case "-h", "--help", "help":
		fmt.Print(usage)
	default:
		runLookup(os.Args[1:])
	}
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
