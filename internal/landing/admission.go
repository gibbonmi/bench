package landing

import (
	"context"
	"errors"
)

// Admission is the publication decision that the installed broker composes around a
// reviewed landing. Check grades the exact composed tree before the gate. Publish decides
// that tree again after the gate and runs publish while that decision still holds.
type Admission interface {
	Check(tree string) error
	Publish(tree string, publish func() error) error
}

// unadmitted is the generic primitive's empty decision. Only LandReviewed selects it.
type unadmitted struct{}

func (unadmitted) Check(string) error                           { return nil }
func (unadmitted) Publish(_ string, publish func() error) error { return publish() }

// LandAdmitted is the installed broker's reviewed landing. A missing admission refuses,
// so no request value can select the unrestricted primitive.
func (o Owner) LandAdmitted(ctx context.Context, r ReviewedRequest, admission Admission) (ReviewedResult, error) {
	if admission == nil {
		return ReviewedResult{}, errors.New("reviewed landing requires a publication admission")
	}
	return o.landReviewed(ctx, r, admission)
}
