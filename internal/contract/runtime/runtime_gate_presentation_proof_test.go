package runtime

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/status"
)

func provePresentation(t *testing.T, variant string) {
	switch variant {
	case "html":
		proveSelfContainedDashboard(t)
	case "axi":
		proveAXIGateCache(t)
	case "hostile-status", "hostile-dashboard", "hostile-roadmap":
		proveHostileProjection(t, strings.TrimPrefix(variant, "hostile-"))
	default:
		proveSignalSeverity(t, variant)
	}
}

func proveSignalSeverity(t *testing.T, state string) {
	p := newProjectionFixture(t, state)
	wantSeverity := map[string]int{"red": 0, "locked-pending": 1, "interrupted-pending": 2, "invalid": 3, "unavailable": 3, "stale": 7}
	wantDetail := map[string]string{"red": "red", "locked-pending": "locked-pending", "interrupted-pending": "interrupted-pending", "invalid": "invalid verdict", "unavailable": "verdict unavailable", "stale": "stale"}
	gateSignals := 0
	for _, signal := range status.Signals(p.f.Root) {
		if signal.Name != "gate" {
			continue
		}
		gateSignals++
		if signal.Severity != wantSeverity[state] || !strings.Contains(signal.Detail, wantDetail[state]) {
			t.Fatalf("%s gate signal = %+v, want severity %d and detail %q", state, signal, wantSeverity[state], wantDetail[state])
		}
	}
	show := state != "absent" && state != "reusable-green"
	if (gateSignals == 1) != show {
		t.Fatalf("%s gate signal count = %d, show = %v", state, gateSignals, show)
	}
	statusOut := runProjectionSurface(t, p, "status")
	dashboardOut := runProjectionSurface(t, p, "dashboard")
	if show {
		if !strings.Contains(statusOut, "gate") || !strings.Contains(dashboardOut, wantDetail[state]) {
			t.Fatalf("%s signal missing from a public human format", state)
		}
	} else if strings.Contains(statusOut, "gate       ") || strings.Contains(dashboardOut, "<td>gate</td>") {
		t.Fatalf("%s rendered a gate signal without a signal", state)
	}
}

func proveSelfContainedDashboard(t *testing.T) {
	p := newProjectionFixture(t, "reusable-green")
	page := runProjectionSurface(t, p, "dashboard")
	doc := parseDashboardDocument(t, page)
	if err := dashboardHTMLViolation(doc); err != nil {
		t.Fatal(err)
	}
	for _, literal := range []string{"<!DOCTYPE html>", "<html", "<head>", "<style>", "</style>", "<body>", "</body>", "</html>"} {
		if !strings.Contains(page, literal) {
			t.Fatalf("dashboard HTML missing %q", literal)
		}
	}
	if !strings.HasPrefix(page, "<!DOCTYPE html>") || !strings.HasSuffix(strings.TrimSpace(page), "</html>") {
		t.Fatal("dashboard is not one complete HTML document")
	}
	for name, mutation := range map[string]string{
		"head content":  strings.Replace(page, "<head>", "<head><div>bad</div>", 1),
		"relative href": strings.Replace(page, `<html `, `<html href="/x" `, 1),
		"protocol src":  strings.Replace(page, `<html `, `<html src="//cdn/x" `, 1),
		"data resource": strings.Replace(page, `<html `, `<html data="data:text/plain,x" `, 1),
		"css url":       strings.Replace(page, "</style>", ".x{background:url(/x)}</style>", 1),
		"css import":    strings.Replace(page, "</style>", "@import '/x';</style>", 1),
	} {
		doc := parseDashboardDocument(t, mutation)
		if err := dashboardHTMLViolation(doc); err == nil {
			t.Fatalf("dashboard HTML validator accepted %s", name)
		}
	}
}

type dashboardNode struct {
	name     string
	attrs    map[string]string
	content  []dashboardContent
	children []*dashboardNode
}

type dashboardContent struct {
	text  string
	child *dashboardNode
}

func parseDashboardDocument(t *testing.T, page string) *dashboardNode {
	t.Helper()
	xmlPage := strings.TrimPrefix(page, "<!DOCTYPE html>\n")
	for start := 0; ; {
		i := strings.Index(xmlPage[start:], "<meta ")
		if i < 0 {
			break
		}
		i += start
		j := strings.IndexByte(xmlPage[i:], '>')
		if j < 0 {
			t.Fatal("dashboard meta element is unterminated")
		}
		j += i
		if xmlPage[j-1] != '/' {
			xmlPage = xmlPage[:j] + "/" + xmlPage[j:]
			j++
		}
		start = j + 1
	}
	root := &dashboardNode{name: "#document"}
	stack := []*dashboardNode{root}
	dec := xml.NewDecoder(strings.NewReader(xmlPage))
	dec.Strict = true
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("dashboard HTML is structurally invalid: %v", err)
		}
		switch x := tok.(type) {
		case xml.StartElement:
			n := &dashboardNode{name: x.Name.Local, attrs: map[string]string{}}
			for _, attr := range x.Attr {
				n.attrs[attr.Name.Local] = attr.Value
			}
			parent := stack[len(stack)-1]
			parent.children = append(parent.children, n)
			parent.content = append(parent.content, dashboardContent{child: n})
			stack = append(stack, n)
		case xml.EndElement:
			stack = stack[:len(stack)-1]
		case xml.CharData:
			stack[len(stack)-1].content = append(stack[len(stack)-1].content, dashboardContent{text: string(x)})
		}
	}
	if len(stack) != 1 || len(root.children) != 1 || root.children[0].name != "html" {
		t.Fatalf("dashboard HTML root = %#v", root.children)
	}
	return root.children[0]
}

