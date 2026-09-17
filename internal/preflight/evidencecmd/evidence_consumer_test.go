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
// then every page's membership in the declared manifest, then exact page coverage.
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

func integer(value any) string {
	number, _ := value.(float64)
	return strconv.Itoa(int(number))
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
	identity, _, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug, false))
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
	identity, _, _ := prepareEvidence(t, preflighttest.LegacyCommitted(t, root, slug, "paged delivery", false))
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
	for _, test := range []struct {
		name, want string
		pages      []map[string]any
	}{
		{"missing page", "never arrived", received.pages[:len(received.pages)-1]},
		{"duplicate page", "arrived twice", append(append([]map[string]any{}, received.pages...), received.pages[len(received.pages)-1])},
		{"final page only", "never arrived", received.pages[len(received.pages)-1:]},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := consume(identity, delivery{manifest: received.manifest, pages: test.pages}); err == nil ||
				!strings.Contains(err.Error(), test.want) {
				t.Fatalf("%s = %v, want %q", test.name, err, test.want)
			}
		})
	}
}
