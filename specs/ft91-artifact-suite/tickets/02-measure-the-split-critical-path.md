# Measure the split critical path

Blocked by: Split artifact contracts by subject

## What to build

Measure the composed artifact package tree and a fresh changed-tree gate, then
record the observed overlap, removed serial share, remaining critical path, and
dormant outer-width-cap verdict in the compiled decision map.

## Acceptance

- [ ] The recursive artifact JSON trace shows overlapping package pass
  intervals rather than cross-package serialization.
- [ ] The focused suite and fresh changed-tree gate commands, wall times, and
  comparison with the 109-second/128-second baseline are recorded.
- [ ] The record identifies repeated process overhead, the remaining critical
  path, and whether the measured run justifies reviving the dormant outer width
  cap.
