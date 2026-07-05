#!/usr/bin/env bash
# Resolves .bench/gate here — extensionless on purpose (the d77063c regression).
# Also names .bench/done-keeper.sh — a legitimately-named sibling the 1d check
# must abstain on; EXPECT pins the flagged list to the gate ref alone, so a
# regression to prefix-matching (which would also flag the done sibling)
# breaks it.
true
