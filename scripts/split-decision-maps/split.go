package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gibbonmi/bench/internal/maps"
)

var rawTicketHeading = regexp.MustCompile(`^## #([1-9][0-9]*):\s+(.+?)\s*$`)

var terminalHeadings = map[string]bool{
	"## Not yet specified":      true,
	"## Spec-writer discretion": true,
	"## Out of scope":           true,
	"## Sources":                true,
}

// rawTicket is one inline ticket kept as the bytes the file holds. The migration
// must not reflow decided prose, and the maps parser normalizes whitespace and
// drops fenced lines, so its DecisionTicket cannot be re-rendered faithfully.
type rawTicket struct {
	id       string
	title    string
	body     []string
	resolved bool
	gist     string
}

// splitDoc is one inline map cut into the three regions the split shape needs.
type splitDoc struct {
	preamble []string
	terminal []string
	tickets  []rawTicket
}

// resolvedAnswer mirrors ticketAnswerState in internal/maps/validation.go, which is
// unexported. The contract ticket deletes this program, so the mirror is temporary;
// a change to the answer-state rule before then must change both.
func resolvedAnswer(answer string) bool {
	answer = strings.TrimSpace(answer)
	switch {
	case strings.HasPrefix(answer, "— (deferred"), strings.Contains(answer, "GRILL DEFERRED"):
		return false
	case answer == "", strings.HasPrefix(answer, "— (open"):
		return false
	default:
		return true
	}
}

// cutDoc slices content into preamble, tickets, and the terminal block. The maps
// parser is the authority on what the tickets are: cutDoc cross-checks its own cut
// against ParseDecisionMap and fails rather than write a silently different map.
func cutDoc(path string, content []byte) (splitDoc, error) {
	parsed, diagnostics := maps.ParseDecisionMap(content)
	if len(diagnostics) > 0 {
		return splitDoc{}, fmt.Errorf("%s: %s", path, diagnostics[0].Message)
	}
	lines := strings.Split(strings.TrimSuffix(string(content), "\n"), "\n")
	var doc splitDoc
	fenced, region := false, "preamble"
	var current *rawTicket
	flush := func() {
		if current != nil {
			doc.tickets = append(doc.tickets, *current)
			current = nil
		}
	}
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			fenced = !fenced
		}
		if !fenced {
			if match := rawTicketHeading.FindStringSubmatch(line); match != nil && region != "terminal" {
				flush()
				region = "ticket"
				current = &rawTicket{id: match[1], title: match[2]}
				continue
			}
			if region != "preamble" && terminalHeadings[line] {
				flush()
				region = "terminal"
			}
		}
		switch region {
		case "preamble":
			doc.preamble = append(doc.preamble, line)
		case "ticket":
			current.body = append(current.body, line)
		default:
			doc.terminal = append(doc.terminal, line)
		}
	}
	flush()
	if err := reconcile(path, doc, parsed); err != nil {
		return splitDoc{}, err
	}
	for i := range doc.tickets {
		answer := parsed.Tickets[i].Answer
		doc.tickets[i].resolved = resolvedAnswer(answer)
		doc.tickets[i].gist = fitGist(doc.tickets[i].title, firstSentence(doc.tickets[i].body))
	}
	return doc, nil
}

func reconcile(path string, doc splitDoc, parsed maps.DecisionMap) error {
	if len(doc.tickets) != len(parsed.Tickets) {
		return fmt.Errorf("%s: cut found %d tickets, the parser found %d", path, len(doc.tickets), len(parsed.Tickets))
	}
	for i, ticket := range doc.tickets {
		if ticket.id != parsed.Tickets[i].ID || ticket.title != parsed.Tickets[i].Title {
			return fmt.Errorf("%s: cut ticket #%s %q does not match parsed #%s %q",
				path, ticket.id, ticket.title, parsed.Tickets[i].ID, parsed.Tickets[i].Title)
		}
	}
	if len(doc.terminal) == 0 {
		return fmt.Errorf("%s: no terminal section found", path)
	}
	return nil
}

