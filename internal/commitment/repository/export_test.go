package repository

// ApproveWith is approve, so an external test can fail one replacement of an approval.
var ApproveWith = Store.approve

// DefaultRevision is defaultRevision, so an external test can read the default-branch tip.
var DefaultRevision = Store.defaultRevision
