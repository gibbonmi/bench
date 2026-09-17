package evidencecmd_test

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	toonlib "github.com/toon-format/toon-go"
)

// This file is the independent consumer. It starts from a trusted expected identity and
// reconstructs the artifact from returned responses alone: its own join, its own SHA-256,
// and the upstream TOON decoder. It calls no production decoder and reads no writer
// constant, so a production reader that trusts the artifact's own claims cannot satisfy it.

// delivery is what one consumer actually received: the joined manifest and the pages it saw.
type delivery struct {
	manifest string
	pages    []map[string]any
}

// consume checks one delivery against the trusted identity. It verifies the manifest digest,
// then every page's membership in the declared manifest, then exact page coverage, then the
// source bodies the received pages reconstruct.
func consume(trusted string, received delivery) error {
	if "sha256:"+sha(received.manifest) != trusted {
		return errors.New("reconstructed manifest does not match the trusted identity")
	}
	document, err := decodeIndependently(received.manifest)
	if err != nil {
		return err
	}
	declared := map[string]map[string]any{}
	for _, page := range document["pages"] {
		declared[page["source"].(string)+"#"+integer(page["index"])] = page
	}
	sources := map[string]bool{}
	for _, source := range document["sources"] {
		sources[source["id"].(string)] = true
	}
	seen := map[string]bool{}
	for _, page := range received.pages {
		if page["stream"] != "source" {
			continue
		}
		key := page["source"].(string) + "#" + integer(page["index"])
		switch {
		case !sources[page["source"].(string)]:
			return fmt.Errorf("page %s names a source the manifest does not declare", key)
		case declared[key] == nil:
			return fmt.Errorf("page %s is outside the declared membership", key)
		case declared[key]["sha256"] != page["sha256"] || sha(page["content"].(string)) != page["sha256"]:
			return fmt.Errorf("page %s does not match its declared digest", key)
		case seen[key]:
			return fmt.Errorf("page %s arrived twice", key)
		}
		seen[key] = true
	}
	for key := range declared {
		if !seen[key] {
			return fmt.Errorf("page %s never arrived", key)
		}
	}
	return reconstruct(document, received.pages)
}

// reconstruct rebuilds every declared source from the pages the consumer received, in the
// order they arrived. A page must open exactly where its source body stands and must hold
// exactly the length it declares, so an out-of-order or an overlapping page set fails here.
// The joined body then matches the length and the digest the manifest declares.
func reconstruct(document map[string][]map[string]any, pages []map[string]any) error {
	bodies := map[string]string{}
	for _, page := range pages {
		if page["stream"] != "source" {
			continue
		}
		id, content := page["source"].(string), page["content"].(string)
		key := id + "#" + integer(page["index"])
		switch {
		case number(page["offset"]) != len(bodies[id]):
			return fmt.Errorf("page %s opens at %d, not where source %s stands", key, number(page["offset"]), id)
		case number(page["bytes"]) != len(content):
			return fmt.Errorf("page %s declares %d bytes and holds %d", key, number(page["bytes"]), len(content))
		}
		bodies[id] += content
	}
	for _, source := range document["sources"] {
		id := source["id"].(string)
		switch {
		case number(source["bytes"]) != len(bodies[id]):
			return fmt.Errorf("source %s reconstructs %d of %d declared bytes",
				id, len(bodies[id]), number(source["bytes"]))
		case sha(bodies[id]) != source["sha256"]:
			return fmt.Errorf("source %s does not reconstruct its declared digest", id)
		}
	}
	return nil
}

// decodeIndependently decodes the manifest tables this consumer needs, through the upstream
// TOON decoder alone.
func decodeIndependently(manifest string) (map[string][]map[string]any, error) {
	parsed, err := toonlib.DecodeString(manifest)
	if err != nil {
		return nil, err
	}
	root, isObject := parsed.(map[string]any)
	if !isObject {
		return nil, errors.New("manifest is not a TOON object")
	}
	document := map[string][]map[string]any{}
	for _, table := range []string{"sources", "pages"} {
		rows, present := root[table].([]any)
		if !present {
			return nil, fmt.Errorf("manifest holds no %s table", table)
		}
		for _, value := range rows {
			row, isRow := value.(map[string]any)
			if !isRow {
				return nil, fmt.Errorf("%s row is not an object", table)
			}
			document[table] = append(document[table], row)
		}
	}
	return document, nil
}

func integer(value any) string { return strconv.Itoa(number(value)) }

func number(value any) int {
	count, _ := value.(float64)
	return int(count)
}

// sourceIndexes lists the positions at which the pages of source id arrived. It fails the test
// when fewer than count pages arrived, so a caller never indexes past the delivery.
func sourceIndexes(t *testing.T, pages []map[string]any, id string, count int) []int {
	t.Helper()
	var found []int
	for i, page := range pages {
		if page["stream"] == "source" && page["source"] == id {
			found = append(found, i)
		}
	}
	if len(found) < count {
		t.Fatalf("source %s delivered %d pages, want at least %d", id, len(found), count)
	}
	return found
}

// restated copies the page at position at, so a caller changes one row and leaves the delivery
// the consumer received untouched.
func restated(pages []map[string]any, at int) map[string]any {
	copied := map[string]any{}
	for name, value := range pages[at] {
		copied[name] = value
	}
	return copied
}

// reordered exchanges the first two pages of source id, so that source arrives out of order.
func reordered(t *testing.T, pages []map[string]any, id string) []map[string]any {
	t.Helper()
	changed := append([]map[string]any{}, pages...)
	at := sourceIndexes(t, changed, id, 2)
	changed[at[0]], changed[at[1]] = changed[at[1]], changed[at[0]]
	return changed
}

