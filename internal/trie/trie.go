package trie

import (
	"encoding/binary"
	"fmt"
	"net"
	"sort"
)

const emptyNode = 0

type node struct {
	children [2]uint32
	provider uint16
}

type Trie struct {
	nodes     []node
	nextFree  uint32
	Providers []string
	provIndex map[string]uint16
	Warnings  []string
}

func New() *Trie {
	return &Trie{
		nodes:     make([]node, 1, 1<<16),
		nextFree:  1,
		Providers: []string{""},
		provIndex: make(map[string]uint16),
	}
}

func Build(cloudData map[string][]string) *Trie {
	t := New()

	providers := make([]string, 0, len(cloudData))
	for p := range cloudData {
		providers = append(providers, p)
	}
	sort.Strings(providers)

	for _, provider := range providers {
		cidrs := cloudData[provider]
		idx := t.providerIdx(provider)
		for _, cidr := range cidrs {
			_, network, err := net.ParseCIDR(cidr)
			if err != nil {
				t.Warnings = append(t.Warnings, fmt.Sprintf("%s: invalid CIDR %q", provider, cidr))
				continue
			}
			ip4 := network.IP.To4()
			if ip4 == nil {
				t.Warnings = append(t.Warnings, fmt.Sprintf("%s: skipping IPv6 CIDR %s", provider, cidr))
				continue
			}
			ones, _ := network.Mask.Size()
			ip := binary.BigEndian.Uint32(ip4)
			t.insert(ip, ones, idx)
		}
	}

	t.nodes = t.nodes[:t.nextFree]
	return t
}

func (t *Trie) providerIdx(name string) uint16 {
	if idx, ok := t.provIndex[name]; ok {
		return idx
	}
	idx := uint16(len(t.Providers))
	t.Providers = append(t.Providers, name)
	t.provIndex[name] = idx
	return idx
}

func (t *Trie) alloc() uint32 {
	id := t.nextFree
	t.nextFree++
	if int(t.nextFree) >= len(t.nodes) {
		t.nodes = append(t.nodes, make([]node, len(t.nodes))...)
	}
	return id
}

func (t *Trie) insert(ip uint32, prefixLen int, provider uint16) {
	cur := uint32(0)
	for i := 31; i >= 32-prefixLen; i-- {
		bit := (ip >> uint(i)) & 1
		child := t.nodes[cur].children[bit]
		if child == emptyNode {
			child = t.alloc()
			t.nodes[cur].children[bit] = child
		}
		cur = child
	}
	t.nodes[cur].provider = provider
}

func (t *Trie) Lookup(ipStr string) string {
	ip, ok := ParseIPv4(ipStr)
	if !ok {
		return ""
	}
	return t.Providers[t.lookupRaw(ip)]
}

func (t *Trie) lookupRaw(ip uint32) uint16 {
	var match uint16
	cur := uint32(0)
	nodes := t.nodes
	for i := 31; i >= 0; i-- {
		bit := (ip >> uint(i)) & 1
		child := nodes[cur].children[bit]
		if child == emptyNode {
			break
		}
		if nodes[child].provider != 0 {
			match = nodes[child].provider
		}
		cur = child
	}
	return match
}

func ParseIPv4(s string) (uint32, bool) {
	var ip uint32
	var octet uint32
	var dots int
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			octet = octet*10 + uint32(c-'0')
			if octet > 255 {
				return 0, false
			}
		} else if c == '.' {
			ip = (ip << 8) | octet
			octet = 0
			dots++
			if dots > 3 {
				return 0, false
			}
		} else {
			return 0, false
		}
	}
	if dots != 3 {
		return 0, false
	}
	ip = (ip << 8) | octet
	return ip, true
}
