package main

import (
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/worktree"
)

// retainsResponse reports a call whose response replaces the retained response of its
// worktree. A record call stages and records the retained response, so it must not replace
// it. A worktree call runs its child in its own process, and the child's response is the one
// to keep. A retiring leaf can remove the scope that would hold the response.
func retainsResponse(definition commandDefinition, leaf commandLeaf) bool {
	return definition.Name != "record" && definition.Name != "worktree" && !leaf.Retires
}

// retainResponse replaces the retained response of the worktree at root with the finished
// response of owner, the typed verb word, its arguments, and its exit. The worktree's commit
// and dirty paths are read after the verb returns, and only when the write follows. The
// dispatcher discards the error, so a failed write changes neither the response nor the exit.
func retainResponse(owner *responsebound.Owner, verb string, args []string, exit int, root string) error {
	return owner.Retain(worktree.Home(), root, func() (responsebound.Retained, error) {
		head, err := git.Output("-C", root, "rev-parse", "HEAD")
		if err != nil {
			return responsebound.Retained{}, err
		}
		entries, err := git.PorcelainStatus(root)
		if err != nil {
			return responsebound.Retained{}, err
		}
		dirty := make([]string, 0, len(entries))
		for _, entry := range entries {
			dirty = append(dirty, entry.Path)
		}
		return responsebound.Retained{Verb: verb, Args: args, Exit: exit, Head: head, Dirty: dirty}, nil
	})
}
