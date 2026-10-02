# CLI and desktop consistency review

C1 author verification is complete.
Independent review remains pending.
No post-review repair cycle has been used.

## Standards

Pending.

## Spec

Pending.

## Coverage

Pending.

```bench-review-record
{
  "version": 1,
  "spec": "specs/cli-desktop-consistency/spec.md",
  "plan_digest": "sha256:8a7ffa8720bf9e2ab3de0c7c3392ee62bdc9e868652deb7adcf02b80bb07fe85",
  "implementation_session": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
  "chunks": [
    {
      "id": "C1",
      "base": "6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69",
      "tip": "1217c75334b173c8a2d039b876039555738beda1",
      "plan_digest": "sha256:8a7ffa8720bf9e2ab3de0c7c3392ee62bdc9e868652deb7adcf02b80bb07fe85",
      "source_digest": "422d7577bfc454c113673a6232d4151a881fbf71",
      "acceptance_rows": [
        "CD01",
        "CD02",
        "CD03",
        "CD04",
        "CD05",
        "CD06",
        "CD07",
        "CD08",
        "CD09",
        "CD10",
        "CD11",
        "CD12",
        "CD13",
        "CD14",
        "CD15",
        "CD16",
        "CD17"
      ],
      "verification": [
        {
          "id": "c1-compatibility-author",
          "performer": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "422d7577bfc454c113673a6232d4151a881fbf71",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
            "digest": "sha256:3856d224fa2fb2b321317277f865600d81a175f23b9f25c751cf7ca19abc4a9f",
            "excerpt": "Author result in Codex chat 01a0fc58-0f59-7f52-9c4c-a25247d5fa21: bench test --package ./internal/compatibility passed with no skips in 5 ms after the effective-config correction. The current chunk retains the verified package bytes. The ledger retains the row mutation results.\n"
          },
          "requirement": "compatibility",
          "command": "bench test --package ./internal/compatibility",
          "exit_code": 0
        },
        {
          "id": "c1-adopt-author",
          "performer": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "422d7577bfc454c113673a6232d4151a881fbf71",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
            "digest": "sha256:f86de95151685cd755c8e6b92b5df722aefd90d7f3295b25a94b0e625c2b67d8",
            "excerpt": "Author result in Codex chat 01a0fc58-0f59-7f52-9c4c-a25247d5fa21: bench test --package ./internal/adopt passed with no skips in 29298 ms after the canonical-payload collector correction. The current chunk retains the verified package bytes.\n"
          },
          "requirement": "adopt",
          "command": "bench test --package ./internal/adopt",
          "exit_code": 0
        },
        {
          "id": "c1-system-author",
          "performer": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "422d7577bfc454c113673a6232d4151a881fbf71",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
            "digest": "sha256:b0d398543c61313251addc5fffd37bc6a9a8958097e65e558452151aea7734e8",
            "excerpt": "Author result in Codex chat 01a0fc58-0f59-7f52-9c4c-a25247d5fa21: bench test --check system passed with no skips in 90110 ms after exact doctor-route source restoration. The current chunk retains the verified system and doctor source bytes.\n"
          },
          "requirement": "system",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "c1-doctor-route-probe-author",
          "performer": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "422d7577bfc454c113673a6232d4151a881fbf71",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
            "digest": "sha256:0c9c050abe0a0e3599f58f06d53809c519fbbbe93a52231c62c8d6e71da2b3db",
            "excerpt": "Author result in Codex chat 01a0fc58-0f59-7f52-9c4c-a25247d5fa21: the spec's manual swap changed Run: adoptCommand(\"doctor\") to Run: adoptCommand(\"setup\"). bench test --check system exited 1 in 72514 ms; TestCompatibilityMissingPath and TestCompatibilityPaths failed, with one existing wrapper-reload failure. The backup was restored and byte/mode identity verified. The restored system suite exited 0 with no skips in 90110 ms. This is the corrected probe; the earlier unquoted-path assertion failure was not accepted.\n"
          },
          "requirement": "doctor-route-probe",
          "command": "doctor-route-system-swap: follow the Doctor-route mutation procedure",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
              "digest": "sha256:0c9c050abe0a0e3599f58f06d53809c519fbbbe93a52231c62c8d6e71da2b3db",
              "excerpt": "Author result in Codex chat 01a0fc58-0f59-7f52-9c4c-a25247d5fa21: the spec's manual swap changed Run: adoptCommand(\"doctor\") to Run: adoptCommand(\"setup\"). bench test --check system exited 1 in 72514 ms; TestCompatibilityMissingPath and TestCompatibilityPaths failed, with one existing wrapper-reload failure. The backup was restored and byte/mode identity verified. The restored system suite exited 0 with no skips in 90110 ms. This is the corrected probe; the earlier unquoted-path assertion failure was not accepted.\n"
            }
          }
        },
        {
          "id": "c1-evidence-command-author",
          "performer": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "422d7577bfc454c113673a6232d4151a881fbf71",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
            "digest": "sha256:c1bc8cf142afb5d93ec141804f72ded1c4205bd408b76ea04615362161a750d2",
            "excerpt": "Native result in this chat: bench test --package ./internal/preflight/evidencecmd passed with no failures or skips in 39242 ms. Before that, the compiled wrong-base mutation in internal/diff/diff.go failed TestReviewFileReconstruction/predecessor_base at the exact-base assertion, and bench probe reported bit and restored yes. The current chunk retains those source bytes.\n"
          },
          "requirement": "evidence-command",
          "command": "bench test --package ./internal/preflight/evidencecmd",
          "exit_code": 0
        }
      ],
      "reviews": []
    }
  ],
  "completion": {
    "state": "pending",
    "source_digest": "",
    "performer": "",
    "reconciliation": {},
    "verification": []
  }
}
```
