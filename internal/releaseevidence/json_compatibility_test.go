package releaseevidence

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/jsonfile"
	"github.com/gibbonmi/bench/internal/jsonfile/jsonfiletest"
)

// releaseJSONBaseline holds the stored results of the strict release readers and the component
// manifest reader for the hostile document family. An ordinary run reads it and never writes it.
const releaseJSONBaseline = "testdata/strict-json-baseline.json"

// releaseJSONResults keys each release result by its reader and its case, and each component
// result by its case.
type releaseJSONResults struct {
	Release   map[string]string `json:"release"`
	Component map[string]string `json:"component"`
	Package   map[string]string `json:"package"`
}

// accepted is the result of a document that its reader accepts.
const accepted = "accepted"

// overflowNumber is a JSON number out of float64 range, which an untyped decode refuses.
const overflowNumber = "1e999"

// releaseReader is one strict release reader and the seed of its baseline family.
type releaseReader struct {
	seed jsonfiletest.Seed
	// framed marks requirement evidence, which requires a final newline before any decode.
	framed bool
	// panics marks an embedded registry loader. Its panic message carries no error chain.
	panics bool
	// raw is a member inside a json.RawMessage value, which only the structural check reads.
	raw string
	// overflow is a document whose raw value holds overflowNumber. An empty overflow omits the
	// overflow cases.
	overflow string
	read     func(t *testing.T, body []byte) error
}

// cases derives the baseline family of r. The release family adds an escaped duplicate key and
// a duplicate key that precedes a syntax error. For a framed reader, each case except unframed
// and empty ends in a final newline.
func (r releaseReader) cases(t *testing.T) map[string]string {
	t.Helper()
	cases := jsonfiletest.Family(r.seed)
	cases["escaped duplicate key"] = strings.Replace(r.seed.Valid, r.seed.Member, r.seed.Member+","+escapedKey(r.seed.Member), 1)
	cases["syntax plus duplicate"] = strings.TrimSuffix(cases["duplicate key"], "}") + ",}"
	if r.raw != "" {
		if !strings.Contains(r.seed.Valid, r.raw) {
			t.Fatalf("raw member %q is not in %q", r.raw, r.seed.Valid)
		}
		cases["raw duplicate key"] = strings.Replace(r.seed.Valid, r.raw, r.raw+","+r.raw, 1)
	}
	if r.overflow != "" {
		cases["raw overflow"] = r.overflow
		cases["duplicate plus raw overflow"] = strings.Replace(r.overflow, r.seed.Member, r.seed.Member+","+r.seed.Member, 1)
	}
	if r.framed {
		for name, body := range cases {
			if name != "unframed" && name != "empty" && !strings.HasSuffix(body, "\n") {
				cases[name] = body + "\n"
			}
		}
	}
	return cases
}

// escapedKey respells the first key byte of member as a JSON \u escape. The decoded key is
// unchanged.
func escapedKey(member string) string {
	return fmt.Sprintf(`"\u%04x`, member[1]) + member[2:]
}

