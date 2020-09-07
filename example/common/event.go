package common

import (
	"context"
	"time"
)

type timeevent struct {
	t   time.Time
	ctx context.Context
	err error
}

func TimeEvent(t time.Time) *timeevent {
	return &timeevent{t: t, ctx: context.Background()}
}

func (t timeevent) Payload() []byte {
	return []byte(t.t.Format(time.RFC3339))
}

func (t timeevent) Context() context.Context {
	return t.ctx
}

func (t timeevent) Error() error {
	return t.err
}
