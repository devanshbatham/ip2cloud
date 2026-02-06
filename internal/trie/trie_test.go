package trie

import (
	"bytes"
	"fmt"
	"testing"
)

var testData = map[string][]string{
	"aws":   {"52.0.0.0/8", "63.32.0.0/14"},
	"azure": {"64.4.8.0/24"},
	"gcp":   {"34.0.0.0/8"},
}

func TestLookup(t *testing.T) {
	tr := Build(testData)

	tests := []struct {
		ip   string
		want string
	}{
		{"52.1.2.3", "aws"},
		{"52.255.255.255", "aws"},
		{"63.32.40.140", "aws"},
		{"63.33.205.240", "aws"},
		{"63.35.255.255", "aws"},
		{"64.4.8.90", "azure"},
		{"64.4.8.0", "azure"},
		{"64.4.8.255", "azure"},
		{"34.100.50.25", "gcp"},
		{"192.168.1.1", ""},
		{"10.0.0.1", ""},
		{"invalid", ""},
		{"", ""},
		{"999.1.2.3", ""},
		{"1.2.3", ""},
		{"1.2.3.4.5", ""},
	}

	for _, tt := range tests {
		got := tr.Lookup(tt.ip)
		if got != tt.want {
			t.Errorf("Lookup(%q) = %q, want %q", tt.ip, got, tt.want)
		}
	}
}

func TestSingleHost(t *testing.T) {
	tr := Build(map[string][]string{
		"test": {"10.0.0.1/32"},
	})
	if got := tr.Lookup("10.0.0.1"); got != "test" {
		t.Errorf("got %q, want \"test\"", got)
	}
	if got := tr.Lookup("10.0.0.2"); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestLongestPrefixMatch(t *testing.T) {
	tr := Build(map[string][]string{
		"broad":  {"10.0.0.0/8"},
		"narrow": {"10.0.0.0/24"},
	})
	if got := tr.Lookup("10.0.0.5"); got != "narrow" {
		t.Errorf("got %q, want \"narrow\"", got)
	}
	if got := tr.Lookup("10.0.1.5"); got != "broad" {
		t.Errorf("got %q, want \"broad\"", got)
	}
}

func TestBoundaryAddresses(t *testing.T) {
	tr := Build(map[string][]string{"test": {"192.168.1.0/24"}})
	cases := []struct{ ip, want string }{
		{"192.168.1.0", "test"},
		{"192.168.1.255", "test"},
		{"192.168.0.255", ""},
		{"192.168.2.0", ""},
	}
	for _, c := range cases {
		if got := tr.Lookup(c.ip); got != c.want {
			t.Errorf("Lookup(%q) = %q, want %q", c.ip, got, c.want)
		}
	}
}

func TestSerializeRoundTrip(t *testing.T) {
	original := Build(testData)

	var buf bytes.Buffer
	if err := original.Encode(&buf); err != nil {
		t.Fatalf("Encode: %v", err)
	}

	loaded, err := Decode(buf.Bytes())
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}

	ips := []string{"52.1.2.3", "63.33.205.240", "64.4.8.90", "34.100.50.25", "192.168.1.1"}
	for _, ip := range ips {
		want := original.Lookup(ip)
		got := loaded.Lookup(ip)
		if got != want {
			t.Errorf("after round-trip, Lookup(%q) = %q, want %q", ip, got, want)
		}
	}
}

func TestParseIPv4(t *testing.T) {
	tests := []struct {
		input string
		want  uint32
		ok    bool
	}{
		{"0.0.0.0", 0, true},
		{"255.255.255.255", 0xFFFFFFFF, true},
		{"10.0.0.1", 0x0A000001, true},
		{"192.168.1.100", 0xC0A80164, true},
		{"", 0, false},
		{"1.2.3", 0, false},
		{"1.2.3.4.5", 0, false},
		{"256.0.0.0", 0, false},
		{"abc", 0, false},
	}
	for _, tt := range tests {
		got, ok := ParseIPv4(tt.input)
		if ok != tt.ok || (ok && got != tt.want) {
			t.Errorf("ParseIPv4(%q) = (%d, %v), want (%d, %v)", tt.input, got, ok, tt.want, tt.ok)
		}
	}
}

func TestBuildIPv6Skipped(t *testing.T) {
	tr := Build(map[string][]string{
		"gcp": {"2600:1900::/28", "34.0.0.0/8"},
	})

	if len(tr.Warnings) == 0 {
		t.Fatal("expected warnings for IPv6 CIDR, got none")
	}

	if got := tr.Lookup("34.1.2.3"); got != "gcp" {
		t.Errorf("Lookup(34.1.2.3) = %q, want %q", got, "gcp")
	}
}

func TestBuildInvalidCIDR(t *testing.T) {
	tr := Build(map[string][]string{
		"bad": {"not-a-cidr"},
	})

	if len(tr.Warnings) == 0 {
		t.Fatal("expected warnings for invalid CIDR, got none")
	}
}

func TestBuildWarningsReported(t *testing.T) {
	tr := Build(map[string][]string{
		"mixed": {"10.0.0.0/8", "not-valid", "2001:db8::/32"},
	})

	if got := tr.Lookup("10.0.0.1"); got != "mixed" {
		t.Errorf("Lookup(10.0.0.1) = %q, want %q", got, "mixed")
	}

	if len(tr.Warnings) != 2 {
		t.Errorf("got %d warnings, want 2: %v", len(tr.Warnings), tr.Warnings)
	}
}

func BenchmarkLookupHit(b *testing.B) {
	tr := Build(testData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Lookup("52.1.2.3")
	}
}

func BenchmarkLookupMiss(b *testing.B) {
	tr := Build(testData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Lookup("192.168.1.1")
	}
}

func BenchmarkLookupParallel(b *testing.B) {
	tr := Build(testData)
	ips := []string{"52.1.2.3", "63.33.205.240", "64.4.8.90", "34.100.50.25", "192.168.1.1"}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			tr.Lookup(ips[i%len(ips)])
			i++
		}
	})
}

func BenchmarkLookupLargeDataset(b *testing.B) {
	data := make(map[string][]string)
	for i := 0; i < 200; i++ {
		var cidrs []string
		for j := 0; j < 50; j++ {
			cidrs = append(cidrs, fmt.Sprintf("%d.%d.0.0/16", i%256, j%256))
		}
		data[fmt.Sprintf("provider_%d", i)] = cidrs
	}
	tr := Build(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Lookup("100.25.1.1")
	}
}

func BenchmarkLoadBinary(b *testing.B) {
	tr := Build(testData)
	var buf bytes.Buffer
	tr.Encode(&buf)
	raw := buf.Bytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(raw)
	}
}

func BenchmarkParseIPv4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseIPv4("192.168.1.100")
	}
}
