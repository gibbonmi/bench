package spec

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/toon"
)

// positionalHistory is the complete `bench spec history <slug>` stdout.
func positionalHistory(t *testing.T, slug string) string {
	t.Helper()
	out, code, _ := Command([]string{"history", slug})
	if code != 0 {
		t.Fatalf("positional history %q exit=%d output=%q", slug, code, out)
	}
	return out
}

func TestSelectedSpecHistories(t *testing.T) {
	f := historyFixture(t)
	out, code := selectedHistory(t, "9", "only-delete", "mixed", "specs/mixed/spec.md", "retire-only", "absent", "mixed-extra")
	if code != 0 {
		t.Fatalf("selected history exit=%d output=%q", code, out)
	}
	var expectedRows [][]any
	for _, selected := range []struct{ slug, baseline string }{
		{"only-delete", "delete-only"}, {"mixed", "mixed"}, {"retire-only", "retire-only"}, {"mixed-extra", "prefix-slug"},
	} {
		for _, row := range f.capturedRows(t, selected.baseline) {
			expectedRows = append(expectedRows, []any{selected.slug, row[0], row[1], row[2], row[3]})
		}
	}
	want, err := toon.TableTyped("history", []string{"slug", "hash", "date", "kind", "subject"}, expectedRows)
	if err != nil {
		t.Fatal(err)
	}
	if got, expected := historyRows(t, out, "history"), historyRows(t, want, "history"); !reflect.DeepEqual(got, expected) {
		t.Fatalf("selected events=%#v; want %#v", got, expected)
	}
	var order []string
	for _, row := range historyRows(t, out, "histories") {
		order = append(order, row["target"].(string))
		if row["error"] != "" {
			t.Fatalf("summary=%#v", row)
		}
	}
	if !reflect.DeepEqual(order, []string{"only-delete", "mixed", "retire-only", "absent", "mixed-extra"}) {
		t.Fatalf("first-occurrence order=%v", order)
	}
}

func TestSelectedSpecLimit(t *testing.T) {
	f := historyFixture(t)
	mixed := []string{f.hash["newest"], f.hash["mixed-delete"], f.hash["both"]}
	for _, limit := range []int{1, 2, 3, 4} {
		t.Run(strconv.Itoa(limit), func(t *testing.T) {
			out, code := selectedHistory(t, strconv.Itoa(limit), "mixed", "retire-only", "absent")
			if code != 0 {
				t.Fatalf("exit=%d output=%q", code, out)
			}
			want := append(append([]string{}, mixed[:min(limit, len(mixed))]...), f.hash["retire"])
			var hashes []string
			for _, row := range historyRows(t, out, "history") {
				hashes = append(hashes, row["hash"].(string))
			}
			if !reflect.DeepEqual(hashes, want) {
				t.Fatalf("bounded event identities=%v; want %v", hashes, want)
			}
		})
	}
}

func TestSelectedSpecOmissions(t *testing.T) {
	historyFixture(t)
	for _, c := range []struct{ limit, omitted string }{{"1", "2"}, {"2", "1"}, {"3", "0"}, {"4", "0"}} {
		out, code := selectedHistory(t, c.limit, "mixed", "absent", "retire-only")
		if code != 0 {
			t.Fatalf("exit=%d output=%q", code, out)
		}
		rows := historyRows(t, out, "histories")
		for i, want := range []struct{ total, omitted string }{{"3", c.omitted}, {"0", "0"}, {"1", "0"}} {
			row := rows[i]
			if fmt.Sprint(row["total_events"]) != want.total || fmt.Sprint(row["omitted_events"]) != want.omitted || row["error"] != "" {
				t.Fatalf("limit=%s summary=%#v; want total=%s omitted=%s", c.limit, row, want.total, want.omitted)
			}
		}
	}
}

