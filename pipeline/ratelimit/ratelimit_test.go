package ratelimit

import (
	"github.com/bartke/tributary"
)

var _ tributary.Pipeline = &limiter{}
