# Restore the installed shim after candidate repair

Blocked by: none
Writes: .agents/commands/bench-final-check.md, CHANGELOG.md
Covers: none

## What to build

Close the broker-rehearsal learnings from the sealed drain.
Add the missing durable-shim restoration step to the existing final-check owner.
Use the primary checkout's wrapper after candidate repair and before release.
Reuse the existing post-merge status duty to verify the installed command.

## Acceptance

- [x] Final-check states when candidate repair can select a temporary shim target.
- [x] Kit-source guidance names the durable primary wrapper and the restoration order.
- [x] The existing status duty verifies the installed command after source release.
- [x] Existing final-check anchors retain their verdicts.

## Verification

The two preceding Git-reader landings exercised this sequence.
Candidate doctor repair selected the temporary wrapper each time.
Primary wrapper repair restored the installed shim before each source release.
Installed `bench status` passed after both releases.
Run the guidance anchor check, prose checks, and landing gate for this addition.
