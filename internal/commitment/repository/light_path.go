package repository

import (
	"errors"
	"fmt"
	"path"
	"slices"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/tickets"
)

// lightPathTicket is the one ticket of a qualifying light-path folder: the folder slug, the
// ticket's tree path, and the tree path that each of its Writes entries names.
type lightPathTicket struct {
	slug, path string
	writes     []string
}

// lightPath grades a commit that readyFor refused as unbound. One qualifying ticket in tree
// must cover every production path. A tree with no qualifying folder keeps unbound.
func (store Store) lightPath(tree string, policy *commitment.Policy, production []string, unbound error) error {
	found, err := store.lightPathTickets(tree, policy)
	if err != nil {
		return err
	}
	if len(found) == 0 {
		return unbound
	}
	return lightPathCover(found, production)
}

// lightPathPublication grades a landing that readyFor refused as unbound. published.Tree
// removes the delivered folder from tree, so the ticket of a delivery is read at its
// reviewed source, and only the folder that the delivery names counts. A spec-less landing
// is never light-path work. When tree holds a ticket that covers every production path,
// its refusal names the --spec route for that ticket's folder.
func (store Store) lightPathPublication(tree string, policy *commitment.Policy, production []string, delivery *Delivery, unbound error) error {
	read := tree
	if delivery != nil {
		read = delivery.Source
	}
	found, err := store.lightPathTickets(read, policy)
	if err != nil {
		return err
	}
	for _, ticket := range found {
		switch {
		case delivery != nil && spec.ClosedFolderPath(ticket.slug) == delivery.Spec:
			return lightPathCover([]lightPathTicket{ticket}, production)
		case delivery == nil && ticket.coversAll(production):
			return fmt.Errorf("%w; land the light-path change with --spec %q", unbound, ticket.slug)
		}
	}
	return unbound
}

// lightPathCover admits production when one ticket covers every path. A path that no
// ticket covers is named against the first ticket. Only when each path has some cover
// does the change span tickets.
func lightPathCover(found []lightPathTicket, production []string) error {
	for _, ticket := range found {
		if ticket.coversAll(production) {
			return nil
		}
	}
	for _, p := range production {
		if !slices.ContainsFunc(found, func(ticket lightPathTicket) bool { return ticket.covers(p) }) {
			return fmt.Errorf("production path %q is outside the Writes line of light-path ticket %q; add the path to that line, or run bench commitment start", p, found[0].path)
		}
	}
	return errors.New("production paths span more than one light-path ticket; a light-path change carries one ticket")
}

func (ticket lightPathTicket) covers(p string) bool {
	return slices.ContainsFunc(ticket.writes, func(entry string) bool { return tickets.Covers(entry, p) })
}

func (ticket lightPathTicket) coversAll(production []string) bool {
	return !slices.ContainsFunc(production, func(p string) bool { return !ticket.covers(p) })
}

// lightPathFolder is what one listing of specs/<slug> tells about its tickets: how many
// tickets.Ext entries lie at any depth below tickets/, the last one seen, and whether any
// of them is not a regular file.
type lightPathFolder struct {
	count     int
	ticket    string
	irregular bool
}

// lightPathTickets returns, in slug order, the ticket of each tickets-only folder of tree
// that qualifies: exactly one regular tickets.Ext entry below tickets/, a ticket that
// parses with no diagnostic, and no milestone of policy that approves the folder as a
// deliverable. One listing of specs gives every folder, count, and mode.
func (store Store) lightPathTickets(tree string, policy *commitment.Policy) ([]lightPathTicket, error) {
	const root = "specs"
	listing, err := git.Output("-C", store.Root, "ls-tree", "-r", "-z", tree, "--", root)
	if err != nil {
		return nil, fmt.Errorf("read light-path folders: %w", err)
	}
	folders := map[string]*lightPathFolder{}
	for _, record := range strings.Split(listing, "\x00") {
		metadata, name, listed := strings.Cut(record, "\t")
		rest, inRoot := strings.CutPrefix(name, root+"/")
		slug, inside, nested := strings.Cut(rest, "/")
		if !listed || !inRoot || !nested {
			continue
		}
		folder := folders[slug]
		if folder == nil {
			folder = &lightPathFolder{}
			folders[slug] = folder
		}
		below, inTickets := strings.CutPrefix(inside, "tickets/")
		if !inTickets || path.Ext(below) != tickets.Ext {
			continue
		}
		folder.count++
		folder.ticket = name
		if fields := strings.Fields(metadata); len(fields) != 3 || !(git.IndexEntry{Mode: fields[0]}).IsRegularFile() {
			folder.irregular = true
		}
	}
	slugs := make([]string, 0, len(folders))
	for slug := range folders {
		slugs = append(slugs, slug)
	}
	sort.Strings(slugs)
	var found []lightPathTicket
	for _, slug := range slugs {
		folder := folders[slug]
		if folder.irregular || folder.count != 1 || approvedDeliverable(policy, spec.ClosedFolderPath(slug)) || !spec.TicketsOnly(spec.CommitTree(store.Root, tree), slug) {
			continue
		}
		data, err := git.ReadTreeFile(store.Root, tree, folder.ticket)
		if err != nil {
			return nil, err
		}
		name := path.Base(folder.ticket)
		parsed, diagnostics := tickets.ParseTicket(name, data, []string{name}, "")
		if len(diagnostics) != 0 {
			continue
		}
		ticket := lightPathTicket{slug: slug, path: folder.ticket}
		for _, entry := range parsed.Writes {
			written, _ := tickets.WritesPath(entry)
			ticket.writes = append(ticket.writes, written)
		}
		found = append(found, ticket)
	}
	return found, nil
}

// approvedDeliverable reports whether any milestone of policy approves the deliverable at
// folder. Approved work keeps its binding, so it is never light-path work.
func approvedDeliverable(policy *commitment.Policy, folder string) bool {
	if policy == nil {
		return false
	}
	for _, milestone := range policy.Milestones {
		for _, outcome := range milestone.Outcomes {
			for _, binding := range outcome.Deliverables {
				if binding.Source.Path == folder {
					return true
				}
			}
		}
	}
	return false
}
