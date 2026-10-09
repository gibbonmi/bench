package responsebound

import "io"

// The size test runs in the external test package, so it can take its assignment checkout
// from responseboundtest, which imports this package. These names expose the in-package
// response helpers and the owner's create seam to that test.

type Tagged = tagged

var (
	PrivateHome       = privateHome
	OutsideRepository = outsideRepository
	RespondWith       = respondWith
	StdoutLines       = stdoutLines
	Numbered          = numbered
	Joined            = joined
	ErrInjected       = errInjected
	ExclusiveCreate   = exclusiveCreate
)

func StderrWrite(data string) Tagged { return tagged{stderr: true, data: data} }

func FaultAfter(file io.WriteCloser, limit int) io.WriteCloser {
	return &faultAfter{file: file, limit: limit}
}

func SetCreate(owner *Owner, create func(path string) (io.WriteCloser, error)) {
	owner.create = create
}
