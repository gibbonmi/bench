package main

import "github.com/gibbonmi/bench/internal/specbuild"

// The lifecycle reaches all three of these capabilities by type assertion on the owner it
// was constructed with — releaseOwnerFrom, abandonOwnerFrom, and the promotion widening in
// assign.go. An assertion that fails does not fail loudly: it returns false and the
// lifecycle proceeds with the capability quietly gone, so losing or renaming one method
// here would downgrade promotion, release, or abandon at runtime with every test still
// green. These pins turn that silent downgrade into a compile error at the wiring site.
var (
	_ specbuild.GateOwner          = productionGateOwner{}
	_ specbuild.PromotionGateOwner = productionGateOwner{}
	_ specbuild.ReleaseOwner       = productionWorktreeOwner{}
	_ specbuild.AbandonOwner       = productionWorktreeOwner{}
)
