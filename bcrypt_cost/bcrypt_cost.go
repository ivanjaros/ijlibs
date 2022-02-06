package bcrypt_cost

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Adjust when needed, then recompile the whole code-base.
// Computing the cost dynamically proved to be problematic so we'll stick with hard-coded value.
const GlobalCost = 14

// computes bcrypt cost that will cause the current machine
// to compute hash in no less than provided target duration.
// note that jump in one cost point can result in quite large
// duration difference.
func Compute(target time.Duration, nearest ...bool) int {
	cost := bcrypt.MinCost
	secret := []byte("foo bar baz")
	var start time.Time
	var last, diff time.Duration

	for {
		start = time.Now()
		_, _ = bcrypt.GenerateFromPassword(secret, cost)
		diff = time.Now().Sub(start)

		if diff >= target || cost >= bcrypt.MaxCost {
			if len(nearest) > 0 && nearest[0] == true && last != 0 && last < diff {
				cost--
			}
			return cost
		}

		last = diff
		cost++
	}
}
