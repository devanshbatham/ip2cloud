package trie

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"testing"
)

func loadEmbeddedTrie(b *testing.B) *Trie {
	b.Helper()
	data := generateRealisticData()
	return Build(data)
}

func generateRealisticData() map[string][]string {
	return map[string][]string{
		"aws": {
			"3.0.0.0/15", "3.2.0.0/15", "3.4.0.0/15", "3.5.0.0/19",
			"13.32.0.0/15", "13.34.0.0/15", "13.36.0.0/14",
			"15.177.0.0/18", "15.193.0.0/20", "15.221.0.0/16",
			"18.64.0.0/14", "18.130.0.0/16", "18.132.0.0/14",
			"35.71.64.0/18", "35.71.128.0/17", "35.72.0.0/13",
			"52.0.0.0/15", "52.2.0.0/15", "52.4.0.0/14",
			"52.8.0.0/13", "52.16.0.0/15", "52.18.0.0/15",
			"52.20.0.0/14", "52.46.0.0/18", "52.56.0.0/16",
			"52.57.0.0/16", "52.58.0.0/15", "52.60.0.0/14",
			"52.64.0.0/17", "52.64.128.0/17", "52.65.0.0/16",
			"54.64.0.0/15", "54.66.0.0/16", "54.67.0.0/16",
			"54.68.0.0/14", "54.72.0.0/15", "54.74.0.0/15",
			"63.32.0.0/14", "63.34.0.0/15",
		},
		"azure": {
			"4.150.0.0/16", "4.151.0.0/16", "4.152.0.0/14",
			"13.64.0.0/11", "13.96.0.0/13", "13.104.0.0/14",
			"20.0.0.0/11", "20.32.0.0/11", "20.64.0.0/10",
			"20.128.0.0/16", "20.150.0.0/15", "20.160.0.0/12",
			"23.96.0.0/13", "40.64.0.0/10", "40.128.0.0/12",
			"51.4.0.0/15", "51.8.0.0/16", "51.10.0.0/15",
			"52.96.0.0/12", "52.112.0.0/14", "52.120.0.0/14",
			"52.224.0.0/11", "64.4.0.0/18",
			"65.52.0.0/14", "70.37.0.0/17", "70.37.128.0/18",
			"104.40.0.0/13", "104.208.0.0/13",
			"137.116.0.0/15", "137.117.0.0/16",
			"168.61.0.0/16", "168.62.0.0/15",
			"191.232.0.0/13",
		},
		"gcp": {
			"8.34.208.0/20", "8.35.192.0/20", "8.35.200.0/21",
			"23.236.48.0/20", "23.251.128.0/19",
			"34.0.0.0/15", "34.2.0.0/16", "34.3.0.0/16",
			"34.4.0.0/14", "34.8.0.0/13", "34.16.0.0/12",
			"34.32.0.0/11", "34.64.0.0/10", "34.128.0.0/10",
			"35.184.0.0/13", "35.192.0.0/14", "35.196.0.0/15",
			"35.198.0.0/16", "35.199.0.0/17", "35.199.128.0/18",
			"35.200.0.0/13", "35.208.0.0/12", "35.224.0.0/12",
			"35.240.0.0/13",
			"104.154.0.0/15", "104.196.0.0/14",
			"107.167.160.0/19", "107.178.192.0/18",
			"108.59.80.0/20", "108.170.192.0/18",
			"130.211.0.0/16", "142.250.0.0/15",
			"146.148.0.0/17", "162.216.148.0/22",
			"162.222.176.0/21", "172.110.32.0/21",
			"172.217.0.0/16", "172.253.0.0/16",
			"199.36.154.0/23", "199.36.156.0/24",
			"199.192.112.0/22", "199.223.232.0/21",
			"209.85.128.0/17", "216.58.192.0/19",
			"216.239.32.0/19",
		},
		"cloudflare": {
			"103.21.244.0/22", "103.22.200.0/22", "103.31.4.0/22",
			"104.16.0.0/13", "104.24.0.0/14",
			"108.162.192.0/18", "131.0.72.0/22",
			"141.101.64.0/18", "162.158.0.0/15",
			"172.64.0.0/13", "173.245.48.0/20",
			"188.114.96.0/20", "190.93.240.0/20",
			"197.234.240.0/22", "198.41.128.0/17",
		},
		"digitalocean": {
			"24.199.64.0/20", "24.199.80.0/20",
			"64.225.0.0/18", "64.225.64.0/18",
			"68.183.0.0/18", "68.183.64.0/18",
			"128.199.0.0/18", "128.199.64.0/18",
			"134.209.0.0/18", "134.209.64.0/18",
			"137.184.0.0/18", "137.184.64.0/18",
			"138.68.0.0/17", "138.197.0.0/17",
			"139.59.0.0/17", "142.93.0.0/17",
			"143.110.128.0/18", "143.198.0.0/17",
			"144.126.192.0/18", "146.190.0.0/17",
			"157.230.0.0/17", "157.245.0.0/17",
			"159.65.0.0/17", "159.89.0.0/17",
			"159.203.0.0/17", "161.35.0.0/17",
			"162.243.0.0/17", "163.47.8.0/22",
			"164.90.128.0/18", "164.92.64.0/18",
			"165.22.0.0/17", "165.227.0.0/17",
			"167.71.0.0/17", "167.172.0.0/17",
			"174.138.0.0/17",
		},
	}
}

