package chargeevidence

import (
	"fmt"
	"strings"
)

// ReferencePath is the shipped format reference that FormatReference generates.
const ReferencePath = ".agents/skills/bench-craft-delegate/references/charge-evidence-format.md"

// FormatReference renders the shipped format reference from the registry. The shipped
// file must equal this projection byte for byte.
func FormatReference() string {
	var b strings.Builder
	b.WriteString("# Charge evidence format\n\n")
	b.WriteString("This reference describes the charge evidence pack.\n")
	b.WriteString("The Bench format registry generates this file, so do not edit it by hand.\n")
	b.WriteString("An independent reader can use this reference without Bench code.\n\n")

	b.WriteString("## Physical pack\n\n")
	b.WriteString("One regular file holds the header, the canonical manifest, and the raw source bodies in manifest order.\n")
	fmt.Fprintf(&b, "The header has %d bytes.\n", HeaderBytes)
	b.WriteString("The reader derives each source offset from the validated preceding lengths.\n")
	b.WriteString("The expected total length must equal the file length, so the reader refuses truncation and trailing data.\n\n")
	b.WriteString("| Header bytes | Meaning |\n| --- | --- |\n")
	for _, field := range Header {
		fmt.Fprintf(&b, "| %d through %d | %s |\n", field.First, field.Last, field.Meaning)
	}

	b.WriteString("\n## Canonical manifest\n\n")
	fmt.Fprintf(&b, "The manifest uses profile %d.\n", ProfileVersion)
	fmt.Fprintf(&b, "The artifact identifier is `%s` followed by the lowercase hexadecimal SHA-256 digest of the manifest bytes.\n", IdentityPrefix)
	b.WriteString("The manifest is a sequence of flat TOON tables in the order below.\n")
	b.WriteString("Every table appears once, and a table without rows keeps its header.\n")
	b.WriteString("The encoding uses two-space indentation, comma delimiters, and one final line feed per table.\n")
	b.WriteString("The reader decodes the manifest, encodes it again, and refuses bytes that differ.\n\n")
	writeBlocks(&b, "Table", ManifestBlocks)

	b.WriteString("\nSource identifiers are `s` followed by the one-based source position.\n")
	fmt.Fprintf(&b, "A source page holds at most %d raw bytes.\n", PageBytes)
	b.WriteString("Each page is the longest UTF-8 prefix within that limit.\n")
	b.WriteString("A page offset is relative to its own source body.\n\n")
	b.WriteString("Pages cover each source exactly once, with no gap and no overlap.\n")
	b.WriteString("An empty optional source has no page.\n")
	b.WriteString("Every digest cell holds 64 lowercase hexadecimal digits as a string.\n")
	b.WriteString("A page index and an argument index both start at zero.\n")

	b.WriteString("\n## Canonical string quoting\n\n")
	b.WriteString("The pinned shared TOON encoder owns every string quoting and escaping rule.\n")
	b.WriteString("`internal/toon/toon_test.go`'s `TestTableCellEscaping` pins its complete trigger inventory.\n")
	b.WriteString("This reference does not restate a partial trigger list.\n")

	b.WriteString("\n## Metadata source\n\n")
	fmt.Fprintf(&b, "The first source has the role `%s`, the kind `%s`, and an empty path.\n", RoleMetadata, KindDerived)
	fmt.Fprintf(&b, "Repository sources have the kind `%s`, and generated sources have the kind `%s`.\n", KindRepository, KindGenerated)
	b.WriteString("The metadata body uses the same canonical table rules as the manifest.\n")
	b.WriteString("Every source cell holds a manifest source identifier.\n")
	fmt.Fprintf(&b, "A build charge row uses access `%s` and names the selected ticket source in its ticket cell.\n", AccessBuild)
	fmt.Fprintf(&b, "A review charge row uses access `%s` and names the spec source in its ticket cell.\n\n", AccessReview)
	writeBlocks(&b, "Table", MetadataBlocks)
	return b.String()
}

func writeBlocks(b *strings.Builder, heading string, blocks []Block) {
	fmt.Fprintf(b, "| %s | Ordered fields | Rows |\n| --- | --- | --- |\n", heading)
	for _, block := range blocks {
		fields := make([]string, len(block.Fields))
		for i, field := range block.Fields {
			fields[i] = fmt.Sprintf("`%s` (%s)", field.Name, field.Type)
		}
		fmt.Fprintf(b, "| `%s` | %s | %s |\n", block.Name, strings.Join(fields, ", "), block.Rows)
	}
}
