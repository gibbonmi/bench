package main

import (
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/env"
	"github.com/gibbonmi/bench/internal/gate/prospectiveartifact"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/releaseevidence"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
	"github.com/gibbonmi/bench/internal/worktree"
)

// processResources is the composition root's protection protocol. Each kind's validator
// comes from its resource owner; a kind without a row refuses admission.
var processResources = protection.Protocol{
	Validators: protection.Validators{
		protection.KindCache:           gocache.ValidateProtection,
		protection.KindProspective:     prospectiveartifact.ValidateProtection,
		protection.KindBinary:          runbinary.ValidateProtection,
		protection.KindKitTest:         env.ValidateProtection,
		protection.KindReleaseEvidence: releaseevidence.ValidateProtection,
		protection.KindWorktree:        worktree.ValidateProtection,
	},
	ReadLimit: bounds.ControlRecordLimit,
}
