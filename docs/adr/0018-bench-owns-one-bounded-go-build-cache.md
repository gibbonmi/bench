# 18. Bench owns one bounded Go build cache

Status: accepted (2026-08-27)

## Decision

Bench owns one shared Go build cache. Gate runs, lane runs, and test runs share
the cache under shared locks. Cache cleaning takes exclusive ownership and
refuses when an active run holds a shared lock.

No run evicts cache data. Bench measures the cache and reports when its
footprint exceeds 10 GiB. An operator can clean the cache only when no run is
active.
