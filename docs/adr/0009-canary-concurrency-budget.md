# The canary's nested concurrency is budgeted as a product, and the budget is the Go runtime's own

The canary sweep runs a complete inner gate for every fixture, so its cost is the product of two fan-outs that nothing previously related: how many fixtures run at once, and how wide each inner gate's own build and tests run. Left alone both default to machine width, which oversubscribes the host by roughly that width squared and turns most of the gate's wall clock into contention rather than work. The accepted model budgets the product rather than either factor alone. Each inner gate is invoked with an explicit `GOMAXPROCS`, the one lever that caps both its build parallelism and the test binary it produces, and the sweep's worker count is derived by dividing the outer process's own `GOMAXPROCS` by that fixed inner width, floored at one worker and capped at the number of fixtures — so workers times inner width lands at roughly the cores actually available. On a sixteen-core host this took the gate from ten to fifteen minutes at a load average around 123 down to about five and a half minutes at a peak load under 35, with no verdict changed and no check weakened.

The budget source is deliberately the Go runtime's own value and nothing else. It is already cgroup- and container-aware and already honours an operator's explicit setting, so a smaller box or a container is accommodated through the standard Go lever rather than through a Bench-specific one. An inherited value is stripped before the pinned one is set, because the exec environment gives no guaranteed precedence to a duplicate key: appending without stripping would leave which width wins to the exec implementation rather than to this policy.

## Considered options

- **An inner width of one** tied on wall clock — within two seconds over a five-and-a-half-minute run — but leaves twice as many inner gates alive at once, and with them twice the memory and temp-directory pressure, for no measured gain.
- **An inner width of four** cost about thirty percent more wall clock by under-parallelizing a fixture count far larger than the worker pool it produced.
- **A Bench-specific concurrency knob** was rejected outright. It would be a second concurrency vocabulary to learn alongside the Go one, and the Go one already covers every case the knob would have served.
- **Capping the outer gate phases as well** buys nothing while wall clock is bound by the conformance phase rather than by the canary sweep. It stays available if contention symptoms return, but it is not part of this decision.

## Consequences

- The sweep's own package tests run nested inside a fixture's inner gate, at the pinned inner width, where the derived worker pool is a single worker. Concurrency expectations in that package are therefore keyed to the derived bound rather than to machine width, and those that genuinely require overlap declare a capability requirement instead of assuming one — a nested expectation keyed to machine width deadlocks until the phase timeout rather than failing cleanly. The project profile carries the commands that prove both directions.
- The inner width is a single named entry in the shared policy registry, and the conformance layer enforces both that the sweep consumes it and that no production code redeclares the same value under a bound-like name.
- This is a semantics-neutral latency fix, which is the only kind ADR 0003 permits for the canary sweep.
