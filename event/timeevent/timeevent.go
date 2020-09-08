package timeevent

import (
	"context"
	"time"
)

type timeevent struct {
	data time.Time
	ctx  context.Context
	err  error
}

func New(t time.Time) *timeevent {
	return &timeevent{data: t, ctx: context.Background()}
}

func (t *timeevent) Payload() (out []byte) {
	out, t.err = t.data.MarshalBinary()
	return
}

func (t timeevent) Context() context.Context {
	return t.ctx
}

func (t timeevent) Error() error {
	return t.err
}
