package chargeevidence_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	ce "github.com/gibbonmi/bench/internal/chargeevidence"
	"golang.org/x/mod/modfile"
)

// The expectations below are written from the spec's format tables, not read from the
// registry, so an omitted or reordered field turns one of them red.

const (
	fixtureBase = "abcdefabcdefabcdefabcdefabcdefabcdefabcd"
	fixtureTip  = "fedcbafedcbafedcbafedcbafedcbafedcbafedc"
)

func fixtureCandidate(ticket string) ce.Candidate {
	return ce.Candidate{
		Selection: ce.Selection{Mode: "build", Spec: "specs/example/spec.md", Ticket: "specs/example/tickets/one.md", Base: fixtureBase, SourceTip: fixtureTip},
		Metadata: ce.Metadata{
			Charge:   []ce.ChargeRow{{Axis: "", Ticket: "s2", Access: "write-within-fence"}},
			Fence:    []string{"internal/example/"},
			Writes:   []string{"internal/example"},
			Coverage: []string{"PF1"},
			Checks:   []string{"s2"},
			Returns:  []string{"s3"},
		},
		Sources: []ce.SourceInput{
			{Role: "ticket", Kind: "repository", Path: "specs/example/tickets/one.md", Required: true, Data: []byte(ticket)},
			{Role: "spec", Kind: "repository", Path: "specs/example/spec.md", Required: true, Data: []byte("# Spec\n")},
		},
	}
}

func mustBuild(t *testing.T, c ce.Candidate) *ce.Pack {
	t.Helper()
	pack, err := ce.Build(c)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	return pack
}

func sum(data string) string {
	s := sha256.Sum256([]byte(data))
	return hex.EncodeToString(s[:])
}

func TestEvidenceDeterminism(t *testing.T) {
	first := mustBuild(t, fixtureCandidate("# One\n"))
	second := mustBuild(t, fixtureCandidate("# One\n"))
	if first.Identity() != second.Identity() || !bytes.Equal(first.Bytes(), second.Bytes()) {
		t.Fatalf("identical candidates differ: %s %s", first.Identity(), second.Identity())
	}
	if want := "sha256:" + sum(string(first.ManifestBytes())); first.Identity() != want {
		t.Fatalf("identity = %s, want the manifest digest %s", first.Identity(), want)
	}
}

func TestEvidenceIdentityInputs(t *testing.T) {
	base := mustBuild(t, fixtureCandidate("# One\n")).Identity()
	for _, test := range []struct {
		name   string
		mutate func(*ce.Candidate)
	}{
		{"CE18 mode", func(c *ce.Candidate) { c.Selection.Mode, c.Selection.Ticket = "review", "" }},
		{"CE19 spec selector", func(c *ce.Candidate) { c.Selection.Spec = "specs/other/spec.md" }},
		{"CE20 ticket selector", func(c *ce.Candidate) { c.Selection.Ticket = "specs/example/tickets/two.md" }},
		{"CE21 base", func(c *ce.Candidate) { c.Selection.Base = strings.Repeat("c", 40) }},
		{"CE22 source tip", func(c *ce.Candidate) { c.Selection.SourceTip = strings.Repeat("d", 40) }},
		{"CE23 role", func(c *ce.Candidate) { c.Sources[1].Role = "decision" }},
		{"CE24 path", func(c *ce.Candidate) { c.Sources[1].Path = "specs/example/other.md" }},
		{"CE25 requiredness", func(c *ce.Candidate) { c.Sources[1].Required = false }},
		{"CE26 source bytes", func(c *ce.Candidate) { c.Sources[1].Data = []byte("# Spek\n") }},
		{"CE27 metadata", func(c *ce.Candidate) { c.Metadata.Coverage = append(c.Metadata.Coverage, "PF2") }},
	} {
		t.Run(test.name, func(t *testing.T) {
			candidate := fixtureCandidate("# One\n")
			test.mutate(&candidate)
			if got := mustBuild(t, candidate).Identity(); got == base {
				t.Fatalf("changed %s kept identity %s", test.name, got)
			}
		})
	}
}

