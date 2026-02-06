package trie

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

var magic = [4]byte{'I', 'P', '2', 'C'}

const (
	version       = 1
	nodeRecordLen = 12
)

type header struct {
	Magic         [4]byte
	Version       uint16
	Reserved1     uint16
	NodeCount     uint32
	ProviderCount uint16
	Reserved2     uint16
}

func (t *Trie) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Encode(f)
}

func (t *Trie) Encode(w io.Writer) error {
	hdr := header{
		Magic:         magic,
		Version:       version,
		NodeCount:     uint32(len(t.nodes)),
		ProviderCount: uint16(len(t.Providers)),
	}
	if err := binary.Write(w, binary.LittleEndian, &hdr); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	for _, p := range t.Providers {
		nameBytes := []byte(p)
		if err := binary.Write(w, binary.LittleEndian, uint16(len(nameBytes))); err != nil {
			return fmt.Errorf("write provider name length: %w", err)
		}
		if _, err := w.Write(nameBytes); err != nil {
			return fmt.Errorf("write provider name: %w", err)
		}
	}

	buf := make([]byte, nodeRecordLen)
	for _, n := range t.nodes {
		binary.LittleEndian.PutUint32(buf[0:4], n.children[0])
		binary.LittleEndian.PutUint32(buf[4:8], n.children[1])
		binary.LittleEndian.PutUint16(buf[8:10], n.provider)
		buf[10] = 0
		buf[11] = 0
		if _, err := w.Write(buf); err != nil {
			return fmt.Errorf("write node: %w", err)
		}
	}

	return nil
}

func Load(path string) (*Trie, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Decode(data)
}

func Decode(data []byte) (*Trie, error) {
	if len(data) < 16 {
		return nil, fmt.Errorf("file too short")
	}

	if data[0] != magic[0] || data[1] != magic[1] || data[2] != magic[2] || data[3] != magic[3] {
		return nil, fmt.Errorf("invalid magic bytes")
	}

	ver := binary.LittleEndian.Uint16(data[4:6])
	if ver != version {
		return nil, fmt.Errorf("unsupported version %d", ver)
	}

	nodeCount := binary.LittleEndian.Uint32(data[8:12])
	providerCount := binary.LittleEndian.Uint16(data[12:14])

	pos := 16

	providers := make([]string, providerCount)
	provIndex := make(map[string]uint16, providerCount)
	for i := uint16(0); i < providerCount; i++ {
		if pos+2 > len(data) {
			return nil, fmt.Errorf("truncated provider table")
		}
		nameLen := int(binary.LittleEndian.Uint16(data[pos : pos+2]))
		pos += 2
		if pos+nameLen > len(data) {
			return nil, fmt.Errorf("truncated provider name")
		}
		name := string(data[pos : pos+nameLen])
		providers[i] = name
		if name != "" {
			provIndex[name] = i
		}
		pos += nameLen
	}

	need := int(nodeCount) * nodeRecordLen
	if pos+need > len(data) {
		return nil, fmt.Errorf("truncated node array")
	}

	nodes := make([]node, nodeCount)
	for i := uint32(0); i < nodeCount; i++ {
		off := pos + int(i)*nodeRecordLen
		nodes[i] = node{
			children: [2]uint32{
				binary.LittleEndian.Uint32(data[off : off+4]),
				binary.LittleEndian.Uint32(data[off+4 : off+8]),
			},
			provider: binary.LittleEndian.Uint16(data[off+8 : off+10]),
		}
	}

	return &Trie{
		nodes:     nodes,
		nextFree:  nodeCount,
		Providers: providers,
		provIndex: provIndex,
	}, nil
}
