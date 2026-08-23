package adopt

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// setup_prompt_test.go pins FT76 story 3: one-at-a-time question sequencing at the
// prompt-I/O unit seam, plus the single-confirm default for an ambiguity-free plan.

// TestAskQuestionsSequencing drives the low-level sequencing function directly with a
// synthetic two-question set. This proves order and per-question application without
// depending on setup ever growing a second real ambiguity source.
func TestAskQuestionsSequencing(t *testing.T) {
	var order []string
	qs := []setupQuestion{
		{
			prompt:  "first question",
			options: []string{"a", "b"},
			apply: func(choice int, f *setupFacts) {
				order = append(order, "first")
				f.gateCommand = []string{"a", "b"}[choice]
			},
		},
		{
			prompt:  "second question",
			options: []string{"c", "d"},
			apply: func(choice int, f *setupFacts) {
				order = append(order, "second")
				f.profileName = []string{"c", "d"}[choice]
			},
		},
	}
	var stdout bytes.Buffer
	reader := bufio.NewReader(strings.NewReader("2\n1\n"))
	facts := setupFacts{}
	if err := askQuestions(qs, &facts, reader, &stdout); err != nil {
		t.Fatalf("askQuestions: %v", err)
	}
	if !slicesEqual(order, []string{"first", "second"}) {
		t.Fatalf("question application order = %v, want [first second]", order)
	}
	out := stdout.String()
	if strings.Index(out, "first question") > strings.Index(out, "second question") {
		t.Fatalf("first question printed after second:\n%s", out)
	}
	if facts.gateCommand != "b" || facts.profileName != "c" {
		t.Fatalf("facts after answers = %+v, want gateCommand=b profileName=c", facts)
	}
}

// TestAskQuestionsRejectsOutOfRangeAnswer proves an invalid answer fails closed rather
// than silently resolving the ambiguity.
func TestAskQuestionsRejectsOutOfRangeAnswer(t *testing.T) {
	qs := []setupQuestion{{prompt: "q", options: []string{"only"}, apply: func(int, *setupFacts) {}}}
	reader := bufio.NewReader(strings.NewReader("9\n"))
	var stdout bytes.Buffer
	facts := setupFacts{}
	if err := askQuestions(qs, &facts, reader, &stdout); err == nil {
		t.Fatal("askQuestions accepted an out-of-range answer")
	}
}

// TestAskQuestionsEmptySetIsANoop proves an ambiguity-free plan asks nothing.
func TestAskQuestionsEmptySetIsANoop(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader(""))
	var stdout bytes.Buffer
	facts := setupFacts{}
	if err := askQuestions(nil, &facts, reader, &stdout); err != nil {
		t.Fatalf("askQuestions with no questions: %v", err)
	}
	if stdout.String() != "" {
		t.Fatalf("askQuestions with no questions printed %q, want nothing", stdout.String())
	}
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestSetupInteractiveSingleConfirm proves an ambiguity-free plan on a real TTY asks
// exactly one confirm and no questions.
func TestSetupInteractiveSingleConfirm(t *testing.T) {
	setupPromptTestRepo(t)
	var stdout, stderr bytes.Buffer
	alwaysTTY := func(io.Reader) bool { return true }
	code := setup(nil, strings.NewReader("y\n"), &stdout, &stderr, "1.0.0", alwaysTTY)
	if code != 0 {
		t.Fatalf("setup exit = %d, stderr:\n%s", code, stderr.String())
	}
	out := stdout.String()
	if got := strings.Count(out, "Proceed with this plan?"); got != 1 {
		t.Fatalf("confirm prompt printed %d times, want exactly 1:\n%s", got, out)
	}
	if strings.Contains(out, "select [1-") {
		t.Fatalf("ambiguity-free plan asked a question:\n%s", out)
	}
}

// TestSetupInteractiveResolvesAmbiguityOneAtATime proves a genuinely ambiguous plan on a
// real TTY asks the open question before the single confirm. The chosen answer lands in
// the written gate.sh.
func TestSetupInteractiveResolvesAmbiguityOneAtATime(t *testing.T) {
	root := setupPromptTestRepo(t)
	if err := os.WriteFile(filepath.Join(root, "package.json"), []byte(`{"name":"fixture","scripts":{"test":"echo ok"}}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	alwaysTTY := func(io.Reader) bool { return true }
	// go.mod sorts before package.json in detectGateCandidates, so answer "1" selects
	// go test ./... - the same command an equivalent --yes ambiguity-free fixture would
	// have gotten, letting the assertion below pin exactly which candidate won.
	code := setup(nil, strings.NewReader("1\ny\n"), &stdout, &stderr, "1.0.0", alwaysTTY)
	if code != 0 {
		t.Fatalf("setup exit = %d, stderr:\n%s", code, stderr.String())
	}
	out := stdout.String()
	questionAt := strings.Index(out, "which one should .bench/gate.sh run?")
	confirmAt := strings.Index(out, "Proceed with this plan?")
	if questionAt < 0 || confirmAt < 0 || questionAt > confirmAt {
		t.Fatalf("question did not precede the single confirm:\n%s", out)
	}
	if got := strings.Count(out, "Proceed with this plan?"); got != 1 {
		t.Fatalf("confirm prompt printed %d times, want exactly 1:\n%s", got, out)
	}
	gate, err := os.ReadFile(filepath.Join(root, ".bench", "gate.sh"))
	if err != nil {
		t.Fatalf("read gate.sh: %v", err)
	}
	if !strings.Contains(string(gate), "go test ./...") {
		t.Fatalf("gate.sh does not embed the answered candidate:\n%s", gate)
	}
}

func setupPromptTestRepo(t *testing.T) string {
	t.Helper()
	root := setupBareTestRepo(t)
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/fixture\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// setupBareTestRepo is setupPromptTestRepo without the go.mod: a zero-signal repository
// carrying no build-system file at all. This is what the seeded gate input manifest has
// to survive alongside the fail-closed gate stub.
func setupBareTestRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	cmd := exec.Command("git", "-C", root, "init", "-q")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_KIT", filepath.Clean(filepath.Join(wd, "..", "..")))
	t.Chdir(root)
	return root
}

// TestSetupPreviewNamesGateInputs pins SD2: the plan preview names
// .bench/gate-inputs.json in both states, in the profile line's voice.
func TestSetupPreviewNamesGateInputs(t *testing.T) {
	root := setupPromptTestRepo(t)
	kit := os.Getenv("BENCH_KIT")

	absent := renderSetupPreview(root, inspectRepo(root, kit))
	if !strings.Contains(absent, "  .bench/gate-inputs.json absent -> will be seeded declaring BENCH_HOME and HOME\n") {
		t.Fatalf("preview does not announce the seed when absent:\n%s", absent)
	}

	if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".bench", "gate-inputs.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	present := renderSetupPreview(root, inspectRepo(root, kit))
	if !strings.Contains(present, "  .bench/gate-inputs.json exists -> left as-is (reviewer-owned content)\n") {
		t.Fatalf("preview does not report the existing file as left as-is:\n%s", present)
	}
	if strings.Contains(present, "will be seeded") {
		t.Fatalf("preview still announces a seed for an existing file:\n%s", present)
	}
}
