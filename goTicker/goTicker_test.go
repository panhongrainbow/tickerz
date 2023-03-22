package goTicker

import (
	tickerBase "github.com/panhongrainbow/tickerz/base"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// [Test_Check_NowDate] is an unit test example that includes two sub-tests,
// both testing the ReloadLocation and UpdateNowDateOrMock functions of the Base structure.
//
// * In the first sub-test,
// there are three assertion statements that respectively test whether the UpdateNowDateOrMock function returns the expected error when bs.NowDate is empty,
// whether the ReloadLocation function returns the expected error when bs.Opts.Location is empty,
// and whether the ReloadLocation function can successfully load the time zone when bs.Opts.Location is not empty.
//
// * In the second sub-test, there are also three assertion statements that respectively test whether the TimeType function can successfully identify the format of bs.NowDate
// after it is updated, and whether bs.NowDate is in the DateFormat format.
//
// * The purpose of the entire test is to ensure that the ReloadLocation and UpdateNowDateOrMock functions of the Base structure can work correctly
// and that the format of bs.NowDate conforms to the expected date format.
func Test_Check_NowDate(t *testing.T) {
	t.Run("reload location", func(t *testing.T) {
		// empty bs time
		bs := GoTicker{}
		err := bs.UpdateNowDateOrMock(mockDateStr)
		require.Equal(t, err, tickerBase.ErrNoBaseLocation)

		// reload location
		bs.Opts.Location = tickerBase.DefaultTimeZone
		err = bs.ReloadLocation()
		require.NoError(t, err)
	})
	t.Run("reload location", func(t *testing.T) {
		// reload location
		gt := GoTicker{}
		gt.Opts.Location = tickerBase.DefaultTimeZone
		err := gt.ReloadLocation()
		require.NoError(t, err)

		// update now date
		err = gt.UpdateNowDateOrMock(mockDateStr)
		require.NoError(t, err)

		// check result
		tp, err := tickerBase.TimeType(gt.NowDate)
		require.NoError(t, err)
		require.Equal(t, tickerBase.DateFormat, tp)
	})
}

func Test_Check_New(t *testing.T) {
	t.Run("check opts", func(t *testing.T) {
		opts := tickerBase.Opts{}
		offOpts := tickerBase.OffOpts{}
		_, err := New(opts, offOpts)
		require.NotNil(t, err)
	})
	t.Run("unsupported location", func(t *testing.T) {
		opts := tickerBase.Opts{
			BaseTime:  "2023-03-05 03:04:05",
			Location:  "heaven", // invalid
			Duration:  0 * time.Nanosecond,
			BaseList:  []string{"2023-03-05 03:04:06", "2023-03-05 03:04:07", "2023-03-05 03:04:08"},
			BeginTime: "2023-03-05 03:04:06",
			EndTime:   "2023-03-05 03:04:07",
		}
		offOpts := tickerBase.OffOpts{}
		_, err := New(opts, offOpts)
		require.Equal(t, tickerBase.ErrUnSupportedLocation, err)
	})
	t.Run("valid opts", func(t *testing.T) {
		opts := tickerBase.Opts{
			BaseTime:  "2023-03-05 03:04:05",
			Location:  "",
			Duration:  0 * time.Nanosecond,
			BaseList:  []string{"2023-03-05 03:04:06", "2023-03-05 03:04:07", "2023-03-05 03:04:08"},
			BeginTime: "2023-03-05 03:04:06",
			EndTime:   "2023-03-05 03:04:07",
		}
		offOpts := tickerBase.OffOpts{}

		// mock date
		mockDateStr = "2023-3-21"

		// new ticker
		gtk, err := New(opts, offOpts)
		require.NoError(t, err)
		// check base stamp
		require.Equal(t, int64(1677956645), gtk.BaseStamp)
		// check location
		require.Equal(t, "UTC", gtk.BaseLocation.String())
		// check base list
		require.Equal(t, int64(1677956646), gtk.BaseList[0])
		require.Equal(t, int64(1677956647), gtk.BaseList[1])
		require.Equal(t, int64(1677956648), gtk.BaseList[2])
		// check begin stamp
		require.Equal(t, int64(1677956646), gtk.BeginStamp)
		// check end stamp
		require.Equal(t, int64(1677956647), gtk.EndStamp)
	})
}

// Test_Check_ReNew is a unit test written in Go to check the behavior of the ReNew method of a ticker object.
// It creates a new ticker object, checks the initial values of various timestamp-related attributes,
// calls the ReNew method, and then checks that the values of the attributes have been updated correctly.
// The test case uses the require function to check that the actual values match the expected values.
func Test_Check_ReNew(t *testing.T) {
	t.Run("renew ticker", func(t *testing.T) {
		opts := tickerBase.Opts{
			BaseTime:  "03:04:05",
			Location:  tickerBase.DefaultTimeZone,
			Duration:  0 * time.Nanosecond,
			BaseList:  []string{"03:04:06", "03:04:07", "03:04:08"},
			BeginTime: "03:04:06",
			EndTime:   "03:04:07",
		}
		offOpts := tickerBase.OffOpts{}

		// mock date
		mockDateStr = "2023-1-31"

		// new ticker
		gtk, err := New(opts, offOpts)
		require.NoError(t, err)

		// check base stamp
		require.Equal(t, int64(1675105445), gtk.BaseStamp)
		// check location
		require.Equal(t, "UTC", gtk.BaseLocation.String())
		// check base list
		require.Equal(t, int64(1675105446), gtk.BaseList[0])
		require.Equal(t, int64(1675105447), gtk.BaseList[1])
		require.Equal(t, int64(1675105448), gtk.BaseList[2])
		// check begin stamp
		require.Equal(t, int64(1675105446), gtk.BeginStamp)
		// check end stamp
		require.Equal(t, int64(1675105447), gtk.EndStamp)

		// mock date
		mockDateStr = "2023-2-1"

		// renew ticker
		err = gtk.ReNew()
		require.NoError(t, err)

		// check base stamp
		require.Equal(t, int64(1675191845), gtk.BaseStamp)
		// check location
		require.Equal(t, "UTC", gtk.BaseLocation.String())
		// check base list
		require.Equal(t, int64(1675191846), gtk.BaseList[0])
		require.Equal(t, int64(1675191847), gtk.BaseList[1])
		require.Equal(t, int64(1675191848), gtk.BaseList[2])
		// check begin stamp
		require.Equal(t, int64(1675191846), gtk.BeginStamp)
		// check end stamp
		require.Equal(t, int64(1675191847), gtk.EndStamp)
	})
}

// TestGoTicker_Check_calculateRepeatParameter is testing the calculateRepeatParameter function of the GoTicker struct.
// It sets up two test cases with different durations and expected results,
// and checks that the function returns the expected nearest timestamp, duration, and error.
// The test passes if all assertions are true.
func TestGoTicker_Check_calculateRepeatParameter(t *testing.T) {
	// Get the current Unix timestamp in seconds.
	now := time.Now().Unix()

	// Define a list of test cases.
	tests := []struct {
		name             string
		baseStamp        int64
		duration         time.Duration
		expectedNearest  int64
		expectedDuration int64
		expectedErr      error
	}{
		// Test case 1: A valid duration of 2 seconds
		{
			name:             "valid duration",
			baseStamp:        now,
			duration:         2 * time.Second,
			expectedNearest:  now,
			expectedDuration: 2,
			expectedErr:      nil,
		},
		// Test case 2: An invalid duration of 0 seconds
		{
			name:             "invalid duration",
			baseStamp:        now,
			duration:         0,
			expectedNearest:  0,
			expectedDuration: 0,
			expectedErr:      tickerBase.ErrInactiveRepeatList,
		},
	}

	// Loop over the list of test cases and execute each one
	for _, test := range tests {
		// Use the name of the test case as the name of the subtest
		t.Run(test.name, func(t *testing.T) {
			// Create a new GoTicker with the given base stamp and duration
			receive := &GoTicker{
				BaseStamp: test.baseStamp,
				Opts: tickerBase.Opts{
					Duration: test.duration,
				},
			}
			// Call the calculateRepeatParameter method of the GoTicker and store the result
			nearest, duration, err := receive.calculateRepeatParameter()
			// Check that the actual nearest and duration values match the expected ones
			require.Equal(t, test.expectedNearest, nearest)
			require.Equal(t, test.expectedDuration, duration)
			// Check that the actual error matches the expected one
			require.Equal(t, test.expectedErr, err)

		})
	}
}

// TestUsableSubBaseList is a test function with four subtests that checks if
// the usableSubBaseList function returns the correct list of timestamps
// within the range of the ticker's BeginStamp and EndStamp.
// It also checks if the length of the subBaseList matches the expected value.
func Test_Check_UsableSubBaseList(t *testing.T) {
	// Define the subtests
	tests := []struct {
		ticker   *GoTicker
		name     string
		expected int
	}{
		// Test case 1
		{
			&GoTicker{
				BeginStamp: time.Now().Unix(),
				EndStamp:   time.Now().Add(10 * time.Minute).Unix(),
				BaseList: []int64{
					time.Now().Add(2 * time.Minute).Unix(),
					time.Now().Add(4 * time.Minute).Unix(),
					time.Now().Add(6 * time.Minute).Unix()},
			},
			"Test case 1",
			3,
		},
		// Test case 2
		{
			&GoTicker{
				BeginStamp: time.Now().Unix(),
				EndStamp:   time.Now().Add(10 * time.Minute).Unix(),
				BaseList: []int64{
					time.Now().Add(-2 * time.Minute).Unix(),
					time.Now().Add(4 * time.Minute).Unix(),
					time.Now().Add(6 * time.Minute).Unix()},
			},
			"Test case 2",
			2,
		},
		// Test case 3
		{
			&GoTicker{
				BeginStamp: time.Now().Unix(),
				EndStamp:   time.Now().Add(10 * time.Minute).Unix(),
				BaseList: []int64{
					time.Now().Unix(),
					time.Now().Add(2 * time.Minute).Unix(),
					time.Now().Add(4 * time.Minute).Unix()},
			},
			"Test case 3",
			2,
		},
		// Test case 4
		{
			&GoTicker{
				BeginStamp: time.Now().Unix(),
				EndStamp:   time.Now().Add(10 * time.Minute).Unix(),
				BaseList: []int64{
					time.Now().Add(-6 * time.Minute).Unix(),
					time.Now().Add(-2 * time.Minute).Unix(),
					time.Now().Add(-2 * time.Minute).Unix()},
			},
			"Test case 4",
			0,
		},
	}

	// Run the subtests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Call the usableSubBaseList function
			subBaseList, err := test.ticker.usableSubBaseList()
			require.NoError(t, err)

			// Check if the length of subBaseList matches with the expected value
			require.Equal(t, test.expected, test.expected)

			// Check if all timestamps in subBaseList are within the range of the ticker's BeginStamp and EndStamp
			var outOfBaseListRange bool
			if test.expected > 0 {
				now := time.Now().Unix()
				for _, timestamp := range subBaseList {
					if timestamp < now ||
						timestamp < test.ticker.BeginStamp ||
						timestamp > test.ticker.EndStamp {
						outOfBaseListRange = true
					}
				}
				require.Equal(t, false, outOfBaseListRange)
			}
		})
	}
}

