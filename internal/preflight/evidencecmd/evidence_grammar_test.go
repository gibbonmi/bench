package evidencecmd_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// The expected grammar forms, response headers, and cell types in this file are stated
// independently of the operation and format registries, so an omission from either turns a
// case red.

func headerLine(out string) string {
	line, _, _ := strings.Cut(out, "\n")
	return line
}

// assertTypedRow proves that every expected field decoded with its declared TOON type.
func assertTypedRow(t *testing.T, row map[string]any, types map[string]string) {
	t.Helper()
	if len(row) != len(types) {
		t.Fatalf("row has %d fields, want %d: %v", len(row), len(types), row)
	}
	for name, want := range types {
		got := ""
		switch row[name].(type) {
		case string:
			got = "string"
		case float64:
			got = "integer"
		case bool:
			got = "boolean"
		}
		if got != want {
			t.Errorf("field %s = %T, want %s", name, row[name], want)
		}
	}
}

// TestEvidencePreparedSchema is CE153.
func TestEvidencePreparedSchema(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	_, row, out := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug, false))
	if got := headerLine(out); got != "prepared[1]{evidence,mode,base,source_tip,assignment,selection,metadata,sources,pages,manifest_bytes,response_complete,delivery,next}:" {
		t.Fatalf("prepared header = %q", got)
	}
	assertTypedRow(t, row, map[string]string{"evidence": "string", "mode": "string", "base": "string", "source_tip": "string",
		"assignment": "string", "selection": "string", "metadata": "string", "sources": "integer", "pages": "integer",
		"manifest_bytes": "integer", "response_complete": "boolean", "delivery": "string", "next": "string"})
	if row["mode"] != "build" || row["selection"] != "manifest:selection" || row["metadata"] != "s1" || row["assignment"] != preflighttest.ChargeFixtureAssignment || row["sources"] != float64(6) {
		t.Fatalf("prepared values = %v", row)
	}
}

// TestEvidenceManifestPageSchema is CE154 and TestEvidenceSourcePageSchema is CE155.
func TestEvidenceManifestPageSchema(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	identity, prepared, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug, false))
	page := traverseEvidence(t, identity)[0]
	assertPageSchema(t, page)
	if page.row["stream"] != "manifest" || page.row["source"] != "" || page.row["total"] != prepared["manifest_bytes"] || page.row["evidence"] != identity {
		t.Fatalf("manifest page = %v", page.row)
	}
}

func TestEvidenceSourcePageSchema(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	identity, _, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug, false))
	for _, page := range traverseEvidence(t, identity) {
		if page.row["stream"] != "source" {
			continue
		}
		assertPageSchema(t, page)
		if page.row["source"] != "s1" || page.row["index"] != float64(0) || page.row["offset"] != float64(0) {
			t.Fatalf("first source page = %v", page.row)
		}
		return
	}
	t.Fatal("no source page")
}

func assertPageSchema(t *testing.T, page evidencePage) {
	t.Helper()
	if got := headerLine(page.raw); got != "page[1]{evidence,stream,source,index,offset,bytes,total,sha256,content,response_complete,stream_end,next}:" {
		t.Fatalf("page header = %q", got)
	}
	assertTypedRow(t, page.row, map[string]string{"evidence": "string", "stream": "string", "source": "string", "index": "integer",
		"offset": "integer", "bytes": "integer", "total": "integer", "sha256": "string", "content": "string",
		"response_complete": "boolean", "stream_end": "boolean", "next": "string"})
}