// TestEvidenceCanonicalProfile pins the exact canonical bytes. The mutated ticket role
// and spec path add a comma-bearing and a colon-bearing string value, so the pinned
// bytes also fix their quoted forms; the encoder module owns the quoting rule itself.
func TestEvidenceCanonicalProfile(t *testing.T) {
	const ticket, spec = "# One\n", "# Spec\n"
	const role, path = "ticket, primary", "specs/example/spec.md:pinned"
	metadata := "charge[1]{axis,ticket,access}:\n  \"\",s2,write-within-fence\n" +
		"fence[1]{path}:\n  internal/example/\n" +
		"writes[1]{path}:\n  internal/example\n" +
		"coverage[1]{row}:\n  PF1\n" +
		"checks[1]{source}:\n  s2\n" +
		"returns[1]{source}:\n  s3\n" +
		"shared_evidence[0]{kind,source}:\n" +
		"completion_evidence[0]{record,source_digest,plan_digest,record_state,detail}:\n"
	manifest := "profile[1]{version,hash,page_bytes}:\n  1,sha256,8192\n" +
		"selection[1]{mode,spec,ticket,base,source_tip}:\n  build,specs/example/spec.md,specs/example/tickets/one.md," + fixtureBase + "," + fixtureTip + "\n" +
		"sources[3]{id,role,kind,path,required,bytes,sha256}:\n" +
		fmt.Sprintf("  s1,metadata,derived,\"\",true,%d,%s\n", len(metadata), sum(metadata)) +
		fmt.Sprintf("  s2,\"%s\",repository,specs/example/tickets/one.md,true,%d,%s\n", role, len(ticket), sum(ticket)) +
		fmt.Sprintf("  s3,spec,repository,\"%s\",true,%d,%s\n", path, len(spec), sum(spec)) +
		"pages[3]{source,index,offset,bytes,sha256}:\n" +
		fmt.Sprintf("  s1,0,0,%d,%s\n", len(metadata), sum(metadata)) +
		fmt.Sprintf("  s2,0,0,%d,%s\n", len(ticket), sum(ticket)) +
		fmt.Sprintf("  s3,0,0,%d,%s\n", len(spec), sum(spec)) +
		"producers[0]{source,name,version,cwd}:\n" +
		"arguments[0]{source,index,value}:\n"
	candidate := fixtureCandidate(ticket)
	candidate.Sources[0].Role = role
	candidate.Sources[1].Path = path
	pack := mustBuild(t, candidate)
	if got := string(pack.ManifestBytes()); got != manifest {
		t.Fatalf("manifest =\n%s\nwant\n%s", got, manifest)
	}
	if body, _ := pack.Source("s1"); string(body) != metadata {
		t.Fatalf("metadata =\n%s\nwant\n%s", body, metadata)
	}
	header := make([]byte, 24)
	copy(header, "BENCHEV\x00")
	binary.LittleEndian.PutUint32(header[8:], 1)
	binary.LittleEndian.PutUint64(header[16:], uint64(len(manifest)))
	want := string(header) + manifest + metadata + ticket + spec
	if got := string(pack.Bytes()); got != want {
		t.Fatalf("pack bytes differ from header, manifest, and bodies in order")
	}
}

// TestEvidenceSchemaRegistry pins the registered field set, order, and types against the
// spec tables, so an omitted or reordered field reds here as well as in the projection.
func TestEvidenceSchemaRegistry(t *testing.T) {
	want := map[string]string{
		"profile":             "version:integer,hash:string,page_bytes:integer",
		"selection":           "mode:string,spec:string,ticket:string,base:string,source_tip:string",
		"sources":             "id:string,role:string,kind:string,path:string,required:boolean,bytes:integer,sha256:string",
		"pages":               "source:string,index:integer,offset:integer,bytes:integer,sha256:string",
		"producers":           "source:string,name:string,version:string,cwd:string",
		"arguments":           "source:string,index:integer,value:string",
		"charge":              "axis:string,ticket:string,access:string",
		"fence":               "path:string",
		"writes":              "path:string",
		"coverage":            "row:string",
		"checks":              "source:string",
		"returns":             "source:string",
		"shared_evidence":     "kind:string,source:string",
		"completion_evidence": "record:string,source_digest:string,plan_digest:string,record_state:string,detail:string",
		"prepared":            "evidence:string,mode:string,base:string,source_tip:string,assignment:string,selection:string,metadata:string,sources:integer,pages:integer,manifest_bytes:integer,response_complete:boolean,delivery:string,next:string",
		"page":                "evidence:string,stream:string,source:string,index:integer,offset:integer,bytes:integer,total:integer,sha256:string,content:string,response_complete:boolean,stream_end:boolean,next:string",
		"verified":            "evidence:string,manifest_verified:boolean,pages_verified:integer,sources_verified:integer,delivery:string",
		"current":             "evidence:string,assignment:string,base:string,source_tip:string,current:boolean,delivery:string",
	}
	order := map[string][]string{
		"manifest": {"profile", "selection", "sources", "pages", "producers", "arguments"},
		"metadata": {"charge", "fence", "writes", "coverage", "checks", "returns", "shared_evidence", "completion_evidence"},
		"response": {"prepared", "page", "verified", "current"},
	}
	for family, blocks := range map[string][]ce.Block{"manifest": ce.ManifestBlocks, "metadata": ce.MetadataBlocks, "response": ce.ResponseBlocks} {
		var names []string
		for _, block := range blocks {
			names = append(names, block.Name)
			fields := make([]string, len(block.Fields))
			for i, field := range block.Fields {
				fields[i] = field.Name + ":" + field.Type.String()
			}
			if got := strings.Join(fields, ","); got != want[block.Name] {
				t.Errorf("%s block %s = %s, want %s", family, block.Name, got, want[block.Name])
			}
		}
		if got := strings.Join(names, ","); got != strings.Join(order[family], ",") {
			t.Errorf("%s blocks = %s, want %s", family, got, strings.Join(order[family], ","))
		}
	}
}

