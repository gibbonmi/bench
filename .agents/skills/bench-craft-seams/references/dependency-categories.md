# Dependency categories and their test strategies

`craft-seams` routes here when a seam crosses a dependency and the question is
what stands on the far side in a test. Classify the dependency first; the
category decides the strategy.

- **In-process** — code this repo owns, running in the same process: your own
  modules, helpers, domain objects. *Strategy:* use the real thing. No double
  of any kind — construct it and test through the public seam. If real
  construction is too painful, the module's shape is wrong; fix the shape.
- **Local-substitutable** — a real dependency with a faithful local
  stand-in. You can run it inside the test: an embedded or in-memory
  database, a temp directory for a filesystem, an in-process queue.
  *Strategy:* run the real substitute, injected through the same seam
  production uses. It exercises the genuine integration at test speed, so
  prefer it to a mock wherever the substitute exists.
- **Remote-owned** — a service your own team runs across a network: an
  internal API, a shared job runner. *Strategy:* mock the seam in unit tests
  with realistic payloads and failure shapes. Keep one thin integration test
  against the real service to catch contract drift. You own both sides, so
  that test is yours to run and fix.
- **True-external** — a third party you neither run nor control: a vendor
  API, a payments provider, an LLM endpoint. *Strategy:* mock the seam for
  the suite, scripting the documented failure modes (timeouts, rate limits,
  malformed replies) alongside success. Verify the mock's shapes against the
  provider's documentation or a recorded real response — never against the
  code under test. Quarantine any live-call check outside the gate, since a
  verdict that depends on a third party is a broken oracle.

Misclassification is the expensive error in both directions. Mocking an
in-process collaborator binds tests to internals (`bench-craft-tdd`'s
`references/mocking.md` owns that rule), while testing against a live
true-external makes the suite nondeterministic. When a dependency sits
between categories — a database accessed through a vendor cloud API —
classify by who controls availability at test time.
