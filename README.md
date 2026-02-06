<h1 align="center">
    ip2cloud
  <br>
</h1>

<h4 align="center">Check IP addresses against known cloud provider IP address ranges</h4>

<p align="center">
  <a href="#installation">Installation</a>  
  <a href="#usage">Usage</a>  
  <a href="#commands">Commands</a>  
  <a href="#adding-custom-providers">Adding Providers</a>  
  <a href="#benchmarks">Benchmarks</a>
  <br>
</p>

![ip2cloud](https://github.com/devanshbatham/ip2cloud/blob/main/static/banner.png?raw=true)

## Installation

```sh
go install github.com/devanshbatham/ip2cloud/cmd/ip2cloud@latest
```

No other setup needed. Cloud provider IP ranges are embedded in the binary and the trie is built automatically on first run.

## Usage

```sh
cat ips.txt | ip2cloud
```

```sh
ip2cloud 1.2.3.4 5.6.7.8
```

```sh
ip2cloud -p aws,azure < ips.txt
```

```sh
ip2cloud -j < ips.txt
```

### Example

Given a file `ips.txt`:

```
59.82.33.201
63.32.40.140
63.33.205.240
64.4.8.90
64.4.8.67
```

Text output:

```
$ cat ips.txt | ip2cloud

[aliyun] : 59.82.33.201
[aws] : 63.32.40.140
[aws] : 63.33.205.240
[azure] : 64.4.8.67
[azure] : 64.4.8.90
```

JSON output:

```
$ cat ips.txt | ip2cloud -j

{
    "aliyun": [
        "59.82.33.201"
    ],
    "aws": [
        "63.32.40.140",
        "63.33.205.240"
    ],
    "azure": [
        "64.4.8.90",
        "64.4.8.67"
    ]
}
```

IPs with no matching cloud provider are omitted from the output.

## Commands

| Command | Description |
|---------|-------------|
| `ip2cloud` | Lookup IPs from stdin or args (default) |
| `ip2cloud build` | Rebuild binary trie (auto-seeds from embedded data if no `-seed` flag) |
| `ip2cloud build -seed ./data` | Seed from a custom directory of `.txt` files |
| `ip2cloud add <provider> [-f file] [cidrs...]` | Add CIDR ranges to a provider |
| `ip2cloud list` | List providers and range counts |
| `ip2cloud version` | Print version |

## Lookup Flags

| Flag | Description |
|------|-------------|
| `-p`, `-provider` | Comma-separated provider filter (e.g. `aws,gcp`) |
| `-j`, `-json` | JSON output |
| `-w` | Worker count (default: NumCPU) |

## Data Storage

Provider data files and the binary trie are stored under `~/.config/ip2cloud/`:

```
~/.config/ip2cloud/
  data/          # provider .txt files (one CIDR per line)
  ip2cloud.bin   # compiled binary trie
```

The trie is built automatically on first lookup. Run `ip2cloud build` to rebuild it manually after modifying provider data.

## Adding Custom Providers

You can add your own cloud provider or update existing ones using the `add` command.

### CIDR format

Provider data files are plain text with one IPv4 CIDR per line. Lines starting with `#` are ignored. Example:

```
# My custom provider ranges
10.100.0.0/16
10.200.0.0/14
172.20.0.0/15
```

> **Note:** Only IPv4 CIDRs are supported. IPv6 entries are silently skipped.

### Adding ranges via CLI

```sh
# Add CIDRs as arguments
ip2cloud add myprovider 10.100.0.0/16 10.200.0.0/14

# Add CIDRs from a file
ip2cloud add myprovider -f ranges.txt

# Add CIDRs from stdin
cat ranges.txt | ip2cloud add myprovider -f -

# Add without auto-rebuilding the trie
ip2cloud add myprovider -build=false 10.100.0.0/16
```

The trie is rebuilt automatically after adding ranges (disable with `-build=false`).

### Adding ranges manually

You can also create or edit provider files directly under `~/.config/ip2cloud/data/`:

```sh
# Create a new provider
echo "10.100.0.0/16" > ~/.config/ip2cloud/data/myprovider.txt

# Rebuild the trie after manual edits
ip2cloud build
```

The file name (without `.txt`) becomes the provider name used in lookup output.

### Seeding from a directory

To replace all provider data from a custom directory:

```sh
ip2cloud build -seed /path/to/my/data/
```

The directory should contain `.txt` files named after each provider (e.g., `aws.txt`, `myprovider.txt`), each with one CIDR per line.

## Benchmarks

See [benchmark.md](benchmark.md) for the full benchmark suite and instructions.

## Limitations

- IPv4 only (IPv6 CIDRs are skipped)
- Output order is nondeterministic when using multiple workers

## Supported Cloud Providers

- [x] Alibaba Cloud (Aliyun)
- [x] Amazon Web Services (AWS)
- [x] Microsoft Azure
- [x] Cloudflare
- [x] DigitalOcean
- [x] Fastly
- [x] Google Cloud
- [x] IBM Cloud
- [x] Linode
- [x] Oracle Cloud
- [x] Tencent Cloud
- [x] UCloud