func TestReleaseJSONCompatibility(t *testing.T) {
	readers := releaseReaders(t)
	names := make([]string, 0, len(readers))
	for name := range readers {
		names = append(names, name)
	}
	sort.Strings(names)

	// The specific contracts run before the baseline family, so each drift fails under its own
	// name first.
	t.Run("duplicate keys reach jsonfile", func(t *testing.T) {
		for _, name := range names {
			reader := readers[name]
			if reader.panics {
				continue
			}
			cases := reader.cases(t)
			for _, label := range []string{"duplicate key", "nested duplicate key", "escaped duplicate key", "raw duplicate key"} {
				body, ok := cases[label]
				if !ok {
					continue
				}
				if err := reader.read(t, []byte(body)); !errors.Is(err, jsonfile.ErrDuplicateKey) {
					t.Errorf("%s / %s: error = %v, want jsonfile.ErrDuplicateKey", name, label, err)
				}
			}
		}
	})
	t.Run("unknown typed fields", func(t *testing.T) {
		want := fmt.Sprintf("unknown field %q", jsonfiletest.UnknownKey)
		for _, name := range names {
			reader := readers[name]
			cases := reader.cases(t)
			for _, label := range []string{"unknown field", "nested unknown field"} {
				body, ok := cases[label]
				if !ok {
					continue
				}
				if err := reader.read(t, []byte(body)); err == nil || !strings.Contains(err.Error(), want) {
					t.Errorf("%s / %s: error = %v, want a refusal that names %s", name, label, err, want)
				}
			}
		}
	})
	t.Run("requirement evidence framing", func(t *testing.T) {
		for _, name := range names {
			reader := readers[name]
			if !reader.framed {
				continue
			}
			cases := reader.cases(t)
			if err := reader.read(t, []byte(cases["framed"])); err != nil {
				t.Errorf("%s / framed: error = %v, want it accepted", name, err)
			}
			if err := reader.read(t, []byte(cases["unframed"])); err == nil || !strings.Contains(err.Error(), "is missing a final newline") {
				t.Errorf("%s / unframed: error = %v, want the final-newline refusal", name, err)
			}
		}
	})

	var stored releaseJSONResults
	jsonfiletest.ReadBaseline(t, releaseJSONBaseline, &stored)
	t.Run("release baseline", func(t *testing.T) { jsonfiletest.Compare(t, releaseResults(t, readers), stored.Release) })
}

// releaseResults reads every case of each reader's family through that reader.
func releaseResults(t *testing.T, readers map[string]releaseReader) map[string]string {
	t.Helper()
	current := map[string]string{}
	for name, reader := range readers {
		for label, body := range reader.cases(t) {
			current[name+" / "+label] = resultText(reader.read(t, []byte(body)))
		}
	}
	return current
}

func resultText(err error) string {
	if err != nil {
		return err.Error()
	}
	return accepted
}

