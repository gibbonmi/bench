package guards

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestManifestField(t *testing.T) {
	out := "name: pre-push\nboundary: git push\ndenies: destructive git — force-push\nwhy: because\n"
	if got := manifestField(out, "name"); got != "pre-push" {
		t.Errorf("name = %q", got)
	}
	if got := manifestField(out, "denies"); got != "destructive git — force-push" {
		t.Errorf("denies = %q", got)
	}
	if got := manifestField(out, "absent"); got != "" {
		t.Errorf("absent = %q, want empty", got)
	}
}

// guardRow classification against small bash stubs — the exec idiom later slices reuse.
func TestGuardRow(t *testing.T) {
	dir := t.TempDir()
	stub := func(name, body string) string {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte("#!/usr/bin/env bash\n"+body), 0o755); err != nil {
			t.Fatal(err)
		}
		return p
	}
	full := stub("full.sh", `[ "$1" = --describe ] && printf 'name: g\nboundary: b\ndenies: d\nwhy: w\n'`)
	if r, emit := guardRow(full, "full"); !emit || !reflect.DeepEqual(r, []string{"g", "b", "d"}) {
		t.Errorf("full manifest row = %v emit=%v", r, emit)
	}

	info := stub("info.sh", `[ "$1" = --describe ] && printf 'name: s\nboundary: -\ndenies: nothing (informational)\nwhy: w\n'`)
	if _, emit := guardRow(info, "info"); emit {
		t.Errorf("informational guard should be excluded")
	}

	stubbed := stub("stub.sh", `cat >/dev/null; exit 0`)
	if r, _ := guardRow(stubbed, "stub"); !reflect.DeepEqual(r, []string{"stub", "", "no manifest"}) {
		t.Errorf("stub row = %v, want no manifest", r)
	}
}
