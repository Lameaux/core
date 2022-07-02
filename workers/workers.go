package workers

import (
	"time"
)

type Worker interface {
	Name() string
	Run() (bool, error)
	SleepTime() time.Duration
}
