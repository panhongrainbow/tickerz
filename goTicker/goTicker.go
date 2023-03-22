package goTicker

import (
	tickerBase "github.com/panhongrainbow/tickerz/base"
	"time"
)

const (
	StatusNewed uint = iota + 1
	StatusReSet
	StatusMockTime
)

var mockDateStr string

type GoTicker tickerBase.Base

// updateNowDateOrMockIfNeeded updates ticker parameters and reloads location information if necessary.
func (receive *GoTicker) updateNowDateOrMockIfNeeded(mockDateStr string) (err error) {
	// Update nowDate with the current date or mocked date if it's set
	err = receive.UpdateNowDateOrMock(mockDateStr)
	// If base location is not set, reload location and update nowDate again
	if err == tickerBase.ErrNoBaseLocation {
		// Reload location
		err = receive.ReloadLocation()
		if err != nil {
			return
		}
		// Update nowDate again
		err = receive.UpdateNowDateOrMock(mockDateStr)
		if err != nil {
			return
		}
	}
	// Return err values
	return
}

// New is a function named New, which takes two arguments:
// Opts and OffOpts, representing the options for the ticker and the off options, respectively.
func New(opts tickerBase.Opts, offOpts tickerBase.OffOpts) (output *GoTicker, err error) {
	// Check if the options are valid
	err = opts.CheckOpts()
	if err != nil {
		return
	}

	// Create a new GoTicker object
	output = new(GoTicker)
	output.Opts = opts
	output.OffOpts = offOpts

	// If base location is not set, try to reload the location and update the nowDate again
	// ( "receive.UpdateNowDateOrMock(mockDateStr)" may be called twice, which may cause code duplication.
	// Extract it into a separate function or variable to reduce duplicated code.)
	err = output.updateNowDateOrMockIfNeeded(mockDateStr)
	if err != nil {
		return
	}

	// Determine the format of the base time and create a corresponding string
	var baseTimeStr string
	output.BaseStampType, err = tickerBase.TimeType(opts.BaseTime)
	if err != nil {
		return
	}

	// If the base time is in the TimeFormat, concatenate the current date and time
	if output.BaseStampType == tickerBase.TimeFormat {
		baseTimeStr = output.NowDate + " " + opts.BaseTime
		// Otherwise, the base time is already in the DatetimeFormat
	} else if output.BaseStampType == tickerBase.DatetimeFormat {
		baseTimeStr = opts.BaseTime
	}

	// Convert the base time string to a Unix timestamp
	output.BaseStamp, err = tickerBase.TimeValue(baseTimeStr, output.BaseLocation)
	if err != nil {
		return
	}

	// Determine the format of the base list time (if provided)
	if len(opts.BaseList) > 0 {
		output.BaseListType, err = tickerBase.TimeType(opts.BaseList[0])
	}

	// Convert each element in the base list to a Unix timestamp
	output.BaseList = make([]int64, 0, len(opts.BaseList))
	for i := 0; i < len(opts.BaseList); i++ {
		var element int64
		// If the base list time is in the TimeFormat, concatenate the current date and time
		if output.BaseListType == tickerBase.TimeFormat {
			baseTimeStr = output.NowDate + " " + opts.BaseList[i]
			element, err = tickerBase.TimeValue(baseTimeStr, output.BaseLocation)
			if err != nil {
				return
			}
			// Otherwise, the base list time is already in the DatetimeFormat
		} else if output.BaseListType == tickerBase.DatetimeFormat {
			element, err = tickerBase.TimeValue(opts.BaseList[i], output.BaseLocation)
			if err != nil {
				return
			}
		}
		output.BaseList = append(output.BaseList, element)
	}

	// Determine the format of the begin time and create a corresponding string
	var beginTimeStr string
	output.BeginStampType, err = tickerBase.TimeType(opts.BeginTime)
	if err != nil {
		return
	}

	// If the begin time is in the TimeFormat, concatenate the current date and time
	if output.BeginStampType == tickerBase.TimeFormat {
		beginTimeStr = output.NowDate + " " + opts.BeginTime
		// Otherwise, the begin time is already in the DatetimeFormat
	} else if output.BeginStampType == tickerBase.DatetimeFormat {
		beginTimeStr = opts.BeginTime
	}

	// Set the BeginStamp value based on the given beginTimeStr and BaseLocation
	output.BeginStamp, err = tickerBase.TimeValue(beginTimeStr, output.BaseLocation)
	if err != nil {
		return
	}

	// Determine the format of the end time and create a corresponding string
	var endTimeStr string
	output.EndStampType, err = tickerBase.TimeType(opts.EndTime)
	if err != nil {
		return
	}

	// If the end time is in the TimeFormat, concatenate the current date and time
	if output.EndStampType == tickerBase.TimeFormat {
		endTimeStr = output.NowDate + " " + opts.EndTime
		// Otherwise, the end time is already in the DatetimeFormat
	} else if output.EndStampType == tickerBase.DatetimeFormat {
		endTimeStr = opts.EndTime
	}

	// Set the EndStamp value based on the given endTimeStr and BaseLocation
	output.EndStamp, err = tickerBase.TimeValue(endTimeStr, output.BaseLocation)
	if err != nil {
		return
	}

	// Set the Status value to StatusNewed
	output.Status = StatusNewed

	// Return the output and err values
	return
}

