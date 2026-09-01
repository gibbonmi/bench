# Cover detected-project consumer hygiene

Blocked by: scaffold-declared-input-hygiene.md
Writes: internal/systemtest/adoption_test.go, specs/roadmap-light-path-fixes/tickets/scaffold-declared-input-hygiene.md
Covers: none

## What to build

Exercise the generated consumer gate after setup detects a project, using a
Git-ignored declared input with a hostile literal path. Demonstrate that
removing the detected-project hygiene branch turns the focused journey red.
Also demonstrate and record the removal mutation that makes the independent
seeded routing-input expectation bite. The system journey authenticates the
selected executable with its retained BENCH_KIT pin.

## Acceptance

- [ ] A detected-project setup journey executes the generated gate and refuses
      an ignored declared path containing a space or glob character.
- [ ] Removing the detected-project hygiene call makes that journey red.
- [ ] Removing a newly seeded routing input makes the seed expectation red.
- [ ] The original LF2 ticket records both focused mutation commands and their
      red output.
