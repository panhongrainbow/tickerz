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
	ErrBaseListDifferentTypes          = Error("base list contains products of different types")
	ErrBaseListOverlaps                = Error("base list overlaps")
	ErrBeginEndTimeTypeNotEqual        = Error("begin time's type is not the same as end time's")
	ErrIncorrectBaseListOrder          = Error("incorrect base list order")
	ErrIncorrectBeginEndTimeOrder      = Error("incorrect begin end time order")
	ErrNegativeDuration                = Error("negative time duration")
	ErrNoBaseLocation                  = Error("no base location")
	ErrUnsupBasetime                   = Error("unsupported basetime")
	ErrUnSupportedBaseList             = Error("unsuported base list")
	ErrUnSupportedBeginOrEndTimeFormat = Error("unsupported begin or end time format")
	// ErrUnSupportedEndTimeFormat        = Error("unsupported end time format")
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

// CheckOpts function is responsible for validating the received Opts structure.
func (receive *Opts) CheckOpts() (err error) {
	// Validate the format of the base time
	var tp uint
	tp, err = TimeType(receive.BaseTime)
	// If the base time format is not supported, return an error
	if err == ErrUnSupportedTimeFormat {
		return
	}
	// Check if the time type is either TimeFormat or DatetimeFormat
	if tp != TimeFormat && tp != DatetimeFormat {
		err = ErrUnsupBasetime
		return
	}

	// If the location is empty, set it to the default time zone
	if receive.Location == "" {
		receive.Location = DefaultTimeZone
	}
	// Load the location to validate it
	_, err = time.LoadLocation(receive.Location)
	// if the location is not supported, return an error
	if err != nil {
		err = ErrUnSupportedLocation
		return
	}

	// Validate that the duration is not negative
	if receive.Duration.Nanoseconds() < 0 {
		err = ErrNegativeDuration
		return
	}

	// Set the current time
	now := time.Now()
	// Set the date format
	date := now.Format("2006-1-2")

	// Check the validity of the baseList
	err = checkOptsBaseList(receive.BaseList, date)
	if err != nil {
		return
	}

	// Declare variables for begin and end time types
	var beginTp, endTp uint
	// Declare variables for begin and end times
	var beginTime, endTime time.Time
	// Check the validity of the beginTime and convert it to a time.Time object
	beginTp, beginTime, err = checkOptsBeginEndTimeToTime(receive.BeginTime, date)
	if err != nil {
		return
	}
	// Check the validity of the endTime and convert it to a time.Time object
	endTp, endTime, err = checkOptsBeginEndTimeToTime(receive.EndTime, date)
	if err != nil {
		return
	}

	// If either the begin time or the end time is empty, return immediately
	if (beginTp == EmptyTimeFormat && endTp == DatetimeFormat) ||
		(beginTp == DatetimeFormat && endTp == EmptyTimeFormat) {
		return
	}

	// If the time types of begin time and end time are not equal, return an error
	if beginTp != endTp {
		// Assign the value of ErrBeginEndTimeTypeNotEqual to err and return
		err = ErrBeginEndTimeTypeNotEqual
		return
	}

	// Check if both the begin time and end time have values
	if receive.BeginTime != "" && receive.EndTime != "" {
		// Compare the begin time and end time to check if they are in the correct order
		if !beginTime.Before(endTime) {
			// If the begin time is after or equal to the end time, return an error
			err = ErrIncorrectBeginEndTimeOrder
			return
		}
	}

	/*
		If the function has reached this point,
		it means the begin and end times are valid and in the correct order,
		so return without errors
	*/
	return
}

/*
checkOptsBaseList checks baseList validity.
It validates each base time format and checks if the time type is TimeFormat or DatetimeFormat.
It also checks if all elements in baseList have the same time format.
*/
func checkOptsBaseList(baseList []string, date string) (err error) {
	var tp uint
	// Initialize the previous time to zero
	var previous time.Time
	// Initialize the previous time type to zero
	var previousType uint
	// Loop through the base list

	for i := 0; i < len(baseList); i++ {
		// Validate the format of each base time
		tp, err = TimeType(baseList[i])
		// If the format is not supported, return an error
		if err != nil {
			return
		}
		// Check if the time type is either TimeFormat or DatetimeFormat
		if tp != TimeFormat && tp != DatetimeFormat {
			err = ErrUnSupportedBaseList
			return
		}

		// Set the datetime string based on the time type
		var datetime string
		if tp == TimeFormat {
			datetime = date + " " + baseList[i]
		}
		if tp == DatetimeFormat {
			datetime = baseList[i]
		}

		// Check if the previous time type matches the current time type
		if previousType == 0 {
			previousType = tp
		} else if previousType != 0 &&
			previousType != tp {
			err = ErrBaseListDifferentTypes
			return
		}

		// Parse the datetime string into a time.Time object
		var tmp time.Time
		tmp, err = time.Parse(DefaultDateTimeFormatStr, datetime)
		if err != nil {
			return
		}
		// Check if the current time is not before the previous time
		if !previous.Before(tmp) {
			err = ErrIncorrectBaseListOrder
			return
		} else if previous.Before(tmp) {
			previous = tmp
		}
	}

	/*
		If the function has reached this point,
		it means the base list is valid and in the correct order,
		so return without errors
	*/
	return
}

/*
checkOptsBeginEndTimeToTime validates a begin or end time string, converts it into a time.Time object, and returns the time type and an error if there is any.
It supports specific time formats and returns an error if the format is unsupported.
*/
func checkOptsBeginEndTimeToTime(beginEndTimeStr string, date string) (beginEndTimeType uint, beginEndTime time.Time, err error) {
	// Validate the format of the begin or end time string and convert it to time type
	beginEndTimeType, err = TimeType(beginEndTimeStr)
	// If the format is not supported, return an error
	if err != nil {
		// Assign the value of ErrUnSupportedBeginOrEndTimeFormat to err and return
		err = ErrUnSupportedBeginOrEndTimeFormat
		return
	}
	// Check if the value of beginTp or endTp does not match any of the supported time formats
	if beginEndTimeType != TimeFormat &&
		beginEndTimeType != DatetimeFormat &&
		beginEndTimeType != EmptyTimeFormat {
		// Assign the value of ErrUnSupportedBeginOrEndTimeFormat to err and return
		err = ErrUnSupportedBeginOrEndTimeFormat
		return
	}

	// Parse the begin or end time string into a time.Time object using the correct format
	var newTimeStr string
	// Check if the begin or end time string is not empty
	if beginEndTimeStr != "" {
		// If the format of the begin or end time is TimeFormat, concatenate the date with the begin time string
		if beginEndTimeType == TimeFormat {
			newTimeStr = date + " " + beginEndTimeStr
		}
		// If the format of the begin or end time is DatetimeFormat, use the begin or end time string as received
		if beginEndTimeType == DatetimeFormat {
			newTimeStr = beginEndTimeStr
		}
		// Parse the begin or end time string into a time.Time object using the DefaultDateTimeFormatStr format
		beginEndTime, err = time.Parse(DefaultDateTimeFormatStr, newTimeStr)
		if err != nil {
			return
		}
	}

	// Return the beginEndTimeType, beginEndTime and err values
	return
}
