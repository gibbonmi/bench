package gate

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

const prospectiveGatePath = ".bench/gate-prospective.sh"

func buildProspectiveSubjectFor(root, identityRoot string) (subject, error) {
	s, err := buildSubjectForPolicy(root, identityRoot, policyVersion)
	if err != nil {
		return subject{}, err
	}
	if _, err := os.Lstat(filepath.Join(root, prospectiveGatePath)); err == nil {
		s.Resolution = Resolution{Kind: ProspectiveGateSh}
	} else if !errors.Is(err, os.ErrNotExist) {
		return subject{}, err
	}
	return s, nil
}

func buildProspectiveSubjectForGeneration(root, identityRoot string, generation *treeGeneration) (subject, error) {
	s, err := buildSubjectForGeneration(root, identityRoot, generation)
	if err != nil {
		return subject{}, err
	}
	if _, err := os.Lstat(filepath.Join(root, prospectiveGatePath)); err == nil {
		s.Resolution = Resolution{Kind: ProspectiveGateSh}
	} else if !errors.Is(err, os.ErrNotExist) {
		return subject{}, err
	}
	return s, nil
}

func hashProspectivePreparation(c *identityCollector, identity io.Writer, root, pathEnv string) error {
	path := filepath.Join(root, prospectiveGatePath)
	if _, err := os.Lstat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}
	frame(identity, "prospective preparation")
	return c.hashExecutable(root, path, pathEnv, true, 0)
}
