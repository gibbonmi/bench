# Gate phase overlap requires a resource census

The serial gate bounds demand without a shared token pool. The token-pool
proposal is closed because no Bench-owned concurrent phase workload remains
for a reserve to protect. A future proposal for phase overlap requires renewed
shaping and a fresh resource census before adoption.

The retained pool design budgets cores, shares tokens through inherited
descriptors, and makes each concurrent parent own its pool. The parent reclaims
grants when children exit. Acquisition deadlines bound waits, with a one-core
fallback and a shortfall diagnostic. Proportional reserved headroom remains
unpriced. Every phase participates, and Go package and thread widths must fit
the grant. The choice of width split also requires measurement.
