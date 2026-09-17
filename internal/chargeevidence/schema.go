// Package chargeevidence owns the immutable charge evidence pack: its format registry,
// the canonical manifest, the in-memory pack writer, and the strict pack reader.
//
// The package accepts prepared values. It never selects sources, runs collectors, or
// infers approval; preflight owns those policies. This file is the one executable format
// registry. The writer, the strict reader, and the generated format reference all read it.
package chargeevidence

// CellType is the declared TOON type of one registered column.
type CellType int

const (
	// String cells always decode as strings, including numeric-looking values.
	String CellType = iota
	// Integer cells hold unsigned decimal values without leading zeroes.
	Integer
	// Boolean cells hold typed TOON booleans.
	Boolean
)

func (t CellType) String() string {
	switch t {
	case Integer:
		return "integer"
	case Boolean:
		return "boolean"
	}
	return "string"
}

// Field is one registered column.
type Field struct {
	Name string
	Type CellType
}

// Block is one registered flat TOON table. Rows describes its cardinality and order for
// the generated reference.
type Block struct {
	Name   string
	Fields []Field
	Rows   string
}

// FieldNames returns the block's column names in their declared order.
func (b Block) FieldNames() []string {
	names := make([]string, len(b.Fields))
	for i, field := range b.Fields {
		names[i] = field.Name
	}
	return names
}

// HeaderField is one registered byte range of the physical pack header.
type HeaderField struct {
	First, Last int
	Meaning     string
}

const (
	// HeaderMarker is the seven ASCII marker bytes followed by one NUL byte.
	HeaderMarker = "BENCHEV\x00"
	// HeaderBytes is the fixed physical header length.
	HeaderBytes = 24
	// ContainerVersion is the only supported physical pack version.
	ContainerVersion = 1
	// ProfileVersion is the only supported canonical manifest profile.
	ProfileVersion = 1
	// HashName names the digest algorithm in the profile block.
	HashName = "sha256"
	// PageBytes is the maximum raw byte count of one source page.
	PageBytes = 8192
	// IdentityPrefix precedes the hexadecimal manifest digest in an artifact identifier.
	IdentityPrefix = "sha256:"
)

// Source roles and kinds that the format itself defines. Preflight owns every canonical
// and generated role name.
const (
	RoleMetadata   = "metadata"
	KindRepository = "repository"
	KindGenerated  = "generated"
	KindDerived    = "derived"
)

// Header is the registered physical header layout.
var Header = []HeaderField{
	{0, 7, "The seven ASCII bytes `BENCHEV`, then one NUL byte."},
	{8, 11, "Unsigned little-endian container version 1."},
	{12, 15, "Reserved unsigned little-endian value 0."},
	{16, 23, "Unsigned little-endian manifest byte length."},
}

func str(name string) Field  { return Field{name, String} }
func num(name string) Field  { return Field{name, Integer} }
func flag(name string) Field { return Field{name, Boolean} }

// Manifest block names in their canonical order.
const (
	blockProfile   = "profile"
	blockSelection = "selection"
	blockSources   = "sources"
	blockPages     = "pages"
	blockProducers = "producers"
	blockArguments = "arguments"
)

// ManifestBlocks is the registered canonical manifest profile 1.
var ManifestBlocks = []Block{
	{blockProfile, []Field{num("version"), str("hash"), num("page_bytes")},
		"One row: integer 1, string sha256, integer 8192."},
	{blockSelection, []Field{str("mode"), str("spec"), str("ticket"), str("base"), str("source_tip")},
		"One row. Review uses an empty ticket string."},
	{blockSources, []Field{str("id"), str("role"), str("kind"), str("path"), flag("required"), num("bytes"), str("sha256")},
		"Source order. The first source is the metadata source."},
	{blockPages, []Field{str("source"), num("index"), num("offset"), num("bytes"), str("sha256")},
		"Source order, then increasing page index."},
	{blockProducers, []Field{str("source"), str("name"), str("version"), str("cwd")},
		"Source order. Only generated sources have a producer row."},
	{blockArguments, []Field{str("source"), num("index"), str("value")},
		"Source order, then argument index."},
}

// Metadata block names in their canonical order.
const (
	blockCharge     = "charge"
	blockFence      = "fence"
	blockWrites     = "writes"
	blockCoverage   = "coverage"
	blockChecks     = "checks"
	blockReturns    = "returns"
	blockShared     = "shared_evidence"
	blockCompletion = "completion_evidence"
)

// MetadataBlocks is the registered schema of the required metadata source.
var MetadataBlocks = []Block{
	{blockCharge, []Field{str("axis"), str("ticket"), str("access")},
		"One build row with an empty axis, or one row per review axis in axis order."},
	{blockFence, []Field{str("path")}, "Declared spec fence order."},
	{blockWrites, []Field{str("path")}, "Selected ticket write order. Review has zero rows."},
	{blockCoverage, []Field{str("row")}, "Selected ticket coverage order. Review has zero rows."},
	{blockChecks, []Field{str("source")}, "Check-source order."},
	{blockReturns, []Field{str("source")}, "Return-source order."},
	{blockShared, []Field{str("kind"), str("source")}, "Diff, consumers, coverage order. Build has zero rows."},
	{blockCompletion, []Field{str("record"), str("source_digest"), str("plan_digest"), str("record_state"), str("detail")},
		"Completion facts. Build has zero rows."},
}
