package stime

import (
	"math"
	"time"
)

var (
	timeOffset time.Duration
	timeDeltas []int64

	ok  = false
	Now = time.Now
)

func now() time.Time {
	return time.Now().Add(time.Nanosecond * timeOffset)
}

func Delta(t1, t2, srvTime int64) (bool, time.Duration) {
	if ok {
		return true, timeOffset
	}
	delta := t2 - t1
	if delta > 100000000 {
		return false, 0
	}
	checkTime := t1 + delta/2
	d := int64(math.Abs(float64(checkTime - srvTime)))
	timeDeltas = append(timeDeltas, d)
	if len(timeDeltas) >= 3 {
		sum := int64(0)
		for _, v := range timeDeltas {
			sum += v
		}
		avg := sum / int64(len(timeDeltas))
		timeOffset = time.Nanosecond * time.Duration(avg)
		Now = now
		ok = true
		return true, timeOffset
	}
	return false, 0
}
