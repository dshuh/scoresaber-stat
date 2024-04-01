package util

import (
	"time"

	"stove-gitlab.sginfra.net/backend/template/framework/util/log"
)

// GetTS2MilliSec 단위가 milliseconds인 현재 timestamp 반환
func GetTS2MilliSec() int64 {
	return time.Now().UnixNano() / 1000000
}

// GetTS2MicroSec 단위가 microseconds인 현재 timestamp 반환
func GetTS2MicroSec() int64 {
	return time.Now().UnixNano() / 1000
}

// GetTS2NanoSec 단위가 nanoseconds인 현재 timestamp 반환
func GetTS2NanoSec() int64 {
	return time.Now().UnixNano()
}

// GetTS2Sec 단위가 seconds인 현재 timestamp 반환
func GetTS2Sec() int64 {
	return time.Now().Unix()
}

// TimeTrack 경과시간을 출력한다.
func TimeTrack(start time.Time, tag string) {
	elapsed := time.Since(start)
	log.Debugf("%s took %v", tag, elapsed)
}
