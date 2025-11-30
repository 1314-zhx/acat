package util

import (
	"time"
)

// 中国时区
var TimeLocation, _ = time.LoadLocation("Asia/Shanghai")

const Timeform = "2006-01-02 15:04:05"
const TimeformShort = "2006-01-02"

// 当前时间的时间戳
func NowUnix() int64 {
	return time.Now().In(TimeLocation).Unix()
}

// 将unix时间戳格式化为yyyymmdd格式字符串
func FormatFromUnixTimeShort(t int64) string {
	if t > 0 {
		return time.Unix(t, 0).Format(TimeformShort)
	} else {
		return time.Now().Format(TimeformShort)
	}
}
