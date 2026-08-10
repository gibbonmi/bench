// Package freshness verifies that a Bench executable was built from the current sources.
package freshness

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/subprocess"
)

const sealSchema = 1
const auxiliaryInputsManifest = "scripts/go-build.inputs"

// publicationBackupPattern names the copies a publication keeps of the pair it replaces.
const publicationBackupPattern = ".bench-publish-backup-*"

// sealTemporaryPattern names the staging file a seal write promotes into path.
func sealTemporaryPattern(path string) string { return filepath.Base(path) + ".tmp-*" }

var replacePublicationFile = os.Rename

type seal struct {
	Schema     int    `json:"schema"`
	Sources    string `json:"sources"`
	Executable string `json:"executable"`
}

// Digest returns the deterministic content digest of Bench's local build inputs.
func Digest(root string) (string, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	paths, err := buildInputs(root)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	for _, path := range paths {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return "", err
		}
		contents, err := regularContents(path)
		if err != nil {
			return "", fmt.Errorf("read build input %q: %w", rel, err)
		}
		name := filepath.ToSlash(rel)
		fmt.Fprintf(hash, "%d:%s%d:", len(name), name, len(contents))
		if _, err := hash.Write(contents); err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// BuildInputs returns the repository-relative, slash-separated, sorted paths that Digest hashes for root.
func BuildInputs(root string) ([]string, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	paths, err := buildInputs(root)
	if err != nil {
		return nil, err
	}
	relative := make([]string, len(paths))
	for i, path := range paths {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil, err
		}
		relative[i] = filepath.ToSlash(rel)
	}
	return relative, nil
}

// SealDigests returns the source and executable digests recorded in executable's published seal.
func SealDigests(executable string) (sources, digest string, err error) {
	data, err := secureContents(sealPath(executable), false)
	if err != nil {
		return "", "", fmt.Errorf("seal %q: %w", sealPath(executable), err)
	}
	stored, err := parseSeal(data)
	if err != nil {
		return "", "", fmt.Errorf("seal %q: %w", sealPath(executable), err)
	}
	return stored.Sources, stored.Executable, nil
}

// ExecutableDigest returns the digest a seal records for the executable at path.
func ExecutableDigest(path string) (string, error) {
	binary, err := secureContents(path, true)
	if err != nil {
		return "", err
	}
	return digestBytes(binary), nil
}

// Publish replaces executable with staged and writes its matching content seal.
func Publish(root, staged, executable string) error {
	if err := publicationTarget(executable); err != nil {
		return fmt.Errorf("publish executable %q: %w", executable, err)
	}
	if err := publicationTarget(sealPath(executable)); err != nil {
		return fmt.Errorf("publish seal %q: %w", sealPath(executable), err)
	}
	encoded, err := sealContents(root, staged)
	if err != nil {
		return err
	}
	transaction, err := beginPublication(executable)
	if err != nil {
		return err
	}
	defer transaction.close()
	if err := transaction.install(staged); err != nil {
		return fmt.Errorf("publish executable: %w", err)
	}
	if err := transaction.sealWith(encoded); err != nil {
		if restoreErr := transaction.rollback(); restoreErr != nil {
			return fmt.Errorf("publish seal: %w; restore prior pair: %v", err, restoreErr)
		}
		return fmt.Errorf("publish seal: %w", err)
	}
	return nil
}

// sealContents encodes the seal that describes staged as built from root's sources.
func sealContents(root, staged string) ([]byte, error) {
	sources, err := Digest(root)
	if err != nil {
		return nil, err
	}
	binary, err := ExecutableDigest(staged)
	if err != nil {
		return nil, fmt.Errorf("stage executable %q: %w", staged, err)
	}
	return json.Marshal(seal{
		Schema:     sealSchema,
		Sources:    sources,
		Executable: binary,
	})
}

// publicationStepGrace bounds how long an arriving termination waits for a publication
// step to finish before restoring anyway. A step is one rename or one small write, so a
// step still running after the grace means the filesystem is wedged, and honoring the
// termination then matters more than waiting for a step that may never return.
const publicationStepGrace = 2 * time.Second

// publication owns replacing an executable and its seal as one outcome. The pair is only
// consistent before the executable moves and after the seal lands, so the transaction
// both restores the prior pair when the seal fails and answers a termination arriving in
// between — the invoking shell holds no rollback state and needs none.
type publication struct {
	executable    string
	backup        string
	hadExecutable bool
	sealBackup    string
	hadSeal       bool
	// step serializes the two renames against the termination restore, so a signal lands
	// strictly between steps rather than inside one.
	step     sync.Mutex
	resolved atomic.Bool
	signals  chan os.Signal
	// sealTemporary holds the staging file of a seal write that has not yet landed. The
	// termination handler reads it while that write is still in flight, so the name lives
	// on the transaction rather than only in the writer's own frame.
	sealTemporary atomic.Pointer[string]
}

