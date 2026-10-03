package sessioninspect

import (
	"context"
	"fmt"
	"io"

	"github.com/gibbonmi/bench/internal/compatibility"
)

func compatibilityInspection(request compatibility.SessionRequest) (string, error) {
	return compatibility.SessionReport(request).Render()
}

func compatibilityPhase(ctx context.Context, stdout, stderr io.Writer, _ string) int {
	if ctx.Err() != nil {
		return 0
	}
	fmt.Fprintln(stdout, "bench: on start or resume, check the active interface through actual chat tools with normal permissions before dependent work. Read Session compatibility in .bench/BENCH-reference.md.")
	output, err := compatibilityInspection(compatibility.SessionRequest{Operation: "diagnose"})
	if err != nil {
		fmt.Fprintln(stderr, "warning: compatibility inspection unavailable; live evidence remains unknown")
		return 0
	}
	fmt.Fprint(stdout, output)
	fmt.Fprintln(stdout, compatibility.CapabilityAction("failed-interface-retest"))
	recovery, err := compatibility.RenderRecovery()
	if err != nil {
		fmt.Fprintln(stderr, "warning: compatibility recovery guidance unavailable")
		return 0
	}
	fmt.Fprint(stdout, recovery)
	return 0
}
