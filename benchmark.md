# Benchmarks

## Running benchmarks

Run the full suite:

```sh
go test -bench=. -benchmem ./internal/trie/
```

Run a specific benchmark:

```sh
go test -bench=BenchmarkLookupRealisticHit -benchmem ./internal/trie/
```

Run with a longer duration for more accurate results:

```sh
go test -bench=. -benchmem -benchtime=5s ./internal/trie/
```

Compare before/after a change using [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat):

```sh
go test -bench=. -benchmem -count=10 ./internal/trie/ > old.txt
# make changes
go test -bench=. -benchmem -count=10 ./internal/trie/ > new.txt
benchstat old.txt new.txt
```

## Benchmark suite

### Trie construction

| Benchmark | Description |
|-----------|-------------|
| `BenchmarkBuildTrie` | Build trie from realistic cloud provider CIDRs (~150 ranges across 5 providers) |
| `BenchmarkBuildTrieLarge` | Build trie from 20,000 CIDRs across 200 providers |

### Lookup

| Benchmark | Description |
|-----------|-------------|
| `BenchmarkLookupRealisticHit` | Lookup IPs that match a cloud provider |
| `BenchmarkLookupRealisticMiss` | Lookup IPs with no cloud match (private/reserved ranges) |
| `BenchmarkLookupRealisticMixed` | Alternating hits and misses |
| `BenchmarkLookupRealisticParallel` | Parallel lookups across all available goroutines |
| `BenchmarkLookupRandomIPs` | Lookup from a pool of 10k random IPs |
| `BenchmarkBulkLookup` | Sequentially look up 100k random IPs per iteration |
| `BenchmarkLookupRawUint32` | Raw trie traversal bypassing IP string parsing |
| `BenchmarkTrieMemorySize` | Lookup with trie node count and byte size reported as custom metrics |

### Serialization

| Benchmark | Description |
|-----------|-------------|
| `BenchmarkSerializeEncode` | Encode trie to binary format |
| `BenchmarkSerializeDecode` | Decode trie from binary format |
| `BenchmarkSerializeRoundTrip` | Full encode → decode → lookup cycle |

### IP parsing

| Benchmark | Description |
|-----------|-------------|
| `BenchmarkParseIPv4Valid` | Parse valid IPv4 strings to uint32 |
| `BenchmarkParseIPv4Invalid` | Parse invalid IPv4 strings (early exit paths) |

## Interpreting results

```
BenchmarkLookupRealisticHit-16    35480817    32.42 ns/op    0 B/op    0 allocs/op
```

- **35480817** — number of iterations run
- **32.42 ns/op** — time per lookup
- **0 B/op** — zero heap allocations per lookup
- **-16** — GOMAXPROCS (number of CPU cores used)

Key things to look for:
- Lookups should be zero-allocation (`0 B/op`)
- Parallel lookup should scale near-linearly with cores
- `BenchmarkTrieMemorySize` reports `trie-bytes`, `nodes`, and `providers` as custom metrics
