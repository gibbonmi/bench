# Canary inventory has no separate concurrency budget

Status: deprecated

Canary inventory validates accepted fixture bindings in process and does not execute
their owners. Direct planted-reason proofs run as ordinary tests under the existing test
runner, so the canary contract defines no separate parallelism or runtime-width policy.
