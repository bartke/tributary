package filter

import (
	"github.com/bartke/tributary/pipeline/interceptor"
	"github.com/bartke/tributary/sink/handler"
)

type Filter interface {
	Create(name string) (interceptor.Fn, error)
	Clean(name string, s int) handler.Fn
}
