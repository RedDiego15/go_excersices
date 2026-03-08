//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds (accumulated)
	mu        sync.Mutex
}

func (user *User) addTimeUsed(timedUsed time.Duration) {
	user.mu.Lock()
	defer user.mu.Unlock()
	user.TimeUsed += int64(timedUsed.Seconds())
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed (no quota left).
// We hold the user lock for the entire request so only one process per user
// runs at a time; otherwise concurrent requests could all pass the quota check
// before any of them add their time.
func HandleRequest(process func(), u *User) bool {
	u.mu.Lock()
	defer u.mu.Unlock()

	if !u.IsPremium && u.TimeUsed >= 10 {
		return false // no quota left, process killed
	}

	start := time.Now()
	process()
	elapsed := time.Since(start)
	u.TimeUsed += int64(elapsed.Seconds())
	return true
}

func main() {
	RunMockServer()
}
