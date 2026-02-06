package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	ip2cloud "github.com/devanshbatham/ip2cloud"
	"github.com/devanshbatham/ip2cloud/internal/store"
)

const batchSize = 4096

type result struct {
	ip       string
	provider string
}

func runLookup(args []string) {
	fs := flag.NewFlagSet("ip2cloud", flag.ExitOnError)
	jsonOutput := fs.Bool("j", false, "Print output in JSON format")
	fs.BoolVar(jsonOutput, "json", false, "Print output in JSON format")
	workers := fs.Int("w", runtime.NumCPU(), "Number of concurrent workers")
	providerFlag := fs.String("provider", "", "Only check against specific providers (comma-separated, e.g., aws,gcp)")
	fs.StringVar(providerFlag, "p", "", "Only check against specific providers (comma-separated, e.g., aws,gcp)")
	fs.Parse(args)

	if *workers < 1 {
		*workers = 1
	}

	allowedProviders := make(map[string]bool)
	if *providerFlag != "" {
		for _, p := range strings.Split(*providerFlag, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				allowedProviders[strings.ToLower(p)] = true
			}
		}
	}

	s, err := store.DefaultStore()
	if err != nil {
		fatal("%v", err)
	}

	embeddedData, err := ip2cloud.EmbeddedData()
	if err != nil {
		fatal("loading embedded data: %v", err)
	}

	trie, err := s.LoadOrBuildTrie(embeddedData)
	if err != nil {
		fatal("loading trie: %v", err)
	}

	for _, w := range trie.Warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", w)
	}

	ipCh := make(chan []string, *workers*2)
	resCh := make(chan []result, *workers*2)

	var wg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range ipCh {
				var results []result
				for _, ip := range batch {
					provider := trie.Lookup(ip)
					if provider == "" {
						continue
					}
					if len(allowedProviders) > 0 && !allowedProviders[strings.ToLower(provider)] {
						continue
					}
					results = append(results, result{ip: ip, provider: provider})
				}
				if len(results) > 0 {
					resCh <- results
				}
			}
		}()
	}

	go func() {
		if positional := fs.Args(); len(positional) > 0 {
			ipCh <- positional
		} else {
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
			batch := make([]string, 0, batchSize)
			for scanner.Scan() {
				ip := strings.TrimSpace(scanner.Text())
				if ip == "" {
					continue
				}
				batch = append(batch, ip)
				if len(batch) >= batchSize {
					ipCh <- batch
					batch = make([]string, 0, batchSize)
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", err)
			}
			if len(batch) > 0 {
				ipCh <- batch
			}
		}
		close(ipCh)
	}()

	go func() {
		wg.Wait()
		close(resCh)
	}()

	if *jsonOutput {
		grouped := make(map[string][]string)
		for batch := range resCh {
			for _, r := range batch {
				grouped[r.provider] = append(grouped[r.provider], r.ip)
			}
		}
		out, err := json.MarshalIndent(grouped, "", "    ")
		if err != nil {
			fatal("marshaling JSON: %v", err)
		}
		os.Stdout.Write(out)
		os.Stdout.Write([]byte("\n"))
	} else {
		w := bufio.NewWriterSize(os.Stdout, 256*1024)
		for batch := range resCh {
			for _, r := range batch {
				fmt.Fprintf(w, "[%s] : %s\n", r.provider, r.ip)
			}
		}
		if err := w.Flush(); err != nil {
			fatal("flushing output: %v", err)
		}
	}
}