func (n *dashboardNode) normalizedText() string {
	var parts []string
	for _, part := range n.content {
		if part.child != nil {
			parts = append(parts, part.child.normalizedText())
		} else {
			parts = append(parts, part.text)
		}
	}
	return strings.Join(strings.Fields(strings.Join(parts, " ")), " ")
}

func dashboardNodeShape(n *dashboardNode) string {
	var b strings.Builder
	b.WriteByte('<')
	b.WriteString(n.name)
	keys := make([]string, 0, len(n.attrs))
	for key := range n.attrs {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		b.WriteByte(' ')
		b.WriteString(key + "=" + strconv.Quote(n.attrs[key]))
	}
	b.WriteByte('>')
	for _, part := range n.content {
		if part.child != nil {
			b.WriteString(dashboardNodeShape(part.child))
		} else if text := strings.Join(strings.Fields(part.text), " "); text != "" {
			b.WriteString(strconv.Quote(text))
		}
	}
	b.WriteString("</" + n.name + ">")
	return b.String()
}

var cssResourcePattern = regexp.MustCompile(`(?i)(url\s*\(|@import\b)`)

func dashboardHTMLViolation(doc *dashboardNode) error {
	if len(doc.children) != 2 || doc.children[0].name != "head" || doc.children[1].name != "body" {
		return fmt.Errorf("dashboard html children must be one ordered head/body pair: %#v", doc.children)
	}
	head := doc.children[0]
	if len(head.children) != 4 || head.children[0].name != "meta" || head.children[1].name != "meta" || head.children[2].name != "title" || head.children[3].name != "style" {
		return fmt.Errorf("dashboard head must contain meta/meta/title/style in order")
	}
	allowed := map[string]map[string]bool{
		"html": {"head": true, "body": true}, "head": {"meta": true, "title": true, "style": true},
		"body": {"div": true}, "div": {"header": true, "section": true}, "header": {"h1": true, "p": true},
		"section": {"h2": true, "h3": true, "p": true, "dl": true, "table": true, "pre": true, "ul": true},
		"p":       {"span": true}, "dl": {"dt": true, "dd": true}, "table": {"thead": true, "tbody": true},
		"thead": {"tr": true}, "tbody": {"tr": true}, "tr": {"th": true, "td": true}, "ul": {"li": true},
	}
	resourceAttrs := map[string]bool{"src": true, "href": true, "srcset": true, "action": true, "formaction": true, "poster": true, "data": true, "cite": true, "background": true, "ping": true, "manifest": true}
	var walk func(*dashboardNode) error
	walk = func(n *dashboardNode) error {
		for name, value := range n.attrs {
			if resourceAttrs[strings.ToLower(name)] && strings.TrimSpace(value) != "" {
				return fmt.Errorf("dashboard %s has resource-bearing %s=%q", n.name, name, value)
			}
			if strings.EqualFold(name, "style") && cssResourcePattern.MatchString(value) {
				return fmt.Errorf("dashboard inline style can load a resource: %q", value)
			}
		}
		if n.name == "style" && cssResourcePattern.MatchString(n.normalizedText()) {
			return fmt.Errorf("dashboard stylesheet can load an external resource")
		}
		if n.name == "meta" && strings.EqualFold(n.attrs["http-equiv"], "refresh") {
			return fmt.Errorf("dashboard meta refresh can navigate externally")
		}
		for _, child := range n.children {
			if !allowed[n.name][child.name] {
				return fmt.Errorf("dashboard %s contains invalid %s element", n.name, child.name)
			}
			if err := walk(child); err != nil {
				return err
			}
		}
		return nil
	}
	return walk(doc)
}

func dashboardElements(n *dashboardNode, name string) []*dashboardNode {
	var got []*dashboardNode
	if n.name == name {
		got = append(got, n)
	}
	for _, child := range n.children {
		got = append(got, dashboardElements(child, name)...)
	}
	return got
}

type dashboardDetail struct{ label, value string }
type dashboardGateProjection struct {
	present        bool
	empty          string
	paragraphs     int
	badgeParagraph string
	details        []dashboardDetail
}

