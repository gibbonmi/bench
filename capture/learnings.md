# Learnings — usage journal

<!-- entries below -->

- 2026-08-22  After `bench commit` reported its first red, the coordinator edited
  the tree before the gate process closed. The coordinator must wait for the
  command's terminal exit before any repair. Proposed rule: treat partial gate
  output as progress, never as a process boundary.