// maxGistWords is the prose lane's sentence bound. A gist is one graded sentence, so
// the seed is cut back to a clause that fits rather than left for a hand edit.
const maxGistWords = 25

// clauseMarks are the cut points a gist may fall back to, longest keep first.
var clauseMarks = []string{". ", "; ", " — ", ": ", ", ", " and ", " so "}

// danglingWords are the trailing conjunctions a clause cut can leave behind.
var danglingWords = []string{" and", " but", " so", " or", " then"}

var codeSpan = regexp.MustCompile("`[^`]*`")

// firstSentence seeds a gist from the answer's opening paragraph. It reads the raw
// body rather than the parsed Answer, because the parser joins the paragraph's lines
// with newlines and a gist is one line.
func firstSentence(body []string) string {
	var paragraph []string
	seen, started := false, false
	for _, line := range body {
		trimmed := strings.TrimSpace(line)
		if !seen {
			seen = trimmed == "### Answer"
			continue
		}
		if trimmed == "" {
			if started {
				break
			}
			continue
		}
		started = true
		paragraph = append(paragraph, trimmed)
	}
	sentence := strings.TrimPrefix(strings.Join(paragraph, " "), "- ")
	if cut := strings.Index(sentence, ". "); cut >= 0 {
		sentence = sentence[:cut+1]
	}
	return strings.TrimSpace(sentence)
}

// fitGist trims the seeded sentence to the longest clause the word bound admits. The
// title is not trimmed, because a gist a reader cannot match to its ticket is worse
// than a short one.
func fitGist(title, sentence string) string {
	budget := maxGistWords - countWords(title)
	best := ""
	for _, candidate := range clauses(sentence) {
		if countWords(candidate) <= budget && len(candidate) > len(best) {
			best = candidate
		}
	}
	if best == "" {
		best = shortest(clauses(sentence))
	}
	return best
}

// clauses returns the sentence and every prefix that ends at a clause mark, each
// closed with a period.
func clauses(sentence string) []string {
	out := []string{sentence}
	for _, mark := range clauseMarks {
		for cut := 0; ; {
			at := strings.Index(sentence[cut:], mark)
			if at < 0 {
				break
			}
			cut += at + len(mark)
			prefix := strings.TrimRight(strings.TrimSpace(sentence[:cut]), ".,;:— ")
			for _, dangling := range danglingWords {
				prefix = strings.TrimSuffix(prefix, dangling)
			}
			prefix = strings.TrimRight(prefix, ".,;:— ")
			if prefix != "" {
				out = append(out, prefix+".")
			}
		}
	}
	return out
}

func shortest(candidates []string) string {
	best := ""
	for _, candidate := range candidates {
		if best == "" || len(candidate) < len(best) {
			best = candidate
		}
	}
	return best
}

// countWords approximates the prose lane's count: a code span is one word, and a token
// with no letter or digit is none.
func countWords(text string) int {
	count := 0
	for _, token := range strings.Fields(codeSpan.ReplaceAllString(text, "span")) {
		if strings.ContainsFunc(token, func(r rune) bool {
			return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
		}) {
			count++
		}
	}
	return count
}

func renderIndex(topic string, doc splitDoc) string {
	var b strings.Builder
	b.WriteString(strings.Join(trimTrailing(doc.preamble), "\n"))
	b.WriteString("\n\n## Notes\n\n## Decisions so far\n")
	for _, ticket := range doc.tickets {
		if !ticket.resolved {
			continue
		}
		fmt.Fprintf(&b, "\n- [%s](%s/tickets/%s.md): %s", ticket.title, topic, ticket.id, ticket.gist)
	}
	b.WriteString("\n\n")
	b.WriteString(strings.Join(trimTrailing(doc.terminal), "\n"))
	b.WriteString("\n")
	return b.String()
}

func renderTicket(ticket rawTicket) string {
	return "# " + ticket.title + "\n" + strings.Join(trimTrailing(ticket.body), "\n") + "\n"
}

func trimTrailing(lines []string) []string {
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}
