package guards

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestManifestFieldReadsLeadingCommentBlock(t *testing.T) {
	header := "#!/usr/bin/env bash\n  \t\n# threat-model prose\n# denies: destructive git operations   \n# name: block-dangerous-git\n# why: protects repository history\n# boundary: PreToolUse Bash\nset -uo pipefail\n# name: too-late\n"
	if got := manifestField(header, "name"); got != "block-dangerous-git" {
		t.Errorf("name = %q", got)
	}
	if got := manifestField(header, "denies"); got != "destructive git operations" {
		t.Errorf("denies = %q", got)
	}
	if got := manifestField(header, "absent"); got != "" {
		t.Errorf("absent = %q, want empty", got)
	}
	if got := manifestField("# name: \n# name: later\n", "name"); got != "" {
		t.Errorf("first empty name lost first-occurrence precedence: %q", got)
	}
}

func TestGuardRowReadsStaticHeader(t *testing.T) {
	dir := t.TempDir()
	write := func(name, body string) string {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
		return p
	}
	full := write("full.sh", "#!/usr/bin/env bash\n# why: w\n# name: g\n# denies: d\n# boundary: b\nexit 99\n")
	if r, emit := guardRow(full, "full"); !emit || !reflect.DeepEqual(r, []string{"g", "b", "d"}) {
		t.Errorf("full manifest row = %v emit=%v", r, emit)
	}

	info := write("info.sh", "#!/usr/bin/env bash\n# name: s\n# boundary: SessionStart\n# denies: nothing (informational)\n# why: w\nexit 99\n")
	if _, emit := guardRow(info, "info"); emit {
		t.Errorf("informational guard should be excluded")
	}

	for _, tc := range []struct {
		name string
		body string
	}{
		{"absent", "#!/usr/bin/env bash\nexit 99\n"},
		{"incomplete", "#!/usr/bin/env bash\n# name: partial\n# boundary: b\n# denies: d\nexit 99\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := write(tc.name+".sh", tc.body)
			if r, emit := guardRow(path, tc.name); !emit || !reflect.DeepEqual(r, []string{tc.name, "", "no manifest"}) {
				t.Errorf("row = %v emit=%v, want no manifest", r, emit)
			}
		})
	}
}