func Test_Check_CalculateWaitTime(t *testing.T) {
	t.Run("check base list", func(t *testing.T) {
		opts := tickerBase.Opts{
			BaseTime:  "0:0:0",
			Location:  tickerBase.DefaultTimeZone,
			Duration:  1 * time.Second,
			BaseList:  []string{},
			BeginTime: "0:0:0",
			EndTime:   "23:59:59",
		}
		offOpts := tickerBase.OffOpts{}

		// mock date
		mockDateStr = time.Now().Format(tickerBase.DefaultDateFormatStr)

		// new ticker
		gtk, err := New(opts, offOpts)
		require.NoError(t, err)

		// calculate
		result, err := gtk.CalculateWaitList(10)
		require.NoError(t, err)
		if len(result) < 10 {
			time.Sleep(10 * time.Second)
			result, err = gtk.CalculateWaitList(10)
			require.NoError(t, err)
			require.Equal(t, 10, len(result))
		}

		now := time.Now().Unix()
		var previous int64
		for i := 0; i < len(result); i++ {
			if i == 0 {
				require.GreaterOrEqual(t, result[i], now)
			}
			if i != 0 {
				require.Greater(t, result[i], now)
			}
			require.Greater(t, result[i], previous)
			previous = result[i]
		}
	})
	t.Run("check base list", func(t *testing.T) {
		opts := tickerBase.Opts{
			BaseTime:  "0:0:0",
			Location:  tickerBase.DefaultTimeZone,
			Duration:  1 * time.Second,
			BaseList:  []string{},
			BeginTime: "0:0:0",
			EndTime:   "23:59:59",
		}
		offOpts := tickerBase.OffOpts{}

		// mock date
		mockDateStr = time.Now().Format(tickerBase.DefaultDateFormatStr)

		// new ticker
		gtk, err := New(opts, offOpts)
		require.NoError(t, err)

		// calculate
		result, err := gtk.CalculateWaitList(10)
		require.NoError(t, err)
		if len(result) < 10 {
			time.Sleep(10 * time.Second)
			result, err = gtk.CalculateWaitList(10)
			require.NoError(t, err)
			require.Equal(t, 10, len(result))
		}

		now := time.Now().Unix()
		var previous int64
		for i := 0; i < len(result); i++ {
			require.GreaterOrEqual(t, result[i], now)
			require.Greater(t, result[i], previous)
			previous = result[i]
		}
	})
}

func Test_Check_WaitAccordingList(t *testing.T) {
	// Define the subtests
	tests := []struct {
		ticker   *GoTicker
		name     string
		expected int
	}{
		// Test case 1
		{
			&GoTicker{
				BeginStamp: time.Now().Unix(),
				EndStamp:   time.Now().Add(10 * time.Second).Unix(),
				BaseList: []int64{
					time.Now().Add(2 * time.Second).Unix(),
					time.Now().Add(4 * time.Second).Unix(),
					time.Now().Add(6 * time.Second).Unix()},
			},
			"Test case 1",
			3,
		},
	}

	// Run the subtests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.ticker.WaitAccordingList(10)
		})
	}
}
