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
	ErrBaseListDifferentTypes     = Error("base list contains products of different types")
	ErrBaseListOverlaps           = Error("base list overlaps")
	ErrBeginEndTimeTypeNotEqual   = Error("begin time's type is not the same as end time's")
	ErrIncorrectBaseListOrder     = Error("incorrect base list order")
	ErrIncorrectBeginEndTimeOrder = Error("incorrect begin end time order")
	ErrNegativeDuration           = Error("negative time duration")
	ErrNoBaseLocation             = Error("no base location")
	ErrUnsupBasetime              = Error("unsupported basetime")
	ErrUnSupportedBaseList        = Error("unsuported base list")
	ErrUnSupportedBeginTimeFormat = Error("unsupported begin time format")
	ErrUnSupportedEndTimeFormat   = Error("unsupported end time format")
	ErrUnSupportedLocation        = Error("unsupported location")
	ErrUnSupportedTimeFormat      = Error("unsupported time format")
	ErrNoOptLocation              = Error("no opt location")
	ErrTimeParsion                = Error("time parsing error")
	ErrNotNewedTicker             = Error("not newed ticker")
)

const (
	HumanMode uint = iota + 1
	CommputerMode
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type Base struct {
	NowDate        string
	BaseLocation   *time.Location
	BaseStamp      int64
	BaseStampType  uint // time type
	BaseList       []int64
	BaseListType   uint // time type
	BeginStamp     int64
	BeginStampType uint // time type
	EndStamp       int64
	EndStampType   uint // time type
	Mode           uint
	Opts           Opts
	OffOpts        OffOpts
	Status         uint
}

type Opts struct {
	BaseTime  string // dateTime or time
	Location  string
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

// TimeValue is a function named TimeValue which takes in two arguments - a string input representing a datetime value,
// and a pointer to a time.Location value named location.
// * The function returns an int64 value representing the Unix time of the input datetime value and an error.
// * The function first checks if the location argument is nil, and if so, sets it to a default time location.
// It then attempts to parse the input string as a datetime value using the time.ParseInLocation function, with the given location.
// If there is an error, the function returns an error message indicating that the input is not a datetime value.
// * If the parsing is successful, the function converts the parsed time.
// Time value to a Unix timestamp using the Unix() method and assigns it to the output variable.
// * The function then returns the output value and a nil error.
func TimeValue(input string, location *time.Location) (output int64, err error) {
	if location == nil {
		location = defaultTimeLocation
	}
	var tmp time.Time
	tmp, err = time.ParseInLocation(DefaultDateTimeFormatStr, input, location)
	if err != nil {
		err = ErrTimeParsion
	}
	output = tmp.Unix()
	return
}

// CheckOpts defines a method CheckOpts() on a struct Opts.
// * The purpose of this method is to validate various properties of the Opts struct,
// including the format and order of various time-related fields.
// * The method first checks the BaseTime field to ensure that it is in a supported format.
// It then checks the Location field to ensure that it contains a valid time zone.
// The Duration field is also checked to ensure that it is a positive value.
// * The method then checks the BaseList field, which is a list of time strings.
// Each string is checked to ensure that it is in a supported format and that all strings in the list have the same format.
// The order of the strings in the list is also checked to ensure that they are in chronological order.
// * Next, the method checks the BeginTime and EndTime fields to ensure that
// they are in a supported format and that they are the same type.
// If both fields are present, the method checks that BeginTime comes before EndTime.
// * Finally, the method returns an error if any of the checks fail, or nil if all checks pass.
func (receive *Opts) CheckOpts() (err error) {
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
