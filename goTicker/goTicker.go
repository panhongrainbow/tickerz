package goTicker

import (
	"fmt"
	tickerBase "github.com/panhongrainbow/tickerz/base"
	"time"
)

// Define status constants
const (
	StatusNewed    uint = iota + 1 // 1: Newed status
	StatusReSet                    // 2: Reset status
	StatusMockTime                 // 3: Mock time status
	StatusInactiv                  // 4: Inactive status
	StatusDead                     // 5: Dead status
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

	// Return err value
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

	// Return err value
	return
}

// mergeSortedBaseListAndRepeat merges a sorted list with a repeated sequence of numbers, producing a new sorted list of a given length.
func mergeSortedBaseListAndRepeat(list []int64, nearest, duration int64, count int) (output []int64) {
	// Create an empty slice to hold the repeated elements
	repeatArr := make([]int64, 0, count)
	// If nearest is not 0 and duration is at least 1,
	// repeat the nearest element for count times with duration as the interval
	if nearest != 0 && duration >= 1 {
		for i := 0; i < count; i++ {
			repeatArr = append(repeatArr, nearest)
			nearest = nearest + duration
		}
	}

	// Create an empty slice to hold the merged sorted output
	output = make([]int64, 0, count)

	// Merge the two sorted slices by comparing elements from both slices and adding the smaller one to the output
	i, j := 0, 0
	for i < len(list) && j < len(repeatArr) {
		// If the element from the base list is smaller than or equal to the repeated element, add the base element to the output
		if list[i] <= repeatArr[j] {
			output = append(output, list[i])
			i++
			// Otherwise, add the repeated element to the output
		} else {
			output = append(output, repeatArr[j])
			j++
		}

		// If we have added enough elements to the output, return the output
		if len(output) >= count {
			return
		}
	}

	// If there are any remaining elements in either slice, add them to the output
	output = append(output, list[i:]...)
	output = append(output, repeatArr[j:]...)

	// Return the output
	return
}

// calculateRepeatParameter calculates the nearest time based on a given duration and
// returns an error if the duration is less than or equal to 0.
func (receive *GoTicker) calculateRepeatParameter() (nearest, duration int64, err error) {
	// Convert duration to seconds
	duration = int64(receive.Opts.Duration.Seconds())

	// Get current Unix time
	now := time.Now().Unix()

	// CalculateWaitList the nearest time based on duration
	if duration >= 1 {
		// CalculateWaitList the time distance between current time and base stamp
		distance := time.Now().Unix() - receive.BaseStamp - ((now - receive.BaseStamp) % duration)
		nearest = receive.BaseStamp + distance
		// Round the nearest time to a multiple of the duration
	}

	// Return an error if the duration is less than or equal to 0
	if duration <= 0 {
		err = tickerBase.ErrInactiveRepeatList
	}

	// Return the output and err values
	return
}

// usableSubBaseList searches for a suitable base timestamp within a specified time range and
// returns a list of available sub-base timestamps.
func (receive *GoTicker) usableSubBaseList() (output []int64, err error) {
	// Get current Unix time
	now := time.Now().Unix()

LOOP:
	// Loop through the BaseList to find a suitable time
	for i := 0; i < len(receive.BaseList); i++ {
		// Check if the current BaseStamp is within the specified time range and in the future
		if receive.BaseList[i] > receive.BeginStamp &&
			receive.BaseList[i] < receive.EndStamp &&
			receive.BaseList[i] > now {
			// If the current BaseStamp meets the criteria,
			// set the output to the remaining BaseStamps in the list
			output = receive.BaseList[i:]
			break LOOP
		}
	}

	// Return the output and err values
	return

}

// CalculateWaitList calculates a wait list by merging sorted base lists and a repeat list based on given parameters.
func (receive *GoTicker) CalculateWaitList(count int) (waitList []int64, err error) {
	// Calculate the nearest and duration of the repeat parameter
	var nearest, duration int64
	var previousErr error
	nearest, duration, previousErr = receive.calculateRepeatParameter()

	// Get the usable subBaseList within the specified time range
	var subBaseList []int64
	subBaseList, err = receive.usableSubBaseList()
	if err == nil && previousErr != nil {
		err = previousErr
	}

	// Merge the sorted baseList and repeat parameter into a waitList of specified count
	waitList = mergeSortedBaseListAndRepeat(subBaseList, nearest, duration, count)
	if len(waitList) == 0 {
		err = tickerBase.ErrInactiveBaseListAndRepeatList
	}

	// Return the waitList and err values
	return
}

func (receive *GoTicker) WaitAccordingList(count int) (err error) {
	var waitLists []int64
	waitLists, err = receive.CalculateWaitList(count)
	if err != nil {
		fmt.Println(err)
	}

	for _, timePoint := range waitLists {
		now := time.Now().Unix()
		wait := timePoint - now
		if timePoint-now > 0 {
			timer := time.NewTimer(time.Duration(wait) * time.Second)
			<-timer.C
			fmt.Println("時間內執行")
			timer.Stop()
		} else {
			fmt.Println("超過時間後執行")
		}
	}
	return
}
