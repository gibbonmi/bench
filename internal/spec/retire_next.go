package spec

import (
	"fmt"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/git"
)

// retireNext renders the last line of a completed retire. The verified delivery of a spec
// closed each roadmap row that it satisfied in its own publication, and a row that it left
// open is residual work. So when the commitment policy at HEAD records that delivery, the
// line schedules no board cleanup. Otherwise it names the board remainder.
func retireNext(base, resolved, roadmapID, slug string) string {
	if deliveredAtHEAD(base, resolved) {
		return fmt.Sprintf(RetireNextPrefix+"promote durable content, commit as `spec-retire: %s`\n", slug)
	}
	return fmt.Sprintf(RetireNextPrefix+"promote durable content, remove the ROADMAP row%s, commit as `spec-retire: %s`\n", roadmapRemainder(base, roadmapID), slug)
}

// deliveredAtHEAD reports whether the commitment policy at HEAD records a verified
// delivery of the spec. An absent or unreadable policy records none.
func deliveredAtHEAD(base, resolved string) bool {
	data, err := headBlob(base, commitment.PolicyPath)
	if err != nil {
		return false
	}
	policy, err := commitment.Parse(data)
	return err == nil && commitment.PathDelivered(policy, filepath.ToSlash(RelTo(base, resolved)))
}

// headBlob reads the repository-relative path at HEAD through git.
func headBlob(base, rel string) ([]byte, error) {
	args := []string{"show", "HEAD:" + rel}
	if base != "" {
		args = append([]string{"-C", base}, args...)
	}
	return git.Raw(args...)
}
