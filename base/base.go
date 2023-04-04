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
	DefaultTestTimeZone      = "Pacific/Honolulu"
)

// time formats
const (
	EmptyTimeFormat uint = iota + 1
	DateFormat
	TimeFormat
	DatetimeFormat
)

const (
	ErrBaseListDifferentTypes        = Error("base list contains products of different types")
	ErrBaseListOverlaps              = Error("base list overlaps")
	ErrBeginEndTimeTypeNotEqual      = Error("begin time's type is not the same as end time's")
	ErrIncorrectBaseListOrder        = Error("incorrect base list order")
	ErrIncorrectBeginEndTimeOrder    = Error("incorrect begin end time order")
	ErrNegativeDuration              = Error("negative time duration")
	ErrNoBaseLocation                = Error("no base location")
	ErrUnsupBasetime                 = Error("unsupported basetime")
	ErrUnSupportedBaseList           = Error("unsuported base list")
	ErrUnSupportedBeginTimeFormat    = Error("unsupported begin time format")
	ErrUnSupportedEndTimeFormat      = Error("unsupported end time format")
	ErrUnSupportedLocation           = Error("unsupported location")
	ErrUnSupportedTimeFormat         = Error("unsupported time format")
	ErrNoOptLocation                 = Error("no opt location")
	ErrTimeParsion                   = Error("time parsing error")
	ErrNotNewedTicker                = Error("not newed ticker")
	ErrTickerDead                    = Error("ticker is is dead ")
	ErrInactiveBaseList              = Error("inactive base list")
	ErrInactiveRepeatList            = Error("inactive repeat list")
	ErrInactiveBaseListAndRepeatList = Error("inactive base list and repeat list")
	ErrUserInterrupted               = Error("user interrupt")
)

const (
	HumanMode uint = iota + 1
	CommputerMode
)

