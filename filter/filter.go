package filter

import (
	"time"

	"github.com/bartke/tributary/pipeline/interceptor"
	"github.com/bartke/tributary/sink/handler"
)

type Filter interface {
	Create(name string) (interceptor.Fn, error)
	Clean(name string, interval time.Duration) handler.Fn
}
