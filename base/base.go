package base

import (
	"errors"
	"log"
	"time"
)

// This is a code snippet that defines several constants.
// These constants define date formats, time formats, date-time formats, and the default time zone.
const (
	DefaultDateFormatStr     = "2006-1-2"
	DefaultTimeFormatStr     = "15:4:5"
	DefaultDateTimeFormatStr = "2006-1-2 15:4:5"
	DefaultTimeZone          = "Asia/Shanghai"
)

// time formats
const (
	EmptyTimeFormat uint = iota + 1
	DateFormat
	TimeFormat
	DatetimeFormat
)

const (
	ErrUnSupportedTimeFormat      = Error("unsupported time format")
	ErrUnSupportedLocation        = Error("unsupported location")
	ErrUnSupportedBeginTimeFormat = Error("unsupported begin time format")
	ErrUnSupportedEndTimeFormat   = Error("unsupported end time format")
	ErrUnsupBasetime              = Error("unsupported basetime")
	ErrNegativeDuration           = Error("negative time duration")
	ErrUnSupportedBaseList        = Error("unsuported base list")
	ErrBaseListOverlaps           = Error("base list overlaps")
	ErrIncorrectBaseListOrder     = Error("incorrect base list order")
	ErrBaseListDifferentTypes     = Error("base list contains products of different types")
	ErrBeginEndTimeTypeNotEqual   = Error("begin time's type is not the same as end time's")
	ErrIncorrectBeginEndTimeOrder = Error("incorrect begin end time order")
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type Base struct {
	BaseStamp    int64
	BaseLocation *time.Location
	Opts         Opts
	OffOpts      OffOpts
	BeginStamp   int64
	EndStamp     int64
	BaseListType uint
}

type Opts struct {
	BaseTime  string // dateTime or time
	Location  string
	Repeat    bool
	Duration  time.Duration
	BaseList  []string
	BeginTime string
	EndTime   string
}

type OffOpts struct {
	EveryDayOff  bool
	MondayOff    bool
	TuesdayOff   bool
	WednesdayOff bool
	ThursdayOff  bool
	FridayOff    bool
	SaturdayOff  bool
	SundayOff    bool
}

func init() {
	var err error
	defaultTimeLocation, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatal(errors.New("default timeLocation error"))
	}
}

var defaultTimeLocation *time.Location

func TimeType(input string) (output uint, err error) {
	if input == "" {
		output = EmptyTimeFormat
		return
	}
	if _, err = time.Parse(DefaultDateFormatStr, input); err == nil {
		output = DateFormat
		return
	}
	if _, err = time.Parse(DefaultTimeFormatStr, input); err == nil {
		output = TimeFormat
		return
	}
	if _, err = time.Parse(DefaultDateTimeFormatStr, input); err == nil {
		output = DatetimeFormat
		return
	}

	err = ErrUnSupportedTimeFormat
	return
}

func TimeValue(input string, location *time.Location) (output int64, err error) {
	if location == nil {
		location = defaultTimeLocation
	}
	t, err := time.ParseInLocation(DefaultDateTimeFormatStr, input, location)
	if err != nil {
		return -1, errors.New("not datetime type")
	}
	output = t.Unix()
	return
}

func (receive Opts) CheckOpts() (err error) {
	// check base time
	var tp uint
	tp, err = TimeType(receive.BaseTime)
	if err == ErrUnSupportedTimeFormat {
		return
	}
	if tp != TimeFormat && tp != DatetimeFormat {
		err = ErrUnsupBasetime
		return
	}

	// check location
	if receive.Location == "" {
		receive.Location = DefaultTimeZone
	}
	_, err = time.LoadLocation(receive.Location)
	if err != nil {
		err = ErrUnSupportedLocation
		return
	}

	// check duration
	if receive.Duration.Nanoseconds() < 0 {
		err = ErrNegativeDuration
		return
	}

	// check base list
	now := time.Now()
	date := now.Format("2006-1-2")
	var previous time.Time
	var previousType uint
	for i := 0; i < len(receive.BaseList); i++ {
		tp, err = TimeType(receive.BaseList[i])
		if err != nil {
			return
		}
		if tp != TimeFormat && tp != DatetimeFormat {
			err = ErrUnSupportedBaseList
			return
		}

		// make new datetime string
		var datetime string
		if tp == TimeFormat {
			datetime = date + " " + receive.BaseList[i]
		}
		if tp == DatetimeFormat {
			datetime = receive.BaseList[i]
		}

		// check if this base list contains products of different types？
		if previousType == 0 {
			previousType = tp
		} else if previousType != 0 &&
			previousType != tp {
			err = ErrBaseListDifferentTypes
			return
		}

		// check base list order
		var tmp time.Time
		tmp, err = time.Parse(DefaultDateTimeFormatStr, datetime)
		if err != nil {
			// return the original error message
			return
		}
		if !previous.Before(tmp) {
			err = ErrIncorrectBaseListOrder
			return
		} else if previous.Before(tmp) {
			previous = tmp
		}
	}

	// check begin time type
	var beginTp uint
	beginTp, err = TimeType(receive.BeginTime)
	if err != nil {
		err = ErrUnSupportedBeginTimeFormat
		return
	}
	if beginTp != TimeFormat &&
		beginTp != DatetimeFormat &&
		beginTp != EmptyTimeFormat {
		err = ErrUnSupportedBeginTimeFormat
		return
	}

	// check end time type
	var endTp uint
	endTp, err = TimeType(receive.EndTime)
	if err != nil {
		err = ErrUnSupportedEndTimeFormat
		return
	}
	if endTp != TimeFormat &&
		endTp != DatetimeFormat &&
		endTp != EmptyTimeFormat {
		err = ErrUnSupportedEndTimeFormat
		return
	}

	// accept EmptyTimeFormat on either side
	if (beginTp == EmptyTimeFormat && endTp == DatetimeFormat) ||
		(beginTp == DatetimeFormat && endTp == EmptyTimeFormat) {
		return
	}

	// receive.BeginTime's type must be the same as receive.EndTime's
	if beginTp != endTp {
		err = ErrBeginEndTimeTypeNotEqual
		return
	}

	// check begin time and end time order
	var begintimeStr string
	var beginTime time.Time
	if receive.BeginTime != "" {
		if beginTp == TimeFormat {
			begintimeStr = date + " " + receive.BeginTime
		}
		if beginTp == DatetimeFormat {
			begintimeStr = receive.BeginTime
		}
		beginTime, err = time.Parse(DefaultDateTimeFormatStr, begintimeStr)
		if err != nil {
			// return the original error message
			return
		}
	}

	var endtimeStr string
	var endTime time.Time
	if receive.EndTime != "" {
		if endTp == TimeFormat {
			endtimeStr = date + " " + receive.EndTime
		}
		if endTp == DatetimeFormat {
			endtimeStr = receive.EndTime
		}
		endTime, err = time.Parse(DefaultDateTimeFormatStr, endtimeStr)
		if err != nil {
			// return the original error message
			return
		}
	}

	if receive.BeginTime != "" && receive.EndTime != "" {
		if !beginTime.Before(endTime) {
			err = ErrIncorrectBeginEndTimeOrder
			return
		}
	}

	return
}
