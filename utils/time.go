package utils

/*
const char* build_time(void)
{
    static const char* psz_build_time = __DATE__ " " __TIME__;
    return psz_build_time;
}
*/
import "C"

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type DateTime time.Time

const (
	DateLayout       = "2006-01-02"
	DateTimeLayout   = "2006-01-02 15:04:05"
	BuildTimeLayout  = "2006.0102.150405"
	CBuildTimeLayout = "Jan  _2 2006 15:04:05"
)

var StartTime = time.Now()

func (dt *DateTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(DateTimeLayout, string(data), time.Local)
	*dt = DateTime(now)
	return
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(DateTimeLayout)+2)
	b = append(b, '"')
	b = time.Time(dt).AppendFormat(b, DateTimeLayout)
	b = append(b, '"')
	return b, nil
}

func (dt DateTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	ti := time.Time(dt)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

func (dt *DateTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*dt = DateTime(value)
		return nil
	}
	return nil
}

func (dt DateTime) String() string {
	return time.Time(dt).Format(DateTimeLayout)
}

func UpTime() time.Duration {
	return time.Since(StartTime)
}

func UpTimeString() string {
	d := UpTime()
	days := d / (time.Hour * 24)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	d -= minutes * time.Minute
	seconds := d / time.Second
	return fmt.Sprintf("%d Days %d Hours %d Mins %d Secs", days, hours, minutes, seconds)
}

func BuildTime() time.Time {
	t, _ := time.Parse(CBuildTimeLayout, C.GoString(C.build_time()))
	return t
}
