# Close classifier edge gaps

Blocked by: 02-refuse-bench-shell-follow-ons.md
Writes: internal/benchguard, internal/systemtest

## What to build

Keep outer shell syntax visible when one supported wrapper contains a Bench call.
Remove leading file-descriptor redirections from command position.
Resolve supported routine prefixes without skipping the Bench executable.

## Acceptance

- [ ] An outer pipeline, `&&`, or redirection after a wrapper-contained Bench call is refused.
- [ ] A leading file-descriptor redirection cannot hide a Bench call.
- [ ] Supported `env`, `command`, `nohup`, `timeout`, and `xargs` forms cannot hide a Bench call.