func beginPublication(executable string) (*publication, error) {
	transaction := &publication{executable: executable}
	backup, hadExecutable, err := publicationBackup(executable, true)
	if err != nil {
		return nil, fmt.Errorf("backup executable: %w", err)
	}
	transaction.backup, transaction.hadExecutable = backup, hadExecutable
	sealBackup, hadSeal, err := publicationBackup(sealPath(executable), false)
	if err != nil {
		transaction.remove(backup)
		return nil, fmt.Errorf("backup seal: %w", err)
	}
	transaction.sealBackup, transaction.hadSeal = sealBackup, hadSeal
	transaction.watch()
	return transaction, nil
}

// watch arms the transaction against the terminations an operator or a supervisor sends.
// The restore has to run before the process ends, and a process cannot honor a signal and
// keep running as if it had not arrived, so the handler exits under the shell convention
// for the signal it received.
func (p *publication) watch() {
	p.signals = make(chan os.Signal, 1)
	signal.Notify(p.signals, subprocess.CancelSignals...)
	go func() {
		received, delivered := <-p.signals
		if !delivered {
			return
		}
		held := p.awaitStep()
		if !p.resolved.Swap(true) {
			_ = p.restore()
		}
		p.removeTemporaries()
		if held {
			p.step.Unlock()
		}
		publicationExit(received)
	}()
}

