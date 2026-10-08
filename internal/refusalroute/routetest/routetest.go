// Package routetest holds the walk of a printed refusal route that the write verbs' tests
// share: it reads the route a face printed, carries the route out step by step, and checks
// the producing fixtures against the registry. How one step runs differs by package, so
// each walk passes its own runner, and this package owns none.
package routetest

import (
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/shellcommand"
)

// ReviewerMarker is the marker a reviewer route opens with. The spec pins it, so the walk
// spells it rather than read it from the renderer it grades.
const ReviewerMarker = "reviewer: "

// CandidateFaces are the commitment faces that only the authorization of a commit's
// candidate raises. The commitment verb authorizes no candidate, so the commit walk proves
// these faces, and the commitment walk proves every other commitment face.
var CandidateFaces = []string{refusalroute.CommitmentUnbound, refusalroute.CommitmentLightPathOutside}

// walked reports whether the walk of verb proves face: a face of verb's own, except that the
// commit walk proves each candidate face in place of the commitment walk.
func walked(verb refusalroute.Verb, face refusalroute.Face) bool {
	if slices.Contains(CandidateFaces, face.Name) {
		return verb == refusalroute.Commit
	}
	return face.Verb == verb
}

// Fixtures requires a producing fixture for each cause of each face that the walk of verb
// proves, and such a face for each fixture. keys holds each fixture's face and cause, with
// an empty cause for a face that one cause raises. It returns the proved faces by name.
func Fixtures(t testing.TB, verb refusalroute.Verb, keys [][2]string) map[string]refusalroute.Face {
	t.Helper()
	faces := map[string]refusalroute.Face{}
	for _, face := range refusalroute.Inventory() {
		if walked(verb, face) {
			faces[face.Name] = face
		}
	}
	produced := map[string]bool{}
	for _, key := range keys {
		if produced[key[0]+"/"+key[1]] {
			t.Fatalf("%s face %q has two producing fixtures for cause %q", verb, key[0], key[1])
		}
		if _, ok := faces[key[0]]; !ok {
			t.Fatalf("fixture %q produces no face that the %s walk proves", key[0], verb)
		}
		produced[key[0]], produced[key[0]+"/"+key[1]] = true, true
	}
	for name := range faces {
		if !produced[name] {
			t.Errorf("face %q that the %s walk proves has no producing fixture", name, verb)
		}
	}
	return faces
}

// Next is the route of the next= line that output prints. It reports false unless output
// prints exactly one such line, because a refusal prints one route.
func Next(output string) (string, bool) {
	var routes []string
	for _, line := range strings.Split(output, "\n") {
		if route, ok := strings.CutPrefix(line, refusalroute.NextField+"="); ok {
			routes = append(routes, route)
		}
	}
	if len(routes) != 1 {
		return "", false
	}
	return routes[0], true
}

// NextCell is the route in the one row of the next table that output ends with, the table
// the commitment verb prints. It reports false unless output prints exactly one such table,
// and that table decodes to one row whose one cell is the route.
func NextCell(output string) (string, bool) {
	lines := strings.SplitAfter(output, "\n")
	start := -1
	for index, line := range lines {
		if strings.HasPrefix(line, refusalroute.NextField+"[") {
			if start >= 0 {
				return "", false
			}
			start = index
		}
	}
	if start < 0 {
		return "", false
	}
	document, err := axitest.DecodeDocument(strings.Join(lines[start:], ""))
	if err != nil || len(document.Blocks) != 1 {
		return "", false
	}
	rows, err := document.Rows(refusalroute.NextField)
	if err != nil || len(rows) != 1 {
		return "", false
	}
	row, _ := rows[0].(map[string]any)
	route, cell := row[refusalroute.NextColumn].(string)
	if !cell || len(row) != 1 {
		return "", false
	}
	return route, true
}

// Steps checks one printed route against the face that printed it and returns its steps.
// The reviewer marker opens the route exactly when the face is the reviewer's, the route
// names each of names, and after the preface it prints the steps the face declares.
func Steps(t testing.TB, face refusalroute.Face, next, preface string, names []string) []string {
	t.Helper()
	route, reviewer := strings.CutPrefix(next, ReviewerMarker)
	if reviewer != (face.Authority == refusalroute.Reviewer) {
		t.Fatalf("%s next = %q, want the reviewer marker exactly when the face is the reviewer's", face.Name, next)
	}
	for _, name := range names {
		if !strings.Contains(route, name) {
			t.Fatalf("%s next = %q, want it to name %q", face.Name, next, name)
		}
	}
	steps := refusalroute.Steps(strings.TrimPrefix(route, preface+"; "))
	if len(steps) != len(face.Route) {
		t.Fatalf("%s next = %q prints %d steps, want the %d the face declares", face.Name, next, len(steps), len(face.Route))
	}
	return steps
}

// Follow carries out a printed route step by step and returns the result of the last step
// it ran. carry answers the fixture's own means for one step, by the step's index in the
// face's route: an instruction, or a step of a reviewer route. run runs every other step
// verbatim. An agent command step always runs verbatim, so the route the agent reads is
// the route the walk proves.
func Follow[R any](t testing.TB, face refusalroute.Face, steps []string, carry func(index int) (func(), bool), run func(step string) R) R {
	t.Helper()
	var last R
	for index, step := range face.Route {
		carried, ok := carry(index)
		switch {
		case ok && face.Authority == refusalroute.Agent && step.IsCommand():
			t.Fatalf("%s step %d is an agent command, which the walk runs verbatim; the fixture may not carry it", face.Name, index+1)
		case ok:
			carried()
		case step.IsCommand():
			last = run(steps[index])
		default:
			t.Fatalf("%s step %d %q is an instruction that the fixture does not carry out", face.Name, index+1, steps[index])
		}
	}
	return last
}

// Words checks that a printed command step reads as the operator would paste it, and
// returns its words: one simple command with no placeholder left, that runs a Bench verb.
func Words(t testing.TB, step string) []string {
	t.Helper()
	stream := shellcommand.Parse(step)
	if stream.Unlexed || len(stream.Commands) != 1 || len(stream.Tokens) != stream.Commands[0].End-stream.Commands[0].Start || strings.Contains(step, "<") {
		t.Fatalf("printed step %q is not one simple command with every slot filled", step)
	}
	words := shellcommand.ProjectCommandWords(stream.Tokens)
	if len(words) < 2 || words[0] != "bench" {
		t.Fatalf("printed step %q runs no Bench verb", step)
	}
	return words
}

// Diagnostic reports whether a step's words are `bench doctor`. Its exit reports the health
// of the whole install, and a fixture repository is never a linked one, so a walk runs the
// step verbatim and does not grade its exit. Every other Bench verb must exit 0.
func Diagnostic(words []string) bool {
	return len(words) == 2 && words[0] == "bench" && words[1] == "doctor"
}
