package conformance

// Each reason applies to a whole class of dispatch names, so the file states the reason once.
const (
	whyPlumbing = "hook- and adapter-driven plumbing: its argv is produced by the kit, never typed by an agent, so there is no misuse for a grammar to report"
	whyNested   = "dispatches a subcommand tree: each leaf owns its grammar"
)

// subcommandRouting is the explicit registry the routing check grades cmd/bench/main.go
// against. It is deliberately exhaustive, not a list of only the interesting cases. A name
// reaching either dispatch surface with no row here turns red. This makes the check fail
// closed against the next subcommand someone adds.
var subcommandRouting = map[string]routingEntry{
	"anchors": routed("cmd/bench"),
	// The `clean` child is dispatched inside the module, whose one grammar declares both
	// the bare verb and the child. So the verb stays routed rather than nested-exempt.
	"cache":     routed("internal/gocache"),
	"commands":  routed("cmd/bench"),
	"consumers": routed("internal/consumers"),
	"commit":    routed("internal/commit"),
	"coverage":  routed("internal/coverage"),
	"dashboard": routed("internal/dashboard"),
	"diff":      routed("internal/diff"),
	"guards":    routed("internal/guards"),
	"handoff":   routed("internal/handoff"),
	"harnesses": routed("internal/harnesses"),
	"idea":      routed("internal/roadmap"),
	"learning":  routed("internal/roadmap"),
	"retro":     routed("internal/roadmap"),
	"learnings": routed("internal/learnings"),
	"maps":      routed("internal/maps"),
	"models":    routed("internal/models"),
	"outline":   routed("internal/outline"),
	"preflight": routed("internal/preflight"),
	// prep-release takes a flat argv with no subcommand tree. It is routed, not exempt
	// like the release commands beside it in the dispatch switch.
	"prep-release": routed("internal/preprelease"),
	"roadmap":      routed("internal/roadmap"),
	"skills-index": routed("internal/skillsindex"),
	"status":       routed("internal/status"),
	"structure":    routed("internal/structure"),
	"test":         routed("internal/testreport"),
	// probe owns its own grammar and hands the focused-run owner only the selection it
	// parsed, so it is routed rather than nested.
	"probe": routed("internal/probe"),

	"check-agent-line":      exempt(whyPlumbing),
	"freshness-check":       exempt(whyPlumbing),
	"freshness-publish":     exempt(whyPlumbing),
	"gate-go":               exempt(whyPlumbing),
	"gate-prose":            exempt(whyPlumbing),
	"gate-phases":           exempt(whyPlumbing),
	"guard-bench-follow-on": exempt(whyPlumbing),
	"guard-file-write":      exempt(whyPlumbing),
	"guard-git":             exempt(whyPlumbing),
	"resolve-model":         exempt(whyPlumbing),
	"stop-verdict":          exempt(whyPlumbing),
	"tree-hash":             exempt(whyPlumbing),
	"worktree-hook":         exempt(whyPlumbing),
	"worktree-lease-file":   exempt(whyPlumbing),
	"worktree-pool":         exempt(whyPlumbing),

	"canary":            exempt(whyNested),
	"doctor":            routed("internal/adopt"),
	"gate":              exempt(whyNested),
	"gate-run":          exempt(whyNested),
	"init":              exempt(whyNested),
	"link":              exempt(whyNested),
	"release":           exempt(whyNested),
	"release-preflight": exempt(whyNested),
	"resume-clean":      exempt(whyNested),
	"session-inspect":   exempt(whyNested),
	"setup":             exempt(whyNested),
	"shift":             exempt(whyNested),
	"spec":              exempt(whyNested),
	"unlink":            exempt(whyNested),
	"upgrade":           exempt(whyNested),
	"worktree":          exempt(whyNested),

	"version": exempt("takes no arguments: the dispatch case prints the build-time version line and returns"),
	"help":    exempt("takes no arguments: the command prints the top-level inventory and returns"),
}
