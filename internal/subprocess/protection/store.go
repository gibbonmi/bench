package protection

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync/atomic"
	"syscall"

	"github.com/gibbonmi/bench/internal/durablefile"
)

// ErrSealed refuses a new user of a generation whose admission a cleaner closed.
var ErrSealed = errors.New("protection: the resource generation is sealed")

// Validator checks a descriptor's recognized store location and existing owner identity.
type Validator func(Descriptor) error

// Validators is the kind-validator table that the CLI composition root supplies.
type Validators map[Kind]Validator

// Protocol holds the composition root's validators and the caller's record read limit. Its
// unset persistence and probe are the operating system's.
type Protocol struct {
	Validators Validators
	ReadLimit  int64
	replacer   durablefile.Replacer
	probe      Probe
}

var installed atomic.Pointer[Protocol]

// Install makes p the protocol of this process's launch adapters.
func Install(p Protocol) { installed.Store(&p) }

// Installed returns the installed protocol, or nil before Install.
func Installed() *Protocol { return installed.Load() }

func (p *Protocol) replace() durablefile.Replacer {
	if p.replacer.Rename == nil {
		return durablefile.Native()
	}
	return p.replacer
}

func (p *Protocol) observer() Probe {
	if p.probe == nil {
		return syscall.Kill
	}
	return p.probe
}

// owner runs the kind's validator. A kind without one refuses, so no table row means no path.
func (p *Protocol) owner(d Descriptor) error {
	var validate Validator
	if p != nil {
		validate = p.Validators[d.Kind]
	}
	if validate == nil {
		return fmt.Errorf("no %s validator is installed, so %s gains no admission", d.Kind, d.ResourcePath)
	}
	return validate(d)
}

// OpenResource validates the owner's identity and joins the store's open generation, or
// publishes a fresh one when none exists. It never rewrites a generation or its users. A
// failed call deletes a store that it created, unless a generation was published there.
func (p *Protocol) OpenResource(d Descriptor) (Descriptor, error) {
	d.ResourceID = ""
	if err := d.checkOwner(); err != nil {
		return Descriptor{}, err
	}
	if err := p.owner(d); err != nil {
		return Descriptor{}, err
	}
	created := false
	switch err := p.makeDir(d.Store); {
	case err == nil:
		created = true
	case !errors.Is(err, fs.ErrExist):
		return Descriptor{}, err
	}
	opened, published, err := p.join(d)
	if err != nil && created && !published {
		err = errors.Join(err, os.RemoveAll(d.Store))
	}
	return opened, err
}

// Reopen is OpenResource for a store whose generation can be sealed. An owner release, a reuse,
// a removal, or a clean seals the generation before it grades the users, so a refused one
// leaves the generation sealed, and OpenResource refuses a sealed generation. Reopen recovers
// the store, which retires the sealed generation when every user is safe, and then opens a
// fresh generation.
func (p *Protocol) Reopen(d Descriptor) (Descriptor, error) {
	opened, err := p.OpenResource(d)
	if errors.Is(err, ErrSealed) {
		if _, err = p.Recover(d.Store, nil); err == nil {
			opened, err = p.OpenResource(d)
		}
	}
	return opened, err
}

// join takes the store's lock and joins or publishes its generation. Published reports whether
// the store holds a generation, so a failed caller knows whether a user can depend on the store.
func (p *Protocol) join(d Descriptor) (_ Descriptor, published bool, _ error) {
	release, err := lock([]string{d.Store})
	if err != nil {
		return Descriptor{}, false, err
	}
	defer release()
	users := filepath.Join(d.Store, usersDir)
	if err := p.makeDir(users); err != nil && !errors.Is(err, fs.ErrExist) {
		return Descriptor{}, false, err
	}
	if err := checkDir(users); err != nil {
		return Descriptor{}, false, err
	}
	var r resourceRecord
	path := filepath.Join(d.Store, resourceFile)
	switch err := readRecord(path, p.ReadLimit, &r); {
	case errors.Is(err, fs.ErrNotExist):
		r = resourceRecord{Schema: schema, ResourceID: newID(), Kind: d.Kind, ResourcePath: d.ResourcePath, Binding: d.Binding, OwnerPID: os.Getpid(), Admission: admissionOpen, UsersDir: usersDir}
		if err := p.publish(path, r); err != nil {
			return Descriptor{}, false, err
		}
	case err != nil:
		return Descriptor{}, true, err
	default:
		if err := r.valid(); err != nil {
			return Descriptor{}, true, err
		}
	}
	d.ResourceID = r.ResourceID
	if err := r.matches(d); err != nil {
		return Descriptor{}, true, err
	}
	if r.Admission == admissionSealed {
		return Descriptor{}, true, ErrSealed
	}
	return d, true, nil
}

