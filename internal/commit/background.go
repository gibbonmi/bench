package commit

import (
	"io"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/gate"
)

// startBackgroundGrade starts the background grade of published, the commit that this call
// published in root. A call that published nothing starts none: a dry run and a refusal each
// leave published empty.
func startBackgroundGrade(root, published string, stderr io.Writer) {
	if published != "" {
		gate.StartBackground(benchhome.Dir(), gate.KitValue(), root, published, stderr)
	}
}
