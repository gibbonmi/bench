//go:build system

// Package systemtest is the one package operand the gate's system phase carries. The
// fixture root becomes its own kit wherever the ambient kit is unset, and a phase whose
// packages the loader cannot expand refuses every citation before the two diagnostics
// this fixture pins.
package systemtest
