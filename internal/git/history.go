package git

import "strings"

// FirstParentsAfter lists the commits on the first-parent chain of destination that the
// marker commit cannot reach, oldest first. When marker is an ancestor of destination, the
// list is the first-parent commits after marker up to destination, and it is empty when the
// two are one commit.
func FirstParentsAfter(root, marker, destination string) ([]string, error) {
	out, err := Output("-C", root, "rev-list", "--first-parent", "--reverse", destination, "^"+marker)
	if err != nil {
		return nil, err
	}
	return strings.Fields(out), nil
}
