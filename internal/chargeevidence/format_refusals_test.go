package chargeevidence_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"testing"

	ce "github.com/gibbonmi/bench/internal/chargeevidence"
)

// repack rebuilds a physical pack around an edited manifest, with the original bodies.
func repack(t *testing.T, pack *ce.Pack, edit func(string) string) ([]byte, string) {
	t.Helper()
	old := pack.ManifestBytes()
	manifest := edit(string(old))
	if manifest == string(old) {
		t.Fatalf("manifest edit changed nothing")
	}
	data := pack.Bytes()
	header := bytes.Clone(data[:24])
	binary.LittleEndian.PutUint64(header[16:], uint64(len(manifest)))
	out := append(append(header, manifest...), data[24+len(old):]...)
	return out, "sha256:" + sum(manifest)
}

// tamperSourceBody rebuilds a pack around one raw source body, replacing old with a new
// body of the same length and updating every manifest occurrence of its digest to match.
// A single-page body's digest occurs once in its source row and once in its page row.
func tamperSourceBody(t *testing.T, pack *ce.Pack, old, new string) ([]byte, string) {
	t.Helper()
	if len(old) != len(new) {
		t.Fatalf("tamper changes body length")
	}
	manifestOld := pack.ManifestBytes()
	oldDigest, newDigest := sum(old), sum(new)
	if strings.Count(string(manifestOld), oldDigest) == 0 {
		t.Fatalf("digest %s does not appear in the manifest", oldDigest)
	}
	manifest := strings.ReplaceAll(string(manifestOld), oldDigest, newDigest)
	data := pack.Bytes()
	if count := bytes.Count(data, []byte(old)); count != 1 {
		t.Fatalf("body %q appears %d times in the pack, want once", old, count)
	}
	data = bytes.Replace(data, []byte(old), []byte(new), 1)
	header := bytes.Clone(data[:24])
	binary.LittleEndian.PutUint64(header[16:], uint64(len(manifest)))
	out := append(append(header, manifest...), data[24+len(manifestOld):]...)
	return out, "sha256:" + sum(manifest)
}

func assertRefusal(t *testing.T, data []byte, identity, class string) {
	t.Helper()
	_, err := ce.Read(data, identity)
	var refusal *ce.Refusal
	if !errors.As(err, &refusal) || refusal.Class != class {
		t.Fatalf("Read = %v, want refusal class %s", err, class)
	}
}

func replaceOnce(t *testing.T, old, new string) func(string) string {
	return func(s string) string {
		t.Helper()
		if strings.Count(s, old) != 1 {
			t.Fatalf("manifest holds %q %d times, want once", old, strings.Count(s, old))
		}
		return strings.Replace(s, old, new, 1)
	}
}

