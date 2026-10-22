package main

import (
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// processResources is the composition root's protection protocol. Each kind's validator
// comes from its resource owner; a kind without a row refuses admission.
var processResources = protection.Protocol{
	Validators: protection.Validators{
		protection.KindCache:  gocache.ValidateProtection,
		protection.KindBinary: runbinary.ValidateProtection,
	},
	ReadLimit: bounds.ControlRecordLimit,
}
