# What a good test looks like

`craft-tdd` routes here when the seam is chosen and the question is what the
test itself should say. Four properties decide it.

- **Public-interface behavior.** The test calls the interface callers use and
  asserts what a caller observes — returned values, visible state, emitted
  errors. It never reaches for a private field, an internal helper, or the
  call pattern between collaborators; a refactor that keeps behavior must
  keep the test green.
- **Independently derived expected values.** The expected value comes from
  the spec, a hand computation, or a trusted reference — written as a
  literal. Never run the implementation and paste back what it returned, and
  never recompute the expectation with the implementation's own algorithm:
  both pass by construction against any bug.
- **One behavior per test.** One scenario, one logical assertion, a name that
  states the behavior. When it fails, the name alone says what broke; a test
  asserting five things reports only that the first one did.
- **Realistic failure shapes.** Error cases assert the failure the caller
  actually meets — the typed error, the message a user sees, the state left
  behind — not merely "an error occurred". Draw the cases from the edge
  inventory, not from what the implementation happens to raise.

## The pair

```go
func TestCloseMarksIssueClosed(t *testing.T) {
	repo := NewIssueRepo(fakeHTTP)
	repo.Close(41)
	if got := repo.Get(41).Status; got != "closed" {
		t.Errorf("status = %q, want %q", got, "closed")
	}
}
```
Good — one behavior, exercised and read back through the public interface,
against a literal expectation the spec supplies.

```go
func TestClose(t *testing.T) {
	repo := NewIssueRepo(fakeHTTP)
	repo.Close(41)
	want := statusLabel(closedState) // the implementation's own derivation
	if repo.Get(41).Status != want {
		t.Error("close failed")
	}
}
```
Bad — the expectation is recomputed with the implementation's algorithm, so
the test passes even when `statusLabel` is wrong, and the failure message
names no observed value.