func TestEvidencePageSplit(t *testing.T) {
	// A three-byte code point straddles the page limit, so the first page ends before it.
	body := strings.Repeat("a", 8191) + "雪" + "tail"
	pack := mustBuild(t, fixtureCandidate(body))
	var sizes []int
	for _, page := range pack.Manifest().Pages {
		if page.Source == "s2" {
			sizes = append(sizes, page.Bytes)
		}
	}
	if fmt.Sprint(sizes) != "[8191 7]" {
		t.Fatalf("page sizes = %v, want [8191 7]", sizes)
	}
	if got, _ := pack.Source("s2"); string(got) != body {
		t.Fatalf("reconstructed body differs")
	}
}

// generatedSourceCandidate extends fixtureCandidate with one generated source that
// declares a producer and two ordered arguments.
func generatedSourceCandidate(ticket string) ce.Candidate {
	c := fixtureCandidate(ticket)
	c.Sources = append(c.Sources, ce.SourceInput{
		Role: "diff", Kind: "generated", Path: "diff", Required: true,
		Data:     []byte("+++ generated diff\n"),
		Producer: &ce.Producer{Name: "collector", Version: "v1", Cwd: ".", Arguments: []string{"alpha", "beta"}},
	})
	return c
}

// TestEvidenceGeneratedSourceProvenance is CV5: a generated source's producer and
// argument rows round-trip through the strict reader unchanged.
func TestEvidenceGeneratedSourceProvenance(t *testing.T) {
	pack := mustBuild(t, generatedSourceCandidate("# One\n"))
	wantProducers := []ce.ProducerRow{{Source: "s4", Name: "collector", Version: "v1", Cwd: "."}}
	wantArguments := []ce.ArgumentRow{{Source: "s4", Index: 0, Value: "alpha"}, {Source: "s4", Index: 1, Value: "beta"}}
	if got := pack.Manifest().Producers; !reflect.DeepEqual(got, wantProducers) {
		t.Fatalf("producers = %+v, want %+v", got, wantProducers)
	}
	if got := pack.Manifest().Arguments; !reflect.DeepEqual(got, wantArguments) {
		t.Fatalf("arguments = %+v, want %+v", got, wantArguments)
	}
	reread, err := ce.Read(pack.Bytes(), pack.Identity())
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if got := reread.Manifest().Producers; !reflect.DeepEqual(got, wantProducers) {
		t.Fatalf("round-trip producers = %+v, want %+v", got, wantProducers)
	}
	if got := reread.Manifest().Arguments; !reflect.DeepEqual(got, wantArguments) {
		t.Fatalf("round-trip arguments = %+v, want %+v", got, wantArguments)
	}
}

