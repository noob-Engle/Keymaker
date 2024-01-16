package pkg

import (
	"io"
	"net/http"
	"strconv"
	"time"
)

type STime struct {
}

// QXGetNowTimestamp 13位时间戳
func (STime) QXGetNowTimestamp() (timestamp int) {
	timestamp = int(time.Now().UnixMilli())
	return
}
func (STime) QXGetNowTimestampSecond() (timestamp int) {
	timestamp = int(time.Now().UnixNano())
	return
}

func (STime) QXGetNowTime() (times time.Time) {
	times = time.Now()
	return
}

func (STime) DToText(times time.Time) (text string) {
	text = times.Format("2006-01-02 15:04:05")
	return
}

func (STime) CCreateTime(year, month, day, hour, minute, second, millisecond int) (value time.Time) {
	value = time.Date(year, time.Month(month), day, hour, minute, second, millisecond, time.Local)
	return
}

func (STime) CLoadText(text string) (times time.Time) {
	times = allType.DToTime(text)
	return
}

func (STime) CReloadTime(timestamp any) (times time.Time) {
	res, _ := strconv.ParseInt(any_to_doc(timestamp), 0, 64)
	times = time.Unix(res, 0)
	return
}

// DToTimestamp 13位
func (STime) DToTimestamp(times time.Time) (timestamp int) {
	timestamp = int(times.UnixMilli())
	return
}

func (STime) QNGetYear(times time.Time) (value int) {
	value = times.Year()
	return
}
func (STime) QYGetMonth(times time.Time) (value int) {
	value = int(times.Month())
	return
}
func (STime) QRGetDay(times time.Time) (value int) {
	value = times.Day()
	return
}
func (STime) QXGetHour(times time.Time) (value int) {
	value = times.Hour()
	return
}
func (STime) QFGetMinute(times time.Time) (value int) {
	value = times.Minute()
	return
}
func (STime) QMGetSecond(times time.Time) (value int) {
	value = times.Second()
	return
}
func (STime) QXGetWeekday(times time.Time) (value int) {
	value = int(times.Weekday())
	return
}
func (STime) QZGetWeek(times time.Time) (value int) {
	_, value = times.ISOWeek()
	return
}

func (STime) QDGetDayOfYear(times time.Time) (value int) {
	value = times.YearDay()
	return
}

// YCDelayProgram 让程序延时执行 再这里等待时间到达收 继续执行 单位 毫秒
func (STime) YCDelayProgram(millisecond int) (value int) {
	time.Sleep(time.Duration(millisecond) * time.Millisecond)
	return
}

func (STime) QBGetChineseTime() (times time.Time, returnerr error) {
	url := "https://www.baidu.com/"
	resp, err := http.Get(url)
	if err != nil {
		returnerr = err
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	times = allType.DToTime(resp.Header.Get("date"), time.RFC1123)
	times = times.Add(time.Hour * 8)
	return
}

// ZJChangeTime 正负=增减 为0时 不变
func (STime) ZJChangeTime(times time.Time, year, month, day, hour, minute, second int) (times1 time.Time) {
	if year != 0 || month != 0 || day != 0 {
		times1 = times.AddDate(year, month, day)
	}
	if hour != 0 || minute != 0 || second != 0 {
		times1 = times1.Add(time.Hour*time.Duration(hour) + time.Minute*time.Duration(minute) + time.Second*time.Duration(second))
	}
	return
}
