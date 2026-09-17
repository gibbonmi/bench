package chargeevidence

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
)

// Selection is the resolved selector row of one preparation.
type Selection struct {
	Mode, Spec, Ticket, Base, SourceTip string
}

// ManifestSource is one source descriptor of the canonical manifest.
type ManifestSource struct {
	ID, Role, Kind, Path string
	Required             bool
	Bytes                int
	SHA256               string
}

// Page is one source page descriptor. Offset is relative to its source body.
type Page struct {
	Source        string
	Index, Offset int
	Bytes         int
	SHA256        string
}

// ProducerRow records the declared collector of one generated source.
type ProducerRow struct {
	Source, Name, Version, Cwd string
}

// ArgumentRow records one ordered semantic collector argument.
type ArgumentRow struct {
	Source string
	Index  int
	Value  string
}

// Manifest is the typed canonical manifest.
type Manifest struct {
	Selection Selection
	Sources   []ManifestSource
	Pages     []Page
	Producers []ProducerRow
	Arguments []ArgumentRow
}

// SourceID returns the manifest identifier of the source at a one-based ordinal. The
// metadata source is always ordinal 1.
func SourceID(ordinal int) string { return "s" + strconv.Itoa(ordinal) }

// InputSourceID returns the manifest identifier of one candidate source input at its
// zero-based position in Candidate.Sources. The metadata source always precedes every
// input, so the first input is ordinal 2.
func InputSourceID(index int) string { return SourceID(index + 2) }

// Identity returns the artifact identifier of canonical manifest bytes.
func Identity(manifest []byte) string {
	return IdentityPrefix + digest(manifest)
}

// ValidIdentity reports whether value is a well-formed artifact identifier.
func ValidIdentity(value string) bool {
	hexPart, ok := strings.CutPrefix(value, IdentityPrefix)
	return ok && validDigest(hexPart)
}

