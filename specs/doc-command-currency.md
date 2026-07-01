# living doc command currency

## Problem

The command surface is now:

- `/bench-setup-repo`
- `/bench-shape-idea`
- `/bench-write-spec`
- `/bench-debug`
- `/bench-implement-spec`
- `/bench-review-implementation`
- `/bench-final-check`
- `/bench-update-kit`
- `/bench-integrate-learnings`

Some living docs still name older command forms. A cold session that greps for a
command can find names that no longer exist, and the gate only checks the AGENTS.md
index plus command adapters.

## Solution

Define the current living-doc surface and make the gate fail on stale command names
there. Refresh the cold-pickup docs at the same time.

Living docs for command currency:

- `README.md`
- `AGENTS.md`
- `.bench/BENCH.md`
- `.bench/learnings.md`
- `CONTEXT.md`
- `HANDOFF.md`
- `specs/`, except files explicitly marked `command-currency: historical`

Archival or deliberation records:

- `decisions/`
- dated release/history entries in `CHANGELOG.md`

`CHANGELOG.md` keeps historical command names inside dated entries, but its header and
Unreleased/current guidance must use the live command names.

## User stories

1. As a cold session, I want setup, workflow, and learning docs to name commands that
   exist now.
2. As a kit maintainer, I want `HANDOFF.md` to be a current pickup doc or not ship at
   all.
3. As a kit maintainer, I want the gate to catch command-name drift in living docs.
4. As a reader of history, I want old command names to remain in historical decision
   and changelog context where they describe what was true at the time.

## Implementation decisions

- **Keep `HANDOFF.md`.** It ships in `package.json`, so rewrite it as a short
  current-state pickup guide that points at the canonical docs and lists the full CLI
  surface.
- **Gate helper:** derive valid slash command names from `.agents/commands/*.md`.
  Scan the living files for slash-command-looking tokens and fail when a token is not
  in the valid set or a small allowed external set.
- **Allowed external tokens:** keep only real, intentional non-Bench references such as
  `/model` if they appear on the living surface. Do not allow historical Bench aliases
  such as bench-learn, bench-update, resynthesize, spec, or build.
- **History stays history:** do not gate `decisions/` or historical changelog entries.
  If stale names in those files become confusing later, rewrite or retire the file as a
  separate doc decision.

## Testing decisions

- **Seam:** `bench gate`, because the behavior is conformance of the repo's shipped and
  working docs.
- Add a gate check that fails on stale slash-command names in living docs.
- Add or update a canary fixture for the command-currency check, so the check must
  prove it bites.

## Out of scope

- Rewriting every historical decision map.
- Renaming commands again to the shorter bench-learn and bench-update forms.
- Policing plain English words like "spec" or "review" when they are not slash
  invocations.
