# Document anchors use a declarative substring registry

Document guidance anchors are ordered declarative data evaluated by one shared matcher, with section-scoped placement and HTML-comment exclusion where required; the agent-facing query reads that registry rather than re-parsing enforcement code. Substring matching is intentionally retained because it gives deterministic, byte-stable diagnostics and cheap omission or relocation oracles, while paraphrased semantic restatements remain review-owned until evidence justifies a stronger mechanism.
