# The kit rule keeps a wrapper form and a Go form

Status: accepted

The kit rule selects the kit that serves a tree.
When the kit is the top level of its own checkout, and the tree is a working tree of the same repository, the tree is the kit.
Every other input keeps the given kit, and a git failure also keeps the given kit.

Two production forms hold this rule.
The shell wrapper holds one form, because the wrapper selects the Bench binary before any Go code runs.
The gate package holds the other form, and a merge uses it to find the lane that grades the merge target.

The one-source rule forbids two production forms of one fact.
Its exception for an independent test expectation does not apply to two production forms.
The reviewer accepts this pair as a recorded exception to the one-source rule.

The wrapper differential test guards the agreement of the two forms.
For each repository shape of the wrapper's kit test with no environment override, the test runs the wrapper and the Go form.
The test requires the same kit from both forms, so a Go form that drops a condition of the rule turns the test red.

## Considered options

- A Go owner of the kit rule that the wrapper consumes: rejected, because the wrapper selects the binary before that owner can run. This owner needs its own design.
