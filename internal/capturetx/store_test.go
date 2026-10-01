package capturetx

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
)

func TestAppendSerializesConcurrentWriters(t *testing.T) {
	root := gittest.Repo(t)
	path := filepath.Join(root, "capture", "IDEAS.md")
	source := Source{Name: "capture/IDEAS.md", Path: path}
	const writers = 20
	var group sync.WaitGroup
	errs := make(chan error, writers)
	for i := 0; i < writers; i++ {
		group.Add(1)
		go func(i int) {
			defer group.Done()
			line := fmt.Sprintf("entry-%02d\n", i)
			errs <- Append(root, source, func(current []byte) ([]byte, error) {
				return append(append([]byte(nil), current...), line...), nil
			})
		}(i)
	}
	group.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < writers; i++ {
		if count := strings.Count(string(data), fmt.Sprintf("entry-%02d\n", i)); count != 1 {
			t.Fatalf("entry %d occurs %d times in %q", i, count, data)
		}
	}
}
