package preflight

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestChargeRefusesLiveLinkedRequiredSource(t *testing.T) {
	root, slug := seedConformant(t)
	if err := os.Remove(buildPhase); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("../../internal/"+slug+"/foo.go", buildPhase); err != nil {
		t.Fatal(err)
	}
	runGit(t, "add", "-A")
	runGit(t, "commit", "-q", "-m", "linked phase")
	args := chargeArgs(t, root, slug, false)
	args[8] = runGit(t, "rev-parse", "HEAD")
	out, code := Command(args)
	if code != 1 || !strings.Contains(out, "source required") || !strings.Contains(out, buildPhase) || strings.Contains(out, "complete") {
		t.Fatalf("live linked source = (%d):\n%s", code, out)
	}
}

func TestChargeRequiredSpecAndCoverageRefuse(t *testing.T) {
	for _, test := range []struct {
		name    string
		prepare func(t *testing.T, slug string)
	}{
		{"spec", func(t *testing.T, slug string) {
			if err := os.Remove("specs/" + slug + "/spec.md"); err != nil {
				t.Fatal(err)
			}
		}},
		{"coverage", func(t *testing.T, slug string) {
			mustWriteFile(t, "specs/"+slug+"/spec.md", strings.Replace(specBody(slug), "### Acceptance coverage map", "### Other map", 1))
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, slug := seedConformant(t)
			test.prepare(t, slug)
			runGit(t, "add", "-A")
			runGit(t, "commit", "-q", "-m", "missing charge input")
			args := chargeArgs(t, root, slug, false)
			args[8] = runGit(t, "rev-parse", "HEAD")
			out, code := Command(args)
			if code != 1 || !strings.Contains(out, "error:") || strings.Contains(out, "complete") {
				t.Fatalf("%s required source = (%d):\n%s", test.name, code, out)
			}
		})
	}
}

func TestChargeRefusesRequiredSourceOutsidePinnedTip(t *testing.T) {
	root, slug := seedConformant(t)
	runGit(t, "rm", "--cached", buildPhase)
	mustWriteFile(t, filepath.Join(root, ".git/info/exclude"), buildPhase+"\n")
	runGit(t, "commit", "-q", "-m", "remove required source from tip")

	out, code := Command(chargeArgs(t, root, slug, true))
	if code != 1 || !strings.Contains(out, "source required") ||
		!strings.Contains(out, buildPhase) || !strings.Contains(out, "source tip") ||
		strings.Contains(out, "complete") {
		t.Fatalf("ignored source outside tip = (%d):\n%s", code, out)
	}
}

func TestChargeDistinguishesAbsentAndEmptyRequiredInputs(t *testing.T) {
	tests := []struct {
		name       string
		absent     func(t *testing.T, slug string)
		empty      func(t *testing.T, slug string)
		absentWant string
		emptyWant  string
	}{
		{
			name: "spec",
			absent: func(t *testing.T, slug string) {
				if err := os.Remove("specs/" + slug + "/spec.md"); err != nil {
					t.Fatal(err)
				}
			},
			empty: func(t *testing.T, slug string) {
				mustWriteFile(t, "specs/"+slug+"/spec.md", "")
			},
			absentWant: "spec folder is missing",
			emptyWant:  "Status:  (want staged)",
		},
		{
			name: "tickets directory",
			absent: func(t *testing.T, slug string) {
				if err := os.RemoveAll("specs/" + slug + "/tickets"); err != nil {
					t.Fatal(err)
				}
			},
			empty: func(t *testing.T, slug string) {
				if err := os.Remove("specs/" + slug + "/tickets/one.md"); err != nil {
					t.Fatal(err)
				}
			},
			absentWant: "tickets directory is absent",
			emptyWant:  "rows-owned",
		},
		{
			name: "selected ticket",
			absent: func(t *testing.T, slug string) {
				mustWriteFile(t, "specs/"+slug+"/tickets/two.md", ticketDoc("Two", "PF1", "PF2"))
				if err := os.Remove("specs/" + slug + "/tickets/one.md"); err != nil {
					t.Fatal(err)
				}
			},
			empty: func(t *testing.T, slug string) {
				mustWriteFile(t, "specs/"+slug+"/tickets/one.md", "")
			},
			absentWant: "selected ticket",
			emptyWant:  "tickets-parse",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			outputs := make([]string, 2)
			for i, state := range []struct {
				name    string
				prepare func(*testing.T, string)
				want    string
			}{{"absent", test.absent, test.absentWant}, {"empty", test.empty, test.emptyWant}} {
				t.Run(state.name, func(t *testing.T) {
					root, slug := seedConformant(t)
					state.prepare(t, slug)
					runGit(t, "add", "-A")
					runGit(t, "commit", "-q", "-m", state.name+" required input")
					out, code := Command(chargeArgs(t, root, slug, false))
					outputs[i] = out
					if code != 1 || !strings.Contains(out, state.want) ||
						!strings.Contains(out, " — ") || strings.Contains(out, "complete") {
						t.Fatalf("%s %s = (%d):\n%s", test.name, state.name, code, out)
					}
				})
			}
			if outputs[0] == outputs[1] {
				t.Fatalf("absent and empty %s produced the same diagnostic", test.name)
			}
		})
	}
}