func TestEvidenceManifestRefusals(t *testing.T) {
	pack := mustBuild(t, fixtureCandidate("# One\n"))
	for _, test := range []struct {
		name, class string
		edit        func(t *testing.T) func(string) string
	}{
		{"CE29 profile version", "unsupported-profile", func(t *testing.T) func(string) string {
			return replaceOnce(t, "  1,sha256,8192\n", "  2,sha256,8192\n")
		}},
		{"CE29 page size", "unsupported-profile", func(t *testing.T) func(string) string {
			return replaceOnce(t, "  1,sha256,8192\n", "  1,sha256,4096\n")
		}},
		{"CE30 duplicate block", "noncanonical", func(t *testing.T) func(string) string {
			return func(s string) string { return s + "arguments[0]{source,index,value}:\n" }
		}},
		{"CE30 duplicate field", "noncanonical", func(t *testing.T) func(string) string {
			return replaceOnce(t, "arguments[0]{source,index,value}:", "arguments[0]{source,index,value,value}:")
		}},
		{"CE30 duplicate identifier", "invalid-manifest", func(t *testing.T) func(string) string {
			return func(s string) string { return strings.Replace(s, "  s3,spec,", "  s2,spec,", 1) }
		}},
		{"CE30 unknown block", "schema", func(t *testing.T) func(string) string {
			return func(s string) string { return s + "extra[0]{value}:\n" }
		}},
		{"CE30 reordered blocks", "noncanonical", func(t *testing.T) func(string) string {
			return replaceOnce(t, "producers[0]{source,name,version,cwd}:\narguments[0]{source,index,value}:\n", "arguments[0]{source,index,value}:\nproducers[0]{source,name,version,cwd}:\n")
		}},
		{"CE30 reordered columns", "noncanonical", func(t *testing.T) func(string) string {
			return replaceOnce(t, "profile[1]{version,hash,page_bytes}:\n  1,sha256,8192", "profile[1]{hash,version,page_bytes}:\n  sha256,1,8192")
		}},
		{"CE30 blank separator", "noncanonical", func(t *testing.T) func(string) string {
			return replaceOnce(t, "\nselection[1]", "\n\nselection[1]")
		}},
		{"CE30 wider indentation", "noncanonical", func(t *testing.T) func(string) string {
			return replaceOnce(t, "  1,sha256,8192", "    1,sha256,8192")
		}},
		{"CE30 needless quotes", "noncanonical", func(t *testing.T) func(string) string {
			return replaceOnce(t, "  s2,ticket,", "  \"s2\",ticket,")
		}},
		{"CE30 missing final newline", "noncanonical", func(t *testing.T) func(string) string {
			return func(s string) string { return strings.TrimSuffix(s, "\n") }
		}},
		{"CE126 quoted boolean", "cell-type", func(t *testing.T) func(string) string {
			return func(s string) string {
				return strings.Replace(s, ",ticket,repository,specs/example/tickets/one.md,true,", ",ticket,repository,specs/example/tickets/one.md,\"true\",", 1)
			}
		}},
		{"CE126 quoted integer", "cell-type", func(t *testing.T) func(string) string {
			return replaceOnce(t, "  1,sha256,8192", "  \"1\",sha256,8192")
		}},
		{"CE126 fractional integer", "cell-type", func(t *testing.T) func(string) string {
			return replaceOnce(t, "  1,sha256,8192", "  1.5,sha256,8192")
		}},
		{"CE126 negative integer", "cell-type", func(t *testing.T) func(string) string {
			return func(s string) string { return strings.Replace(s, "  s2,0,0,", "  s2,-1,0,", 1) }
		}},
		{"CE126 numeric digest", "cell-type", func(t *testing.T) func(string) string {
			ticket := sum("# One\n")
			return func(s string) string { return strings.Replace(s, ","+ticket+"\n", ",1234\n", 1) }
		}},
		{"CV4 integer upper bound", "cell-type", func(t *testing.T) func(string) string {
			return replaceOnce(t, "  1,sha256,8192", "  1,sha256,9007199254740992")
		}},
		{"CV4 exponent form", "noncanonical", func(t *testing.T) func(string) string {
			return replaceOnce(t, "  1,sha256,8192", "  1,sha256,8192e0")
		}},
		{"CV4 negative zero", "noncanonical", func(t *testing.T) func(string) string {
			return func(s string) string { return strings.Replace(s, "  s2,0,0,", "  s2,-0,0,", 1) }
		}},
		{"CV4 fractional-form integer", "noncanonical", func(t *testing.T) func(string) string {
			return replaceOnce(t, "  1,sha256,8192", "  1.0,sha256,8192")
		}},
		{"CV4 leading zero", "cell-type", func(t *testing.T) func(string) string {
			return replaceOnce(t, "  1,sha256,8192", "  1,sha256,08192")
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			data, identity := repack(t, pack, test.edit(t))
			assertRefusal(t, data, identity, test.class)
		})
	}
}

