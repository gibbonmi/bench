# One owner bounds every public Bench response

Every public Bench command prints a bounded response. One response owner
applies a line value of 10 and a byte value of 4096. The production policy
registry holds both values, and no other package states them. The owner
derives its line cut from the two values.

A response over either value prints its first lines, one spill line, and its
last lines. A private spill file holds the complete output of both streams, in
arrival order. The spill line names the totals, the omitted lines, the cut
lines, and the file. A long line keeps a prefix no longer than the line cut,
and the cut never splits a well-formed character.

Each public command declares a bound disposition. The set of exempt commands
is closed, and each exemption states its reason. Hooks and internal plumbing
are outside the bound. The spill store is private to the Bench home, and the
retirement of an assignment removes that assignment's spill files.

This budget belongs to the response owner only. A command view that selects
rows keeps its selection, and it does not state a second budget.
