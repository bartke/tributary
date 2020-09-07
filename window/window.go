package window

import (
	"github.com/bartke/tributary/pipeline/injector"
	"github.com/bartke/tributary/pipeline/multiinjector"
)

type Windower interface {
	Create(v interface{}) injector.Injector
	Query(q string) multiinjector.MultiInjector
}
