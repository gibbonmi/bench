package chargeevidence

import (
	"bytes"
	"encoding/binary"
	"math"
	"slices"
	"unicode/utf8"

	"github.com/gibbonmi/bench/internal/toon"
)

// Producer is the declared provenance of one generated source.
type Producer struct {
	Name, Version, Cwd string
	Arguments          []string
}

// SourceInput is one prepared canonical or generated source. The writer assigns its
// identifier from its position: the first input is ordinal 2, after the metadata source.
type SourceInput struct {
	Role, Kind, Path string
	Required         bool
	Data             []byte
	Producer         *Producer
}

// Candidate is one complete prepared evidence set before publication.
type Candidate struct {
	Selection Selection
	Metadata  Metadata
	Sources   []SourceInput
}

// Pack is a strictly validated evidence pack. Its accessors return copies, so a caller
// cannot change validated bytes.
type Pack struct {
	identity string
	data     []byte
	manifest Manifest
	metadata Metadata
	bodies   map[string][]byte
}

// Identity returns the artifact identifier.
func (p *Pack) Identity() string { return p.identity }

// Bytes returns the complete physical pack.
func (p *Pack) Bytes() []byte { return bytes.Clone(p.data) }

// ManifestBytes returns the canonical manifest bytes.
func (p *Pack) ManifestBytes() []byte {
	length := binary.LittleEndian.Uint64(p.data[16:HeaderBytes])
	return bytes.Clone(p.data[HeaderBytes : HeaderBytes+length])
}

// Manifest returns the decoded manifest.
func (p *Pack) Manifest() Manifest {
	m := p.manifest
	m.Sources = slices.Clone(m.Sources)
	m.Pages = slices.Clone(m.Pages)
	m.Producers = slices.Clone(m.Producers)
	m.Arguments = slices.Clone(m.Arguments)
	return m
}

// Metadata returns the decoded metadata source.
func (p *Pack) Metadata() Metadata {
	m := p.metadata
	m.Charge = slices.Clone(m.Charge)
	m.Fence = slices.Clone(m.Fence)
	m.Writes = slices.Clone(m.Writes)
	m.Coverage = slices.Clone(m.Coverage)
	m.Checks = slices.Clone(m.Checks)
	m.Returns = slices.Clone(m.Returns)
	m.Shared = slices.Clone(m.Shared)
	m.Completion = slices.Clone(m.Completion)
	return m
}

// Source returns the exact body of one declared source.
func (p *Pack) Source(id string) ([]byte, bool) {
	body, ok := p.bodies[id]
	return bytes.Clone(body), ok
}

// Build writes a candidate into an in-memory pack and returns it only after the strict
// reader accepts the complete result.
func Build(c Candidate) (*Pack, error) {
	metadata, err := EncodeMetadata(c.Metadata)
	if err != nil {
		return nil, err
	}
	inputs := append([]SourceInput{{Role: RoleMetadata, Kind: KindDerived, Required: true, Data: metadata}}, c.Sources...)
	m := Manifest{Selection: c.Selection}
	var bodies bytes.Buffer
	for i, input := range inputs {
		id := SourceID(i + 1)
		if err := validateInput(i, input); err != nil {
			return nil, err
		}
		m.Sources = append(m.Sources, ManifestSource{id, input.Role, input.Kind, input.Path, input.Required, len(input.Data), digest(input.Data)})
		offset := 0
		for index, size := range pageSizes(input.Data) {
			page := input.Data[offset : offset+size]
			m.Pages = append(m.Pages, Page{id, index, offset, size, digest(page)})
			offset += size
		}
		if input.Producer != nil {
			p := input.Producer
			m.Producers = append(m.Producers, ProducerRow{id, p.Name, p.Version, p.Cwd})
			for index, value := range p.Arguments {
				m.Arguments = append(m.Arguments, ArgumentRow{id, index, value})
			}
		}
		bodies.Write(input.Data)
	}
	if err := validateManifest(m); err != nil {
		return nil, refuse(RefuseCandidate, "%v", err)
	}
	manifest, err := encodeManifest(m)
	if err != nil {
		return nil, err
	}
	data := make([]byte, HeaderBytes, HeaderBytes+len(manifest)+bodies.Len())
	copy(data, HeaderMarker)
	binary.LittleEndian.PutUint32(data[8:12], ContainerVersion)
	binary.LittleEndian.PutUint64(data[16:24], uint64(len(manifest)))
	data = append(append(data, manifest...), bodies.Bytes()...)
	return Read(data, Identity(manifest))
}

func validateInput(index int, input SourceInput) error {
	if err := validateDescriptor(index, input.Role, input.Kind, input.Path, input.Required, len(input.Data)); err != nil {
		return refuse(RefuseCandidate, "%v", err)
	}
	if (input.Kind == KindGenerated) != (input.Producer != nil) {
		return refuse(RefuseCandidate, "source %d: exactly the generated sources declare a producer", index+1)
	}
	if !supportedText(input.Data) {
		return refuse(RefuseCandidate, "source %d holds bytes the shared TOON adapter cannot represent", index+1)
	}
	return nil
}

