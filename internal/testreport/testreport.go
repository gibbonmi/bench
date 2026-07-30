// Package testreport runs a fresh Go test invocation and renders its observed package terminals.
package testreport

import (
	"encoding/json"
	"os/exec"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

var grammar = usage.Grammar{
	Cmd:     "bench test [--full] [package]",
	Help:    "usage: bench test [--full] [package]",
	Flags:   []usage.Flag{{Name: "--full"}},
	MaxArgs: 1,
}

// Command runs Go from root and renders one stable row for each terminal package event.
func Command(root string, args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	packageExpr := "./..."
	if len(parsed.Positionals) == 1 {
		packageExpr = parsed.Positionals[0]
	}

	cmd := exec.Command("go", "test", "-json", "-count=1", packageExpr)
	cmd.Dir = root
	raw, err := cmd.Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			return toon.Errorf("go test failed to start", err.Error()) + "\n", 1
		}
	}
	packages := packageRows(raw)
	out, renderErr := render(packages)
	if renderErr != nil {
		return toon.RenderError(renderErr) + "\n", 1
	}
	if err != nil {
		return out, 1
	}
	return out, 0
}

type event struct {
	Action  string
	Package string
	Test    string
	Output  string
}

func packageRows(raw []byte) [][]string {
	statuses := map[string]string{}
	for _, line := range strings.Split(string(raw), "\n") {
		var e event
		if json.Unmarshal([]byte(line), &e) != nil || e.Package == "" {
			continue
		}
		if e.Test == "" && strings.Contains(e.Output, "[no test files]") {
			statuses[e.Package] = "no-tests"
		}
		if e.Action == "pass" {
			if statuses[e.Package] != "no-tests" {
				statuses[e.Package] = "pass"
			}
		}
		if e.Action == "fail" {
			statuses[e.Package] = "fail"
		}
	}
	packages := make([]string, 0, len(statuses))
	for pkg := range statuses {
		packages = append(packages, pkg)
	}
	sort.Strings(packages)
	rows := make([][]string, 0, len(packages))
	for _, pkg := range packages {
		rows = append(rows, []string{pkg, statuses[pkg]})
	}
	return rows
}

func render(packages [][]string) (string, error) {
	packageBlock, err := toon.Table("packages", []string{"package", "status"}, packages)
	if err != nil {
		return "", err
	}
	failureBlock, err := toon.Table("failures", []string{"package", "test", "line"}, nil)
	if err != nil {
		return "", err
	}
	skipBlock, err := toon.Table("skips", []string{"package", "test", "reason"}, nil)
	if err != nil {
		return "", err
	}
	return packageBlock + failureBlock + skipBlock, nil
}
