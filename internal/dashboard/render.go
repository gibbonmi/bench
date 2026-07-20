// render.go is the pure half of the dashboard: the sanitized view projection, the
// HTML renderer, and the one self-contained page template. No IO, no clock, no git —
// Command (dashboard.go) gathers the Snapshot; everything here is a function of it.
package dashboard

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/sanitize"
)

// view is the sanitized, render-ready projection of a Snapshot. Every git- and file-sourced
// string is run through sanitize.Controls before it lands here, so the template only ever
// escapes control-byte-free text — control-rune escaping is the one step html/template does
// not do. RoadmapText and Sequence are the exception: they render inside <pre>, where
// html/template's markup neutralization is enough on its own, so they route through
// sanitize.Preformatted instead, which keeps newline and tab literal so the panel's layout
// survives while still escaping every other control rune.
type view struct {
	GeneratedAt    string
	HasGate        bool
	Gate           gateView
	Signals        []signalView
	RoadmapPresent bool
	RoadmapText    string
	Sequence       string
	Ideas          []string
	OpenLearnings  int
	Worktrees      []worktreeView
	WorktreesErr   string
}

type gateView struct {
	State, PendingStatus, Status, CachedTree, WorkTree, Timestamp, Age string
	Stale                                                              bool
}

type signalView struct{ Name, Detail, Action string }

type worktreeView struct{ Class, Path string }

// Render turns a Snapshot into the complete self-contained HTML document. It is pure: it
// reads nothing but its argument. Escaping is contextual (html/template neutralizes markup
// and quote injection in every interpolated field) plus a control-rune pass the template
// cannot do itself — sanitize.Controls for every field, sanitize.Preformatted for the two
// <pre>-rendered fields so their layout survives. A template-execution error is a
// template-source bug, unreachable from repo data, so the seam stays a total function.
func Render(s Snapshot) string {
	v := view{
		GeneratedAt:    s.GeneratedAt.Format(time.RFC3339),
		HasGate:        s.Gate.Present,
		RoadmapPresent: s.RoadmapPresent,
		RoadmapText:    sanitize.Preformatted(s.RoadmapText),
		Sequence:       sanitize.Preformatted(s.Sequence),
		OpenLearnings:  s.OpenLearnings,
		WorktreesErr:   sanitize.Controls(s.WorktreesErr),
	}
	if s.Gate.Present {
		v.Gate = gateView{
			State:         sanitize.Controls(s.Gate.State),
			PendingStatus: sanitize.Controls(s.Gate.PendingStatus),
			Status:        sanitize.Controls(s.Gate.Status),
			CachedTree:    sanitize.Controls(s.Gate.CachedTree),
			WorkTree:      sanitize.Controls(s.Gate.WorkTree),
			Timestamp:     sanitize.Controls(s.Gate.Timestamp),
			Age:           gateAge(s.GeneratedAt, s.Gate.Timestamp),
			Stale:         s.Gate.Stale,
		}
	}
	for _, sig := range s.Signals {
		v.Signals = append(v.Signals, signalView{
			Name:   sanitize.Controls(sig.Name),
			Detail: sanitize.Controls(sig.Detail),
			Action: sanitize.Controls(sig.Action),
		})
	}
	for _, idea := range s.Ideas {
		v.Ideas = append(v.Ideas, sanitize.Controls(idea))
	}
	for _, wt := range s.Worktrees {
		v.Worktrees = append(v.Worktrees, worktreeView{
			Class: sanitize.Controls(string(wt.Class)),
			Path:  sanitize.Controls(wt.Path),
		})
	}
	var b strings.Builder
	if err := pageTemplate.Execute(&b, v); err != nil {
		// Unreachable from repo data: the template is a compiled-in constant, so an error
		// here is a source bug, not an input condition. Return whatever was written.
		return b.String()
	}
	return b.String()
}

// gateAge renders how long before generation the gate ran, from the cache timestamp. It is
// pure — the "now" is the injected generation time, never the wall clock. An absent or
// unparseable timestamp, or a future one, yields the empty string, which the template omits.
func gateAge(generatedAt time.Time, timestamp string) string {
	if timestamp == "" {
		return ""
	}
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return ""
	}
	d := generatedAt.Sub(t)
	switch {
	case d < 0:
		return ""
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}

// pageTemplate is the one self-contained document: inline <style> only (light + a
// prefers-color-scheme dark palette), neutral typography, no external request, no
// JavaScript. Every interpolated field is auto-escaped by html/template in its context.
var pageTemplate = template.Must(template.New("dashboard").Parse(pageHTML))

const pageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Bench dashboard</title>
<style>
:root {
  color-scheme: light dark;
  --bg: #f7f7f5;
  --card: #ffffff;
  --fg: #1a1a1a;
  --muted: #666666;
  --border: #dddddd;
  --accent: #2a5d84;
  --warn: #9a5b00;
  --bad: #a11;
}
@media (prefers-color-scheme: dark) {
  :root {
    --bg: #14161a;
    --card: #1c1f24;
    --fg: #e6e6e6;
    --muted: #9aa0a6;
    --border: #2c313a;
    --accent: #7db4dd;
    --warn: #d9a35b;
    --bad: #e06c75;
  }
}
* { box-sizing: border-box; }
body {
  margin: 0;
  padding: 2rem 1rem 4rem;
  background: var(--bg);
  color: var(--fg);
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  line-height: 1.5;
}
.wrap { max-width: 60rem; margin: 0 auto; }
header { margin-bottom: 1.5rem; }
h1 { font-size: 1.5rem; margin: 0 0 .25rem; }
h2 { font-size: 1.05rem; margin: 0 0 .75rem; }
h3 { font-size: .95rem; margin: 1rem 0 .5rem; color: var(--muted); }
.generated { color: var(--muted); font-size: .85rem; margin: 0; }
.card {
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: .5rem;
  padding: 1rem 1.25rem;
  margin-bottom: 1rem;
}
table { width: 100%; border-collapse: collapse; font-size: .9rem; }
th, td { text-align: left; padding: .35rem .5rem; border-bottom: 1px solid var(--border); vertical-align: top; }
th { color: var(--muted); font-weight: 600; }
td:first-child { white-space: nowrap; }
dl { display: grid; grid-template-columns: max-content 1fr; gap: .25rem 1rem; margin: 0; font-size: .9rem; }
dt { color: var(--muted); }
dd { margin: 0; word-break: break-all; }
pre {
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: .35rem;
  padding: .75rem;
  overflow-x: auto;
  font-size: .85rem;
  white-space: pre-wrap;
  word-break: break-word;
}
ul { margin: 0; padding-left: 1.25rem; font-size: .9rem; }
.empty { color: var(--muted); font-style: italic; margin: 0; }
.error { color: var(--bad); margin: 0; }
.badge { display: inline-block; padding: .05rem .5rem; border-radius: .25rem; font-size: .8rem; font-weight: 600; }
.badge.green { color: var(--accent); }
.badge.red { color: var(--bad); }
.badge.stale { color: var(--warn); border: 1px solid var(--warn); margin-left: .5rem; }
.count { font-size: 1.75rem; font-weight: 700; margin: 0; }
</style>
</head>
<body>
<div class="wrap">
<header>
  <h1>Bench dashboard</h1>
  <p class="generated">Generated {{.GeneratedAt}}</p>
</header>

<section class="card">
  <h2>Gate</h2>
  {{if .HasGate}}
  <p>
    <span class="badge {{if .Gate.Stale}}stale{{else}}{{.Gate.Status}}{{end}}">{{if .Gate.PendingStatus}}{{.Gate.PendingStatus}}{{else if .Gate.Status}}{{.Gate.Status}}{{else}}{{.Gate.State}}{{end}}</span>
    {{if .Gate.Stale}}<span class="badge stale">stale</span>{{end}}
  </p>
  <dl>
    <dt>cached tree</dt><dd>{{.Gate.CachedTree}}</dd>
    <dt>work tree</dt><dd>{{.Gate.WorkTree}}</dd>
    {{if .Gate.Timestamp}}<dt>gated at</dt><dd>{{.Gate.Timestamp}}{{if .Gate.Age}} ({{.Gate.Age}}){{end}}</dd>{{end}}
  </dl>
  {{else}}
  <p class="empty">No gate cache yet — run bench gate.</p>
  {{end}}
</section>

<section class="card">
  <h2>Signals</h2>
  {{if .Signals}}
  <table>
    <thead><tr><th>signal</th><th>detail</th><th>next action</th></tr></thead>
    <tbody>
    {{range .Signals}}<tr><td>{{.Name}}</td><td>{{.Detail}}</td><td>{{.Action}}</td></tr>
    {{end}}</tbody>
  </table>
  {{else}}
  <p class="empty">No signals — nothing pending.</p>
  {{end}}
</section>

<section class="card">
  <h2>Roadmap</h2>
  {{if .RoadmapPresent}}
  <pre>{{.RoadmapText}}</pre>
  {{if .Sequence}}
  <h3>Recommended sequence</h3>
  <pre>{{.Sequence}}</pre>
  {{end}}
  {{else}}
  <p class="empty">No ROADMAP.md — run /bench-what-next to create the working roadmap.</p>
  {{end}}
</section>

<section class="card">
  <h2>Ideas</h2>
  {{if .Ideas}}
  <ul>{{range .Ideas}}<li>{{.}}</li>{{end}}</ul>
  {{else}}
  <p class="empty">No parked ideas.</p>
  {{end}}
</section>

<section class="card">
  <h2>Open learnings</h2>
  <p class="count">{{.OpenLearnings}}</p>
</section>

<section class="card">
  <h2>Worktrees</h2>
  {{if .WorktreesErr}}
  <p class="error">git worktree list failed: {{.WorktreesErr}}</p>
  {{else if .Worktrees}}
  <table>
    <thead><tr><th>class</th><th>path</th></tr></thead>
    <tbody>
    {{range .Worktrees}}<tr><td>{{.Class}}</td><td>{{.Path}}</td></tr>
    {{end}}</tbody>
  </table>
  {{else}}
  <p class="empty">No out-of-pool, leased, or warm worktrees.</p>
  {{end}}
</section>
</div>
</body>
</html>
`