// TestEvidenceTypedCells is CE11: a numeric-looking source tip and every digest decode as
// strings.
func TestEvidenceTypedCells(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	t.Setenv("GIT_COMMITTER_DATE", "2026-01-02T03:04:05Z")
	tip := ""
	for i := 0; i < 512; i++ {
		preflighttest.RunGit(t, "commit", "--amend", "-q", "-m", "numeric tip "+strconv.Itoa(i))
		tip = preflighttest.RunGit(t, "rev-parse", "HEAD")
		if tip[0] == '0' && tip[1] >= '0' && tip[1] <= '9' {
			break
		}
	}
	identity, row, out := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug, false))
	if row["source_tip"] != tip || !strings.Contains(out, ",\""+tip+"\",") {
		t.Fatalf("source tip = %#v, want string %s", row["source_tip"], tip)
	}
	manifest, _ := reconstructEvidence(t, identity, traverseEvidence(t, identity))
	for _, page := range preflighttest.TableRows(t, preflighttest.DecodeMap(t, manifest), "pages") {
		if _, ok := page["sha256"].(string); !ok {
			t.Fatalf("page digest = %#v", page["sha256"])
		}
	}
}

// TestEvidenceBuildGrammar is CE150.
func TestEvidenceBuildGrammar(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	valid := preflighttest.ChargeArgs(t, root, slug, false)
	base, tip := valid[6], valid[8]
	for _, form := range [][]string{
		valid,
		append(append([]string{}, valid...), "--max-store-bytes", "1073741824"),
		{"build", slug, "--source-tip", tip, "--base", base, "--ticket", "one.md", "--charge"},
	} {
		if out, code := preflight.Command(form); code != 0 || !strings.HasPrefix(out, "prepared[1]") {
			t.Fatalf("accepted form %v = (%d):\n%s", form, code, out)
		}
	}
	for _, test := range []struct {
		name string
		args []string
		want string
	}{
		{"missing base", []string{"build", slug, "--charge", "--ticket", "one.md", "--source-tip", tip}, "--charge requires build, --ticket, --base, and --source-tip"},
		{"missing tip", []string{"build", slug, "--charge", "--ticket", "one.md", "--base", base}, "--charge requires build, --ticket, --base, and --source-tip"},
		{"missing ticket", []string{"build", slug, "--charge", "--base", base, "--source-tip", tip}, "--charge requires build, --ticket, --base, and --source-tip"},
		{"read-only flag", append(append([]string{}, valid...), "--cursor", "v1"), "unknown argument: --cursor"},
		{"duplicate quota", append(append([]string{}, valid...), "--max-store-bytes", "1", "--max-store-bytes", "2"), "unknown argument: --max-store-bytes"},
		{"extra operand", append(append([]string{}, valid...), "extra"), "unknown argument: extra"},
		{"zero quota", append(append([]string{}, valid...), "--max-store-bytes", "0"), "needs a positive decimal byte count"},
		{"signed quota", append(append([]string{}, valid...), "--max-store-bytes", "+5"), "needs a positive decimal byte count"},
		{"padded quota", append(append([]string{}, valid...), "--max-store-bytes", "05"), "needs a positive decimal byte count"},
		{"overflowing quota", append(append([]string{}, valid...), "--max-store-bytes", "18446744073709551616"), "needs a positive decimal byte count"},
		{"missing quota value", append(append([]string{}, valid...), "--max-store-bytes"), "missing argument: --max-store-bytes"},
		// CE173. The build guidance migration retires the legacy full charge, so the exact
		// form the phase used to run now refuses through the bounded usage path.
		{"CE173 build charge full", append(append([]string{}, valid...), "--full"), "--charge and --full requires review"},
		{"CE173 build full alone", []string{"build", slug, "--full"}, "--full requires"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if out, code := preflight.Command(test.args); code != 2 || !strings.Contains(out, test.want) {
				t.Fatalf("%s = (%d):\n%s", test.name, code, out)
			}
		})
	}
}

