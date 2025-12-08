package util

import (
	"fmt"
	"time"
)

// 中国时区
var TimeLocation, _ = time.LoadLocation("Asia/Shanghai")

const TimeformShort = "2006-01-02"
const DateTimeLocalFormat = "2006-01-02T15:04"

type DateTimeLocal time.Time

func (dt DateTimeLocal) Time() time.Time {
	return time.Time(dt)
}

// UnmarshalJSON 实现自定义 JSON 反序列化,shouldBindJSON，调用这个方法，因为这个自定义类型，实现了这个接口，自动检查时发现有自定义类型实现了
// 接口就调用这个方法，且仅对自定义的 DateTimeLocal类型有效
func (dt *DateTimeLocal) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == `""` || string(data) == `null` {
		return nil
	}

	s := string(data[1 : len(data)-1]) // 去掉引号

	// 使用 InLocation 解析，指定为上海时区
	t, err := time.ParseInLocation(DateTimeLocalFormat, s, TimeLocation)
	if err != nil {
		return fmt.Errorf("cannot parse %q as datetime-local in Asia/Shanghai: %w", s, err)
	}

	*dt = DateTimeLocal(t)
	return nil
}

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