func TestSelectedSpecDetailRoute(t *testing.T) {
	f := historyFixture(t)
	targets := []string{"mixed", "a space", "quote' name", "café", "--limit", "help", "$(touch sentinel)", "[x]*?", "x.md.md", ".md"}
	for i, target := range targets[1:] {
		if target != ".md" {
			f.commit(t, target, i+9, "spec-retire: "+SlugOf(target))
		}
	}
	out, code := selectedHistory(t, "1", targets...)
	if code != 0 {
		t.Fatalf("exit=%d output=%q", code, out)
	}
	rows := historyRows(t, out, "histories")
	if len(rows) != len(targets) {
		t.Fatalf("summaries=%#v", rows)
	}
	for i, row := range rows {
		t.Run(targets[i], func(t *testing.T) {
			// The shell function prints each word of the recovery command NUL-framed, so
			// the test reads the argv that a shell gives the command.
			detail := row["detail"].(string)
			raw, err := exec.Command("sh", "-c", "bench() { printf '%s\\000' \"$@\"; }; "+detail).Output()
			if err != nil {
				t.Fatalf("detail %q: %v", detail, err)
			}
			argv := strings.Split(strings.TrimSuffix(string(raw), "\x00"), "\x00")
			if len(argv) < 2 || argv[0] != "spec" || argv[1] != "history" {
				t.Fatalf("invalid recovery route %q: argv=%q", detail, argv)
			}
			full, exit, _ := Command(argv[1:])
			if exit != 0 {
				t.Fatalf("recovery %q exit=%d output=%q", detail, exit, full)
			}
			count, newest := 1, f.hash[targets[i]]
			switch targets[i] {
			case "mixed":
				count, newest = 3, f.hash["newest"]
			case ".md":
				count = 0
			}
			recovered := historyRows(t, full, "history")
			if len(recovered) != count || fmt.Sprint(row["total_events"]) != strconv.Itoa(count) {
				t.Fatalf("recovered %s through %q count=%d summary=%v want=%d", targets[i], detail, len(recovered), row["total_events"], count)
			}
			if count > 0 && recovered[0]["hash"] != newest {
				t.Fatalf("recovered another history for %s: %#v", targets[i], recovered)
			}
		})
	}
	if _, err := os.Stat(filepath.Join(f.root, "sentinel")); !os.IsNotExist(err) {
		t.Fatalf("shell-shaped target executed: %v", err)
	}
}

func TestSelectedHistoriesExcludeOldOutput(t *testing.T) {
	f := historyFixture(t)
	out, code := selectedHistory(t, "1", "mixed")
	if code != 0 {
		t.Fatalf("exit=%d output=%q", code, out)
	}
	document, err := axitest.DecodeDocument(out)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(document.Blocks, []string{"histories", "history"}) || strings.Count(out, "\nhistory[") != 1 {
		t.Fatalf("unexpected output bodies: %s", out)
	}
	rows := historyRows(t, out, "history")
	if len(rows) != 1 || rows[0]["hash"] != f.hash["newest"] {
		t.Fatalf("unexpected selected body=%#v", rows)
	}
	for _, omitted := range []string{"mixed-delete", "both", "prefix", "delete", "retire"} {
		if strings.Contains(out, f.hash[omitted]) {
			t.Fatalf("unrequested %s body survived: %s", omitted, out)
		}
	}
}

func TestSelectedHistoryPreservesProducer(t *testing.T) {
	f := historyFixture(t)
	var want []HistoryFact
	for _, row := range f.capturedRows(t, "mixed") {
		want = append(want, HistoryFact{Slug: "mixed", Hash: row[0], Date: row[1], Kind: row[2], Subject: row[3]})
	}
	before, err := History("mixed")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before, want) {
		t.Fatalf("complete producer=%#v want baseline=%#v", before, want)
	}
	out, code := selectedHistory(t, "1", "mixed", "absent")
	if code != 0 {
		t.Fatalf("selected exit=%d output=%q", code, out)
	}
	after, err := History("mixed")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(after, want) {
		t.Fatalf("producer after selection=%#v want complete=%#v", after, want)
	}
}