func requireExactDashboardGateProjection(t *testing.T, page string, p projectionFixture) {
	t.Helper()
	doc := parseDashboardDocument(t, page)
	var gateSections []*dashboardNode
	for _, section := range dashboardElements(doc, "section") {
		for _, child := range section.children {
			if child.name == "h2" && child.normalizedText() == "Gate" {
				gateSections = append(gateSections, section)
			}
		}
	}
	if len(gateSections) != 1 {
		t.Fatalf("dashboard Gate sections = %d, want 1", len(gateSections))
	}
	section := gateSections[0]
	got := dashboardGateProjection{}
	for _, child := range section.children {
		switch child.name {
		case "h2":
		case "p":
			got.paragraphs++
			if child.attrs["class"] == "empty" {
				if len(child.attrs) != 1 || len(child.children) != 0 {
					t.Fatal("dashboard Gate empty paragraph has extra attributes or descendants")
				}
				got.empty = child.normalizedText()
				continue
			}
			got.present = true
			got.badgeParagraph = dashboardNodeShape(child)
		case "dl":
			for i := 0; i < len(child.children); i += 2 {
				if i+1 >= len(child.children) || child.children[i].name != "dt" || child.children[i+1].name != "dd" {
					t.Fatal("dashboard Gate details are not exact dt/dd pairs")
				}
				got.details = append(got.details, dashboardDetail{child.children[i].normalizedText(), child.children[i+1].normalizedText()})
			}
		default:
			t.Fatalf("dashboard Gate section has unexpected %s element", child.name)
		}
	}
	want := dashboardGateProjection{empty: "No gate cache yet — run bench gate.", paragraphs: 1}
	if p.state != "absent" {
		want = dashboardGateProjection{present: true, paragraphs: 1, details: []dashboardDetail{{"cached tree", p.cachedTree}, {"work tree", p.workTree}}}
		badge := map[string][2]string{
			"reusable-green": {"badge green", "green"}, "red": {"badge red", "red"},
			"locked-pending": {"badge ", "locked-pending"}, "interrupted-pending": {"badge ", "interrupted-pending"},
			"invalid": {"badge ", "invalid"}, "unavailable": {"badge ", "unavailable"},
		}[p.state]
		want.badgeParagraph = expectedBadgeParagraph(badge)
		if p.state == "stale" {
			want.badgeParagraph = expectedBadgeParagraph([2]string{"badge stale", "green"}, [2]string{"badge stale", "stale"})
		}
		if p.state == "invalid" || p.state == "unavailable" {
			want.details[0].value = ""
		}
		needsTimestamp := p.state == "reusable-green" || p.state == "red" || p.state == "stale"
		if needsTimestamp && len(got.details) != 3 {
			t.Fatalf("%s dashboard Gate timestamp details = %#v, want one gated-at field", p.state, got.details)
		}
		if needsTimestamp {
			stamp := strings.SplitN(got.details[2].value, " (", 2)[0]
			if _, err := time.Parse(time.RFC3339, stamp); err != nil {
				t.Fatalf("dashboard gated-at timestamp = %q: %v", got.details[2].value, err)
			}
			want.details = append(want.details, got.details[2])
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s dashboard Gate projection\nwant: %#v\ngot:  %#v", p.state, want, got)
	}
}

func expectedBadgeParagraph(badges ...[2]string) string {
	p := &dashboardNode{name: "p", attrs: map[string]string{}}
	for _, badge := range badges {
		span := &dashboardNode{name: "span", attrs: map[string]string{"class": badge[0]}, content: []dashboardContent{{text: badge[1]}}}
		p.children = append(p.children, span)
		p.content = append(p.content, dashboardContent{child: span})
	}
	return dashboardNodeShape(p)
}

func proveAXIGateCache(t *testing.T) {
	p := newProjectionFixture(t, "locked-pending")
	out := runProjectionSurface(t, p, "roadmap")
	want := "gate_cache[1]{present,state,pending_status,status,cached_tree,work_tree,timestamp,stale}:"
	if strings.Count(out, want) != 1 {
		t.Fatalf("AXI gate-cache schema count = %d, want 1\n%s", strings.Count(out, want), out)
	}
	if !strings.Contains(out, "true,pending,locked-pending") {
		t.Fatalf("AXI gate-cache row lost typed pending fields:\n%s", out)
	}
}

func proveHostileProjection(t *testing.T, surface string) {
	p := newProjectionFixture(t, "invalid")
	hostile := "SECRET-FT78\x1b<script>alert(1)</script>\x00\x7f\n"
	if err := os.WriteFile(p.cache, []byte(hostile), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(p.cache, 0o600); err != nil {
		t.Fatal(err)
	}
	out := runProjectionSurface(t, p, surface)
	for _, raw := range []string{"SECRET-FT78", "<script>", "alert(1)", "\x1b", "\x00", "\x7f"} {
		if strings.Contains(out, raw) {
			t.Fatalf("%s leaked hostile cache bytes %q", surface, raw)
		}
	}
}
