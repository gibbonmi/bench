# Canary inventory has no separate latency budget

Status: deprecated

Canary inventory is bounded in-process validation of accepted fixture bindings and does
not execute their owners. Direct planted-reason proofs run as ordinary tests and inherit
their packages' normal performance expectations, so the canary contract defines no
separate latency policy.
