package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	ip2cloud "github.com/devanshbatham/ip2cloud"
	"github.com/devanshbatham/ip2cloud/internal/store"
)

func runBuild(args []string) {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	seedDir := fs.String("seed", "", "Seed data from a directory of .txt files (e.g., ./data)")
	fs.Parse(args)

	s, err := store.DefaultStore()
	if err != nil {
		fatal("%v", err)
	}

	if err := s.Init(); err != nil {
		fatal("creating data dir: %v", err)
	}

	if *seedDir != "" {
		if err := seedFromDir(s, *seedDir); err != nil {
			fatal("seeding: %v", err)
		}
	} else {
		embeddedData, err := ip2cloud.EmbeddedData()
		if err != nil {
			fatal("loading embedded data: %v", err)
		}
		if err := s.SeedFromFS(embeddedData); err != nil {
			fatal("seeding embedded data: %v", err)
		}
	}

	t, err := s.Build()
	if err != nil {
		fatal("build: %v", err)
	}

	for _, w := range t.Warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", w)
	}

	providers := t.Providers[1:]
	fmt.Printf("Built trie: %d providers, saved to %s\n", len(providers), s.BinPath)
}

func seedFromDir(s *store.Store, dir string) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".txt") {
			return err
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		dst := filepath.Join(s.DataDir, d.Name())
		return os.WriteFile(dst, src, 0644)
	})
}
