package httputil

import (
	"math"
	"time"
)

func ExponentialBackoff(min, max time.Duration, num int) time.Duration {
	mult := math.Pow(2, float64(num)) * float64(min)
	sleep := time.Duration(mult)
	if float64(sleep) != mult || sleep > max {
		sleep = max
	}
	return sleep
}

func LinearWithJitterBackoff(min, max, num int) time.Duration {
	panic("not implemented")
}
