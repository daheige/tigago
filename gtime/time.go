package gtime

import (
	"time"
)

const (
	tmFmtWithMS = "2006-01-02 15:04:05.999"
	tmFmtMissMS = "2006-01-02 15:04:05"
)

// TimeZone 默认时区设置，可以是Local本地时区
var TimeZone = "Asia/Shanghai"

// SetTimeZone set time zone.
func SetTimeZone(zone string) {
	TimeZone = zone
}

// FormatTime format a time.Time to string as 2006-01-02 15:04:05.999
func FormatTime(t time.Time) string {
	return t.Format(tmFmtWithMS)
}

// FormatTime19 format a time.Time to string as 2006-01-02 15:04:05
// 将time.time转换为日期格式
func FormatTime19(t time.Time) string {
	return t.Format(tmFmtMissMS)
}

// GetLoc 获取时区loc
func GetLoc(zone string) *time.Location {
	loc, _ := time.LoadLocation(TimeZone)
	return loc
}

// GetCurrentLocalTime 当前本地时间
func GetCurrentLocalTime() string {
	return GetTimeByTimeZone(TimeZone)
}

// GetTimeByTimeZone get time by zone.
func GetTimeByTimeZone(zone string) string {
	loc, _ := time.LoadLocation(zone)
	return time.Now().In(loc).Format(tmFmtMissMS)
}

// FormatNow format time.Now() use FormatTime
func FormatNow() string {
	return FormatTime(time.Now())
}

// FormatUTC format time.Now().UTC() use FormatTime
func FormatUTC() string {
	return FormatTime(time.Now().UTC())
}

// ParseTime parse a string to time.Time
func ParseTime(s string) (time.Time, error) {
	if len(s) == len(tmFmtMissMS) {
		return time.ParseInLocation(tmFmtMissMS, s, time.Local)
	}
	return time.ParseInLocation(tmFmtWithMS, s, time.Local)
}

// ParseTimeUTC parse a string as "2006-01-02 15:04:05.999" to time.Time
func ParseTimeUTC(s string) (time.Time, error) {
	if len(s) == len(tmFmtMissMS) {
		return time.ParseInLocation(tmFmtMissMS, s, time.UTC)
	}
	return time.ParseInLocation(tmFmtWithMS, s, time.UTC)
}

// NumberTime format a time.Time to number as 20060102150405999
func NumberTime(t time.Time) uint64 {
	y, m, d := t.Date()
	h, M, s := t.Clock()
	ms := t.Nanosecond() / 1000000
	return uint64(ms+s*1000+M*100000+h*10000000+d*1000000000) +
		uint64(m)*100000000000 + uint64(y)*10000000000000
}

// NumberNow format time.Now() use NumberTime
func NumberNow() uint64 {
	return NumberTime(time.Now())
}

// NumberUTC format time.Now().UTC() use NumberTime
func NumberUTC() uint64 {
	return NumberTime(time.Now().UTC())
}

// parseNumber parse a uint64 as 20060102150405999 to time.Time
func parseNumber(t uint64, tl *time.Location) (time.Time, error) {
	ns := int((t % 1000) * 1000000)
	t /= 1000
	s := int(t % 100)
	t /= 100
	M := int(t % 100)
	t /= 100
	h := int(t % 100)
	t /= 100
	d := int(t % 100)
	t /= 100
	m := time.Month(t % 100)
	y := int(t / 100)

	return time.Date(y, m, d, h, M, s, ns, tl), nil
}

// ParseNumber parse t to time.Time
func ParseNumber(t uint64) (time.Time, error) {
	return parseNumber(t, time.Local)
}

// ParseNumberUTC parset t to utc time
func ParseNumberUTC(t uint64) (time.Time, error) {
	return parseNumber(t, time.UTC)
}

// StrToTime like php strtotime()
// StrToTime("02/01/2006 15:04:05", "02/01/2016 15:04:05") == 1451747045
// StrToTime("3 04 PM", "8 41 PM") == -62167144740
func StrToTime(format, strTime string) (int64, error) {
	t, err := time.Parse(format, strTime)
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}

// Time time()
func Time() int64 {
	return time.Now().Unix()
}

// Date date()
// Date("02/01/2006 15:04:05 PM", 1524799394)
func Date(format string, timestamp int64) string {
	return time.Unix(timestamp, 0).Format(format)
}

// Sleep sleep()
func Sleep(t int64) {
	time.Sleep(time.Duration(t) * time.Second)
}

// Usleep usleep()
func Usleep(t int64) {
	time.Sleep(time.Duration(t) * time.Microsecond)
}