// releaseReaders holds the strict release readers. Each reads a body through its real validator.
// The native proof reader runs checkNativeProof with fixed digests and no release plan.
func releaseReaders(t *testing.T) map[string]releaseReader {
	t.Helper()
	commit, version := "abc123", fixtureVersion
	producer := ft87Requirement(t)
	envelope := func(payload, sum string) string {
		return fmt.Sprintf(`{"schema_version":1,"key":%q,"owner":%q,"schema":%q,"identity":{"source_commit":%q,"package_version":%q},"status":"satisfied","reason":"","payload":%s,"digest":%q}`,
			producer.Key, producer.Owner, producer.Schema, commit, version, payload, sum)
	}
	payload := `{"flag":"BENCH_OFFLINE=1","journeys":["direct"]}`
	// The overflow envelope carries the digest of empty bytes, so only the overflow refusal stops it.
	overflow := envelope(strings.TrimSuffix(payload, "}")+`,"n":`+overflowNumber+`}`, digest(nil))
	rootSeed := func(valid string) jsonfiletest.Seed {
		return jsonfiletest.Seed{Valid: valid, Member: `"schema_version":1`, Wrong: `"schema_version":"1"`}
	}
	spdx := jsonfiletest.Seed{
		Valid:  `{"SPDXID":"SPDXRef-DOCUMENT","SPDXVersion":"SPDX-2.3","creationInfo":{"creators":["Tool: bench"]},"dataLicense":"CC0-1.0","name":"redbench","documentNamespace":"https://example.invalid/redbench","packages":[{"SPDXID":"SPDXRef-Package","name":"redbench","versionInfo":"` + version + `","downloadLocation":"NOASSERTION","licenseConcluded":"MIT","licenseDeclared":"MIT"}],"relationships":[{"spdxElementId":"SPDXRef-DOCUMENT","relationshipType":"DESCRIBES","relatedSpdxElement":"SPDXRef-Package"}]}`,
		Member: `"dataLicense":"CC0-1.0"`, Wrong: `"dataLicense":1`, Nested: `"versionInfo":"` + version + `"`,
	}
	reproducibility := jsonfiletest.Seed{Valid: `{"schema_version":1,"status":"green","builds":2,"artifacts":[{"name":"bench.tgz","size":1,"sha256":"aa","match":true}],"evidence":[]}`, Member: `"builds":2`, Wrong: `"builds":"2"`, Nested: `"match":true`}
	proofTarget := targetEvidence{OS: "linux", Arch: "amd64", Runner: "native"}
	producerSeed := rootSeed(envelope(payload, digest(mustCanonicalPayload([]byte(payload)))))
	producerSeed.Nested = fmt.Sprintf(`"source_commit":%q`, commit)

	readers := map[string]releaseReader{
		"embedded phase registry": {
			seed:   jsonfiletest.Seed{Valid: strings.TrimSpace(string(registryJSON)), Member: `"publish_only": ["changelog"]`, Wrong: `"publish_only": "changelog"`, Nested: `"handler":"external"`},
			panics: true,
			read:   embeddedReader(&registryJSON, func() { loadRegistry() }),
		},
		"embedded requirement registry": {
			seed:   jsonfiletest.Seed{Valid: strings.TrimSpace(string(requirementsJSON)), Member: `"schema_version": 1`, Wrong: `"schema_version": "1"`, Nested: `"path": "component-manifest.json"`},
			panics: true,
			read:   embeddedReader(&requirementsJSON, func() { loadRequirements() }),
		},
		"producer envelope": {seed: producerSeed, framed: true, raw: `"flag":"BENCH_OFFLINE=1"`, overflow: overflow, read: func(t *testing.T, body []byte) error {
			return validateRequirementBytes(producer, body, Identity{SourceCommit: &commit, PackageVersion: &version})
		}},
		"SPDX document": {seed: spdx, read: func(t *testing.T, body []byte) error {
			return validateSPDXDocument(body, "redbench", version)
		}},
		"rollback policy": {seed: rootSeed(governanceSeeds["core.policy.recovery_rollback"]), read: func(t *testing.T, body []byte) error {
			root := t.TempDir()
			writeReleaseFile(t, root, "governance/policies/recovery-rollback.json", body)
			_, err := readRollbackTarget(root)
			return err
		}},
		"reproducibility record": {seed: reproducibility, read: func(t *testing.T, body []byte) error {
			root := t.TempDir()
			writeReleaseFile(t, root, "dist/reproducibility.json", body)
			_, err := readReproducibility(root, []artifactEvidence{{Name: "bench.tgz", Size: 1, SHA256: "aa"}})
			return err
		}},
		"native proof": {seed: rootSeed(`{"schema_version":1,"target":"linux-amd64","runner":"native","status":"green","rebuilt_sha256":"bin","binary_sha256":"bin","package_sha256":"pkg","archive_sha256":"arc","musl_status":"green","operations_status":"green","strip_status":"green","tools_status":"green"}`), read: func(t *testing.T, body []byte) error {
			_, err := checkNativeProof("native proof", body, proofTarget, func() (nativeProofDigests, error) {
				return nativeProofDigests{binary: "bin", pkg: "pkg", archive: "arc"}, nil
			})
			// The baseline holds the bare decode cause, as decodeStrict reports it.
			if cause := errors.Unwrap(err); cause != nil {
				return cause
			}
			return err
		}},
	}
	for key, valid := range governanceSeeds {
		record := requirementByKey(t, key)
		seed := rootSeed(valid)
		if key == "core.policy.security_response" {
			seed.Nested = `"critical":7`
		}
		readers["governance "+key] = releaseReader{seed: seed, framed: true, read: func(t *testing.T, body []byte) error {
			return validateRequirementBytes(record, body, Identity{})
		}}
	}
	return readers
}

// governanceSeeds holds one valid policy document for each governance policy schema.
var governanceSeeds = map[string]string{
	"core.policy.supported_versions":        `{"schema_version":1,"policy":"supported-versions","support":"latest minor","previous_minor_days":90,"eol":"announced"}`,
	"core.policy.security_response":         `{"schema_version":1,"policy":"security-response","private_route":"security advisory","acknowledge_business_days":2,"triage_business_days":5,"mitigation_business_days":{"critical":7,"high":14,"medium":30,"low":"next release"}}`,
	"core.policy.dependency_license_change": `{"schema_version":1,"policy":"dependency-license-change","review_required":true,"license_change_requires":"review","dependency_change_requires":"review"}`,
	"core.policy.threat_model":              `{"schema_version":1,"policy":"threat-model","release_input":true,"trust_boundary":"local","primary_threats":["tampering"]}`,
	"core.policy.recovery_rollback":         `{"schema_version":1,"policy":"recovery-rollback","rollback_target":"previous","preserve_prior_generation":true,"npm_rollback":"deprecate"}`,
	"core.policy.support":                   `{"schema_version":1,"policy":"support","route":"issues","personal_email":false,"non_personal":true}`,
}

