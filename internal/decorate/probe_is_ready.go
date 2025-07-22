package decorate

import "time"

var probeDelay = (time.Millisecond * 15).Milliseconds()

func probeIsReady(startTime int64) bool {
	var (
		now       = time.Now().UnixMilli()
		readyTime = startTime + probeDelay
	)

	return startTime > 0 && readyTime < now
}
