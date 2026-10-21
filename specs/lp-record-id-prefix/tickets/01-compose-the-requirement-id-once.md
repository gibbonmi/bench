# Compose the requirement ID with the prefix once

Blocked by: none
Writes: internal/reviewrecord/recordcmd/ticket.go, internal/reviewrecord/recordcmd/ticket_test.go
Covers: none

## What to build

`bench record ticket <slug> --id-prefix <p>` writes one verification entry for each owed
requirement. The evidence ID of each entry joins `<p>`, a hyphen, and the requirement ID. A
plan can give a requirement an ID that already starts with `<p>-`. The command then writes
the prefix two times, for example `t2-t2-preflight-packages` for the requirement
`t2-preflight-packages` under the prefix `t2`.

When the requirement ID starts with `<p>-`, use the requirement ID as the evidence ID. When
the requirement ID does not start with `<p>-`, keep the current ID: `<p>`, a hyphen, and the
requirement ID.

## Acceptance

- [ ] A ticket call with the prefix `p` writes the evidence ID `p-tests` for an owed requirement with the ID `p-tests`.
- [ ] A ticket call with the prefix `p` writes the evidence ID `p-additional` for an owed requirement with the ID `additional`.
- [ ] The focused test of the first item is red before the fix.
