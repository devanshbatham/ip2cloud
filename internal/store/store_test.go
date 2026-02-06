package store

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestAddAndBuildRoundTrip(t *testing.T) {
	tmp := t.TempDir()
	s := &Store{
		DataDir: filepath.Join(tmp, "data"),
		BinPath: filepath.Join(tmp, "ip2cloud.bin"),
	}

	if err := s.AddRanges("testprov", []string{"10.0.0.0/8"}); err != nil {
		t.Fatalf("AddRanges: %v", err)
	}

	if _, err := s.Build(); err != nil {
		t.Fatalf("Build: %v", err)
	}

	tr, err := s.LoadTrie()
	if err != nil {
		t.Fatalf("LoadTrie: %v", err)
	}

	if got := tr.Lookup("10.0.0.1"); got != "testprov" {
		t.Errorf("Lookup(10.0.0.1) = %q, want %q", got, "testprov")
	}
	if got := tr.Lookup("192.168.1.1"); got != "" {
		t.Errorf("Lookup(192.168.1.1) = %q, want empty", got)
	}
}

func TestListProviders(t *testing.T) {
	tmp := t.TempDir()
	s := &Store{
		DataDir: filepath.Join(tmp, "data"),
		BinPath: filepath.Join(tmp, "ip2cloud.bin"),
	}

	if err := s.AddRanges("zeta", []string{"10.0.0.0/8"}); err != nil {
		t.Fatalf("AddRanges zeta: %v", err)
	}
	if err := s.AddRanges("alpha", []string{"192.168.0.0/16", "172.16.0.0/12"}); err != nil {
		t.Fatalf("AddRanges alpha: %v", err)
	}

	providers, err := s.ListProviders()
	if err != nil {
		t.Fatalf("ListProviders: %v", err)
	}

	if len(providers) != 2 {
		t.Fatalf("got %d providers, want 2", len(providers))
	}

	if providers[0].Name != "alpha" || providers[0].RangeCount != 2 {
		t.Errorf("providers[0] = %+v, want {Name:alpha RangeCount:2}", providers[0])
	}
	if providers[1].Name != "zeta" || providers[1].RangeCount != 1 {
		t.Errorf("providers[1] = %+v, want {Name:zeta RangeCount:1}", providers[1])
	}
}

func TestSeedFromFSDoesNotOverwrite(t *testing.T) {
	tmp := t.TempDir()
	dataDir := filepath.Join(tmp, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		t.Fatal(err)
	}

	existing := "10.0.0.0/8\n"
	if err := os.WriteFile(filepath.Join(dataDir, "aws.txt"), []byte(existing), 0644); err != nil {
		t.Fatal(err)
	}

	s := &Store{
		DataDir: dataDir,
		BinPath: filepath.Join(tmp, "ip2cloud.bin"),
	}

	seedFS := fstest.MapFS{
		"aws.txt": &fstest.MapFile{Data: []byte("192.168.0.0/16\n")},
	}

	if err := s.SeedFromFS(seedFS); err != nil {
		t.Fatalf("SeedFromFS: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(dataDir, "aws.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != existing {
		t.Errorf("file content = %q, want %q", got, existing)
	}
}

func TestSeedFromFSWritesMissing(t *testing.T) {
	tmp := t.TempDir()
	dataDir := filepath.Join(tmp, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		t.Fatal(err)
	}

	s := &Store{
		DataDir: dataDir,
		BinPath: filepath.Join(tmp, "ip2cloud.bin"),
	}

	seedFS := fstest.MapFS{
		"aws.txt": &fstest.MapFile{Data: []byte("10.0.0.0/8\n")},
	}

	if err := s.SeedFromFS(seedFS); err != nil {
		t.Fatalf("SeedFromFS: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dataDir, "aws.txt")); err != nil {
		t.Errorf("expected aws.txt to exist: %v", err)
	}
}