// TestEvidenceProvenanceRefusals is CV5's refusal half: a producer row naming a
// repository source, and an out-of-order argument row, both refuse.
func TestEvidenceProvenanceRefusals(t *testing.T) {
	t.Run("producer names repository source", func(t *testing.T) {
		pack := mustBuild(t, fixtureCandidate("# One\n"))
		data, identity := repack(t, pack, replaceOnce(t,
			"producers[0]{source,name,version,cwd}:\n",
			"producers[1]{source,name,version,cwd}:\n  s2,fake,v1,.\n"))
		assertRefusal(t, data, identity, "invalid-manifest")
	})
	t.Run("out-of-order arguments", func(t *testing.T) {
		pack := mustBuild(t, generatedSourceCandidate("# One\n"))
		data, identity := repack(t, pack, replaceOnce(t,
			"  s4,0,alpha\n  s4,1,beta\n",
			"  s4,1,beta\n  s4,0,alpha\n"))
		assertRefusal(t, data, identity, "invalid-manifest")
	})
}

func TestEvidenceCandidateRefusals(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*ce.Candidate)
	}{
		{"control byte", func(c *ce.Candidate) { c.Sources[0].Data = []byte("bad \x1b\n") }},
		{"invalid UTF-8", func(c *ce.Candidate) { c.Sources[0].Data = []byte("bad \xff\n") }},
		{"empty required", func(c *ce.Candidate) { c.Sources[0].Data = nil }},
		{"undeclared reference", func(c *ce.Candidate) { c.Metadata.Checks = []string{"s9"} }},
		{"generated without producer", func(c *ce.Candidate) { c.Sources[0].Kind = "generated" }},
		{"unknown kind", func(c *ce.Candidate) { c.Sources[0].Kind = "other" }},
		{"build without ticket", func(c *ce.Candidate) { c.Selection.Ticket = "" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			candidate := fixtureCandidate("# One\n")
			test.mutate(&candidate)
			if pack, err := ce.Build(candidate); err == nil {
				t.Fatalf("Build accepted %s: %s", test.name, pack.Identity())
			}
		})
	}
}

// TestEvidenceEncoderModuleRequired ties the reference's named encoder module to the
// encoder the shared TOON adapter actually imports, and to the go.mod requirement that
// pins its version.
func TestEvidenceEncoderModuleRequired(t *testing.T) {
	adapter := filepath.Join("..", "toon", "toon.go")
	file, err := parser.ParseFile(token.NewFileSet(), adapter, nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parse %s: %v", adapter, err)
	}
	imported := false
	for _, spec := range file.Imports {
		if path, _ := strconv.Unquote(spec.Path.Value); path == ce.EncoderModule {
			imported = true
		}
	}
	if !imported {
		t.Fatalf("the shared TOON adapter does not import encoder module %s", ce.EncoderModule)
	}
	path := filepath.Join("..", "..", "go.mod")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	parsed, err := modfile.Parse(path, data, nil)
	if err != nil {
		t.Fatalf("parse go.mod: %v", err)
	}
	for _, req := range parsed.Require {
		if req.Mod.Path == ce.EncoderModule {
			return
		}
	}
	t.Fatalf("go.mod has no require for encoder module %s", ce.EncoderModule)
}

// TestEvidenceFormatProjection is CE148: the shipped reference equals the registry
// projection, and every block's own projected row names every one of its registered
// fields. A row-scoped check catches an omitted field even when another block's row
// carries the same label.
func TestEvidenceFormatProjection(t *testing.T) {
	shipped, err := os.ReadFile(filepath.Join("..", "..", filepath.FromSlash(ce.ReferencePath)))
	if err != nil {
		t.Fatalf("read shipped reference: %v", err)
	}
	generated := ce.FormatReference()
	if string(shipped) != generated {
		t.Fatalf("shipped %s differs from the registry projection:\n%s", ce.ReferencePath, generated)
	}
	lines := strings.Split(generated, "\n")
	for _, block := range append(append(append([]ce.Block{}, ce.ManifestBlocks...), ce.MetadataBlocks...), ce.ResponseBlocks...) {
		prefix := "| `" + block.Name + "` |"
		var row string
		count := 0
		for _, line := range lines {
			if strings.HasPrefix(line, prefix) {
				row = line
				count++
			}
		}
		if count != 1 {
			t.Fatalf("projection holds %d rows for block %s, want 1", count, block.Name)
		}
		for _, field := range block.Fields {
			label := fmt.Sprintf("`%s` (%s)", field.Name, field.Type)
			if !strings.Contains(row, label) {
				t.Errorf("%s row omits field %s:\n%s", block.Name, label, row)
			}
		}
	}
}
