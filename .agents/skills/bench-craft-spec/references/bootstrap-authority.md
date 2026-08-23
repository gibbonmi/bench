# Bootstrap authority before execution

Charged from `craft-spec` when a story claims trusted or authenticated execution,
or refusal before execution.

Trace the real process from the raw OS entrypoint to the first
candidate-controlled instruction, and through every executable hop. At each
hop, name the already-trusted validator and how it authenticates the next
executable before it launches that executable. A path, record, digest, or
executable cannot authenticate itself. Without an independent trust root, the
design is incomplete unless a reviewer-visible trust assumption says otherwise.

Coverage places markers in candidate-controlled executables, corrupts or
replaces the next authority, and asserts no marker runs before refusal.
Slicing names who publishes, locates, validates, and invokes the first
trusted executable; no complete owner is a pre-build slicing defect.

Two guards apply to the exclusions. First, no **amputated callers**: a
**Won't handle** line needs one surviving in-scope caller. Second,
**compatibility proven, not promised**: divergence from a named external
format is a reviewer decision, never a silent promise.
