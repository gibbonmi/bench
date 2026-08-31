// Tests for the ticket-file grammar: required fields, field mechanics, the
// blocker graph, the citation grammar, and the hostile inputs.
package tickets

import (
	"reflect"
	"strings"
	"testing"
)

// wellFormed renders a ticket that yields no diagnostic, so each case states only its own fault.
func wellFormed() string {
	return strings.Join([]string{
		"# Build the ticket-file parser",
		"",
		"Blocked by: none",
		"Writes: internal/tickets/tickets.go (new)",
		"Covers: TG1, TG2",
		"",
		"## What to build",
		"",
		"One parser owns the schema.",
		"",
		"## Acceptance",
		"",
		"- [ ] TG1 — the parser names each absent required field.",
		"",
		"## Delegate charge",
		"",
		"Build it.",
		"",
	}, "\n")
}

func TestParseTicketRequiredFields(t *testing.T) {
	cases := []struct {
		name     string
		document string
		want     []string
	}{
		{
			name:     "well formed ticket is clean",
			document: wellFormed(),
		},
		{
			name:     "every required field absent",
			document: "Nothing here.\n",
			want: []string{
				"missing title",
				"missing Blocked by",
				"missing Writes",
				"missing Covers",
				"missing What to build section",
				"missing Acceptance section",
			},
		},
		{
			name:     "one absent field is named alone",
			document: strings.Replace(wellFormed(), "Covers: TG1, TG2\n", "", 1),
			want:     []string{"missing Covers"},
		},
		{
			name:     "an absent heading is named",
			document: strings.Replace(wellFormed(), "## Acceptance", "## Acceptance criteria", 1),
			want:     []string{"missing Acceptance section"},
		},
		{
			name:     "an empty field value is missing",
			document: strings.Replace(wellFormed(), "Writes: internal/tickets/tickets.go (new)", "Writes:", 1),
			want:     []string{"missing Writes"},
		},
		{
			name:     "the delegate charge heading stays optional",
			document: strings.Replace(wellFormed(), "## Delegate charge", "", 1),
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			_, got := ParseTicket("parser.md", []byte(testCase.document), []string{"parser.md"}, "TG")
			assertDiagnostics(t, got, testCase.want)
		})
	}
}

func TestParseTicketFieldMechanics(t *testing.T) {
	fenced := strings.Replace(wellFormed(), "One parser owns the schema.",
		"```\nBlocked by: quoted.md\nCovers: XX9\n```", 1)
	cases := []struct {
		name     string
		document string
		want     []string
		check    func(t *testing.T, ticket Ticket)
	}{
		{
			name:     "a second Blocked by line is a duplicate field",
			document: strings.Replace(wellFormed(), "Blocked by: none\n", "Blocked by: none\nBlocked by: other.md\n", 1),
			want:     []string{"duplicate Blocked by"},
			check: func(t *testing.T, ticket Ticket) {
				if ticket.Blockers != nil {
					t.Fatalf("Blockers = %#v, want the first value", ticket.Blockers)
				}
			},
		},
		{
			name:     "a second heading is a duplicate section",
			document: wellFormed() + "\n## Acceptance\n",
			want:     []string{"duplicate Acceptance section"},
		},
		{
			name:     "a fenced field line parses as no field",
			document: fenced,
			check: func(t *testing.T, ticket Ticket) {
				if ticket.Blockers != nil {
					t.Fatalf("Blockers = %#v, want none", ticket.Blockers)
				}
				if !reflect.DeepEqual(ticket.Covers, []string{"TG1", "TG2"}) {
					t.Fatalf("Covers = %#v, want the unfenced value", ticket.Covers)
				}
			},
		},
		{
			name:     "the parsed fields carry their values",
			document: wellFormed(),
			check: func(t *testing.T, ticket Ticket) {
				if ticket.Title != "Build the ticket-file parser" {
					t.Fatalf("Title = %q", ticket.Title)
				}
				if !reflect.DeepEqual(ticket.Writes, []string{"internal/tickets/tickets.go (new)"}) {
					t.Fatalf("Writes = %#v", ticket.Writes)
				}
			},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			ticket, got := ParseTicket("parser.md", []byte(testCase.document), []string{"parser.md", "other.md"}, "TG")
			assertDiagnostics(t, got, testCase.want)
			if testCase.check != nil {
				testCase.check(t, ticket)
			}
		})
	}
}

