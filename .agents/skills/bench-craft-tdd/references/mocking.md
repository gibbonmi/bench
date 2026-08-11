# Mocking at system seams

`craft-tdd` routes here when a test needs a double and the question is whether
and how to mock. Two rules decide the whether; the honesty rule decides the
how.

- **Mock only at system seams** — the places the process meets the world:
  another process, the network, the clock, the filesystem, randomness, an
  external API. These are slow, nondeterministic, or unavailable in a test,
  which is what earns a double.
- **Never mock what you own in-process.** An in-process collaborator you can
  construct is used for real; a mock there pins the current internal wiring,
  so a behavior-preserving refactor kills the test. If a real construction is
  too painful for a test, that pain is a shape problem in the module — fix
  the shape (`craft-seams`), don't paper over it with a mock.
- **Keep stubs honest.** A stub scripted to return exactly the success shape
  hollows the test out — it passes while the real integration is never
  exercised. Return realistic payloads, and script the failure behaviors the
  real seam produces — timeouts, refusals, malformed replies — drawn from
  the edge inventory, not only the happy reply.

The mapping from a dependency's category to its whole test strategy — when the
real thing beats any double — is `craft-seams`'s
`references/dependency-categories.md`.

## The pair

```python
fake_http.respond("/issues/41/close", status=503, body="maintenance")
repo = IssueRepo(fake_http)
with pytest.raises(IssueUnavailable):
    repo.close(41)
```
Good — the double stands at a real system seam (the injected network client)
and scripts a failure shape the live service actually produces; the assertion
reads the caller-visible outcome.

```python
repo = IssueRepo(real_http)
repo.cache = Mock()
repo.close(41)
repo.cache.invalidate.assert_called_once_with(41)
```
Bad — the mock replaces an owned in-process collaborator and asserts the call
pattern, so the test survives a broken close and dies on a harmless refactor.
