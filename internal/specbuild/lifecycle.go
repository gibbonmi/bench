// Package specbuild owns the durable lifecycle of a reviewed spec build.
package specbuild

import "context"

// GateOwner validates retained exact green evidence for a working subject.
type GateOwner interface {
	Bootstrap(context.Context, string, string, string) error
}

// WorktreeOwner creates a request-idempotent owned worktree at start.
type WorktreeOwner interface {
	Create(context.Context, string, string, string, string) (OwnedWorktree, error)
}

// OwnedWorktree identifies a worktree created by the existing ownership owner.
type OwnedWorktree struct {
	ID, Path string
}

// Service coordinates spec build transitions from one working checkout.
type Service struct {
	root      string
	gate      GateOwner
	worktrees WorktreeOwner
}

// New constructs a lifecycle service rooted at one working checkout.
func New(root string, gate GateOwner, worktrees WorktreeOwner) *Service {
	return &Service{root: root, gate: gate, worktrees: worktrees}
}

// Status is the compact public projection of one spec build.
type Status struct {
	Slug, State, Subject, Next string
}

// Assignment records the externally useful ownership binding for one ticket.
type Assignment struct {
	ID, Path, Base string
	Rows           []string
	Fence          []string
	Assumptions    []string
}
