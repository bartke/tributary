package window

import (
	"github.com/bartke/tributary/pipeline/injector"
	"github.com/bartke/tributary/pipeline/interceptor"
)

type Windower interface {
	Create(v interface{}) interceptor.Fn
	Query(q string) injector.Fn
}
