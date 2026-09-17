# Charge evidence format

This reference describes the charge evidence pack.
The Bench format registry generates this file, so do not edit it by hand.
An independent reader can use this reference without Bench code.

## Physical pack

One regular file holds the header, the canonical manifest, and the raw source bodies in manifest order.
The header has 24 bytes.
The reader derives each source offset from the validated preceding lengths.
The expected total length must equal the file length, so the reader refuses truncation and trailing data.

| Header bytes | Meaning |
| --- | --- |
| 0 through 7 | The seven ASCII bytes `BENCHEV`, then one NUL byte. |
| 8 through 11 | Unsigned little-endian container version 1. |
| 12 through 15 | Reserved unsigned little-endian value 0. |
| 16 through 23 | Unsigned little-endian manifest byte length. |

## Canonical manifest

The manifest uses profile 1.
The artifact identifier is `sha256:` followed by the lowercase hexadecimal SHA-256 digest of the manifest bytes.
The manifest is a sequence of flat TOON tables in the order below.
Every table appears once, and a table without rows keeps its header.
The encoding uses two-space indentation, comma delimiters, and one final line feed per table.
The reader decodes the manifest, encodes it again, and refuses bytes that differ.

| Table | Ordered fields | Rows |
| --- | --- | --- |
| `profile` | `version` (integer), `hash` (string), `page_bytes` (integer) | One row: integer 1, string sha256, integer 8192. |
| `selection` | `mode` (string), `spec` (string), `ticket` (string), `base` (string), `source_tip` (string) | One row. Review uses an empty ticket string. |
| `sources` | `id` (string), `role` (string), `kind` (string), `path` (string), `required` (boolean), `bytes` (integer), `sha256` (string) | Source order. The first source is the metadata source. |
| `pages` | `source` (string), `index` (integer), `offset` (integer), `bytes` (integer), `sha256` (string) | Source order, then increasing page index. |
| `producers` | `source` (string), `name` (string), `version` (string), `cwd` (string) | Source order. Only generated sources have a producer row. |
| `arguments` | `source` (string), `index` (integer), `value` (string) | Source order, then argument index. |

Source identifiers are `s` followed by the one-based source position.
A source page holds at most 8192 raw bytes.
Each page is the longest UTF-8 prefix within that limit.
A page offset is relative to its own source body.

Pages cover each source exactly once, with no gap and no overlap.
An empty optional source has no page.
Every digest cell holds 64 lowercase hexadecimal digits as a string.
A page index and an argument index both start at zero.

## Canonical string quoting

The pinned shared TOON encoder owns every string quoting and escaping rule.
`internal/toon/toon_test.go`'s `TestTableCellEscaping` pins its complete trigger inventory.
This reference does not restate a partial trigger list.

## Metadata source

The first source has the role `metadata`, the kind `derived`, and an empty path.
Repository sources have the kind `repository`, and generated sources have the kind `generated`.
The metadata body uses the same canonical table rules as the manifest.
Every source cell holds a manifest source identifier.
A build charge row uses access `write-within-fence` and names the selected ticket source in its ticket cell.
A review charge row uses access `read-only` and names the spec source in its ticket cell.

| Table | Ordered fields | Rows |
| --- | --- | --- |
| `charge` | `axis` (string), `ticket` (string), `access` (string) | One build row with empty axis, or the existing review-axis inventory order. |
| `fence` | `path` (string) | Declared spec fence order. |
| `writes` | `path` (string) | Selected ticket write order. Review has zero rows. |
| `coverage` | `row` (string) | Selected ticket coverage order. Review has zero rows here. |
| `checks` | `source` (string) | Existing build check-source order or review skill source. |
| `returns` | `source` (string) | Existing return-source order. |
| `shared_evidence` | `kind` (string), `source` (string) | Diff, consumers, coverage order. Build has zero rows. |
| `completion_evidence` | `record` (string), `source_digest` (string), `plan_digest` (string), `record_state` (string), `detail` (string) | Existing completion facts. Build has zero rows. |