func TestParseTicketBlockerGraph(t *testing.T) {
	cases := []struct {
		name     string
		blockers string
		siblings []string
		want     []string
	}{
		{
			name:     "none parses to zero edges",
			blockers: "none",
			siblings: []string{"parser.md", "lift.md"},
		},
		{
			name:     "a sibling basename resolves",
			blockers: "lift.md",
			siblings: []string{"parser.md", "lift.md"},
		},
		{
			name:     "a dangling blocker names both basenames",
			blockers: "retitled.md",
			siblings: []string{"parser.md", "lift.md"},
			want:     []string{"parser.md: dangling blocker retitled.md"},
		},
		{
			name:     "a repeated blocker is a duplicate",
			blockers: "lift.md, lift.md",
			siblings: []string{"parser.md", "lift.md"},
			want:     []string{"parser.md: duplicate blocker lift.md"},
		},
		{
			name:     "a self blocker is named",
			blockers: "parser.md",
			siblings: []string{"parser.md", "lift.md"},
			want:     []string{"parser.md: self blocker parser.md"},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			document := strings.Replace(wellFormed(), "Blocked by: none", "Blocked by: "+testCase.blockers, 1)
			ticket, got := ParseTicket("parser.md", []byte(document), testCase.siblings, "TG")
			assertDiagnostics(t, got, testCase.want)
			if testCase.blockers == "none" && ticket.Blockers != nil {
				t.Fatalf("Blockers = %#v, want none", ticket.Blockers)
			}
		})
	}

	cycles := []struct {
		name   string
		parsed []Ticket
		want   []string
	}{
		{
			name: "an acyclic chain reports no cycle",
			parsed: []Ticket{
				{Name: "a.md", Blockers: []string{"b.md"}},
				{Name: "b.md", Blockers: []string{"c.md"}},
				{Name: "c.md"},
			},
		},
		{
			name: "a cycle names one edge",
			parsed: []Ticket{
				{Name: "a.md", Blockers: []string{"b.md"}},
				{Name: "b.md", Blockers: []string{"a.md"}},
			},
			want: []string{"cycle edge b.md -> a.md"},
		},
		{
			name: "a self edge is no cycle",
			parsed: []Ticket{
				{Name: "a.md", Blockers: []string{"a.md"}},
			},
		},
	}
	for _, testCase := range cycles {
		t.Run(testCase.name, func(t *testing.T) {
			assertDiagnostics(t, Cycles(testCase.parsed), testCase.want)
		})
	}
}