// overlapped restates the second page of source id at the first page's offset, so the two
// declared ranges overlap.
func overlapped(t *testing.T, pages []map[string]any, id string) []map[string]any {
	t.Helper()
	changed := append([]map[string]any{}, pages...)
	at := sourceIndexes(t, changed, id, 2)
	moved := restated(changed, at[1])
	moved["offset"] = changed[at[0]]["offset"]
	changed[at[1]] = moved
	return changed
}

// shortened understates the first page of source id by one byte, so that page declares a length
// the content it carries does not hold. The page keeps its content and its digest.
func shortened(t *testing.T, pages []map[string]any, id string) []map[string]any {
	t.Helper()
	changed := append([]map[string]any{}, pages...)
	at := sourceIndexes(t, changed, id, 1)
	trimmed := restated(changed, at[0])
	trimmed["bytes"] = float64(number(trimmed["bytes"]) - 1)
	changed[at[0]] = trimmed
	return changed
}

// receive traverses the default stream and returns what the consumer received.
func receive(t *testing.T, identity string) delivery {
	t.Helper()
	received := delivery{}
	for _, page := range traverseEvidence(t, identity) {
		if page.row["stream"] == "manifest" {
			received.manifest += page.row["content"].(string)
		}
		received.pages = append(received.pages, page.row)
	}
	return received
}

// TestEvidenceIndependentReader is CE7 and CE8.
func TestEvidenceIndependentReader(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	identity, _, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug))
	received := receive(t, identity)
	if err := consume(identity, received); err != nil {
		t.Fatalf("complete delivery = %v, want acceptance", err)
	}
	t.Run("CE7 manifest hash differs from the trusted identity", func(t *testing.T) {
		if err := consume("sha256:"+strings.Repeat("0", 64), received); err == nil ||
			!strings.Contains(err.Error(), "trusted identity") {
			t.Fatalf("foreign identity = %v, want a manifest digest rejection", err)
		}
	})
	t.Run("CE8 page outside declared membership", func(t *testing.T) {
		content := "outside the manifest\n"
		outside := map[string]any{"stream": "source", "source": "s2", "index": float64(99),
			"sha256": sha(content), "content": content}
		tampered := delivery{manifest: received.manifest, pages: append(append([]map[string]any{}, received.pages...), outside)}
		if err := consume(identity, tampered); err == nil || !strings.Contains(err.Error(), "outside the declared membership") {
			t.Fatalf("undeclared page = %v, want a membership rejection", err)
		}
	})
	t.Run("CE8 page of an undeclared source", func(t *testing.T) {
		content := "another source\n"
		foreign := map[string]any{"stream": "source", "source": "s99", "index": float64(0),
			"sha256": sha(content), "content": content}
		tampered := delivery{manifest: received.manifest, pages: append(append([]map[string]any{}, received.pages...), foreign)}
		if err := consume(identity, tampered); err == nil || !strings.Contains(err.Error(), "does not declare") {
			t.Fatalf("undeclared source page = %v, want a membership rejection", err)
		}
	})
}

// TestEvidenceDeliveryCoverage is CE106: the delivery check rejects missing, duplicate, and
// final-only page sets, so a terminal response alone proves no earlier delivery.
func TestEvidenceDeliveryCoverage(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	ticket := preflighttest.TicketDoc("One", "PF1", "PF2") + strings.Repeat("paged ticket\n", 1000)
	preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/one.md", ticket)
	identity, _, _ := prepareEvidence(t, preflighttest.LegacyCommitted(t, root, slug, "paged delivery"))
	received := receive(t, identity)
	if err := consume(identity, received); err != nil {
		t.Fatalf("complete delivery = %v, want acceptance", err)
	}
	sourcePages := 0
	for _, page := range received.pages {
		if page["stream"] == "source" {
			sourcePages++
		}
	}
	if sourcePages < 3 {
		t.Fatalf("delivery fixture held %d source pages, want at least 3", sourcePages)
	}
	short := shortened(t, received.pages, "s2")
	held := len(short[sourceIndexes(t, short, "s2", 1)[0]]["content"].(string))
	for _, test := range []struct {
		name, want string
		pages      []map[string]any
	}{
		{"missing page", "never arrived", received.pages[:len(received.pages)-1]},
		{"duplicate page", "arrived twice", append(append([]map[string]any{}, received.pages...), received.pages[len(received.pages)-1])},
		{"final page only", "never arrived", received.pages[len(received.pages)-1:]},
		{"out-of-order pages", "not where source s2 stands", reordered(t, received.pages, "s2")},
		{"overlapping pages", "not where source s2 stands", overlapped(t, received.pages, "s2")},
		{"page shorter than it declares", fmt.Sprintf("declares %d bytes and holds %d", held-1, held), short},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := consume(identity, delivery{manifest: received.manifest, pages: test.pages}); err == nil ||
				!strings.Contains(err.Error(), test.want) {
				t.Fatalf("%s = %v, want %q", test.name, err, test.want)
			}
		})
	}
	t.Run("source digest the pages do not produce", func(t *testing.T) {
		// The crafted manifest keeps every page digest and every declared range, so the
		// delivery fails only where the joined body meets the declared source digest.
		crafted := strings.Replace(received.manifest, sha(ticket), strings.Repeat("c", 64), 1)
		if crafted == received.manifest {
			t.Fatalf("manifest holds no digest of the %d-byte ticket", len(ticket))
		}
		err := consume("sha256:"+sha(crafted), delivery{manifest: crafted, pages: received.pages})
		if err == nil || !strings.Contains(err.Error(), "does not reconstruct its declared digest") {
			t.Fatalf("crafted source digest = %v, want a reconstruction rejection", err)
		}
	})
}
