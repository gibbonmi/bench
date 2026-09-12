package preflight

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestChargePreservesUnicodeTicketWithoutFinalNewline(t *testing.T) {
	root, slug := seedConformant(t)
	ticket := strings.Replace(ticketDoc("One", "PF1", "PF2"), "Writes: specs", "Writes: specs, internal/example", 1)
	ticket = strings.TrimSuffix(ticket, "\n") + "\n\nRésumé 雪"
	mustWriteFile(t, "specs/"+slug+"/tickets/one.md", ticket)
	runGit(t, "add", "specs/"+slug+"/tickets/one.md")
	runGit(t, "commit", "-q", "-m", "unicode ticket without final newline")

	out, code := Command(chargeArgs(t, root, slug, true))
	wantIdentity := (chargeSource{path: "specs/" + slug + "/tickets/one.md", data: []byte(ticket)}).identity()
	if code != 0 || !strings.Contains(out, "Résumé 雪\"") ||
		!strings.Contains(out, wantIdentity) || !strings.Contains(out, "\"specs, internal/example\"") ||
		!strings.Contains(out, "PF1") || !strings.Contains(out, "PF2") ||
		!strings.Contains(out, "does x") || !strings.Contains(out, "cli seam") ||
		!strings.Contains(out, "catches z") {
		t.Fatalf("unicode ticket without final newline = (%d):\n%s", code, out)
	}
}

func TestChargeTextSinksUseTOONEscapingAndRefusal(t *testing.T) {
	t.Run("permitted controls", func(t *testing.T) {
		root, slug := seedConformant(t)
		ticket := ticketDoc("One", "PF1", "PF2") + "marker:\tvalue\nreturn:\rvalue\n"
		mustWriteFile(t, "specs/"+slug+"/tickets/one.md", ticket)
		runGit(t, "add", "specs/"+slug+"/tickets/one.md")
		runGit(t, "commit", "-q", "-m", "ticket with permitted controls")
		out, code := Command(chargeArgs(t, root, slug, true))
		if code != 0 || !strings.Contains(out, "marker:\\tvalue\\nreturn:\\rvalue\\n") {
			t.Fatalf("permitted controls = (%d):\n%s", code, out)
		}
	})
	for _, control := range []struct {
		name string
		byte byte
	}{{"ESC", 0x1b}, {"BEL", 0x07}} {
		t.Run(control.name, func(t *testing.T) {
			root, slug := seedConformant(t)
			mustWriteFile(t, buildPhase, "unsafe "+string(control.byte)+" source\n")
			runGit(t, "add", buildPhase)
			runGit(t, "commit", "-q", "-m", "source with "+control.name)
			out, code := Command(chargeArgs(t, root, slug, true))
			if code != 1 || !strings.Contains(out, "source required") ||
				!strings.Contains(out, buildPhase) ||
				!strings.Contains(out, "cannot represent") || strings.Contains(out, "complete,next}") {
				t.Fatalf("%s source = (%d):\n%s", control.name, code, out)
			}
		})
	}
}

func TestChargeQuotesNumericLookingSourceTip(t *testing.T) {
	root, slug := seedConformant(t)
	t.Setenv("GIT_COMMITTER_DATE", "2026-01-02T03:04:05Z")
	tip := ""
	for i := 0; i < 512; i++ {
		runGit(t, "commit", "--amend", "-q", "-m", "numeric tip "+strconv.Itoa(i))
		tip = runGit(t, "rev-parse", "HEAD")
		if tip[0] == '0' && tip[1] >= '0' && tip[1] <= '9' {
			break
		}
	}
	if tip[0] != '0' || tip[1] < '0' || tip[1] > '9' {
		t.Fatalf("fixture did not produce a numeric-looking tip: %s", tip)
	}
	out, code := Command(chargeArgs(t, root, slug, false))
	if code != 0 || !strings.Contains(out, ",\""+tip+"\",") {
		t.Fatalf("numeric-looking source tip = (%d, %s):\n%s", code, tip, out)
	}
}

func TestChargeHostileTicketPathRemainsData(t *testing.T) {
	root, slug := seedConformant(t)
	name := "one [*] $(touch sentinel).md"
	rel := "specs/" + slug + "/tickets/dir /" + name
	mustWriteFile(t, "specs/"+slug+"/spec.md", specBody(slug, "- `specs/"+slug+"/` (ticket fixtures)"))
	mustWriteFile(t, rel, ticketDoc("Hostile path", "PF1", "PF2"))
	runGit(t, "add", "specs/"+slug)
	runGit(t, "commit", "-q", "-m", "hostile ticket path")
	args := chargeArgs(t, root, slug, false)
	args[4] = name
	out, code := Command(args)
	if code != 0 || !strings.Contains(out, rel) ||
		!strings.Contains(out, "'"+name+"'") || strings.Contains(out, "complete,true") {
		t.Fatalf("hostile ticket path = (%d):\n%s", code, out)
	}
	if _, err := os.Stat(filepath.Join(root, "sentinel")); !os.IsNotExist(err) {
		t.Fatalf("hostile ticket path created sentinel: %v", err)
	}
}
