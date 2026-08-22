# Diagnose the partial SessionStart environment

Blocked by: none
Writes: internal/sessioninspect/sessioninspect.go, internal/sessioninspect/sessioninspect_test.go, internal/bounds/bounds.go, internal/bounds/bounds_test.go, internal/systemtest/session_start_test.go, CONTEXT.md, projects/benchkit.md, CHANGELOG.md
Line: `gpt-5.6-sol` / high — SessionStart guidance is cross-harness and the login-shell discovery carries a trust assumption.

## What to build

In a Go repository whose harness PATH cannot resolve Go, bounded session
inspection asks the user's clean Bash login where Go resolves. A safe absolute
existing executable produces an informational diagnosis and a shell-quoted
command that prepends its directory to the existing PATH. Discovery runs through
the process-group owner under the two-second environment-discovery bound; timeout
or nonzero exit prints no recovery command, continues SessionStart, and exits
zero. The reported binary is never executed and no conventional home path is
searched. The system fixture compares partial and healthy shapes over one
repository and fake Go artifact, while phase-close evidence separately records
actual Codex-client and CLI invocations against one WSL repository/toolchain
before any portability claim. The glossary and project cold-session guidance
record environment closure and PATH preservation.

## Acceptance

- [ ] The reproduced loaded-marker/no-PATH-Go shape prints the discovered clean-login Go diagnosis (covers TE5).
- [ ] The recovery shell-quotes and prepends the discovered directory to literal `"$PATH"` (covers TE6 and TE13).
- [ ] SessionStart does not execute the discovered Go artifact (covers TE7).
- [ ] Partial and missing-tool cases exit zero (covers TE8).
- [ ] The same repository with Go already on harness PATH prints no recovery warning (covers TE9).
- [ ] No clean-login Go produces an honest diagnosis with no path-bearing PATH assignment (covers TE10).
- [ ] Outside a repository remains completely silent (covers TE11).
- [ ] Relative, nonexistent, multiline, and control-bearing discovery output cannot enter a recovery command (covers TE14).
- [ ] Discovery exceeding `bounds.EnvironmentDiscoveryTimeout` is killed as a process group, prints no recovery command, and SessionStart continues with exit zero (covers TE15).
- [ ] A timed-out discovery shell with a sentinel-writing child leaves no live child and no sentinel after SessionStart returns (covers TE16).
- [ ] Phase-close evidence records one actual Codex-client invocation and one actual CLI invocation against the same WSL repository, Go executable, and initialization files before making a portability claim (manual evidence for story 12).
