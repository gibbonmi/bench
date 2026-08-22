# Roadmap recurrence uses owner-scoped incident keys

Each roadmap item owns a sorted set of unique incident keys, and its recurrence
count is derived from that set rather than stored independently. Capture evidence
cites one current owner and one stable incident key. Reviewed maintenance records a
new key before removing its source, while malformed or ambiguous evidence makes the
sequence untrusted instead of inviting inference. Recurrence influences priority
only after severity, actionability, dependencies, and explicit reviewer pricing
remain tied, so no command automatically reorders the roadmap.
