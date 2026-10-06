# Retire the delivered FT391 spec

Blocked by: none
Writes: specs/light-path-commitment-exemption, reviews/light-path-commitment-exemption.md
Covers: none

## What to build

The light-path-commitment-exemption spec is implemented, and FT391 is delivered. ADR 0028 and ADR 0023 hold its durable decisions, and the guide states its rules.
Then `bench spec retire light-path-commitment-exemption` removes the spec folder and its review record.

## Acceptance

- [ ] The spec folder and its review record are absent, and the landing commit subject ends with `spec-retire: light-path-commitment-exemption`.
