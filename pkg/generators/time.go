package generators

import (
	"time"
)

type TimeGenerator struct{}

func NewTimeGenerator() *TimeGenerator {
	return &TimeGenerator{}
}

func (n *TimeGenerator) NowDate() time.Time {
	return time.Now()
}