func TestParseTicketCoversGrammar(t *testing.T) {
	cases := []struct {
		name   string
		covers string
		tag    string
		want   []string
		parsed []string
	}{
		{
			name:   "none parses to zero citations",
			covers: "none",
			tag:    "TG",
		},
		{
			name:   "a row-ID list parses in order",
			covers: "TG3, TG1",
			tag:    "TG",
			parsed: []string{"TG3", "TG1"},
		},
		{
			name:   "a repeated token is a duplicate",
			covers: "TG1, TG1",
			tag:    "TG",
			want:   []string{"duplicate Covers token TG1"},
			parsed: []string{"TG1", "TG1"},
		},
		{
			name:   "a foreign tag is named",
			covers: "TG1, XX2",
			tag:    "TG",
			want:   []string{"foreign Covers token XX2: spec tag is TG"},
			parsed: []string{"TG1", "XX2"},
		},
		{
			name:   "a malformed token is named",
			covers: "TG1, TG-2",
			tag:    "TG",
			want:   []string{`malformed Covers token "TG-2"`},
			parsed: []string{"TG1", "TG-2"},
		},
		{
			name:   "an empty spec tag skips the tag rule",
			covers: "XX2",
			parsed: []string{"XX2"},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			document := strings.Replace(wellFormed(), "Covers: TG1, TG2", "Covers: "+testCase.covers, 1)
			ticket, got := ParseTicket("parser.md", []byte(document), []string{"parser.md"}, testCase.tag)
			assertDiagnostics(t, got, testCase.want)
			if !reflect.DeepEqual(ticket.Covers, testCase.parsed) && !(len(ticket.Covers) == 0 && len(testCase.parsed) == 0) {
				t.Fatalf("Covers = %#v, want %#v", ticket.Covers, testCase.parsed)
			}
		})
	}

	// TG20: a row ID in body prose lands in no parsed Covers set.
	prose := strings.Replace(wellFormed(), "One parser owns the schema.", "The parser also serves TG9 and TG8.", 1)
	ticket, diagnostics := ParseTicket("parser.md", []byte(prose), []string{"parser.md"}, "TG")
	assertDiagnostics(t, diagnostics, nil)
	if !reflect.DeepEqual(ticket.Covers, []string{"TG1", "TG2"}) {
		t.Fatalf("prose tokens reached Covers = %#v", ticket.Covers)
	}
}

func TestParseTicketHostileInput(t *testing.T) {
	const nbsp = "\u00a0"
	cases := []struct {
		name     string
		document string
		want     []string
		check    func(t *testing.T, ticket Ticket)
	}{
		{
			name:     "a non-breaking space in a field prefix reads as no field",
			document: strings.Replace(wellFormed(), "Blocked by: none", nbsp+"Blocked by: none", 1),
			want:     []string{"missing Blocked by"},
		},
		{
			name:     "a non-breaking space inside a field prefix reads as no field",
			document: strings.Replace(wellFormed(), "Blocked by: none", "Blocked"+nbsp+"by: none", 1),
			want:     []string{"missing Blocked by"},
		},
		{
			name:     "an unterminated fence grades nothing after the opening marker",
			document: "# Title\n\nBlocked by: none\n\n```\nWrites: swallowed.go\nCovers: TG1\n\n## What to build\n\n## Acceptance\n",
			want: []string{
				"unterminated fence",
				"missing Writes",
				"missing Covers",
				"missing What to build section",
				"missing Acceptance section",
			},
		},
		{
			name:     "a closed fence reports no unterminated fence",
			document: strings.Replace(wellFormed(), "One parser owns the schema.", "```\nquoted\n```", 1),
		},
		{
			name:     "a required field on the last line without a newline parses",
			document: strings.TrimSuffix(reorderedLastField(), "\n"),
			check: func(t *testing.T, ticket Ticket) {
				if !reflect.DeepEqual(ticket.Covers, []string{"TG1", "TG2"}) {
					t.Fatalf("Covers = %#v, want the last-line value", ticket.Covers)
				}
			},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			ticket, got := ParseTicket("parser.md", []byte(testCase.document), []string{"parser.md"}, "TG")
			assertDiagnostics(t, got, testCase.want)
			if testCase.check != nil {
				testCase.check(t, ticket)
			}
		})
	}
}

// reorderedLastField renders a well-formed ticket whose last line is a required field.
func reorderedLastField() string {
	return strings.Join([]string{
		"# Build the ticket-file parser",
		"",
		"Blocked by: none",
		"Writes: internal/tickets/tickets.go (new)",
		"",
		"## What to build",
		"",
		"One parser owns the schema.",
		"",
		"## Acceptance",
		"",
		"- [ ] TG1 — the parser names each absent required field.",
		"",
		"Covers: TG1, TG2",
		"",
	}, "\n")
}

func assertDiagnostics(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) == 0 && len(want) == 0 {
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("diagnostics = %#v, want %#v", got, want)
	}
}