// view is one store's strictly read generation and users.
type view struct {
	resource resourceRecord
	users    map[string]userRecord
}

// read strictly reads the store of d. Any unreadable or invalid record is an error, because it
// retains the resource and refuses admission.
func (p *Protocol) read(d Descriptor) (view, error) {
	v := view{users: map[string]userRecord{}}
	if err := readRecord(filepath.Join(d.Store, resourceFile), p.ReadLimit, &v.resource); err != nil {
		return v, err
	}
	if err := v.resource.valid(); err != nil {
		return v, err
	}
	if err := v.resource.matches(d); err != nil {
		return v, err
	}
	users := filepath.Join(d.Store, usersDir)
	if err := checkDir(users); err != nil {
		return v, err
	}
	entries, err := os.ReadDir(users)
	if err != nil {
		return v, err
	}
	for _, entry := range entries {
		dir := filepath.Join(users, entry.Name())
		if err := checkDir(dir); err != nil {
			return v, err
		}
		var u userRecord
		if err := readRecord(filepath.Join(dir, userFile), p.ReadLimit, &u); err != nil {
			return v, err
		}
		if err := u.valid(entry.Name()); err != nil {
			return v, err
		}
		if u.ResourceID != v.resource.ResourceID {
			return v, fmt.Errorf("user %s names generation %s in %s", u.UserID, u.ResourceID, d.Store)
		}
		v.users[u.UserID] = u
	}
	return v, nil
}

// lock takes each store's protocol lock in the given order and returns the release of all. It
// waits for each lock that another open file holds.
func lock(stores []string) (func(), error) {
	return lockStores(stores, syscall.LOCK_EX)
}

// tryLock takes the protocol lock of store only when no other open file holds it, so its
// caller never waits on another process. A held lock is an error.
func tryLock(store string) (func(), error) {
	return lockStores([]string{store}, syscall.LOCK_EX|syscall.LOCK_NB)
}

// lockStores takes each store's protocol lock with the flock operation how and returns the
// release of all. An open file description owns each lock, so in-process users exclude each
// other too.
func lockStores(stores []string, how int) (func(), error) {
	var held []*os.File
	release := func() {
		for _, file := range held {
			_ = file.Close()
		}
	}
	for _, store := range stores {
		if err := checkDir(store); err != nil {
			release()
			return nil, err
		}
		file, err := os.OpenFile(filepath.Join(store, lockFile), os.O_RDWR|os.O_CREATE|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, recordMode)
		if err == nil {
			held = append(held, file)
			err = lockFileOf(file, how)
		}
		if err != nil {
			release()
			return nil, err
		}
	}
	return release, nil
}

func lockFileOf(file *os.File, how int) error {
	if info, err := file.Stat(); err != nil || !info.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular lock file", file.Name())
	}
	return syscall.Flock(int(file.Fd()), how)
}

// publish durably replaces the record at path.
func (p *Protocol) publish(path string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return p.replace().Replace(path, data, recordMode)
}

// makeDir durably creates one private directory.
func (p *Protocol) makeDir(path string) error {
	if err := os.Mkdir(path, dirMode); err != nil {
		return err
	}
	if err := os.Chmod(path, dirMode); err != nil {
		return err
	}
	return p.syncDir(filepath.Dir(path))
}

// syncDir makes a directory's entries durable.
func (p *Protocol) syncDir(path string) error {
	dir, err := p.replace().OpenDirectory(path)
	if err != nil {
		return err
	}
	return errors.Join(dir.Sync(), dir.Close())
}

// remove deletes the given entries in order, then makes their directory's change durable.
func (p *Protocol) remove(dir string, entries ...string) error {
	for _, entry := range entries {
		if err := p.replace().Remove(entry); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	return p.syncDir(dir)
}

// removeUser deletes one user record and its directory durably.
func (p *Protocol) removeUser(store, user string) error {
	users := filepath.Join(store, usersDir)
	return p.remove(users, filepath.Join(users, user, userFile), filepath.Join(users, user))
}