// byteFixture commits three deletions of one flat spec. Each positional row is 2 indent
// bytes, an 8-byte hash, a 10-byte date, the 6-byte kind `delete`, a 1-byte subject,
// 3 commas, and a newline: 31 bytes. The header `history[3]{hash,date,kind,subject}:`
// and its newline are 36 bytes. The commit days fix the hashes. With the gittest identity,
// these days give no delete hash that TOON quotes, such as one with a leading zero digit.
func byteFixture(t *testing.T) historyFixtureData {
	t.Helper()
	f := emptyHistoryRepo(t)
	runGit(t, f.root, "config", "core.abbrev", "8")
	flat := filepath.Join(f.root, "specs", "bytes.md")
	for i, subject := range []string{"a", "b", "c"} {
		// Each removal takes the emptied specs directory with it.
		if err := os.MkdirAll(filepath.Dir(flat), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(flat, []byte("spec\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		f.commit(t, "add-"+subject, 2*i+11, "add bytes spec")
		runGit(t, f.root, "rm", "-q", "specs/bytes.md")
		f.commit(t, subject, 2*i+12, subject)
	}
	return f
}

func TestSelectedHistoryTrueBytes(t *testing.T) {
	const header, row, events = 36, 31, 3
	byteFixture(t)
	positional := positionalHistory(t, "bytes")
	lines := strings.SplitAfter(strings.TrimSuffix(positional, "\n"), "\n")
	if len(lines) != 1+events || len(lines[0]) != header {
		t.Fatalf("positional table=%q; want a %d-byte header and %d rows", positional, header, events)
	}
	for _, line := range lines[1:] {
		if len(strings.TrimSuffix(line, "\n"))+1 != row || !strings.Contains(line, ",delete,") {
			t.Fatalf("row %q; want a %d-byte delete row", line, row)
		}
	}
	if len(positional) != header+events*row {
		t.Fatalf("positional bytes=%d; want %d + %d × %d = %d", len(positional), header, events, row, header+events*row)
	}
	for _, limit := range []string{"1", "3", "4"} {
		out, code := selectedHistory(t, limit, "bytes")
		if code != 0 {
			t.Fatalf("exit=%d output=%q", code, out)
		}
		if got := fmt.Sprint(historyRows(t, out, "histories")[0]["total_bytes"]); got != strconv.Itoa(len(positional)) {
			t.Fatalf("limit=%s bytes=%s; want the complete positional count %d", limit, got, len(positional))
		}
	}
}

func TestSelectedHistoryTrueBytesCountsUTF8(t *testing.T) {
	historyFixture(t)
	for _, limit := range []string{"1", "3"} {
		out, code := selectedHistory(t, limit, "mixed", "absent")
		if code != 0 {
			t.Fatalf("exit=%d output=%q", code, out)
		}
		for i, operand := range []string{"mixed", "absent"} {
			want := len(positionalHistory(t, operand))
			if got := fmt.Sprint(historyRows(t, out, "histories")[i]["total_bytes"]); got != strconv.Itoa(want) {
				t.Fatalf("limit=%s %s bytes=%s; want complete UTF-8 bytes %d", limit, operand, got, want)
			}
		}
	}
}

func TestSelectedHistoryHostileSubject(t *testing.T) {
	f := historyFixture(t)
	f.commit(t, "refused", 9, "bad\x1b subject spec-retire: tainted")
	f.commit(t, "visible", 10, "visible clean subject spec-retire: tainted")
	f.commit(t, "permitted", 11, "tab\treturn\rcafé spec-retire: representable")
	f.commit(t, "folded", 12, "line\ncontinued spec-retire: representable-newline")
	out, code := selectedHistory(t, "1", "tainted", "mixed", "representable", "representable-newline", "absent")
	if code != 1 {
		t.Fatalf("hidden refused subject exit=%d output=%q", code, out)
	}
	rows := historyRows(t, out, "histories")
	if len(rows) != 5 || rows[0]["error"] != selectedHistoryUnrepresentable || rows[0]["detail"] != historyDetail("tainted", SlugOf("tainted")) || rows[0]["total_events"] != nil || rows[0]["total_bytes"] != nil || rows[0]["omitted_events"] != nil {
		t.Fatalf("per-target refused history=%#v", rows)
	}
	for _, row := range rows[1:] {
		if row["error"] != "" {
			t.Fatalf("lost representable history=%#v", row)
		}
	}
	events := historyRows(t, out, "history")
	if len(events) != 3 || events[0]["slug"] != "mixed" ||
		events[1]["subject"] != "tab\treturn\rcafé spec-retire: representable" ||
		events[2]["subject"] != "line continued spec-retire: representable-newline" {
		t.Fatalf("representable producer subjects=%#v", events)
	}
}
