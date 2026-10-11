# The completion plan section

Charged from the `craft-spec` template. Every build preflight mode refuses a
spec that has no valid completion plan fence. Each chunk checkpoint reads the
same fence, as `.bench/BENCH-reference.md` states.

Put the fence under a `### Completion plan` heading at the end of
`## Implementation chunks`. The opening line is exactly the fence marker and
`bench-completion-plan`, with nothing after it. Write one fence in the spec.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "<stable chunk ID>",
      "tickets": ["<ticket basename>"],
      "verification": [
        {"id": "<requirement ID>", "command": "<exact command>"}
      ]
    }
  ],
  "final_verification": [
    {"id": "<requirement ID>", "command": "<exact command>"}
  ]
}
```

The keys have these meanings:

- `version`: Write `1` in the authored plan.
- `chunks`: The plan holds one entry for each row of the chunk table.
- `id` (chunk): The stable chunk ID from the chunk table. `--chunk` takes this
  value at the chunk checkpoint.
- `tickets`: The ticket basenames under `specs/<slug>/tickets/`. Each ticket is
  in exactly one chunk.
- `verification`: The requirements that the chunk checkpoint grades. Each
  requirement needs a passing result for the chunk source.
- `id` (requirement): A requirement name that is unique in its list.
- `command`: The exact command whose passing result the checkpoint requires.
- `probe` (optional): A named mutation. The checkpoint then also requires a
  `bench probe` result that bit, with a passed restore.
- `final_verification`: The requirements that the `--complete` checkpoint
  grades. These requirements name no ticket.

The parser refuses an unknown key, a duplicate key, and an empty list. A chunk
must not come before a chunk that holds one of its `Blocked by:` tickets. Do not
write `execution` or a requirement `ticket`; the version 2 amendment adds them.
The `check-floor` row of build preflight names each floor command that a chunk
omits.
