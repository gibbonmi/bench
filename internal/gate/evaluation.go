package gate

import (
	"errors"
	"os/exec"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

type executionEvaluation interface {
	acceptPre() (subject, error)
	validatePre() (subject, error)
	capturePost() (subject, error)
}

type gateEvaluation struct {
	runtimeRoot     string
	identityRoot    string
	prospective     bool
	preSource       treeSource
	validateTree    func(string) error
	postSource      treeSource
	pre             *treeGeneration
	post            *treeGeneration
	acceptedSubject subject
}

type engineEvaluation struct {
	root   string
	engine gateEngine
	pre    *treeGeneration
	post   *treeGeneration
}

func newEngineEvaluationAtKit(root, _ string, engine gateEngine) *engineEvaluation {
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

func newWorkingTreeEvaluation(root string) *gateEvaluation {
	return newWorkingTreeEvaluationAtKit(root, root)
}

func newWorkingTreeEvaluationAtKit(root, _ string) *gateEvaluation {
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
	return plan, nil
}

func (e *gateEvaluation) build(generation *treeGeneration) (subject, error) {
	if e.prospective {
		return buildProspectiveSubjectForGeneration(e.runtimeRoot, e.identityRoot, generation)
	}
	return buildSubjectForGeneration(e.runtimeRoot, e.identityRoot, generation)
}