func supportedText(data []byte) bool {
	return utf8.Valid(data) && toon.Representable(string(data))
}

// pageSizes splits valid UTF-8 into the longest code-point-aligned prefixes of at most
// PageBytes bytes. An empty body has no pages.
func pageSizes(data []byte) []int {
	var sizes []int
	for len(data) > 0 {
		size := min(len(data), PageBytes)
		for size < len(data) && size > 0 && !utf8.RuneStart(data[size]) {
			size--
		}
		if size == 0 {
			size = min(len(data), PageBytes)
		}
		sizes = append(sizes, size)
		data = data[size:]
	}
	return sizes
}

// Read strictly validates a complete physical pack against the trusted expected identity.
// It derives every source offset from validated preceding lengths.
func Read(data []byte, expected string) (*Pack, error) {
	if !ValidIdentity(expected) {
		return nil, refuse(RefuseIdentifier, "expected identity is not sha256 followed by 64 lowercase hexadecimal digits")
	}
	if len(data) < HeaderBytes {
		return nil, refuse(RefuseTruncated, "pack holds %d bytes, fewer than the %d-byte header", len(data), HeaderBytes)
	}
	if string(data[:8]) != HeaderMarker {
		return nil, refuse(RefuseMarker, "header marker is not BENCHEV")
	}
	if version := binary.LittleEndian.Uint32(data[8:12]); version != ContainerVersion {
		return nil, refuse(RefuseContainerVersion, "container version %d is not supported", version)
	}
	if reserved := binary.LittleEndian.Uint32(data[12:16]); reserved != 0 {
		return nil, refuse(RefuseReserved, "reserved header value is %d", reserved)
	}
	length := binary.LittleEndian.Uint64(data[16:24])
	if length > math.MaxUint64-HeaderBytes {
		return nil, refuse(RefuseOverflow, "manifest length %d overflows the pack offset", length)
	}
	end := HeaderBytes + length
	if end > uint64(len(data)) {
		return nil, refuse(RefuseTruncated, "manifest ends at %d beyond the %d-byte pack", end, len(data))
	}
	manifestBytes := data[HeaderBytes:end]
	if Identity(manifestBytes) != expected {
		return nil, refuse(RefuseIdentity, "manifest digest differs from the expected identity")
	}
	m, err := decodeManifest(manifestBytes)
	if err != nil {
		return nil, err
	}
	total := end
	for _, s := range m.Sources {
		if uint64(s.Bytes) > math.MaxUint64-total {
			return nil, refuse(RefuseOverflow, "source %s length overflows the pack length", s.ID)
		}
		total += uint64(s.Bytes)
	}
	if total > uint64(len(data)) {
		return nil, refuse(RefuseTruncated, "sources end at %d beyond the %d-byte pack", total, len(data))
	}
	if total < uint64(len(data)) {
		return nil, refuse(RefuseTrailing, "pack holds %d bytes after its declared end %d", uint64(len(data))-total, total)
	}
	bodies, err := readBodies(data[end:], m)
	if err != nil {
		return nil, err
	}
	metadata, err := DecodeMetadata(bodies[SourceID(1)])
	if err != nil {
		return nil, err
	}
	for _, ref := range metadata.references() {
		if _, ok := bodies[ref]; !ok {
			return nil, refuse(RefuseInvalid, "metadata names undeclared source %q", ref)
		}
	}
	return &Pack{identity: expected, data: bytes.Clone(data), manifest: m, metadata: metadata, bodies: bodies}, nil
}

func readBodies(region []byte, m Manifest) (map[string][]byte, error) {
	bodies := make(map[string][]byte, len(m.Sources))
	start, next := 0, 0
	for _, s := range m.Sources {
		body := region[start : start+s.Bytes]
		start += s.Bytes
		var sizes []int
		for ; next < len(m.Pages) && m.Pages[next].Source == s.ID; next++ {
			p := m.Pages[next]
			if digest(body[p.Offset:p.Offset+p.Bytes]) != p.SHA256 {
				return nil, refuse(RefusePageDigest, "source %s page %d digest differs", s.ID, p.Index)
			}
			sizes = append(sizes, p.Bytes)
		}
		if digest(body) != s.SHA256 {
			return nil, refuse(RefuseSourceDigest, "source %s digest differs", s.ID)
		}
		if !supportedText(body) {
			return nil, refuse(RefuseSourceBytes, "source %s holds bytes the shared TOON adapter cannot represent", s.ID)
		}
		if !slices.Equal(sizes, pageSizes(body)) {
			return nil, refuse(RefusePageSplit, "source %s pages are not the canonical UTF-8 split", s.ID)
		}
		bodies[s.ID] = body
	}
	return bodies, nil
}
