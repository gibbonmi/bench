package main

import "testing"

// The idiom-setting table test for the module: pure logic, table-driven, no process
// boundary. Acceptance for the version line lives at the shell seam (the gate's
// version-routing contract); this pins the format one layer down.
func TestVersionLine(t *testing.T) {
	cases := []struct {
		v, goos, goarch, want string
	}{
		{"0.2.0", "linux", "amd64", "benchkit 0.2.0 (linux/amd64)"},
		{"dev", "darwin", "arm64", "benchkit dev (darwin/arm64)"},
		{"1.0.0", "linux", "arm64", "benchkit 1.0.0 (linux/arm64)"},
	}
	for _, c := range cases {
		if got := versionLine(c.v, c.goos, c.goarch); got != c.want {
			t.Errorf("versionLine(%q,%q,%q) = %q, want %q", c.v, c.goos, c.goarch, got, c.want)
		}
	}
}

func TestRunVersionExits0(t *testing.T) {
	if rc := run([]string{"version"}, nil, nil); rc != 0 {
		t.Errorf("run version exit = %d, want 0", rc)
	}
}

func TestRunUnknownExits2(t *testing.T) {
	if rc := run([]string{"nope"}, nil, nil); rc != 2 {
		t.Errorf("run nope exit = %d, want 2", rc)
	}
}