func digest(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func validDigest(value string) bool {
	if len(value) != 64 {
		return false
	}
	for i := 0; i < len(value); i++ {
		c := value[i]
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}

func encodeManifest(m Manifest) ([]byte, error) {
	values := map[string][][]any{
		blockProfile:   {{ProfileVersion, HashName, PageBytes}},
		blockSelection: {{m.Selection.Mode, m.Selection.Spec, m.Selection.Ticket, m.Selection.Base, m.Selection.SourceTip}},
	}
	for _, s := range m.Sources {
		values[blockSources] = append(values[blockSources], []any{s.ID, s.Role, s.Kind, s.Path, s.Required, s.Bytes, s.SHA256})
	}
	for _, p := range m.Pages {
		values[blockPages] = append(values[blockPages], []any{p.Source, p.Index, p.Offset, p.Bytes, p.SHA256})
	}
	for _, p := range m.Producers {
		values[blockProducers] = append(values[blockProducers], []any{p.Source, p.Name, p.Version, p.Cwd})
	}
	for _, a := range m.Arguments {
		values[blockArguments] = append(values[blockArguments], []any{a.Source, a.Index, a.Value})
	}
	return encodeBlocks(ManifestBlocks, values)
}

// decodeManifest strictly decodes and semantically validates canonical manifest bytes.
func decodeManifest(data []byte) (Manifest, error) {
	values, err := decodeBlocks(ManifestBlocks, data)
	if err != nil {
		return Manifest{}, err
	}
	profile := values[blockProfile]
	if len(profile) != 1 {
		return Manifest{}, refuse(RefuseInvalid, "profile has %d rows, want 1", len(profile))
	}
	if profile[0][0].(int) != ProfileVersion {
		return Manifest{}, refuse(RefuseProfile, "manifest profile %d is not supported", profile[0][0].(int))
	}
	if profile[0][1].(string) != HashName || profile[0][2].(int) != PageBytes {
		return Manifest{}, refuse(RefuseProfile, "profile hash %q and page size %d are not supported", profile[0][1], profile[0][2])
	}
	selection := values[blockSelection]
	if len(selection) != 1 {
		return Manifest{}, refuse(RefuseInvalid, "selection has %d rows, want 1", len(selection))
	}
	var m Manifest
	s := selection[0]
	m.Selection = Selection{s[0].(string), s[1].(string), s[2].(string), s[3].(string), s[4].(string)}
	for _, r := range values[blockSources] {
		m.Sources = append(m.Sources, ManifestSource{r[0].(string), r[1].(string), r[2].(string), r[3].(string), r[4].(bool), r[5].(int), r[6].(string)})
	}
	for _, r := range values[blockPages] {
		m.Pages = append(m.Pages, Page{r[0].(string), r[1].(int), r[2].(int), r[3].(int), r[4].(string)})
	}
	for _, r := range values[blockProducers] {
		m.Producers = append(m.Producers, ProducerRow{r[0].(string), r[1].(string), r[2].(string), r[3].(string)})
	}
	for _, r := range values[blockArguments] {
		m.Arguments = append(m.Arguments, ArgumentRow{r[0].(string), r[1].(int), r[2].(string)})
	}
	return m, validateManifest(m)
}

func validateManifest(m Manifest) error {
	if err := validateSelection(m.Selection); err != nil {
		return err
	}
	if len(m.Sources) == 0 {
		return refuse(RefuseInvalid, "manifest declares no sources")
	}
	generated := map[string]bool{}
	for i, s := range m.Sources {
		if s.ID != SourceID(i+1) {
			return refuse(RefuseInvalid, "source %d has identifier %q, want %q", i+1, s.ID, SourceID(i+1))
		}
		if err := validateDescriptor(i, s.Role, s.Kind, s.Path, s.Required, s.Bytes); err != nil {
			return err
		}
		if !validDigest(s.SHA256) {
			return refuse(RefuseInvalid, "source %s digest is not 64 lowercase hexadecimal digits", s.ID)
		}
		if s.Kind == KindGenerated {
			generated[s.ID] = true
		}
	}
	if err := validatePages(m.Sources, m.Pages); err != nil {
		return err
	}
	return validateProvenance(m, generated)
}

func validateSelection(s Selection) error {
	switch {
	case s.Mode != "build" && s.Mode != "review":
		return refuse(RefuseInvalid, "selection mode %q is not build or review", s.Mode)
	case s.Spec == "" || s.Base == "" || s.SourceTip == "":
		return refuse(RefuseInvalid, "selection needs a spec, a base, and a source tip")
	case (s.Mode == "build") != (s.Ticket != ""):
		return refuse(RefuseInvalid, "build selection needs a ticket and review selection has none")
	}
	return nil
}

// validateDescriptor holds the source-kind rules the writer and the reader share.
func validateDescriptor(index int, role, kind, path string, required bool, size int) error {
	metadata := role == RoleMetadata
	switch {
	case (index == 0) != metadata:
		return refuse(RefuseInvalid, "source %d role %q: the metadata source is exactly the first source", index+1, role)
	case metadata && (kind != KindDerived || path != "" || !required):
		return refuse(RefuseInvalid, "metadata source must be required, derived, and pathless")
	case !metadata && kind != KindRepository && kind != KindGenerated:
		return refuse(RefuseInvalid, "source %d kind %q is not repository or generated", index+1, kind)
	case role == "" || path == "" && !metadata:
		return refuse(RefuseInvalid, "source %d needs a role and a path", index+1)
	case required && size == 0:
		return refuse(RefuseInvalid, "required source %d is empty", index+1)
	}
	return nil
}

func validatePages(sources []ManifestSource, pages []Page) error {
	next := 0
	for _, s := range sources {
		covered := 0
		for index := 0; next < len(pages) && pages[next].Source == s.ID; index, next = index+1, next+1 {
			p := pages[next]
			switch {
			case p.Index != index:
				return refuse(RefuseInvalid, "source %s page %d has index %d", s.ID, index, p.Index)
			case p.Offset < covered:
				return refuse(RefusePageOverlap, "source %s page %d starts at %d inside covered bytes %d", s.ID, index, p.Offset, covered)
			case p.Offset > covered:
				return refuse(RefusePageGap, "source %s page %d starts at %d after a gap from %d", s.ID, index, p.Offset, covered)
			case p.Bytes == 0 || p.Bytes > PageBytes:
				return refuse(RefuseInvalid, "source %s page %d holds %d bytes", s.ID, index, p.Bytes)
			case p.Offset+p.Bytes > s.Bytes:
				return refuse(RefusePageRange, "source %s page %d ends at %d beyond source length %d", s.ID, index, p.Offset+p.Bytes, s.Bytes)
			case !validDigest(p.SHA256):
				return refuse(RefuseInvalid, "source %s page %d digest is malformed", s.ID, index)
			}
			covered += p.Bytes
		}
		if covered != s.Bytes {
			return refuse(RefusePageGap, "source %s pages cover %d of %d bytes", s.ID, covered, s.Bytes)
		}
	}
	if next != len(pages) {
		return refuse(RefuseInvalid, "page %d names source %q outside source order", next, pages[next].Source)
	}
	return nil
}

func validateProvenance(m Manifest, generated map[string]bool) error {
	produced := 0
	for _, s := range m.Sources {
		if !generated[s.ID] {
			continue
		}
		if produced >= len(m.Producers) || m.Producers[produced].Source != s.ID {
			return refuse(RefuseInvalid, "generated source %s has no producer row in source order", s.ID)
		}
		produced++
	}
	if produced != len(m.Producers) {
		return refuse(RefuseInvalid, "producer rows name sources that are not generated")
	}
	order := map[string]int{}
	for i, s := range m.Sources {
		order[s.ID] = i
	}
	last, index := -1, 0
	for _, a := range m.Arguments {
		position, ok := order[a.Source]
		switch {
		case !ok || !generated[a.Source]:
			return refuse(RefuseInvalid, "argument row names source %q that is not generated", a.Source)
		case position < last:
			return refuse(RefuseInvalid, "argument rows leave source order at %s", a.Source)
		case position > last:
			last, index = position, 0
		}
		if a.Index != index {
			return refuse(RefuseInvalid, "source %s argument %d has index %d", a.Source, index, a.Index)
		}
		index++
	}
	return nil
}
