package commit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const malformedSpacing = "package sample\n\nfunc Value( )int{return 1}\n"
const formattedSpacing = "package sample\n\nfunc Value() int { return 1 }\n"
const malformedChange = "package sample\n\nfunc Value( )int{return 2}\n"
const formattedChange = "package sample\n\nfunc Value() int { return 2 }\n"

func TestCommandFormatsNamedGoFileBeforeAuthorization(t *testing.T) {
	root, _ := landingRepo(t, 0, func(t *testing.T, root string) {
		mustWrite(t, filepath.Join(root, "named.go"), formattedSpacing, 0o644)
		mustWrite(t, filepath.Join(root, ".bench", "gate.sh"), "#!/bin/sh\nset -eu\n[ \"$(sed -n '3p' named.go)\" = 'func Value() int { return 2 }' ]\n", 0o755)
	})
	mustWrite(t, filepath.Join(root, "named.go"), malformedChange, 0o644)

	code, stdout, stderr := runCommand(t, root, "-m", "format", "named.go")
	if code != 0 {
		t.Fatalf("formatted commit = (%d, %q, %q), want success", code, stdout, stderr)
	}
	got, err := os.ReadFile(filepath.Join(root, "named.go"))
	if err != nil || string(got) != formattedChange {
		t.Fatalf("named.go = %q, %v, want formatted bytes", got, err)
	}
	if !strings.Contains(stdout, "formatted Go paths: named.go\n") {
		t.Fatalf("stdout = %q, want formatted-path disclosure", stdout)
	}
}

func TestCommandFormatsOnlyNamedGoFiles(t *testing.T) {
	root, _ := landingRepo(t, 0, func(t *testing.T, root string) {
		mustWrite(t, filepath.Join(root, "named.go"), formattedSpacing, 0o644)
	})
	mustWrite(t, filepath.Join(root, "named.go"), malformedChange, 0o644)
	mustWrite(t, filepath.Join(root, "unnamed.go"), malformedSpacing, 0o644)

	code, stdout, stderr := runCommand(t, root, "-m", "format named", "named.go")
	if code != 0 {
		t.Fatalf("named commit = (%d, %q, %q), want success", code, stdout, stderr)
	}
	unnamed, err := os.ReadFile(filepath.Join(root, "unnamed.go"))
	if err != nil || string(unnamed) != malformedSpacing {
		t.Fatalf("unnamed.go = %q, %v, want untouched", unnamed, err)
	}
	if strings.Contains(stdout, "unnamed.go") {
		t.Fatalf("stdout disclosed an unnamed path: %q", stdout)
	}
}

func TestCommandFormatsNamedDirectoryChangedGoDescendants(t *testing.T) {
	const hostile = "pkg/-[x] name.go"
	root, _ := landingRepo(t, 0, func(t *testing.T, root string) {
		mustMkdirAll(t, filepath.Join(root, "pkg"))
		mustWrite(t, filepath.Join(root, "pkg", "tracked.go"), formattedSpacing, 0o644)
		mustWrite(t, filepath.Join(root, "pkg", "deleted.go"), formattedSpacing, 0o644)
	})
	mustWrite(t, filepath.Join(root, "pkg", "tracked.go"), malformedChange, 0o644)
	mustWrite(t, filepath.Join(root, filepath.FromSlash(hostile)), "package sample\n\nfunc Other( )int{return 3}\n", 0o644)
	if err := os.Remove(filepath.Join(root, "pkg", "deleted.go")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("tracked.go", filepath.Join(root, "pkg", "linked.go")); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(root, "outside.go"), malformedSpacing, 0o644)

	code, stdout, stderr := runCommand(t, root, "-m", "format directory", "pkg")
	if code != 0 {
		t.Fatalf("directory commit = (%d, %q, %q), want success", code, stdout, stderr)
	}
	for path, want := range map[string]string{
		"pkg/tracked.go": formattedChange,
		hostile:          "package sample\n\nfunc Other() int { return 3 }\n",
		"outside.go":     malformedSpacing,
	} {
		got, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil || string(got) != want {
			t.Fatalf("%s = %q, %v, want %q", path, got, err, want)
		}
	}
	if !strings.Contains(stdout, hostile) || strings.Contains(stdout, "outside.go") {
		t.Fatalf("formatted disclosure = %q", stdout)
	}
}

func TestCommandFormatFailureIsAtomic(t *testing.T) {
	root, _ := landingRepo(t, 0, func(t *testing.T, root string) {
		mustMkdirAll(t, filepath.Join(root, "pkg"))
		mustWrite(t, filepath.Join(root, "pkg", "valid.go"), formattedSpacing, 0o644)
	})
	mustWrite(t, filepath.Join(root, "pkg", "valid.go"), malformedChange, 0o644)
	invalid := "package sample\n\nfunc broken(\n"
	mustWrite(t, filepath.Join(root, "pkg", "invalid.go"), invalid, 0o644)

	code, stdout, stderr := runCommand(t, root, "-m", "format failure", "pkg")
	if code != 1 || stdout != "" || !strings.Contains(stderr, "format named Go files") {
		t.Fatalf("format failure = (%d, %q, %q), want pre-publication refusal", code, stdout, stderr)
	}
	valid, _ := os.ReadFile(filepath.Join(root, "pkg", "valid.go"))
	if string(valid) != malformedChange {
		t.Fatalf("valid.go changed despite atomic format refusal: %q", valid)
	}
}

func TestCommandDryRunDoesNotFormat(t *testing.T) {
	root, _ := landingRepo(t, 0, func(t *testing.T, root string) {
		mustWrite(t, filepath.Join(root, "named.go"), formattedSpacing, 0o644)
	})
	mustWrite(t, filepath.Join(root, "named.go"), malformedChange, 0o644)

	code, stdout, stderr := runCommand(t, root, "--dry-run", "-m", "dry", "named.go")
	if code != 0 || stderr != "" || !strings.Contains(stdout, "dry run:") {
		t.Fatalf("dry run = (%d, %q, %q), want success", code, stdout, stderr)
	}
	got, _ := os.ReadFile(filepath.Join(root, "named.go"))
	if string(got) != malformedChange {
		t.Fatalf("dry run formatted named.go: %q", got)
	}
}
