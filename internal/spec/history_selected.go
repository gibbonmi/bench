package spec

import (
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// SelectedHistoryUsage is the selected history grammar, which parsing and public
// discovery share.
const SelectedHistoryUsage = historyCmd + " --spec <slug-or-path> [--spec <slug-or-path>]... --limit <positive-count>"

var selectedHistoryGrammar = usage.Grammar{
	Cmd:  historyCmd,
	Help: "usage: " + SelectedHistoryUsage,
	Flags: []usage.Flag{
		{Name: "--spec", HasValue: true, NoEmptyValue: true, Required: true, Repeatable: true},
		{Name: "--limit", HasValue: true, NoEmptyValue: true, Required: true},
	},
}

// selectedTargetControls is the error cell of a spec operand that no output line can
// carry. The row names the request ordinal in place of the operand.
const selectedTargetControls = "spec operand contains control characters"

// selectedHistoryUnrepresentable is the error cell of a history whose complete table the
// TOON owner refuses, such as one with a control byte in a commit subject.
const selectedHistoryUnrepresentable = "history is not representable"

// selectsHistories reports whether args choose the selected view. The check runs before
// specArg parses the positional operand, because that grammar refuses every flag.
func selectsHistories(args []string) bool {
	return usage.FlagPresent(selectedHistoryGrammar, args, "--spec") || usage.FlagPresent(selectedHistoryGrammar, args, "--limit")
}

// selectedHistories answers several spec histories in one table pair. The complete
// History producer runs for each distinct slug; the count, the byte size, and the
// recovery command describe that complete history, and at most limit newest events join
// the event table. Each failed target keeps its own summary row with unknown counts.
func selectedHistories(args []string) (string, int) {
	parsed, line, code := usage.Parse(selectedHistoryGrammar, args)
	if line != "" {
		return line + "\n", code
	}
	limit, err := strconv.Atoi(parsed.Flags["--limit"])
	if err != nil || limit <= 0 || parsed.EndedFlags {
		return selectedHistoryGrammar.Help + "\n", 2
	}
	if _, err := git.Root(); err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	var summaries, rows [][]any
	seenTargets, seenSlugs := map[string]bool{}, map[string]bool{}
	exit := 0
	for i, target := range parsed.Repeated["--spec"] {
		if seenTargets[target] {
			continue
		}
		seenTargets[target] = true
		if !sanitize.LineSafe(target) {
			summaries = append(summaries, []any{sanitize.TargetPointer(i + 1), "", nil, nil, nil, "", selectedTargetControls})
			exit = 1
			continue
		}
		slug := SlugOf(target)
		if seenSlugs[slug] {
			continue
		}
		seenSlugs[slug] = true
		detail := historyDetail(target, slug)
		events, err := History(slug)
		if err != nil {
			summaries = append(summaries, []any{target, slug, nil, nil, nil, detail, historyDerivationFailed})
			exit = 1
			continue
		}
		complete, err := renderHistory(events)
		if err != nil {
			summaries = append(summaries, []any{target, slug, nil, nil, nil, detail, selectedHistoryUnrepresentable})
			exit = 1
			continue
		}
		total := len(events)
		if total > limit {
			events = events[:limit]
		}
		summaries = append(summaries, []any{target, slug, total, len(complete), total - len(events), detail, ""})
		for _, event := range events {
			rows = append(rows, []any{event.Slug, event.Hash, event.Date, event.Kind, event.Subject})
		}
	}
	summary, err := toon.TableTyped("histories", []string{"target", "slug", "total_events", "total_bytes", "omitted_events", "detail", "error"}, summaries)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	history, err := toon.TableTyped("history", []string{"slug", "hash", "date", "kind", "subject"}, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return summary + history, exit
}

// historyDetail is the complete-history command of one selected spec. It names the slug
// when the positional path normalizes that slug to itself, and the original operand
// otherwise, as for `x.md.md` and `.md`. The end-of-options marker keeps a slug that
// reads as help or a flag a positional operand.
func historyDetail(target, slug string) string {
	operand := slug
	if SlugOf(slug) != slug {
		operand = target
	}
	command := historyCmd + " "
	if strings.HasPrefix(operand, "-") || operand == "help" {
		command += "-- "
	}
	return command + axi.ShellQuote(operand)
}