func BenchmarkBuildTrie(b *testing.B) {
	data := generateRealisticData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Build(data)
	}
}

func BenchmarkBuildTrieLarge(b *testing.B) {
	data := make(map[string][]string)
	for i := 0; i < 200; i++ {
		var cidrs []string
		for j := 0; j < 100; j++ {
			cidrs = append(cidrs, fmt.Sprintf("%d.%d.0.0/16", i%256, j%256))
		}
		data[fmt.Sprintf("provider_%d", i)] = cidrs
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Build(data)
	}
}

func BenchmarkLookupRealisticHit(b *testing.B) {
	tr := loadEmbeddedTrie(b)
	ips := []string{
		"52.1.2.3", "34.100.50.25", "20.40.50.60",
		"104.16.5.100", "68.183.10.20", "35.200.1.1",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Lookup(ips[i%len(ips)])
	}
}

func BenchmarkLookupRealisticMiss(b *testing.B) {
	tr := loadEmbeddedTrie(b)
	ips := []string{
		"192.168.1.1", "10.0.0.1", "172.16.0.1",
		"8.8.8.8", "1.1.1.1", "127.0.0.1",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Lookup(ips[i%len(ips)])
	}
}

func BenchmarkLookupRealisticMixed(b *testing.B) {
	tr := loadEmbeddedTrie(b)
	ips := []string{
		"52.1.2.3", "192.168.1.1", "34.100.50.25",
		"10.0.0.1", "104.16.5.100", "8.8.8.8",
		"20.40.50.60", "172.16.0.1", "68.183.10.20",
		"1.1.1.1",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Lookup(ips[i%len(ips)])
	}
}

func BenchmarkLookupRealisticParallel(b *testing.B) {
	tr := loadEmbeddedTrie(b)
	ips := []string{
		"52.1.2.3", "192.168.1.1", "34.100.50.25",
		"10.0.0.1", "104.16.5.100", "8.8.8.8",
		"20.40.50.60", "172.16.0.1", "68.183.10.20",
		"1.1.1.1",
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			tr.Lookup(ips[i%len(ips)])
			i++
		}
	})
}

func BenchmarkSerializeEncode(b *testing.B) {
	tr := loadEmbeddedTrie(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		tr.Encode(&buf)
	}
}

func BenchmarkSerializeDecode(b *testing.B) {
	tr := loadEmbeddedTrie(b)
	var buf bytes.Buffer
	tr.Encode(&buf)
	raw := buf.Bytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(raw)
	}
}

func BenchmarkSerializeRoundTrip(b *testing.B) {
	tr := loadEmbeddedTrie(b)
	var buf bytes.Buffer
	tr.Encode(&buf)
	raw := buf.Bytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loaded, _ := Decode(raw)
		loaded.Lookup("52.1.2.3")
	}
}

func BenchmarkParseIPv4Valid(b *testing.B) {
	ips := []string{"192.168.1.100", "10.0.0.1", "255.255.255.255", "0.0.0.0"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseIPv4(ips[i%len(ips)])
	}
}

func BenchmarkParseIPv4Invalid(b *testing.B) {
	ips := []string{"not-an-ip", "256.0.0.0", "1.2.3", ""}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseIPv4(ips[i%len(ips)])
	}
}

func BenchmarkLookupRawUint32(b *testing.B) {
	tr := loadEmbeddedTrie(b)
	ip := binary.BigEndian.Uint32([]byte{52, 1, 2, 3})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.lookupRaw(ip)
	}
}

func BenchmarkLookupRandomIPs(b *testing.B) {
	tr := loadEmbeddedTrie(b)
	rng := rand.New(rand.NewSource(42))
	ips := make([]string, 10000)
	for i := range ips {
		ips[i] = fmt.Sprintf("%d.%d.%d.%d",
			rng.Intn(256), rng.Intn(256), rng.Intn(256), rng.Intn(256))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Lookup(ips[i%len(ips)])
	}
}

func BenchmarkTrieMemorySize(b *testing.B) {
	tr := loadEmbeddedTrie(b)
	nodeBytes := len(tr.nodes) * 10
	b.ReportMetric(float64(nodeBytes), "trie-bytes")
	b.ReportMetric(float64(len(tr.nodes)), "nodes")
	b.ReportMetric(float64(len(tr.Providers)-1), "providers")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Lookup("52.1.2.3")
	}
}

func BenchmarkBulkLookup(b *testing.B) {
	tr := loadEmbeddedTrie(b)
	rng := rand.New(rand.NewSource(42))
	ips := make([]string, 100000)
	for i := range ips {
		ips[i] = fmt.Sprintf("%d.%d.%d.%d",
			rng.Intn(256), rng.Intn(256), rng.Intn(256), rng.Intn(256))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, ip := range ips {
			tr.Lookup(ip)
		}
	}
}
