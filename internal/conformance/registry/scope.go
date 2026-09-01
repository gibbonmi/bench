package registry

// ConformanceScopeEnv transports one named conformance check to the root entry point.
const ConformanceScopeEnv = "BENCH_CONFORMANCE_SCOPE"

// ConsumerOnlyEnv limits load-validity-metadata to the rule shipped in linked
// consumer gates. The scaffolded gate is the only producer of this intent.
const ConsumerOnlyEnv = "BENCH_CONFORMANCE_CONSUMER_ONLY"