func TestEvidencePackRefusals(t *testing.T) {
	big := strings.Repeat("x", 10000)
	pack := mustBuild(t, fixtureCandidate(big))
	identity := pack.Identity()
	header := func(offset int, put func([]byte)) []byte {
		data := pack.Bytes()
		put(data[offset:])
		return data
	}
	secondPage := fmt.Sprintf("  s2,1,8192,1808,%s\n", sum(big[8192:]))
	for _, test := range []struct {
		name, class string
		data        func(t *testing.T) ([]byte, string)
	}{
		{"CE31 container version", "unsupported-container-version", func(*testing.T) ([]byte, string) {
			return header(8, func(b []byte) { binary.LittleEndian.PutUint32(b, 2) }), identity
		}},
		{"CE32 header marker", "invalid-marker", func(*testing.T) ([]byte, string) {
			return header(0, func(b []byte) { copy(b, "BENCHEX\x00") }), identity
		}},
		{"CE33 reserved bits", "nonzero-reserved", func(*testing.T) ([]byte, string) {
			return header(12, func(b []byte) { binary.LittleEndian.PutUint32(b, 1) }), identity
		}},
		{"CE34 manifest length overflow", "length-overflow", func(*testing.T) ([]byte, string) {
			return header(16, func(b []byte) { binary.LittleEndian.PutUint64(b, ^uint64(0)) }), identity
		}},
		{"CE35 overlapping pages", "page-overlap", func(t *testing.T) ([]byte, string) {
			return repack(t, pack, replaceOnce(t, secondPage, fmt.Sprintf("  s2,1,8191,1808,%s\n", sum(big[8192:]))))
		}},
		{"CE36 gapped pages", "page-gap", func(t *testing.T) ([]byte, string) {
			return repack(t, pack, replaceOnce(t, secondPage, fmt.Sprintf("  s2,1,8193,1808,%s\n", sum(big[8192:]))))
		}},
		{"CE37 truncated body", "truncated", func(*testing.T) ([]byte, string) {
			data := pack.Bytes()
			n := len(data) - 1
			return data[:n:n], identity
		}},
		{"CE37 truncated header", "truncated", func(*testing.T) ([]byte, string) {
			return pack.Bytes()[:23:23], identity
		}},
		{"CE37 truncated manifest", "truncated", func(*testing.T) ([]byte, string) {
			return pack.Bytes()[:30:30], identity
		}},
		{"CE38 trailing data", "trailing-data", func(*testing.T) ([]byte, string) {
			return append(pack.Bytes(), 'x'), identity
		}},
		{"CE125 page beyond source", "page-range", func(t *testing.T) ([]byte, string) {
			return repack(t, pack, replaceOnce(t, secondPage, fmt.Sprintf("  s2,1,8192,1809,%s\n", sum(big[8192:]))))
		}},
		{"CV1 page split", "page-split", func(t *testing.T) ([]byte, string) {
			firstPage := fmt.Sprintf("  s2,0,0,8192,%s\n", sum(big[:8192]))
			newFirst := fmt.Sprintf("  s2,0,0,8000,%s\n", sum(big[:8000]))
			newSecond := fmt.Sprintf("  s2,1,8000,2000,%s\n", sum(big[8000:]))
			return repack(t, pack, func(s string) string {
				if strings.Count(s, firstPage) != 1 || strings.Count(s, secondPage) != 1 {
					t.Fatalf("manifest holds an unexpected page-row count")
				}
				return strings.Replace(strings.Replace(s, firstPage, newFirst, 1), secondPage, newSecond, 1)
			})
		}},
		{"CV2 source digest", "source-digest", func(t *testing.T) ([]byte, string) {
			p := mustBuild(t, fixtureCandidate("# One\n"))
			ticket := sum("# One\n")
			return repack(t, p, func(s string) string { return strings.Replace(s, ","+ticket+"\n", ","+strings.Repeat("a", 64)+"\n", 1) })
		}},
		{"CV2 zero-byte optional source digest", "source-digest", func(t *testing.T) ([]byte, string) {
			c := fixtureCandidate("# One\n")
			c.Sources = append(c.Sources, ce.SourceInput{Role: "optional", Kind: "repository", Path: "specs/example/optional.md"})
			p := mustBuild(t, c)
			return repack(t, p, replaceOnce(t, ","+sum("")+"\n", ","+strings.Repeat("f", 64)+"\n"))
		}},
		{"CV3 control byte source", "source-bytes", func(t *testing.T) ([]byte, string) {
			p := mustBuild(t, fixtureCandidate("safe\n"))
			return tamperSourceBody(t, p, "safe\n", "s\x1bfe\n")
		}},
		{"CV3 invalid UTF-8 source", "source-bytes", func(t *testing.T) ([]byte, string) {
			p := mustBuild(t, fixtureCandidate("safe\n"))
			return tamperSourceBody(t, p, "safe\n", "s\x80fe\n")
		}},
		{"changed page body", "page-digest", func(*testing.T) ([]byte, string) {
			data := pack.Bytes()
			data[len(data)-10] = 'y'
			return data, identity
		}},
		{"foreign identity", "identity-mismatch", func(*testing.T) ([]byte, string) {
			return pack.Bytes(), "sha256:" + strings.Repeat("0", 64)
		}},
		{"malformed identity", "invalid-identifier", func(*testing.T) ([]byte, string) {
			return pack.Bytes(), "sha256:" + strings.Repeat("A", 64)
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			data, id := test.data(t)
			assertRefusal(t, data, id, test.class)
		})
	}
	if _, err := ce.Read(pack.Bytes(), identity); err != nil {
		t.Fatalf("unchanged pack refused: %v", err)
	}
}