func TestChargeDistinguishesAbsentAndEmptyGuidanceSources(t *testing.T) {
	for _, source := range []string{delegateSkill, delegateProcedure, buildPhase} {
		t.Run(source, func(t *testing.T) {
			for _, state := range []string{"absent", "empty"} {
				t.Run(state, func(t *testing.T) {
					root, slug := seedConformant(t)
					if state == "absent" {
						if err := os.Remove(source); err != nil {
							t.Fatal(err)
						}
					} else {
						mustWriteFile(t, source, "")
					}
					runGit(t, "add", "-A")
					runGit(t, "commit", "-q", "-m", state+" guidance source")
					out, code := Command(chargeArgs(t, root, slug, false))
					if code != 1 || !strings.Contains(out, "source required") ||
						!strings.Contains(out, source+" is "+state) ||
						!strings.Contains(out, "restore the named canonical source") ||
						strings.Contains(out, "complete") {
						t.Fatalf("%s %s = (%d):\n%s", source, state, code, out)
					}
				})
			}
		})
	}
}

func TestChargeGrammarBoundariesRefuse(t *testing.T) {
	root, slug := seedConformant(t)
	valid := chargeArgs(t, root, slug, false)
	base, tip := valid[6], valid[8]
	tests := []struct {
		name string
		args []string
		code int
		want string
	}{
		{"duplicate option", append(append([]string{}, valid...), "--base", base), 2, "unknown argument: --base"},
		{"unknown flag", append(append([]string{}, valid...), "--unknown"), 2, "unknown argument"},
		{"missing value", []string{"build", slug, "--charge", "--ticket", "one.md", "--base", base, "--source-tip"}, 2, "missing argument: --source-tip"},
		{"review charge", []string{"review", slug, "--charge", "--ticket", "one.md", "--base", base, "--source-tip", tip}, 2, "--charge requires build"},
		{"full without charge", []string{"build", slug, "--full"}, 2, "--ticket and --full require --charge"},
		{"ticket without charge", []string{"build", slug, "--ticket", "one.md"}, 2, "--ticket and --full require --charge"},
		{"missing ticket", []string{"build", slug, "--charge", "--base", base, "--source-tip", tip}, 2, "--charge requires build"},
		{"ticket traversal", []string{"build", slug, "--charge", "--ticket", "../one.md", "--base", base, "--source-tip", tip}, 1, "selected ticket"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out, code := Command(test.args)
			if code != test.code || !strings.Contains(out, test.want) ||
				strings.Contains(out, "complete") {
				t.Fatalf("grammar case = (%d):\n%s", code, out)
			}
		})
	}
}
