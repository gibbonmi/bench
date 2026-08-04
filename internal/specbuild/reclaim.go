package specbuild

import (
	"errors"
	"fmt"
	"strings"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

// ReclamationPlan is the classified provisional residue of one spec-build run and the
// fingerprint an apply must match. Applied distinguishes a receipt for refs that are gone
// from a plan for refs that are still there.
type ReclamationPlan struct {
	Slug        string
	Fingerprint string
	Applied     bool
	Refs        []ReclamationRef
}

// ReclamationRef is one provisional ref, the object it held when it was classified, and
// the disposition that decides whether an apply may delete it.
type ReclamationRef struct{ Name, Object, Assignment, Disposition string }

// ReclamationClass is one disposition and how many refs the plan puts in it.
type ReclamationClass struct {
	Name  string
	Count int
}

// Classes reports every disposition the enumeration can carry, in a fixed order and
// including the empty ones, so a receipt cannot hide a class by omitting it — a maintainer
// reading zero retained refs has to be able to tell that from the class not being counted.
func (p ReclamationPlan) Classes() []ReclamationClass {
	classes := make([]ReclamationClass, 0, 4)
	for _, disposition := range []provisionalDisposition{refReclaimable, refActive, refUnclassified, refAmbiguous} {
		class := ReclamationClass{Name: string(disposition)}
		for _, ref := range p.Refs {
			if ref.Disposition == class.Name {
				class.Count++
			}
		}
		classes = append(classes, class)
	}
	return classes
}

// Reclaim returns the read-only inventory a reclaim apply must match. It mutates no ref.
func (s *Service) Reclaim(slug string) (ReclamationPlan, error) {
	release, run, err := s.openRun(slug)
	if err != nil {
		return ReclamationPlan{}, err
	}
	defer release()
	return s.reclamationPlan(slug, run)
}

// ApplyReclaim deletes exactly the refs the named fingerprint's plan classified reclaimable
// and leaves every other class for a maintainer to judge. Each deletion carries the object
// the enumeration observed, so a ref that moved between plan and apply fails closed rather
// than being clobbered. Unlike promotion's reclamation this refuses a failed deletion
// instead of absorbing it: the fingerprint asserts the exact refs exist, so a deletion that
// does not take means the tree drifted under a maintainer who asked for a specific act.
//
// Git deletes one ref at a time and gives no transaction across the set, so a refusal
// partway through has already spent some deletions. The outcome is made honest rather than
// atomic: the returned receipt reports exactly the refs that are gone and carries the
// fingerprint of that set alone — never the spent plan's, which no longer describes
// anything — and the error names the ref that drifted. Convergence is a fresh plan, which
// re-reads the tree and covers only the remainder.
func (s *Service) ApplyReclaim(slug, fingerprint string) (ReclamationPlan, error) {
	if fingerprint == "" {
		return ReclamationPlan{}, errors.New("spec build reclaim fingerprint is required; request a fresh plan")
	}
	if !canonicalDigest(fingerprint) {
		return ReclamationPlan{}, errors.New("spec build reclaim fingerprint is malformed; request a fresh plan")
	}
	release, run, err := s.openRun(slug)
	if err != nil {
		return ReclamationPlan{}, err
	}
	defer release()
	plan, err := s.reclamationPlan(slug, run)
	if err != nil {
		return ReclamationPlan{}, err
	}
	if plan.Fingerprint != fingerprint {
		return ReclamationPlan{}, errors.New("spec build reclaim plan drifted; request a fresh plan")
	}
	if err := s.faultAt("reclaim/apply"); err != nil {
		return ReclamationPlan{}, err
	}
	deleted := ReclamationPlan{Slug: slug, Applied: true, Refs: []ReclamationRef{}}
	for _, ref := range plan.Refs {
		if ref.Disposition != string(refReclaimable) {
			continue
		}
		if err := benchgit.DeleteBranchExact(s.root, ref.Name, ref.Object); err != nil {
			deleted.Fingerprint = digest(reclamationFacts(deleted))
			return deleted, fmt.Errorf("reclaim %s: %w; request a fresh plan to finish the remainder", ref.Name, err)
		}
		deleted.Refs = append(deleted.Refs, ref)
	}
	plan.Applied = true
	return plan, nil
}

// openRun holds the run's lock and hands back its record. It reads the record rather than
// resolving the spec because the residue outlives the spec: the runs that stranded refs
// promoted long ago, and requiring a spec file would put exactly those runs out of reach.
func (s *Service) openRun(slug string) (func(), record, error) {
	if err := requireSlug(slug); err != nil {
		return nil, record{}, err
	}
	release, err := s.lock(slug)
	if err != nil {
		return nil, record{}, err
	}
	run, found, err := s.load(slug)
	if err != nil || !found {
		release()
		if err == nil {
			err = errors.New("spec build run does not exist")
		}
		return nil, record{}, err
	}
	return release, run, nil
}

func (s *Service) reclamationPlan(slug string, run record) (ReclamationPlan, error) {
	refs, err := s.provisionalResidue(run)
	if err != nil {
		return ReclamationPlan{}, err
	}
	plan := ReclamationPlan{Slug: slug, Refs: make([]ReclamationRef, 0, len(refs))}
	for _, ref := range refs {
		plan.Refs = append(plan.Refs, ReclamationRef{Name: ref.Name, Object: ref.Object, Assignment: ref.Assignment, Disposition: string(ref.Disposition)})
	}
	plan.Fingerprint = digest(reclamationFacts(plan))
	return plan, nil
}

// reclamationFacts commits the fingerprint to the run it was planned for and to every ref's
// name, object, and classification. The domain tag keeps it unusable as any other
// operation's authorization, and the NUL separator keeps a ref name containing the
// separator from forging a different inventory with the same encoding.
func reclamationFacts(plan ReclamationPlan) string {
	var facts strings.Builder
	facts.WriteString("bench-specbuild-reclaim/v1\x00" + plan.Slug)
	for _, ref := range plan.Refs {
		facts.WriteString("\x00ref\x00" + ref.Name + "\x00" + ref.Object + "\x00" + ref.Disposition)
	}
	return facts.String()
}
