package store

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devanshbatham/ip2cloud/internal/trie"
)

type Store struct {
	DataDir string
	BinPath string
}

func DefaultStore() (*Store, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	base := filepath.Join(configDir, "ip2cloud")
	return &Store{
		DataDir: filepath.Join(base, "data"),
		BinPath: filepath.Join(base, "ip2cloud.bin"),
	}, nil
}

func (s *Store) Init() error {
	return os.MkdirAll(s.DataDir, 0755)
}

func (s *Store) ReadProviderRanges(provider string) ([]string, error) {
	path := filepath.Join(s.DataDir, provider+".txt")
	return readLines(path)
}

func (s *Store) ProviderExists(provider string) bool {
	path := filepath.Join(s.DataDir, provider+".txt")
	_, err := os.Stat(path)
	return err == nil
}

func (s *Store) AddRanges(provider string, cidrs []string) error {
	if err := s.Init(); err != nil {
		return err
	}
	path := filepath.Join(s.DataDir, provider+".txt")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, cidr := range cidrs {
		w.WriteString(cidr)
		w.WriteByte('\n')
	}
	return w.Flush()
}

func (s *Store) OverwriteRanges(provider string, cidrs []string) error {
	if err := s.Init(); err != nil {
		return err
	}
	path := filepath.Join(s.DataDir, provider+".txt")
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, cidr := range cidrs {
		w.WriteString(cidr)
		w.WriteByte('\n')
	}
	return w.Flush()
}

func (s *Store) ListProviders() ([]ProviderInfo, error) {
	entries, err := os.ReadDir(s.DataDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var result []ProviderInfo
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".txt") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".txt")
		ranges, err := readLines(filepath.Join(s.DataDir, e.Name()))
		if err != nil {
			continue
		}
		result = append(result, ProviderInfo{Name: name, RangeCount: len(ranges)})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result, nil
}

type ProviderInfo struct {
	Name       string
	RangeCount int
}

func (s *Store) Build() (*trie.Trie, error) {
	entries, err := os.ReadDir(s.DataDir)
	if err != nil {
		return nil, fmt.Errorf("reading data dir: %w", err)
	}

	cloudData := make(map[string][]string)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".txt") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".txt")
		ranges, err := readLines(filepath.Join(s.DataDir, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", e.Name(), err)
		}
		cloudData[name] = ranges
	}

	t := trie.Build(cloudData)
	if err := t.Save(s.BinPath); err != nil {
		return nil, fmt.Errorf("saving binary trie: %w", err)
	}
	return t, nil
}

func (s *Store) LoadTrie() (*trie.Trie, error) {
	return trie.Load(s.BinPath)
}

func (s *Store) LoadOrBuildTrie(seedFS fs.FS) (*trie.Trie, error) {
	t, err := trie.Load(s.BinPath)
	if err == nil {
		return t, nil
	}
	if err := s.Init(); err != nil {
		return nil, fmt.Errorf("creating data dir: %w", err)
	}
	if err := s.SeedFromFS(seedFS); err != nil {
		return nil, fmt.Errorf("seeding data: %w", err)
	}
	return s.Build()
}

func (s *Store) SeedFromFS(seedFS fs.FS) error {
	return fs.WalkDir(seedFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".txt") {
			return err
		}
		src, err := fs.ReadFile(seedFS, path)
		if err != nil {
			return err
		}
		dst := filepath.Join(s.DataDir, d.Name())
		if _, err := os.Stat(dst); err == nil {
			return nil
		}
		return os.WriteFile(dst, src, 0644)
	})
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, sc.Err()
}
