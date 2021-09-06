//go:generate mockgen -source timer.go -destination ../mock/timer.go -package mock
package timer

import "time"

// Timer is an abstraction layer over time.Now
type Timer interface {
	Now() time.Time
}

// New return Timer implementaion
func New() traveler {
	return traveler{}
}

type traveler struct {
	now *time.Time
}

// Now return the current local time
func (t traveler) Now() time.Time {
	if t.now != nil {
		return *t.now
	}
	return time.Now()
}

// Travel fix NOW value
func (t *traveler) Travel(now *time.Time) {
	t.now = now
}