func (p *publication) awaitStep() bool {
	deadline := time.Now().Add(publicationStepGrace)
	for {
		if p.step.TryLock() {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func (p *publication) install(staged string) error {
	p.step.Lock()
	defer p.step.Unlock()
	return replacePublicationFile(staged, p.executable)
}

func (p *publication) sealWith(encoded []byte) error {
	p.step.Lock()
	defer p.step.Unlock()
	if err := writeSeal(sealPath(p.executable), encoded, p.trackSealTemporary); err != nil {
		return err
	}
	// Resolving under the same lock the restore takes is what stops a termination
	// arriving on the heels of a landed seal from undoing a complete publication.
	p.resolved.Store(true)
	return nil
}

func (p *publication) rollback() error {
	p.step.Lock()
	defer p.step.Unlock()
	if p.resolved.Swap(true) {
		return nil
	}
	return p.restore()
}

func (p *publication) restore() error {
	if err := restorePriorFile(p.backup, p.executable, p.hadExecutable); err != nil {
		return err
	}
	return restorePriorFile(p.sealBackup, sealPath(p.executable), p.hadSeal)
}

func (p *publication) close() {
	signal.Stop(p.signals)
	close(p.signals)
	p.removeTemporaries()
}

// removeTemporaries deletes every temporary the transaction owns. Both the termination
// handler and close call it, and either may run after the other, so a removal of what is
// already gone is the normal case rather than a failure.
func (p *publication) removeTemporaries() {
	p.remove(p.backup)
	p.remove(p.sealBackup)
	if temporary := p.sealTemporary.Load(); temporary != nil {
		p.remove(*temporary)
	}
}

func (p *publication) trackSealTemporary(path string) {
	if path == "" {
		p.sealTemporary.Store(nil)
		return
	}
	p.sealTemporary.Store(&path)
}

func (p *publication) remove(backup string) {
	if backup != "" {
		_ = os.Remove(backup)
	}
}

// restorePriorFile puts back what publication found at path: the backed-up bytes, or
// nothing at all where there was nothing before.
func restorePriorFile(backup, path string, existed bool) error {
	if existed {
		return os.Rename(backup, path)
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func publicationExit(received os.Signal) {
	code := 128
	if signaled, ok := received.(syscall.Signal); ok {
		code += int(signaled)
	}
	os.Exit(code)
}

func publicationBackup(path string, executable bool) (string, bool, error) {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	contents, err := secureContents(path, executable)
	if err != nil {
		return "", false, err
	}
	backup, err := os.CreateTemp(filepath.Dir(path), publicationBackupPattern)
	if err != nil {
		return "", false, err
	}
	name := backup.Name()
	if _, err := backup.Write(contents); err != nil {
		_ = backup.Close()
		_ = os.Remove(name)
		return "", false, err
	}
	if err := backup.Chmod(info.Mode().Perm()); err != nil {
		_ = backup.Close()
		_ = os.Remove(name)
		return "", false, err
	}
	if err := backup.Close(); err != nil {
		_ = os.Remove(name)
		return "", false, err
	}
	return name, true, nil
}

func publicationTarget(path string) error {
	if err := rejectSymlinkComponents(filepath.Dir(path)); err != nil {
		return err
	}
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("is a symbolic link")
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("is not a regular file")
	}
	return nil
}

type listedPackage struct {
	Dir          string
	GoFiles      []string
	CgoFiles     []string
	CFiles       []string
	CXXFiles     []string
	MFiles       []string
	HFiles       []string
	FFiles       []string
	SFiles       []string
	SwigFiles    []string
	SwigCXXFiles []string
	SysoFiles    []string
	EmbedFiles   []string
}

func buildInputs(root string) ([]string, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	command := exec.Command("go", "list", "-buildvcs=false", "-json", "-deps", "./cmd/bench")
	command.Dir = root
	output, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("resolve Bench build inputs: %w", err)
	}
	decoder := json.NewDecoder(strings.NewReader(string(output)))
	paths := map[string]struct{}{}
	for {
		var pkg listedPackage
		if err := decoder.Decode(&pkg); err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("decode resolved Bench build inputs: %w", err)
		}
		for _, name := range packageFiles(pkg) {
			path := filepath.Join(pkg.Dir, name)
			if isWithinRoot(root, path) {
				paths[path] = struct{}{}
			}
		}
	}
	for _, name := range []string{
		"go.mod",
	} {
		path := filepath.Join(root, filepath.FromSlash(name))
		if _, err := os.Lstat(path); err != nil {
			return nil, fmt.Errorf("required build input %q: %w", name, err)
		}
		paths[path] = struct{}{}
	}
	auxiliary, err := auxiliaryBuildInputs(root)
	if err != nil {
		return nil, err
	}
	for _, path := range auxiliary {
		paths[path] = struct{}{}
	}
	if path := filepath.Join(root, "go.sum"); exists(path) {
		paths[path] = struct{}{}
	}
	ordered := make([]string, 0, len(paths))
	for path := range paths {
		ordered = append(ordered, path)
	}
	sort.Slice(ordered, func(i, j int) bool {
		return filepath.ToSlash(ordered[i]) < filepath.ToSlash(ordered[j])
	})
	return ordered, nil
}

func auxiliaryBuildInputs(root string) ([]string, error) {
	manifest := filepath.Join(root, filepath.FromSlash(auxiliaryInputsManifest))
	data, err := regularContents(manifest)
	if err != nil {
		return nil, fmt.Errorf("read auxiliary build-input manifest %q: %w", auxiliaryInputsManifest, err)
	}
	paths := []string{manifest}
	keys := map[string]struct{}{}
	for number, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		key, name, ok := strings.Cut(line, "=")
		if !ok || key == "" || name == "" || strings.TrimSpace(key) != key || strings.TrimSpace(name) != name {
			return nil, fmt.Errorf("malformed auxiliary build input at %s:%d", auxiliaryInputsManifest, number+1)
		}
		if _, exists := keys[key]; exists {
			return nil, fmt.Errorf("duplicate auxiliary build input key %q", key)
		}
		keys[key] = struct{}{}
		path := filepath.Join(root, filepath.FromSlash(name))
		if !isWithinRoot(root, path) || filepath.IsAbs(name) {
			return nil, fmt.Errorf("auxiliary build input %q leaves the source root", name)
		}
		if _, err := os.Lstat(path); err != nil {
			return nil, fmt.Errorf("required build input %q: %w", name, err)
		}
		paths = append(paths, path)
	}
	if len(paths) == 1 {
		return nil, fmt.Errorf("auxiliary build-input manifest %q is empty", auxiliaryInputsManifest)
	}
	return paths, nil
}

func packageFiles(pkg listedPackage) []string {
	var files []string
	for _, group := range [][]string{
		pkg.GoFiles, pkg.CgoFiles, pkg.CFiles, pkg.CXXFiles, pkg.MFiles, pkg.HFiles,
		pkg.FFiles, pkg.SFiles, pkg.SwigFiles, pkg.SwigCXXFiles, pkg.SysoFiles, pkg.EmbedFiles,
	} {
		files = append(files, group...)
	}
	return files
}

func isWithinRoot(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func exists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

func regularContents(path string) ([]byte, error) {
	return secureContents(path, false)
}

func secureContents(path string, executable bool) ([]byte, error) {
	if err := rejectSymlinkComponents(path); err != nil {
		return nil, err
	}
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("is a symbolic link")
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("is not a regular file")
	}
	if executable && info.Mode()&0o111 == 0 {
		return nil, fmt.Errorf("is not executable")
	}
	if executable && info.Size() == 0 {
		return nil, fmt.Errorf("is empty")
	}
	// A path replacement after Lstat stays untrusted instead of redirecting this read.
	fd, err := syscall.Open(path, syscall.O_RDONLY|syscall.O_CLOEXEC|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil, err
	}
	file := os.NewFile(uintptr(fd), path)
	defer file.Close()
	opened, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if !os.SameFile(info, opened) {
		return nil, fmt.Errorf("changed while opening")
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if executable && len(data) == 0 {
		return nil, fmt.Errorf("is empty")
	}
	return data, nil
}

func rejectSymlinkComponents(path string) error {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	volume := filepath.VolumeName(absolute)
	current := volume + string(filepath.Separator)
	relative := strings.TrimPrefix(absolute, current)
	for _, component := range strings.Split(relative, string(filepath.Separator)) {
		if component == "" {
			continue
		}
		current = filepath.Join(current, component)
		info, err := os.Lstat(current)
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("path component %q is a symbolic link", current)
		}
	}
	return nil
}

// Verify reports whether executable has a matching content seal for root.
func Verify(root, executable string) error {
	stored, err := verifiedExecutable(executable)
	if err != nil {
		return refusal(root, executable, err)
	}
	sources, err := Digest(root)
	if err != nil {
		return refusal(root, executable, err)
	}
	decision := Select(SelectionInput{
		StoredSource: stored.Sources, CurrentSource: sources,
		StoredExecutable: stored.Executable, CurrentExecutable: stored.Executable,
	})
	if !decision.Accepted {
		return refusal(root, executable, errors.New(decision.Reason))
	}
	return nil
}

// VerifyExecutable checks that executable and its adjacent seal are an intact published
// pair without rediscovering source inputs. A run owner has already verified the source
// digest; descendants use this narrower check so inherited selection cannot mutate the
// subject through Go's cache discovery.
func VerifyExecutable(executable string) error {
	_, err := verifiedExecutable(executable)
	return err
}

func verifiedExecutable(executable string) (seal, error) {
	binary, err := secureContents(executable, true)
	if err != nil {
		return seal{}, err
	}
	sealData, err := secureContents(sealPath(executable), false)
	if err != nil {
		return seal{}, fmt.Errorf("seal %q: %w", sealPath(executable), err)
	}
	stored, err := parseSeal(sealData)
	if err != nil {
		return seal{}, fmt.Errorf("seal %q: %w", sealPath(executable), err)
	}
	decision := Select(SelectionInput{
		StoredSource: stored.Sources, CurrentSource: stored.Sources,
		StoredExecutable: stored.Executable, CurrentExecutable: digestBytes(binary),
	})
	if !decision.Accepted {
		return seal{}, errors.New(decision.Reason)
	}
	return stored, nil
}

// Check verifies an executable from current sources, then requires its freshness subcommand.
func Check(root, executable string) error {
	if err := Verify(root, executable); err != nil {
		return err
	}
	command := exec.Command(executable, "freshness-check", root)
	command.Stdout = io.Discard
	command.Stderr = io.Discard
	if err := command.Run(); err != nil {
		return refusal(root, executable, fmt.Errorf("freshness-check failed"))
	}
	return nil
}

func parseSeal(data []byte) (seal, error) {
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	var stored seal
	if err := decoder.Decode(&stored); err != nil {
		return seal{}, fmt.Errorf("malformed: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return seal{}, fmt.Errorf("malformed trailing data")
	}
	if stored.Schema != sealSchema || !isDigest(stored.Sources) || !isDigest(stored.Executable) {
		return seal{}, fmt.Errorf("malformed contents")
	}
	return stored, nil
}

// writeSeal promotes data into path through a staging file, reporting that file's name to
// track for as long as it is in flight: a termination answered mid-write ends the process
// before any deferred cleanup here can run, so the staging file needs an owner that
// outlives this call.
func writeSeal(path string, data []byte, track func(string)) (err error) {
	temporary, err := os.CreateTemp(filepath.Dir(path), sealTemporaryPattern(path))
	if err != nil {
		return err
	}
	track(temporary.Name())
	defer func() {
		if err != nil {
			_ = os.Remove(temporary.Name())
		}
		track("")
	}()
	if _, err := temporary.Write(data); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Chmod(0o644); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return replacePublicationFile(temporary.Name(), path)
}

func sealPath(executable string) string { return executable + ".seal" }

func digestBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func isDigest(value string) bool {
	if len(value) != sha256.Size*2 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

func refusal(root, executable string, cause error) error {
	return fmt.Errorf("bench binary %q is untrusted: %v; rebuild with %s", executable, cause, RebuildAction(root))
}

// RebuildAction returns the copy-paste command that republishes root's Bench binary.
func RebuildAction(root string) string {
	return fmt.Sprintf("cd %s && bash scripts/go-build.sh %s %s", shellQuote(root), shellQuote(root), shellQuote(filepath.Join(root, "dist", "bench")))
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}
