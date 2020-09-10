package standardevent

import (
	"context"

	"github.com/bartke/tributary"
)

type standardevent struct {
	payload []byte
	ctx     context.Context
	err     error
}

func New(ctx context.Context, p []byte, err error) tributary.Event {
	return &standardevent{payload: p, ctx: ctx, err: err}
}

func (m standardevent) Payload() []byte {
	return m.payload
}

func (m standardevent) Context() context.Context {
	return m.ctx
}

func (m standardevent) Error() error {
	return m.err
}
