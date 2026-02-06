package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/devanshbatham/ip2cloud/internal/store"
)

func runList() {
	s, err := store.DefaultStore()
	if err != nil {
		fatal("%v", err)
	}

	providers, err := s.ListProviders()
	if err != nil {
		fatal("listing providers: %v", err)
	}

	if len(providers) == 0 {
		fmt.Println("No providers found. Run 'ip2cloud build' to import embedded data.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "PROVIDER\tRANGES")
	total := 0
	for _, p := range providers {
		fmt.Fprintf(w, "%s\t%d\n", p.Name, p.RangeCount)
		total += p.RangeCount
	}
	fmt.Fprintf(w, "\t\nTOTAL\t%d\n", total)
	w.Flush()
}
