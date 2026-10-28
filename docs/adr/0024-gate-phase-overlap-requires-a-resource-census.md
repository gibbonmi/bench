# Gate phase overlap requires a resource census

The background grade is the one Bench-owned workload that runs beside a
foreground gate. In a kit worktree, a lane pass can queue a background grade of
its tip. One grade runs at a time on the machine. Its Go test phases run at two
packages and two parallel tests, at nice 19, and in the idle I/O class. A linked
repository starts no grade.

A resource census measured this bound beside the foreground gate. The census
ratio of beside to alone was 1.258, above its bound of 1.2. The reviewer accepted
that ratio and admitted the grade. A shared token pool is not necessary, because
this one bounded grade is the only concurrent workload. Any other proposal for
phase overlap requires renewed shaping and a fresh resource census before
adoption.

The retained pool design budgets cores, shares tokens through inherited
descriptors, and makes each concurrent parent own its pool. The parent reclaims
grants when children exit. Acquisition deadlines bound waits, with a one-core
fallback and a shortfall diagnostic. Proportional reserved headroom remains
unpriced. Every phase participates, and Go package and thread widths must fit
the grant. The choice of width split also requires measurement.