// ReNew updates various timestamp-related values based on the current date,
// and updates the status of the ticker.
func (receive *GoTicker) ReNew() (err error) {
	// Check if the ticker has been initialized with New() function
	if receive.Status != StatusNewed {
		err = tickerBase.ErrNotNewedTicker
		return
	}

	// Update nowDate with the current date or mocked date if it's set
	err = receive.UpdateNowDateOrMock(mockDateStr)
	if err != tickerBase.ErrNoBaseLocation && err != nil {
		return
	}

	// Update baseList if baseListType is TimeFormat
	if receive.BaseStampType == tickerBase.TimeFormat {
		receive.BaseStamp, err = tickerBase.TimeValue(receive.NowDate+" "+receive.Opts.BaseTime, receive.BaseLocation)
		if err != nil {
			return
		}
	}

	// Update beginStamp if beginStampType is TimeFormat
	if receive.BaseListType == tickerBase.TimeFormat {
		receive.BaseList = make([]int64, 0, len(receive.Opts.BaseList))
		for i := 0; i < len(receive.Opts.BaseList); i++ {
			var baseListElement int64
			baseListElement, err = tickerBase.TimeValue(receive.NowDate+" "+receive.Opts.BaseList[i], receive.BaseLocation)
			if err != nil {
				return
			}
			receive.BaseList = append(receive.BaseList, baseListElement)
		}
	}

	// Update beginStamp if beginStampType is TimeFormat
	if receive.BeginStampType == tickerBase.TimeFormat {
		receive.BeginStamp, err = tickerBase.TimeValue(receive.NowDate+" "+receive.Opts.BeginTime, receive.BaseLocation)
		if err != nil {
			return
		}
	}

	// Update endStamp if endStampType is TimeFormat
	if receive.EndStampType == tickerBase.TimeFormat {
		receive.EndStamp, err = tickerBase.TimeValue(receive.NowDate+" "+receive.Opts.EndTime, receive.BaseLocation)
		if err != nil {
			return
		}
	}

	// Return the output and err values
	return
}

// UpdateNowDateOrMock is a method named UpdateNowDateOrMock that belongs to a struct type GoTicker.
// This method can take a string as input and update the date and time in the GoTicker struct based on that input.
func (receive *GoTicker) UpdateNowDateOrMock(input string) (err error) {
	if input != "" {
		receive.NowDate = input
		return
	}

	if receive.BaseLocation == nil {
		err = tickerBase.ErrNoBaseLocation
		return
	}

	now := time.Now().In(receive.BaseLocation)
	receive.NowDate = now.Format(tickerBase.DefaultDateFormatStr)

	return
}

// ReloadLocation is a method named ReloadLocation for the Base structure in Golang,
// which is intended to reload the time zone.
func (receive *GoTicker) ReloadLocation() (err error) {
	if receive.Opts.Location == "" {
		receive.Opts.Location = tickerBase.DefaultTimeZone
	}

	// reload time location
	receive.BaseLocation, err = time.LoadLocation(receive.Opts.Location)
	if err != nil {
		err = tickerBase.ErrUnSupportedLocation
		return
	}

	return
}

func (receive *GoTicker) Calculate(count int) (output []int64, err error) {
	output = make([]int64, 0, count)
	repeat := make([]int64, 0, count)

	duration := int64(receive.Opts.Duration.Seconds())
	now := time.Now().Unix()
	if duration >= 1 {
		distance := time.Now().Unix() - receive.BaseStamp - ((now - receive.BaseStamp) % duration)
		next := receive.BaseStamp + distance
		var loopCount int
	LOOP:
		for {
			if next > receive.BeginStamp &&
				next < receive.EndStamp &&
				next > now {
				repeat = append(repeat, next)
				if next > receive.EndStamp ||
					len(repeat) > count ||
					loopCount > 86400 {
					break LOOP
				}
			}
			next = next + duration
		}
	}

	// for test
	output = repeat

	return
}
