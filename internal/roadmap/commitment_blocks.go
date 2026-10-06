package roadmap

import "github.com/gibbonmi/bench/internal/commitment"

// tableBlock is one TOON table of a roadmap response.
type tableBlock struct {
	name   string
	fields []string
	rows   [][]any
}

// OutlookSource returns the commitment outlook for a repository root. The command
// adapter supplies the one projection that status and the dashboard also render, so the
// roadmap reads no policy and decides no eligibility itself.
type OutlookSource func(root string) commitment.Outlook

// commitmentBlocks renders the outlook table and its blocker table. Before adoption the
// outlook state is adoption-required, which marks the recommended sequence as unapproved
// input rather than a commitment.
func commitmentBlocks(outlook commitment.Outlook) []tableBlock {
	return []tableBlock{
		{commitment.OutlookTable, commitment.OutlookFields, stringRows([][]string{outlook.Cells()})},
		{commitment.BlockerTable, commitment.BlockerFields, stringRows(outlook.BlockerCells())},
	}
}
