package www

import (
	"math/rand"
	"time"
)

// DelayFunc is a fucntion used to generate the teh duration to pause,
// before attempting a new retry.
// parameter: attempt num
// return time to delay
type DelayFunc func(int) time.Duration

// NoDelay always returns no delay
var NoDelay DelayFunc = func(attempt int) time.Duration {
	return 0
}

// LinearDelay returns 1sec * number of attempts
var LinearDelay DelayFunc = func(attempt int) time.Duration {
	return time.Duration(attempt) * time.Second
}

// Linear500Delay returns 500ms * number of attempts
var Linear500Delay DelayFunc = func(attempt int) time.Duration {
	return time.Duration(attempt) * (500 * time.Millisecond)
}

// LinearJitterDelay returns 1sec * number of attempts +- 125ms of jitter jitter
var LinearJitterDelay DelayFunc = func(attempt int) time.Duration {
	r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	return time.Duration(attempt)*time.Second + time.Duration((r.Float64()*250)-125)*time.Millisecond
}
