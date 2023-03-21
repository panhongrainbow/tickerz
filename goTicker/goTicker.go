package goTicker

import (
	tickerBase "github.com/panhongrainbow/tickerz/base"
	"time"
)

type GoTicker tickerBase.Base

// New is a function named New, which takes two arguments:
// Opts and OffOpts, representing the options for the ticker and the off options, respectively.
// * The function returns GoTicker and error. In the function body,
// it first calls the CheckOpts method of opts to check if the options are correct.
// * If they are correct, it creates a GoTicker struct and assigns Opts and OffOpts to the struct's Opts and OffOpts fields, respectively.
// Next, it calls the UpdateNowDate method to reload the current date and updates BaseStamp based on the BaseTime in the options.
// * If the location for BaseTime is not specified in the options,
// it calls the ReloadLocation method to reload the location and the current date and BaseStamp.
// * Then, it generates BaseList, BeginStamp, and EndStamp based on the BaseList, BeginTime, and EndTime in the options, respectively.
// * Finally, it returns GoTicker and error.
func New(opts tickerBase.Opts, offOpts tickerBase.OffOpts) (output *GoTicker, err error) {
	// check opts
	err = opts.CheckOpts()
	if err != nil {
		return
	}

	// create goTicker
	output = new(GoTicker)
	output.Opts = opts
	output.OffOpts = offOpts

	// reload current date
	err = output.UpdateNowDate()
	if err != tickerBase.ErrNoBaseLocation && err != nil {
		return
	}
	if err == tickerBase.ErrNoBaseLocation {
		err = output.ReloadLocation()
		if err != nil {
			return
		}
		err = output.UpdateNowDate()
		if err != nil {
			return
		}
	}

	// reform base time
	var baseTimeStr string
	var tp uint
	tp, err = tickerBase.TimeType(opts.BaseTime)
	if err != nil {
		return
	}
	if tp == tickerBase.TimeFormat {
		baseTimeStr = output.NowDate + " " + opts.BaseTime
	} else if tp == tickerBase.DatetimeFormat {
		baseTimeStr = opts.BaseTime
	}

	// renew base time
	output.BaseStamp, err = tickerBase.TimeValue(baseTimeStr, output.BaseLocation)
	if err != nil {
		return
	}

	// check base list element type
	if len(opts.BaseList) > 0 {
		tp, err = tickerBase.TimeType(opts.BaseList[0])
	}

	// make base list
	output.BaseList = make([]int64, 0, len(opts.BaseList))
	for i := 0; i < len(opts.BaseList); i++ {
		var element int64
		if tp == tickerBase.TimeFormat {
			baseTimeStr = output.NowDate + " " + opts.BaseList[i]
			element, err = tickerBase.TimeValue(baseTimeStr, output.BaseLocation)
			if err != nil {
				return
			}
		} else if tp == tickerBase.DatetimeFormat {
			element, err = tickerBase.TimeValue(opts.BaseList[i], output.BaseLocation)
			if err != nil {
				return
			}
		}
		output.BaseList = append(output.BaseList, element)
	}

	// make begin time
	var beginTimeStr string
	tp, err = tickerBase.TimeType(opts.BeginTime)
	if err != nil {
		return
	}
	if tp == tickerBase.TimeFormat {
		beginTimeStr = output.NowDate + " " + opts.BeginTime
	} else if tp == tickerBase.DatetimeFormat {
		beginTimeStr = opts.BeginTime
	}
	output.BeginStamp, err = tickerBase.TimeValue(beginTimeStr, output.BaseLocation)
	if err != nil {
		return
	}

	// make end time
	var endTimeStr string
	tp, err = tickerBase.TimeType(opts.EndTime)
	if err != nil {
		return
	}
	if tp == tickerBase.TimeFormat {
		endTimeStr = output.NowDate + " " + opts.EndTime
	} else if tp == tickerBase.DatetimeFormat {
		endTimeStr = opts.EndTime
	}
	output.EndStamp, err = tickerBase.TimeValue(endTimeStr, output.BaseLocation)
	if err != nil {
		return
	}

	return
}

// UpdateNowDate is a method named UpdateNowDate for the Base structure in Golang,
// which is intended to update the current date and time.
// * This method receives a pointer of the Base structure named receive as the receiver,
// and the type of the receiver is *Base,
// indicating that this method is a pointer method of the Base structure
// and can be directly called on a pointer variable of the Base structure.
// * In the implementation of the method, it first checks whether the BaseLocation attribute of the Base structure is nil.
// If it is nil, it returns an ErrNoBaseLocation error. If BaseLocation is not nil,
// the Now function in the time package is called to get the current date and time,
// and then the In method of the time.Time object is called to convert the current date and
// time to the time zone of the Base structure.
// * Finally, the Format method of the time.Time object is called to format the current date and
// time as a string in the DefaultDateFormatStr format, and
// the formatted string is stored in the NowDate attribute of the Base structure.
// * Finally, the method returns a variable named err, which can be ErrNoBaseLocation or nil,
// indicating whether the method was executed successfully.
func (receive *GoTicker) UpdateNowDate() (err error) {
	if receive.BaseLocation == nil {
		err = tickerBase.ErrNoBaseLocation
		return
	}

	now := time.Now().In(receive.BaseLocation)
	receive.NowDate = now.Format(tickerBase.DefaultDateFormatStr)

	return
}

// ReloadLocation is a method named ReloadLocation for the Base structure in Golang, which is intended to reload the time zone.
// * This method receives a pointer of the Base structure named receive as the receiver,
// and the type of the receiver is *Base, indicating that this method is a pointer method of the Base structure
// and can be directly called on a pointer variable of the Base structure.
// * In the implementation of the method, it first checks whether the Opts.Location attribute of the Base structure is empty.
// If it is empty, it returns an ErrNoOptLocation error.
// If Opts.Location is not empty, the LoadLocation function in the time package is called to reload the time zone.
// If an error occurs during the reloading process, it returns an ErrUnSupportedLocation error.
// * Finally, the method returns a variable named err, which can be ErrNoOptLocation, ErrUnSupportedLocation, or nil,
// indicating whether the method was executed successfully.
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
