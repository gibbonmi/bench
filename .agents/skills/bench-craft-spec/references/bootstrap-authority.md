# Bootstrap authority before execution

Charged from `craft-spec` when a story claims trusted or authenticated execution,
or refusal before execution.

Trace the real process from the raw OS entrypoint to the first
candidate-controlled instruction and through every executable hop, naming at
each hop the already-trusted validator and how it authenticates the next
executable before launching the next executable. A path, record, digest, or
executable cannot authenticate itself. Without an independent trust root, the
design is incomplete unless a reviewer-visible trust assumption says otherwise.

Coverage places markers in candidate-controlled executables, corrupts or
replaces the next authority, and asserts no marker runs before refusal.
Slicing names who publishes, locates, validates, and invokes the first
trusted executable; no complete owner is a pre-build slicing defect.

Two guards on the exclusions: no **amputated callers** (a **Won't handle**
line needs one surviving in-scope caller) and **compatibility proven, not
promised** (divergence from a named external format is a reviewer decision,
never a silent promise).
