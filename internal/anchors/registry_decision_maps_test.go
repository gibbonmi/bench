package anchors

import "testing"

// TestDecisionMapAuthoringAnchorsRedOnRemoval holds the decision-map authoring steps and
// the template's asset path. The shaping phase file must send the author to one ready
// map's Sources block, must name both first-skeleton verbs, and must keep the unnamed
// prose-check term out. The rendered template must spell the topic-folder asset path.
// TestDecisionMapSplitAnchorsRedOnRemoval holds the two command files' half of that
// path. Each needle and diagnostic is written here independently of the registry.
func TestDecisionMapAuthoringAnchorsRedOnRemoval(t *testing.T) {
	const (
		shapeIdea = ".agents/commands/bench-shape-idea.md"
		schema    = "internal/maps/schema.go"
	)
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:   shapeIdea,
				needle: "Read one ready decision map's `## Sources` block before the first write.",
				want:   ".agents/commands/bench-shape-idea.md dropped the ready-map Sources read before the first decision-map write",
			},
			{
				file:   shapeIdea,
				needle: "Run `bench maps` and `bench gate-prose` on the first skeleton.",
				want:   ".agents/commands/bench-shape-idea.md dropped the bench-maps and bench-gate-prose checks on the first decision-map skeleton",
			},
			{
				file:      shapeIdea,
				needle:    "prose preflight",
				want:      ".agents/commands/bench-shape-idea.md writes an unnamed prose check; `bench gate-prose` is the one handle",
				forbidden: true,
			},
			{
				file:   schema,
				needle: "decisions/<topic>/assets/",
				want:   "internal/maps/schema.go dropped the decisions/<topic>/assets/ path for a map-owned asset from the decision-map template",
			},
		},
		templates: map[string]string{schema: "package maps\n%s"},
	}.check(t)
}

// TestDecisionMapSplitAnchorsRedOnRemoval holds the split shape in the two phase files.
// A map is an index, and each decision lives in one ticket file the index links. The
// shaping file must teach that shape, the ticket-naming rule, the worktree-lease claim,
// and the current-session evidence rule for a grill recommendation. The spec file must
// move the topic folder as one unit and must retire it whole. Neither file may keep the
// flat asset path, so the two forbidden rows red a command that survived the move
// unedited. Each needle and diagnostic is written here independently of the registry, so
// a rule reworded in the command cannot define itself green.
func TestDecisionMapSplitAnchorsRedOnRemoval(t *testing.T) {
	const (
		shapeIdea = ".agents/commands/bench-shape-idea.md"
		writeSpec = ".agents/commands/bench-write-spec.md"
	)
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:   shapeIdea,
				needle: "The map is an index: it lists the decisions made and links the ticket that holds each one.",
				want:   ".agents/commands/bench-shape-idea.md dropped the index role of the map file",
			},
			{
				file:   shapeIdea,
				needle: "A decision ticket is one file under the map's tickets folder, named by its number.",
				want:   ".agents/commands/bench-shape-idea.md dropped the one-file-per-decision-ticket rule",
			},
			{
				file:   shapeIdea,
				needle: "The answer lives only in the ticket file.",
				want:   ".agents/commands/bench-shape-idea.md dropped the single home of a decision answer",
			},
			{
				file:   shapeIdea,
				needle: "Decisions so far holds one gist line per resolved ticket, with a link to its file.",
				want:   ".agents/commands/bench-shape-idea.md dropped the gist-line contents of Decisions so far",
			},
			{
				file:   shapeIdea,
				needle: "Notes holds the domain, the skills a session consults, and the standing preferences of that map.",
				want:   ".agents/commands/bench-shape-idea.md dropped the contents of the map's Notes section",
			},
			{
				file:   shapeIdea,
				needle: "The shaping worktree lease is the claim, and no owner field enters a ticket.",
				want:   ".agents/commands/bench-shape-idea.md dropped the worktree-lease claim or admitted an owner field into a ticket",
			},
			{
				file:   shapeIdea,
				needle: "Name a ticket by its title, with its number beside it.",
				want:   ".agents/commands/bench-shape-idea.md dropped the title-first naming rule for a decision ticket",
			},
			{
				file:   shapeIdea,
				needle: "A grill recommendation that asserts current-code behavior names the evidence read in the current session.",
				want:   ".agents/commands/bench-shape-idea.md dropped the current-session evidence rule for a grill recommendation",
			},
			{
				file:      shapeIdea,
				needle:    "decisions/assets/",
				want:      ".agents/commands/bench-shape-idea.md keeps the flat decisions/assets/ path; a map-owned asset stays under the topic folder",
				forbidden: true,
			},
			{
				file:   writeSpec,
				needle: "Move the topic folder, its tickets and assets included, into the spec's decisions folder as one unit.",
				want:   ".agents/commands/bench-write-spec.md dropped the compile-time decision-map move",
			},
			{
				file:   writeSpec,
				needle: "Whole-folder retirement removes the compiled topic folders, their tickets and assets included.",
				want:   ".agents/commands/bench-write-spec.md dropped the whole-folder compiled decision-map retirement",
			},
			{
				file:      writeSpec,
				needle:    "decisions/assets/",
				want:      ".agents/commands/bench-write-spec.md keeps the flat decisions/assets/ path; a map-owned asset moves with the topic folder",
				forbidden: true,
			},
		},
	}.check(t)
}

// TestReadmeSplitShapeAnchorRedsOnRemoval holds the README sentence that teaches the
// split shape to a first reader: the map file is an index, and one ticket file under
// the map's tickets folder holds each decision. The needle and the diagnostic are
// written here independently of the registry, so a README rewrite that dropped the
// sentence cannot define itself green.
func TestReadmeSplitShapeAnchorRedsOnRemoval(t *testing.T) {
	anchorHarness{
		group: AfterSpecAuthorization,
		rules: []anchorRule{
			{
				file:   "README.md",
				needle: "each decision lives in one ticket file under the map's tickets folder",
				want:   "README.md dropped the split decision-map shape; each decision lives in one ticket file under the map's tickets folder",
			},
		},
	}.check(t)
}

// TestContextMapTermAnchorsRedOnRemoval holds the four map-term glossary entries in
// CONTEXT.md. The coverage map and the decision map are two artifacts, and each Avoid
// list names the bare word "map". The reader sweep entry reserves "census". The gist
// entry names the one line that carries a resolved decision into the map index. Each
// needle and diagnostic is written here independently of the registry, so a glossary
// that merged the two map entries cannot define itself green.
func TestContextMapTermAnchorsRedOnRemoval(t *testing.T) {
	const file = "CONTEXT.md"
	anchorHarness{
		group: AfterSpecAuthorization,
		rules: []anchorRule{
			{
				file:   file,
				needle: "Not \"map\", not \"traceability matrix\" — coverage map.",
				want:   "CONTEXT.md dropped the coverage-map glossary entry with the Avoid list that names the bare word map",
			},
			{
				file:   file,
				needle: "Not \"PRD\", not \"design doc\", not \"map\" — decision map.",
				want:   "CONTEXT.md decision-map entry dropped the Avoid list that names the bare word map",
			},
			{
				file:   file,
				needle: "Not \"census\", not \"consumer audit\" — reader sweep.",
				want:   "CONTEXT.md dropped the reader-sweep glossary entry with the Avoid list that reserves census",
			},
			{
				file:   file,
				needle: "one line in Decisions so far that links a resolved decision ticket",
				want:   "CONTEXT.md dropped the gist glossary entry that names the Decisions so far line linking a resolved decision ticket",
			},
		},
	}.check(t)
}
