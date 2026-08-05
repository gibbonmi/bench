package gate

import (
	"errors"
	"os/exec"
	"time"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

type executionEvaluation interface {
	acceptPre() (subject, error)
	validatePre() (subject, error)
	scope(Resolution, runMode, time.Time) componentScoping
	capturePost() (subject, error)
	postStrippedSubject() (subject, error)
}

type gateEvaluation struct {
	runtimeRoot           string
	identityRoot          string
	prospective           bool
	preSource             treeSource
	validateTree          func(string) error
	postSource            treeSource
	pre                   *treeGeneration
	post                  *treeGeneration
	acceptedSubject       subject
	acceptedStripped      subject
	acceptedStrippedErr   error
	acceptedStrippedReady bool
	postStripped          subject
	postStrippedErr       error
	postStrippedReady     bool
	scoping               componentScoping
}

type engineEvaluation struct {
	root   string
	engine gateEngine
	pre    *treeGeneration
	post   *treeGeneration
}

func newEngineEvaluation(root string, engine gateEngine) *engineEvaluation {
	return &engineEvaluation{root: root, engine: engine}
}

func (e *engineEvaluation) acceptPre() (subject, error) {
	generation, err := captureTreeGeneration(workingTreeSource{root: e.root})
	if err != nil {
		return subject{}, err
	}
	e.pre = generation
	return e.engine.BuildSubject(e.root)
}

func (e *engineEvaluation) validatePre() (subject, error) { return e.engine.BuildSubject(e.root) }

func (e *engineEvaluation) capturePost() (subject, error) {
	generation, err := captureTreeGeneration(workingTreeSource{root: e.root})
	if err != nil {
		return subject{}, err
	}
	e.post = generation
	return e.engine.PostRunSubject(e.root)
}

func (e *engineEvaluation) postStrippedSubject() (subject, error) {
	if e.post == nil {
		return subject{}, errors.New("post generation unavailable")
	}
	return buildStrippedSubjectForGeneration(e.root, e.post)
}

func (e *engineEvaluation) scope(res Resolution, mode runMode, now time.Time) componentScoping {
	if e.pre == nil {
		return componentScoping{}
	}
	return scopeComponentsForIdentityGenerations(e.root, res, mode, now, e.pre, e.pre, e.pre)
}

func newWorkingTreeEvaluation(root string) *gateEvaluation {
	return &gateEvaluation{
		runtimeRoot:  root,
		identityRoot: root,
		preSource:    workingTreeSource{root: root},
		validateTree: func(want string) error {
			got, err := (workingTreeSource{root: root}).tree()
			if err != nil {
				return err
			}
			if got != want {
				return errors.New("tree changed")
			}
			return nil
		},
		postSource: workingTreeSource{root: root},
	}
}

func newProspectiveTreeEvaluation(runtimeRoot, identityRoot, tree string) *gateEvaluation {
	return &gateEvaluation{
		runtimeRoot:  runtimeRoot,
		identityRoot: identityRoot,
		prospective:  true,
		preSource:    prospectiveTreeSource{root: identityRoot, treeID: tree},
		validateTree: func(want string) error { return validateProspectiveCheckout(runtimeRoot, want) },
		postSource:   workingTreeSource{root: runtimeRoot},
	}
}

func validateProspectiveCheckout(root, tree string) error {
	if err := exec.Command("git", "-C", root, "diff-index", "--quiet", tree, "--").Run(); err != nil {
		return errors.New("tree changed")
	}
	untracked, err := benchgit.Output("-C", root, "ls-files", "--others", "--exclude-standard", "-z")
	if err != nil {
		return err
	}
	if untracked != "" {
		return errors.New("tree changed")
	}
	return nil
}

func (e *gateEvaluation) acceptPre() (subject, error) {
	generation, err := captureTreeGeneration(e.preSource)
	if err != nil {
		return subject{}, err
	}
	e.pre = generation
	plan, err := e.build(generation)
	if err != nil {
		return subject{}, err
	}
	if !e.prospective {
		e.acceptedStripped, e.acceptedStrippedErr = buildStrippedSubjectForGeneration(e.runtimeRoot, generation)
		if e.acceptedStrippedErr != nil {
			return subject{}, e.acceptedStrippedErr
		}
		e.acceptedStrippedReady = true
	}
	e.acceptedSubject = plan
	return plan, nil
}

func (e *gateEvaluation) validatePre() (subject, error) {
	if e.pre == nil {
		return subject{}, errors.New("pre generation unavailable")
	}
	if err := e.validateTree(e.pre.tree); err != nil {
		return subject{}, err
	}
	return e.build(e.pre)
}

func (e *gateEvaluation) scope(res Resolution, mode runMode, now time.Time) componentScoping {
	if e.prospective || e.pre == nil {
		return componentScoping{}
	}
	e.scoping = scopeComponentsForIdentityGenerations(e.runtimeRoot, res, mode, now, e.pre, e.pre, e.pre)
	return e.scoping
}

func (e *gateEvaluation) capturePost() (subject, error) {
	generation, err := captureTreeGeneration(e.postSource)
	if err != nil {
		return subject{}, err
	}
	e.post = generation
	plan, err := e.build(generation)
	if err != nil {
		return subject{}, err
	}
	if !e.prospective {
		e.postStripped, e.postStrippedErr = buildStrippedSubjectForGeneration(e.runtimeRoot, generation)
		e.postStrippedReady = e.postStrippedErr == nil
	}
	return plan, nil
}

func (e *gateEvaluation) acceptedStrippedSubject() (subject, error) {
	if e.pre == nil || !e.acceptedStrippedReady {
		if e.acceptedStrippedErr != nil {
			return subject{}, e.acceptedStrippedErr
		}
		return subject{}, errors.New("pre generation unavailable")
	}
	return e.acceptedStripped, nil
}

func (e *gateEvaluation) postStrippedSubject() (subject, error) {
	if e.post == nil || !e.postStrippedReady {
		if e.postStrippedErr != nil {
			return subject{}, e.postStrippedErr
		}
		return subject{}, errors.New("post generation unavailable")
	}
	return e.postStripped, nil
}

func (e *gateEvaluation) build(generation *treeGeneration) (subject, error) {
	if e.prospective {
		return buildProspectiveSubjectForGeneration(e.runtimeRoot, e.identityRoot, generation)
	}
	return buildSubjectForGeneration(e.runtimeRoot, e.identityRoot, generation)
}
