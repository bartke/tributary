package interceptor

import "github.com/bartke/tributary"

var _ tributary.Pipeline = &interceptor{}
