# Read a stack file and its overrides

Blocked by: none
Writes: internal/worktree/land_stack_input.go (new), internal/worktree/land_stack_input_test.go (new), internal/worktree/parallel_census_test.go
Covers: LS54, LS55, LS56, LS57, LS58, LS59, LS60, LS61, LS62, LS64, LS65, LS96, LS97, LS98

## What to build

Add the stack reader to `internal/worktree`. It is a function that takes the stack path, the caller's working directory, and the override values. It returns the members in file order and the parsed overrides, or one refusal:

- Each member carries its file fields, its file tip, and its effective tip, which is the override tip when an override names its request.
- The parsed overrides are the list of `<request>=<commit>` pairs in argument order, keyed by request.

It reads no ledger and no repository state, so it lands green with no caller. Ticket 07 calls it from the verb, maps each refusal to exit 2, and gives the parsed overrides to its rerun composer. No other code splits an override value.

The reader does these steps:

1. It resolves a relative path against the working directory. It reads the file only when the path resolves to a regular file, so a FIFO, a device, or a directory refuses before any open. An absent path refuses.
2. It decodes the file with `jsonfile.DecodeExactDocument` into `{"members":[{"path","request","base","source_tip","spec","message"}]}`. The `spec` field is optional, and each other field is required. So an empty file, a duplicate key, an unknown field, and trailing content refuse.
3. It refuses fewer than two members and more than sixteen members.
4. It maps each member to the argument vector of one landing and parses that vector with `landGrammar`. So a member accepts exactly what one landing accepts, and a refusal names the member's 1-based index and the grammar line.
5. It parses each override value `<request>=<commit>`. It refuses a value with no `=`, an empty side, a request that names no member, and a second override for one request. Each valid override sets the effective tip of its member.

The members keep the path value as written. Ticket 07 resolves each path to a canonical worktree path.

## Acceptance

- [ ] The reader refuses an absent stack file (LS54).
- [ ] The reader refuses an empty stack file (LS55).
- [ ] The reader refuses a stack file with a duplicate key (LS56).
- [ ] The reader refuses a stack file with an unknown field (LS57).
- [ ] The reader refuses a stack file with content after the object (LS58).
- [ ] The reader refuses a stack file with one member (LS59).
- [ ] The reader refuses a stack file with seventeen members (LS60).
- [ ] The reader refuses a member that has no `message` field, and it names the member index (LS61).
- [ ] The reader accepts a file with sixteen members and no final newline (LS62).
- [ ] The reader refuses a member with an empty `message` with its index and the grammar line of the single landing (LS64).
- [ ] The reader refuses a FIFO at the stack path as not a regular file, and the call returns with no writer on the FIFO (LS65).
- [ ] The reader refuses an override whose request names no member (LS96).
- [ ] The reader refuses two overrides for one request (LS97).
- [ ] The reader refuses an override with no `=`, an empty request, or an empty commit (LS98).
- [ ] A valid override gives its member the override tip as its effective tip, keeps the file tip beside it, and appears in the parsed overrides keyed by its request.

## Checkpoint verification

- `bench test --package ./internal/worktree --run 'TestStackReader'`
- `bench test --check root-conformance`

## Mutation proof

Lower the member bound from sixteen to fifteen. The LS62 test must fail. Spell the probe against the code that the author writes, and require `bench probe` to report `bit` with a byte-exact restore.
