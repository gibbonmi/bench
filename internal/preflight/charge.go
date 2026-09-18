package preflight

import (
	"bytes"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/chargeevidence"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/tickets"
	"github.com/gibbonmi/bench/internal/toon"
)

// buildChargePack applies every build charge refusal in its fixed order and returns the
// validated in-memory pack, or the refusal that stopped it.
func buildChargePack(root string, facts Facts, verdict Verdict, name string, policy []buildSourceDescriptor) (*chargeevidence.Pack, string) {
	if refusal := preparationCheckoutRefusal(root, facts, "charge"); refusal != "" {
		return nil, refusal
	}
	if verdict.Red {
		return nil, chargeVerdictRefusal(verdict)
	}
	selected, parsed, detail, selectionNext := preparationTicket(root, facts, name)
	if detail != "" {
		return nil, chargeRefusal("ticket", detail, selectionNext)
	}
	pack, failure, err := prepareBuildPack(root, facts, selected, parsed, policy)
	if failure != "" {
		return nil, chargeRefusal("source", failure, "restore the named canonical source and rerun the exact charge")
	}
	if err != nil {
		return nil, chargeRefusal("evidence", err.Error(), "repair the prepared evidence input and rerun the exact charge")
	}
	return pack, ""
}

// loadChargeSource reads one required source at the pinned tip and refuses bytes the
// shared TOON adapter cannot represent.
func loadChargeSource(root, sourceTip, path string) ([]byte, string) {
	data, failure := readChargeSource(root, sourceTip, path)
	if failure != "" {
		return nil, failure
	}
	if !toon.Representable(string(data)) {
		return nil, path + " contains a byte spec-TOON cannot represent"
	}
	return data, ""
}

func selectedTicket(entries []tickets.Entry, name string) *tickets.Entry {
	for i := range entries {
		if entries[i].Name == name {
			return &entries[i]
		}
	}
	return nil
}

func selectedParsedTicket(facts Facts, name string) *tickets.Ticket {
	for i := range facts.Tickets {
		if facts.Tickets[i].Name == name {
			return &facts.Tickets[i]
		}
	}
	return nil
}

func readChargeSource(root, sourceTip, rel string) ([]byte, string) {
	read := bounds.ClassifyNoFollow(filepath.Join(root, filepath.FromSlash(rel)))
	if read.State != bounds.StateParsed {
		return nil, rel + " is " + string(read.State) + ": " + read.Reason
	}
	pinned, err := git.Raw("-C", root, "show", sourceTip+":"+rel)
	if err != nil {
		return nil, rel + " is absent or unreadable at source tip " + sourceTip
	}
	if !bytes.Equal(read.Data, pinned) {
		return nil, rel + " does not match source tip " + sourceTip
	}
	return pinned, ""
}

func chargeVerdictRefusal(verdict Verdict) string {
	for _, check := range verdict.Checks {
		if check.Verdict == verdictRed {
			next := check.Next
			if next == "" {
				next = "repair " + check.Check + " and rerun the exact charge"
			}
			return chargeRefusal("preflight", check.Check+": "+check.Detail, next)
		}
	}
	return chargeRefusal("preflight", "required preflight checks are red", "repair the reported check and rerun the exact charge")
}

func rows(values []string) [][]string {
	result := make([][]string, len(values))
	for i, value := range values {
		result[i] = []string{value}
	}
	return result
}

func chargeRefusal(input, detail, next string) string {
	return toon.Errorf(input+" required: "+detail, next) + "\n"
}