// embeddedReader loads body as the embedded registry at data through load, then restores data.
// It reports the panic of load as an error.
func embeddedReader(data *[]byte, load func()) func(*testing.T, []byte) error {
	return func(t *testing.T, body []byte) (err error) {
		previous := *data
		*data = body
		defer func() {
			*data = previous
			if recovered := recover(); recovered != nil {
				err = fmt.Errorf("%v", recovered)
			}
		}()
		load()
		return nil
	}
}

func writeReleaseFile(t *testing.T, root, rel string, body []byte) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, body, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestComponentJSONCompatibility(t *testing.T) {
	labels := requirements.ComponentManifest
	component := releaseReader{
		seed: jsonfiletest.Seed{
			Valid:  componentManifestJSON(t, map[string][]byte{"package.json": []byte(identityPackage)}),
			Member: strconv.Quote(labels.RootFields.SchemaVersion) + ":1",
			Wrong:  strconv.Quote(labels.RootFields.SchemaVersion) + `:"1"`,
			Nested: strconv.Quote(labels.TargetFields.OS) + `:"all"`,
		},
		read: func(t *testing.T, body []byte) error {
			return ValidateTarballForTesting(packageTarball(t, packageMember{labels.Path, body}))
		},
	}
	cases := component.cases(t)
	read := func(body string) error { return component.read(t, []byte(body)) }

	// The specific contracts run before the baseline family, so each drift fails under its own
	// name first.
	t.Run("declared key set", func(t *testing.T) {
		if err := read(component.seed.Valid); err != nil {
			t.Fatalf("manifest with the declared keys: error = %v, want it accepted", err)
		}
		renamed := RequirementsRegistry()
		renamed.ComponentManifest.TargetFields.OS = "platform"
		restore := SetRequirementsForTesting(renamed)
		defer restore()
		if err := read(componentManifestJSON(t, nil)); err != nil {
			t.Errorf("manifest with the renamed registry key: error = %v, want it accepted", err)
		}
		if err := read(component.seed.Valid); err == nil {
			t.Error("manifest with the retired registry key is accepted, want it refused")
		}
	})
	t.Run("extra key", func(t *testing.T) {
		if err := read(cases["unknown field"]); err == nil || !strings.Contains(err.Error(), "fields do not match schema") {
			t.Errorf("manifest with an undeclared root key: error = %v, want the schema mismatch", err)
		}
	})
	t.Run("missing key", func(t *testing.T) {
		body := strings.Replace(component.seed.Valid, component.seed.Nested, `"platform":"all"`, 1)
		want := "field " + labels.TargetFields.OS + " is missing"
		if err := read(body); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("manifest with a replaced target key: error = %v, want %q", err, want)
		}
	})
	t.Run("recursive duplicate keys", func(t *testing.T) {
		for _, label := range []string{"duplicate key", "nested duplicate key", "escaped duplicate key"} {
			if err := read(cases[label]); !errors.Is(err, jsonfile.ErrDuplicateKey) {
				t.Errorf("%s: error = %v, want jsonfile.ErrDuplicateKey", label, err)
			}
		}
	})

	var stored releaseJSONResults
	jsonfiletest.ReadBaseline(t, releaseJSONBaseline, &stored)
	t.Run("component baseline", func(t *testing.T) { jsonfiletest.Compare(t, componentResults(t, component), stored.Component) })
}

// componentResults reads every case of the component family through the archive reader.
func componentResults(t *testing.T, component releaseReader) map[string]string {
	t.Helper()
	current := map[string]string{}
	for label, body := range component.cases(t) {
		current[label] = resultText(component.read(t, []byte(body)))
	}
	return current
}
