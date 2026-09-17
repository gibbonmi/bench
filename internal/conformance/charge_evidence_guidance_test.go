package conformance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/anchors"
)

const boundedBuildActionDiagnosticPrefix = "bounded build action: "

// buildActionSentenceLead opens every guidance sentence that states one bounded build
// action prerequisite. A paragraph that carries the lead is pinned prose.
const buildActionSentenceLead = "Build action requires"

// The retired build full charge form is a flag pair, not one spelling. The guidance may
// name the preflight build command, and it may name the run control, but one line may not
// carry both.
const (
	preflightBuildCommand = "bench preflight build"
	retiredBuildFullFlag  = "--full"
)

func boundedBuildActionAnchors() []anchors.Anchor {
	return anchorsWithDiagnosticPrefix(boundedBuildActionDiagnosticPrefix)
}

// boundedBuildActionGuidance returns the one canonical reader the family pins. The rows'
// File fields are that path's single source, so a row that names a second file fails here.
func boundedBuildActionGuidance(t *testing.T, family []anchors.Anchor) string {
	t.Helper()
	files := map[string]bool{}
	for _, anchor := range family {
		files[anchor.File] = true
	}
	if len(files) != 1 {
		t.Fatalf("bounded build action rows name %d guidance files, want one", len(files))
	}
	for file := range files {
		return file
	}
	return ""
}

// TestEvidenceBuildGuidance grades the bounded build action rows for CE101, CE102, CE103,
// CE140, and CE141. The expectation table below is the independent half of this pair: it
// states each row's kind and diagnostic apart from the registry, so a reworded or deleted
// row bites. The count comes from that table's length, and the mismatch reports rather than
// stops, so a deleted row also reaches the table.
func TestEvidenceBuildGuidance(t *testing.T) {
	actionAnchors := boundedBuildActionAnchors()
	want := []struct {
		kind       anchors.Kind
		diagnostic string
	}{
		{anchors.Require, "permits build action without verified delivery"},
		{anchors.Require, "permits build action without available required context"},
		{anchors.Require, "permits build action without a current binding"},
		{anchors.Require, "permits build action without reviewer approval"},
		{anchors.Require, "permits build action without the complete task supplement"},
		{anchors.Forbid, "re-advertises the retired build full charge form"},
	}
	if got := len(actionAnchors); got != len(want) {
		t.Errorf("bounded build action anchor count = %d, want %d", got, len(want))
	}
	for _, expected := range want {
		found := false
		for _, anchor := range actionAnchors {
			if anchor.Kind == expected.kind && containsDiagnostic([]string{anchor.Diagnostic}, expected.diagnostic) {
				found = true
			}
		}
		if !found {
			t.Errorf("no row of the expected kind states %q", expected.diagnostic)
		}
	}
	runAnchorBites(t, actionAnchors, func(anchor anchors.Anchor) string { return anchor.Diagnostic })
}

// TestEvidenceBuildActionSentencesArePinned binds the guidance to the registry. A paragraph
// that states one bounded build action prerequisite is pinned prose, so exactly one Require
// needle covers each of its sentences. A needle that keeps only its lead clause, a needle
// that repeats another row's text, and an unpinned sentence added to that paragraph all
// bite here.
func TestEvidenceBuildActionSentencesArePinned(t *testing.T) {
	family := boundedBuildActionAnchors()
	needles := map[string]int{}
	for _, anchor := range family {
		if anchor.Kind == anchors.Require {
			needles[anchor.Needle]++
		}
	}
	for needle, rows := range needles {
		if rows != 1 {
			t.Errorf("%d Require rows share the needle %q, want one row each", rows, needle)
		}
	}
	guidance := boundedBuildActionGuidance(t, family)
	sentences := buildActionSentences(guidanceText(t, filepath.Join(NewHarness(t).KitRoot, filepath.FromSlash(guidance))))
	if len(sentences) == 0 {
		t.Fatalf("%s states no bounded build action prerequisite", guidance)
	}
	for _, sentence := range sentences {
		pinned := 0
		for needle := range needles {
			if strings.Contains(needle, sentence) {
				pinned++
			}
		}
		if pinned != 1 {
			t.Errorf("guidance sentence %q is pinned by %d Require rows, want one", sentence, pinned)
		}
	}
}

// TestEvidenceBuildGuidanceRejectsThePreflightBuildFullPair grades the guidance half of
// CE173. The retired form is the pair, so a guidance line that names the preflight build
// command may not also carry the run-control flag. Any flag order re-advertises the retired
// form and bites here. The phase's own full-run control names no preflight command, so it
// stays legal.
func TestEvidenceBuildGuidanceRejectsThePreflightBuildFullPair(t *testing.T) {
	guidance := boundedBuildActionGuidance(t, boundedBuildActionAnchors())
	lines := strings.Split(guidanceText(t, filepath.Join(NewHarness(t).KitRoot, filepath.FromSlash(guidance))), "\n")
	for index, line := range lines {
		if strings.Contains(line, preflightBuildCommand) && strings.Contains(line, retiredBuildFullFlag) {
			t.Errorf("%s line %d names %q and carries %q, which re-advertises the retired build full charge form",
				guidance, index+1, preflightBuildCommand, retiredBuildFullFlag)
		}
	}
}

// guidanceText reads one canonical reader with its line endings normalized.
func guidanceText(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return strings.ReplaceAll(string(data), "\r\n", "\n")
}

// buildActionSentences returns every sentence of every paragraph that carries the bounded
// build action lead. A paragraph runs between blank lines, and a sentence runs to the first
// period that a space, a line end, or the paragraph end follows.
func buildActionSentences(text string) []string {
	var found []string
	for _, paragraph := range strings.Split(text, "\n\n") {
		if !strings.Contains(paragraph, buildActionSentenceLead) {
			continue
		}
		found = append(found, paragraphSentences(paragraph)...)
	}
	return found
}

func paragraphSentences(paragraph string) []string {
	var found []string
	for rest := strings.TrimSpace(paragraph); rest != ""; {
		end := len(rest)
		for index := 0; index < len(rest); index++ {
			if rest[index] != '.' {
				continue
			}
			if index+1 == len(rest) || rest[index+1] == ' ' || rest[index+1] == '\n' {
				end = index + 1
				break
			}
		}
		found = append(found, strings.TrimSpace(rest[:end]))
		rest = strings.TrimSpace(rest[end:])
	}
	return found
}
