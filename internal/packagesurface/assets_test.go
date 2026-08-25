package packagesurface

import "testing"

func TestRequiredPackAssetsIncludeFollowOnGuard(t *testing.T) {
	const want = ".bench/hooks/block-bench-follow-on.sh"
	for _, asset := range RequiredPackAssets {
		if asset == want {
			return
		}
	}
	t.Fatalf("RequiredPackAssets omit %q", want)
}
