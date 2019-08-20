package trait

import (
	"strconv"
	"time"
)

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func Int64ToString(n int64) string {
	return strconv.FormatInt(n, 10)
}