// TestEvidenceReadGrammar is CE151.
func TestEvidenceReadGrammar(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	identity, _, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug, false))
	cursor := "v1." + strings.TrimPrefix(identity, "sha256:") + ".m.0.0"
	for _, form := range [][]string{{"evidence", identity}, {"evidence", identity, "--cursor", cursor}, {"evidence", "--cursor", cursor, identity}} {
		if out, code := preflight.Command(form); code != 0 || !strings.HasPrefix(out, "page[1]") {
			t.Fatalf("accepted form %v = (%d):\n%s", form, code, out)
		}
	}
	for _, test := range []struct {
		name string
		args []string
		want string
	}{
		{"unknown flag", []string{"evidence", identity, "--page"}, "unknown argument: --page"},
		{"extra operand", []string{"evidence", identity, "extra"}, "unknown argument: extra"},
		{"missing operand", []string{"evidence"}, "missing argument"},
		{"missing cursor value", []string{"evidence", identity, "--cursor"}, "missing argument: --cursor"},
		{"duplicate cursor", []string{"evidence", identity, "--cursor", cursor, "--cursor", cursor}, "unknown argument: --cursor"},
		{"base flag", []string{"evidence", identity, "--base", "main"}, "unknown argument: --base"},
		{"quota flag", []string{"evidence", identity, "--max-store-bytes", "5"}, "unknown argument: --max-store-bytes"},
		{"charge flag", []string{"evidence", identity, "--charge"}, "--charge requires"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if out, code := preflight.Command(test.args); code != 2 || !strings.Contains(out, test.want) {
				t.Fatalf("%s = (%d):\n%s", test.name, code, out)
			}
		})
	}
}

// TestEvidenceCursorRefusals is CE63. A malformed or foreign cursor refuses before any
// store access; a well-formed cursor beyond the artifact refuses at the read.
func TestEvidenceCursorRefusals(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	identity, _, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug, false))
	hex := strings.TrimPrefix(identity, "sha256:")
	foreign := strings.Repeat("0", 64)
	for _, test := range []struct{ name, cursor string }{
		{"version", "v2." + hex + ".m.0.0"},
		{"short identity", "v1.abc.m.0.0"},
		{"foreign identity", "v1." + foreign + ".m.0.0"},
		{"stream", "v1." + hex + ".x.0.0"},
		{"manifest ordinal", "v1." + hex + ".m.1.0"},
		{"source ordinal zero", "v1." + hex + ".s.0.0"},
		{"leading zero", "v1." + hex + ".s.01.0"},
		{"signed index", "v1." + hex + ".s.1.-1"},
		{"extra field", "v1." + hex + ".s.1.0.0"},
		{"missing field", "v1." + hex + ".s.1"},
	} {
		t.Run(test.name, func(t *testing.T) {
			out, code := preflight.Command([]string{"evidence", identity, "--cursor", test.cursor})
			if code != 2 || !strings.Contains(out, "invalid cursor bytes=") || strings.Contains(out, test.cursor) {
				t.Fatalf("%s cursor = (%d):\n%s", test.name, code, out)
			}
		})
	}
	for _, cursor := range []string{"v1." + hex + ".m.0.999", "v1." + hex + ".s.99.0", "v1." + hex + ".s.1.999"} {
		out, code := preflight.Command([]string{"evidence", identity, "--cursor", cursor})
		if code != 1 || !strings.Contains(out, "invalid-cursor") {
			t.Fatalf("out-of-range cursor %s = (%d):\n%s", cursor, code, out)
		}
	}
}

