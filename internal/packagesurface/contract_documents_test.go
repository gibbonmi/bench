package packagesurface

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestContractDocumentInputsRefusesUnresolvableAllowlists pins this consumer to the
// canonical payload parser. The inventory it returns decides which documents the
// lifecycle contracts grade, so a row the allowlist forbids — an escaping source, an
// unknown audience, a source claimed twice — must stop the read rather than seed a
// document set nobody consented to. Absent is a defect here too: this check's whole
// subject is the allowlist.
func TestContractDocumentInputsRefusesUnresolvableAllowlists(t *testing.T) {
	for name, refusal := range map[string]struct {
		allowlist string
		want      string
	}{
		// Absent and empty want distinct wordings: both name the same path, so an
		// assertion on the path alone would still pass if a present-but-empty
		// allowlist were reclassified as absent.
		"absent":           {want: consumerPayloadPath + " is absent: the consumer payload inventory has no source"},
		"empty":            {want: consumerPayloadPath + " unreadable: empty"},
		"invalid JSON":     {allowlist: "{not json", want: consumerPayloadPath},
		"unknown audience": {allowlist: `[{"source":"AGENTS.md","audience":"everyone"}]`, want: consumerPayloadPath},
		"empty source":     {allowlist: `[{"source":"","audience":"kit-only"}]`, want: consumerPayloadPath},
		"unsafe source":    {allowlist: `[{"source":"../escape","audience":"kit-only"}]`, want: consumerPayloadPath},
		"duplicate source": {allowlist: `[{"source":"AGENTS.md","audience":"consumer"},{"source":"AGENTS.md","audience":"kit-only"}]`, want: consumerPayloadPath},
	} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
				t.Fatal(err)
			}
			if name != "absent" {
				if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(consumerPayloadPath)), []byte(refusal.allowlist), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			_, err := ContractDocumentInputs(root)
			if err == nil {
				t.Fatalf("ContractDocumentInputs with a %s allowlist returned no error", name)
			}
			if !strings.Contains(err.Error(), refusal.want) {
				t.Fatalf("ContractDocumentInputs with a %s allowlist = %v, want an error containing %q", name, err, refusal.want)
			}
		})
	}
}

// TestContractDocumentInputsResolvesAValidAllowlist is the control: the refusals above
// must not be a reader that fails on everything.
func TestContractDocumentInputsResolvesAValidAllowlist(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("# agents\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(consumerPayloadPath)), []byte(`[{"source":"AGENTS.md","mode":"0644","audience":"consumer"}]`), 0o644); err != nil {
		t.Fatal(err)
	}
	paths, err := ContractDocumentInputs(root)
	if err != nil {
		t.Fatalf("ContractDocumentInputs on a valid allowlist = %v", err)
	}
	want := []string{".bench/consumer-payload.json", "AGENTS.md"}
	if strings.Join(paths, ",") != strings.Join(want, ",") {
		t.Fatalf("ContractDocumentInputs = %v, want %v", paths, want)
	}
}
