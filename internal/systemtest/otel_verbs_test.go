//go:build system

package systemtest

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/otelrecord"
)

// TestOtelCommitAndLandingRecordJourney drives the built binary's two publication verbs
// against a scaffolded repository and reads the record each one left behind. The commit
// runs under its own BENCH_HOME, because the two verbs grade different roots and the
// record is addressed by root.
//
// Rows OT18 and OT19.
func TestOtelCommitAndLandingRecordJourney(t *testing.T) {
	// The commit subject is the objective text the record must never copy. It is a token
	// no path, digest, or seam name can hold, so a match is the record carrying it.
	const subject = "objective subject text 8f3a"
	fixture := scaffoldRecordedPublicationRepo(t)
	source := systemCreateLandingWorktree(t, fixture.root, fixture.home, "otel-record", "otel record")
	base := systemGitOutput(t, fixture.root, "rev-parse", "main")

	if err := os.WriteFile(filepath.Join(source.path, "landed.txt"), []byte("landed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	commitHome := filepath.Join(t.TempDir(), "commit-home")
	committed := systemSelected(t, source.path, fixture.environment(commitHome), "commit", "-m", subject, "landed.txt")
	if committed.code != 0 {
		t.Fatalf("bench commit = (%d, %q, %q)", committed.code, committed.stdout, committed.stderr)
	}
	source.tip = systemGitOutput(t, source.path, "rev-parse", "HEAD")

	commitSpans := readSeamRecord(t, commitHome)
	commit, present := commitSpans["commit"]
	if !present {
		t.Fatalf("the record has no commit span; it holds %v", spanNames(commitSpans))
	}
	// OT18: the commit span carries the subject digest, the outcome, and the size of the
	// write set it composed.
	if commit.Attributes[otelrecord.AttrSubjectID] != source.tip {
		t.Fatalf("commit subject id = %q, want the published commit %s",
			commit.Attributes[otelrecord.AttrSubjectID], source.tip)
	}
	if commit.Attributes[otelrecord.AttrOutcome] != "green" {
		t.Fatalf("commit outcome = %q, want the verb's own exit", commit.Attributes[otelrecord.AttrOutcome])
	}
	if commit.Attributes[otelrecord.AttrMeasurePathCount] != "1" {
		t.Fatalf("commit path count = %q, want the one attributed path",
			commit.Attributes[otelrecord.AttrMeasurePathCount])
	}
	// The commit subject is objective text, and the record is a third durable place it
	// must not reach. No span of the run carries it, in any attribute.
	for name, span := range commitSpans {
		for key, value := range span.Attributes {
			if strings.Contains(value, subject) {
				t.Fatalf("span %s carries the commit subject text in %s: %q", name, key, value)
			}
		}
	}

	landed := systemSelected(t, fixture.root, fixture.environment(fixture.home),
		"worktree", "land", "--request", source.request, "--base", base, "--source-tip", source.tip,
		"-m", "land the recorded source", source.path)
	if landed.code != 0 {
		t.Fatalf("bench worktree land = (%d, %q, %q)", landed.code, landed.stdout, landed.stderr)
	}
	published := systemGitOutput(t, fixture.root, "rev-parse", "main")

	landingSpans := readSeamRecord(t, fixture.home)
	landing, present := landingSpans["worktree.land"]
	if !present {
		t.Fatalf("the record has no landing span; it holds %v", spanNames(landingSpans))
	}
	// OT19: the census raw-call count survives the release step that drops the records.
	// The landing states the same count on its own envelope, so the record and the verb
	// answer for one landing.
	wantCensus := envelopeField(t, landed.stdout, "census")
	if landing.Attributes[otelrecord.AttrMeasureCensusRawCalls] != wantCensus {
		t.Fatalf("landing census raw calls = %q, want the landing's own count %q",
			landing.Attributes[otelrecord.AttrMeasureCensusRawCalls], wantCensus)
	}
	// The landing span also carries the reviewed write set's size and the published
	// subject. One commit touching one path is the reviewed name-only diff.
	if landing.Attributes[otelrecord.AttrMeasurePathCount] != "1" {
		t.Fatalf("landing path count = %q, want the one reviewed path",
			landing.Attributes[otelrecord.AttrMeasurePathCount])
	}
	if landing.Attributes[otelrecord.AttrSubjectID] != published {
		t.Fatalf("landing subject id = %q, want the published commit %s",
			landing.Attributes[otelrecord.AttrSubjectID], published)
	}

	// Story 19: no span of either run carries an attribute outside the declared set.
	for _, spans := range []map[string]recordedSpan{commitSpans, landingSpans} {
		for name, span := range spans {
			for key := range span.Attributes {
				if !slices.Contains(otelrecord.DeclaredAttributes, key) {
					t.Fatalf("span %s carries the undeclared attribute %q", name, key)
				}
			}
		}
	}
}

// envelopeField answers one name=value field of a terminal envelope line.
func envelopeField(t *testing.T, out, name string) string {
	t.Helper()
	for _, field := range strings.Split(strings.Trim(strings.TrimSpace(out), "}"), ",") {
		if key, value, found := strings.Cut(field, "="); found && strings.HasSuffix(key, name) {
			return value
		}
	}
	t.Fatalf("no %s field in the envelope %q", name, out)
	return ""
}

// recordedPublicationRepo is a scaffolded repository whose gate is green, so the two
// publication verbs reach their own publication rather than a refusal.
type recordedPublicationRepo struct {
	root, home string
}

func scaffoldRecordedPublicationRepo(t *testing.T) recordedPublicationRepo {
	t.Helper()
	root, err := os.MkdirTemp(owner.root, "otel-publication [root]-")
	if err != nil {
		t.Fatal(err)
	}
	// The home holds the record the test reads back, and the reader globs for it, so the
	// home carries no path characters a glob pattern would read as a character class.
	home := filepath.Join(t.TempDir(), "bench-home")
	if result := owner.runAt(root, nil, "git", "init", "-q", "-b", "main"); result.code != 0 {
		t.Fatalf("git init = (%d, %q)", result.code, result.stderr)
	}
	// Both verbs create their own commits here, so the identity belongs in the config.
	for _, identity := range [][]string{{"user.email", "bench@local"}, {"user.name", "bench"}} {
		systemGit(t, root, "config", identity[0], identity[1])
	}
	if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, file := range []string{"gate.sh", "gate-prospective.sh"} {
		if err := os.WriteFile(filepath.Join(root, ".bench", file), []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	inputs := "{\"schema\":1,\"closure\":\"local\",\"environment\":[],\"paths\":[],\"tools\":[]}\n"
	if err := os.WriteFile(filepath.Join(root, ".bench", "gate-inputs.json"), []byte(inputs), 0o644); err != nil {
		t.Fatal(err)
	}
	systemGit(t, root, "add", ".")
	systemGit(t, root, "commit", "-qm", "publication record base")
	systemGit(t, root, "update-ref", "refs/bench/green/main", systemGitOutput(t, root, "rev-parse", "HEAD"))
	return recordedPublicationRepo{root: root, home: home}
}

// environment answers the overrides a run recording under benchHome needs.
func (r recordedPublicationRepo) environment(benchHome string) []string {
	return []string{
		"BENCH_RUN_BINARY=" + owner.selected.path,
		"BENCH_KIT=" + owner.kit,
		"BENCH_COMMAND_OBSERVE=1",
		"BENCH_SYSTEM_ROOT=" + r.root,
		"BENCH_HOME=" + benchHome,
	}
}