// TestEvidenceSourceGrammar is CE174: the source-selected read accepts exactly its declared
// forms, including valid cursor pairing, and refuses every other combination.
func TestEvidenceSourceGrammar(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	identity, _, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug, false))
	hex := strings.TrimPrefix(identity, "sha256:")
	for _, form := range [][]string{
		{"evidence", identity, "--source", "s2"},
		{"evidence", identity, "--source", "s2", "--cursor", "v1." + hex + ".s.2.0"},
		{"evidence", "--source", "s2", identity},
	} {
		if out, code := preflight.Command(form); code != 0 || !strings.HasPrefix(out, "page[1]") {
			t.Fatalf("accepted form %v = (%d):\n%s", form, code, out)
		}
	}
	for _, test := range []struct {
		name string
		args []string
		code int
		want string
	}{
		{"missing source value", []string{"evidence", identity, "--source"}, 2, "missing argument: --source"},
		{"duplicate source", []string{"evidence", identity, "--source", "s2", "--source", "s3"}, 2, "unknown argument: --source"},
		{"extra operand", []string{"evidence", identity, "--source", "s2", "extra"}, 2, "unknown argument: extra"},
		{"malformed source", []string{"evidence", identity, "--source", "specs/example/spec.md"}, 2, "invalid source identifier bytes="},
		{"zero source", []string{"evidence", identity, "--source", "s0"}, 2, "invalid source identifier bytes="},
		{"padded source", []string{"evidence", identity, "--source", "s02"}, 2, "invalid source identifier bytes="},
		{"cursor of another source", []string{"evidence", identity, "--source", "s2", "--cursor", "v1." + hex + ".s.3.0"}, 2, "names another source than s2"},
		{"manifest cursor", []string{"evidence", identity, "--source", "s2", "--cursor", "v1." + hex + ".m.0.0"}, 2, "names another source than s2"},
		{"verify combination", []string{"evidence", identity, "--source", "s2", "--verify"}, 2, "cannot be combined"},
		{"current combination", []string{"evidence", identity, "--source", "s2", "--check-current"}, 2, "cannot be combined"},
		{"undeclared source", []string{"evidence", identity, "--source", "s99"}, 1, "unknown-source"},
	} {
		t.Run(test.name, func(t *testing.T) {
			out, code := preflight.Command(test.args)
			if code != test.code || !strings.Contains(out, test.want) {
				t.Fatalf("%s = (%d):\n%s", test.name, code, out)
			}
		})
	}
}

// TestEvidenceVerifyGrammar is CE158 and TestEvidenceCurrentGrammar is CE159: each mode
// accepts only its exclusive declared form and refuses before any store access.
func TestEvidenceVerifyGrammar(t *testing.T) {
	assertExclusiveMode(t, "--verify", "verified[1]")
}

func TestEvidenceCurrentGrammar(t *testing.T) {
	assertExclusiveMode(t, "--check-current", "current[1]")
}

// assertExclusiveMode proves that mode accepts only `evidence <id> <mode>` and that every
// other operand combination refuses with exit 2.
func assertExclusiveMode(t *testing.T, mode, header string) {
	t.Helper()
	root, slug := preflighttest.SeedConformant(t)
	identity, _, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug, false))
	hex := strings.TrimPrefix(identity, "sha256:")
	for _, form := range [][]string{{"evidence", identity, mode}, {"evidence", mode, identity}} {
		if out, code := preflight.Command(form); code != 0 || !strings.HasPrefix(out, header) {
			t.Fatalf("accepted form %v = (%d):\n%s", form, code, out)
		}
	}
	other := "--verify"
	if mode == other {
		other = "--check-current"
	}
	for _, test := range []struct{ name, flag, value, want string }{
		{"with cursor", "--cursor", "v1." + hex + ".m.0.0", "unknown argument: --cursor"},
		{"with source", "--source", "s2", "cannot be combined"},
		{"with base", "--base", "main", "unknown argument: --base"},
		{"with quota", "--max-store-bytes", "5", "unknown argument: --max-store-bytes"},
		{"with the other mode", other, "", "cannot be combined"},
		{"duplicate", mode, "", "unknown argument: " + mode},
	} {
		t.Run(test.name, func(t *testing.T) {
			args := []string{"evidence", identity, mode, test.flag}
			if test.value != "" {
				args = append(args, test.value)
			}
			out, code := preflight.Command(args)
			if code != 2 || !strings.Contains(out, test.want) {
				t.Fatalf("%s %s = (%d):\n%s", mode, test.name, code, out)
			}
		})
	}
	// An invalid identifier refuses before any store access.
	if out, code := preflight.Command([]string{"evidence", "sha256:nope", mode}); code != 2 || !strings.Contains(out, "invalid evidence identifier bytes=") {
		t.Fatalf("%s with an invalid identifier = (%d):\n%s", mode, code, out)
	}
}
