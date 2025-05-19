package utils

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"time"
)

func GetAgeWithIdentificationNumber(identificationNumber string) (string, error) {
	reg := regexp.MustCompile(`^[1-9]\d{5}(18|19|20)(\d{2})((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`)
	params := reg.FindStringSubmatch(identificationNumber)
	if len(params) == 0 {
		return "", errors.New("身份证有误")
	}
	birYear, _ := strconv.Atoi(params[1] + params[2])
	birMonth, _ := strconv.Atoi(params[3])
	age := time.Now().Year() - birYear
	if int(time.Now().Month()) < birMonth {
		age--
	}
	return fmt.Sprintf("%d", age), nil
}

func GetSexWithIdentificationNumber(identificationNumber string) (string, error) {
	if len(identificationNumber) != 18 {
		return "", errors.New("身份证有误")
	}
	sexStr := identificationNumber[16:17]
	sexCode, err := strconv.Atoi(sexStr)
	if err != nil {
		return "", errors.New("身份证有误")
	}
	if sexCode%2 == 0 {
		return "女", nil
	} else {
		return "男", nil
	}
}

type TimeSort []string

func (t TimeSort) Len() int      { return len(t) }
func (t TimeSort) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t TimeSort) Less(i, j int) bool {

	var p1Int int64
	var p2Int int64

	if len(t[i]) == 10 {
		tm1, _ := time.Parse("2006-01-02", t[i])
		tm2, _ := time.Parse("2006-01-02", t[j])
		p1Int = tm1.Unix()
		p2Int = tm2.Unix()
	} else {
		tm1, _ := time.Parse("2006-01-02 15:04:05", t[i])
		tm2, _ := time.Parse("2006-01-02 15:04:05", t[j])
		p1Int = tm1.Unix()
		p2Int = tm2.Unix()
	}

	return p1Int > p2Int
}

func SortByTime(times []string) string {
	sort.Sort(TimeSort(times))

	if len(times) > 0 {
		return times[0]
	}

	return ""
}

// toString convert any value to string
func AnyToString(value interface{}) (d string) {
	val := reflect.ValueOf(value)
	switch value.(type) {
	case int:
		d = strconv.Itoa(int(val.Int()))
	case int64:
		d = strconv.FormatInt(val.Int(), 10)
	case float64, float32:
		d = strconv.FormatFloat(val.Float(), 'f', -1, 64)
	case int8, int16, int32:
		d = strconv.Itoa(int(val.Int()))
	case uint, uint32, uint64:
		d = strconv.FormatUint(val.Uint(), 10)
	case uint8, uint16:
		d = strconv.Itoa(int(val.Int()))
	case []byte:
		d = string(val.Bytes())
	case string:
		d = val.String()
	case bool:
		if val.Bool() {
			d = "true"
		} else {
			d = "false"
		}
	}
	return
}

func FormatDate(dateStr string) time.Time {
	t, e := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if e != nil {
	}
	return t
}
func FormatDateToUnix(dateStr string) int64 {
	t, e := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if e != nil {
	}
	return t.Unix()
}

// StringToInt ...
func StringToInt(s string) int {
	v, e := strconv.Atoi(s)
	if e != nil {
	}
	return v
}

var local, _ = time.LoadLocation("Local")

// 纠正日期字符串为2006-01-02 15:04:05
func RedressDateString(t string) string {
	if len(t) == 0 || t == "0001-01-01 00:00:00 +0000 UTC" {
		return time.Now().In(local).Format("2006-01-02 15:04:05")
	}
	_, err := time.ParseInLocation("2006-01-02 15:04:05", t, local)
	if err != nil {
		return time.Now().In(local).Format("2006-01-02 15:04:05")
	}
	if len(t) == 10 {
		return fmt.Sprintf("%s 00:00:00", t)
	}
	return t
}

func RedressDate(t time.Time) string {
	if t.IsZero() {
		t = time.Now().In(local)
	}
	return t.Format("2006-01-02 15:04:05")
}

/**字符串->时间对象*/
func Str2Time(formatTimeStr string) time.Time {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, formatTimeStr, loc) //使用模板在对应时区转化为time.time类型

	return theTime

}

func Str2TimePointer(formatTimeStr string) *time.Time {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, formatTimeStr, loc) //使用模板在对应时区转化为time.time类型

	if theTime.IsZero() {
		return nil
	}

	return &theTime

}

/**字符串->时间戳*/
func Str2Stamp(formatTimeStr string) int64 {
	timeStruct := Str2Time(formatTimeStr)
	millisecond := timeStruct.UnixNano() / 1e6
	return millisecond
}

/**时间对象->字符串*/
func Time2Str(t time.Time) string {
	const shortForm = "2006-01-01 15:04:05"
	temp := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	str := temp.Format(shortForm)
	return str
}

/**时间对象->字符串*/
func Time2StrPointer(t *time.Time) string {
	const shortForm = "2006-01-01 15:04:05"
	if t == nil {
		return ""
	}

	if t.IsZero() {
		return ""
	}

	temp := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	str := temp.Format(shortForm)
	return str
}

/*时间对象->时间戳*/
func Time2Stamp() int64 {
	t := time.Now()
	millisecond := t.UnixNano() / 1e6
	return millisecond
}

/*时间戳->字符串*/
func Stamp2Str(stamp int64) string {
	timeLayout := "2006-01-02 15:04:05"
	str := time.Unix(stamp/1000, 0).Format(timeLayout)
	return str
}

/*时间戳->时间对象*/
func Stamp2Time(stamp int64) time.Time {
	stampStr := Stamp2Str(stamp)
	timer := Str2Time(stampStr)
	return timer
}