const (
	SignalOnTime uint = iota + 1
	SignalDelay
	SignalUserInterrupt
	SignalWaitForTomorrow
	SignalAbandonPrevious
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type Base struct {
	SerialHandler  func(int64) uint
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
	SignalChan     chan TickerSignal
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

type TickerSignal struct {
	SignalStatus uint
	SerialNumber uint
	DelaySeconds int64
}

func init() {
	var err error
	defaultTimeLocation, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatal(errors.New("default timeLocation error"))
	}
}

var defaultTimeLocation *time.Location

/*
TimeType determines the type of time format based on the string str,
and returns the format code tType and error err.
*/
func TimeType(str string) (tType uint, err error) {
	// If the str string is empty, set the tType format code to EmptyTimeFormat and return
	if str == "" {
		tType = EmptyTimeFormat
		return
	}
	/*
		If the str string matches the default date format,
		set the tType format code to DateFormat and return
	*/
	if _, err = time.Parse(DefaultDateFormatStr, str); err == nil {
		tType = DateFormat
		return
	}
	/*
		If the str string matches the default time format,
		set the tType format code to TimeFormat and return
	*/
	if _, err = time.Parse(DefaultTimeFormatStr, str); err == nil {
		tType = TimeFormat
		return
	}
	/*
		If the str string matches the default date-time format,
		set the tType format code to DatetimeFormat and return
	*/
	if _, err = time.Parse(DefaultDateTimeFormatStr, str); err == nil {
		tType = DatetimeFormat
		return
	}
	/*
		If none of the above conditions are met,
		set the error to ErrUnSupportedTimeFormat and return
	*/
	err = ErrUnSupportedTimeFormat
	return
}

/*
TimeValue parses the dateTimeStr string dateTimeStr as a date-time in the specified location,
and returns its Unix timestamp in seconds.
*/
func TimeValue(dateTimeStr string, location *time.Location) (timeStamp int64, err error) {
	// If no location is specified, use the default time zone
	if location == nil {
		location = defaultTimeLocation
	}
	// Parse the dateTimeStr string dateTimeStr as a date-time in the specified location
	var tmp time.Time
	tmp, err = time.ParseInLocation(DefaultDateTimeFormatStr, dateTimeStr, location)
	// If an error occurs during parsing, set the error to ErrTimeParsion
	if err != nil {
		err = ErrTimeParsion
	}
	// Set the timeStamp value to the Unix timestamp in seconds
	timeStamp = tmp.Unix()
	// Return the timeStamp and err values
	return
}

// CheckOpts is
func (receive *Opts) CheckOpts() (err error) {
	// check base time
	var tp uint
	tp, err = TimeType(receive.BaseTime)
	//
	if err == ErrUnSupportedTimeFormat {
		return
	}
	//
	if tp != TimeFormat && tp != DatetimeFormat {
		err = ErrUnsupBasetime
		return
	}

	//
	if receive.Location == "" {
		receive.Location = DefaultTimeZone
	}
	//
	_, err = time.LoadLocation(receive.Location)
	//
	if err != nil {
		err = ErrUnSupportedLocation
		return
	}

	//
	if receive.Duration.Nanoseconds() < 0 {
		err = ErrNegativeDuration
		return
	}

	//
	now := time.Now()
	//
	date := now.Format("2006-1-2")
	//
	var previous time.Time
	//
	var previousType uint
	//
	for i := 0; i < len(receive.BaseList); i++ {
		//
		tp, err = TimeType(receive.BaseList[i])
		//
		if err != nil {
			return
		}
		//
		if tp != TimeFormat && tp != DatetimeFormat {
			err = ErrUnSupportedBaseList
			return
		}

		//
		var datetime string
		//
		if tp == TimeFormat {
			datetime = date + " " + receive.BaseList[i]
		}
		//
		if tp == DatetimeFormat {
			datetime = receive.BaseList[i]
		}

		//
		if previousType == 0 {
			//
			previousType = tp
			//
		} else if previousType != 0 &&
			previousType != tp {
			//
			err = ErrBaseListDifferentTypes
			return
		}

		//
		var tmp time.Time
		tmp, err = time.Parse(DefaultDateTimeFormatStr, datetime)
		//
		if err != nil {
			return
		}
		//
		if !previous.Before(tmp) {
			err = ErrIncorrectBaseListOrder
			return
			//
		} else if previous.Before(tmp) {
			//
			previous = tmp
		}
	}

	//
	var beginTp uint
	beginTp, err = TimeType(receive.BeginTime)
	//
	if err != nil {
		//
		err = ErrUnSupportedBeginTimeFormat
		return
	}
	//
	if beginTp != TimeFormat &&
		beginTp != DatetimeFormat &&
		beginTp != EmptyTimeFormat {
		//
		err = ErrUnSupportedBeginTimeFormat
		return
	}

	//
	var endTp uint
	endTp, err = TimeType(receive.EndTime)
	//
	if err != nil {
		//
		err = ErrUnSupportedEndTimeFormat
		return
	}
	//
	if endTp != TimeFormat &&
		endTp != DatetimeFormat &&
		endTp != EmptyTimeFormat {
		//
		err = ErrUnSupportedEndTimeFormat
		return
	}

	//
	if (beginTp == EmptyTimeFormat && endTp == DatetimeFormat) ||
		(beginTp == DatetimeFormat && endTp == EmptyTimeFormat) {
		//
		return
	}

	//
	if beginTp != endTp {
		//
		err = ErrBeginEndTimeTypeNotEqual
		return
	}

	//
	var begintimeStr string
	var beginTime time.Time
	//
	if receive.BeginTime != "" {
		//
		if beginTp == TimeFormat {
			//
			begintimeStr = date + " " + receive.BeginTime
		}
		//
		if beginTp == DatetimeFormat {
			//
			begintimeStr = receive.BeginTime
		}
		//
		beginTime, err = time.Parse(DefaultDateTimeFormatStr, begintimeStr)
		if err != nil {
			//
			return
		}
	}

	//
	var endtimeStr string
	//
	var endTime time.Time
	//
	if receive.EndTime != "" {
		//
		if endTp == TimeFormat {
			//
			endtimeStr = date + " " + receive.EndTime
		}
		//
		if endTp == DatetimeFormat {
			//
			endtimeStr = receive.EndTime
		}
		//
		endTime, err = time.Parse(DefaultDateTimeFormatStr, endtimeStr)
		//
		if err != nil {
			//
			return
		}
	}

	//
	if receive.BeginTime != "" && receive.EndTime != "" {
		//
		if !beginTime.Before(endTime) {
			//
			err = ErrIncorrectBeginEndTimeOrder
			return
		}
	}

	//
	return
}
