package consumer

import tg "example.com/fixturemod/target"

// Direct is the plain call-position consumer.
func Direct() { tg.Symbol() }

// registry is the value-position consumer inside a package-level var.
var registry = []tg.T{{Run: tg.Symbol}}

// Holder gives the fixture a method consumer.
type Holder struct{}

// Use is the method consumer.
func (h Holder) Use() { tg.Symbol() }

var _ = registry
