package base

import (
	"github.com/stretchr/testify/require"
	_ "net/http/pprof"
	"testing"
	"time"
)

// [Test_Check_TimeValue_Timestamps] is used to verify whether the TimeValue function can correctly calculate timestamps.
func Test_Check_TimeValue_Timestamps(t *testing.T) {
	// TimeValue of 1970-1-1 8:0:0 CST is -28800
	shanghai, err := time.LoadLocation("Asia/Shanghai")
	past, err := TimeValue("1970-1-1 0:0:0", shanghai)
	require.NoError(t, err)
	require.Equal(t, int64(-28800), past)

	// TimeValue of 1970-1-1 0:0:0 UCT is 0
	utc, err := time.LoadLocation("UTC")
	require.NoError(t, err)
	past, err = TimeValue("1970-1-1 0:0:0", utc)
	require.NoError(t, err)
	require.Equal(t, int64(0), past)

	// TimeValue of 1970-1-1 8:2:5 CST is 125
	tests := []struct {
		timeStr string
	}{
		{"1970-01-01 08:02:05"},
		{"1970-1-1 8:2:5"},
		{"1970-01-1 08:2:05"},
		{"1970-01-01 8:2:5"},
		{"1970-1-1 08:02:05"},
	}
	for i := 0; i < len(tests); i++ {
		past, err = TimeValue(tests[i].timeStr, nil)
		require.NoError(t, err)
		require.Equal(t, int64(125), past)
	}
}

// [Test_TimeType] is a unit test function to verify if the TimeValue function can calculate timestamps correctly.
// It includes multiple test cases for different time formats, such as date, time, and datetime formats.
// The test cases check whether the TimeType function returns the expected time type and error,
// and the require.Equal function is used to compare the actual result with the expected result.
func Test_Check_TimeType_Time(t *testing.T) {
	// test all time types
	tests := []struct {
		timeStr  string
		timeType uint
		err      error
	}{
		// valid cases
		{
			"",
			EmptyTimeFormat,
			nil,
		},
		{
			"3:4:5",
			TimeFormat,
			nil,
		},
		{
			"03:04:05",
			TimeFormat,
			nil,
		},
		{
			"03:04:05",
			TimeFormat,
			nil,
		},
		{
			"3:4:05",
			TimeFormat,
			nil,
		},
		// invalid cases
		{
			"3:4:61",
			0,
			ErrUnSupportedTimeFormat,
		},
		{
			"3:4:61",
			0,
			ErrUnSupportedTimeFormat,
		},
		{
			"3:61:5",
			0,
			ErrUnSupportedTimeFormat,
		},
		{
			"25:4:5",
			0,
			ErrUnSupportedTimeFormat,
		},
		{
			"25:4:5",
			0,
			ErrUnSupportedTimeFormat,
		},
	}
	// verify
	for i := 0; i < len(tests); i++ {
		tp, err := TimeType(tests[i].timeStr)
		require.Equal(t, tp, tests[i].timeType)
		require.Equal(t, err, tests[i].err)
	}
}

// [Test_TimeType] is a unit test function to verify if the TimeValue function can calculate timestamps correctly.
// It includes multiple test cases for different time formats, such as date, time, and datetime formats.
// The test cases check whether the TimeType function returns the expected time type and error,
// and the require.Equal function is used to compare the actual result with the expected result.
func Test_Check_TimeType_Date(t *testing.T) {
	// test all time types
	tests := []struct {
		timeStr  string
		timeType uint
		err      error
	}{
		// valid cases
		{
			"2023-3-5",
			DateFormat,
			nil,
		},
		{
			"2023-03-5",
			DateFormat,
			nil,
		},
		{
			"2023-3-05",
			DateFormat,
			nil,
		},
		{
			"2023-03-05",
			DateFormat,
			nil,
		},
		// invalid cases
		{
			"23-03-05",
			0,
			ErrUnSupportedTimeFormat,
		},
		{
			"2023-0-05",
			0,
			ErrUnSupportedTimeFormat,
		},
		{
			"2023-003-05",
			0,
			ErrUnSupportedTimeFormat,
		},
		{
			"2023-13-05",
			0,
			ErrUnSupportedTimeFormat,
		},
		{
			"2023-03-005",
			0,
			ErrUnSupportedTimeFormat,
		},
		{
			"2023-03-32",
			0,
			ErrUnSupportedTimeFormat,
		},
	}
	// verify
	for i := 0; i < len(tests); i++ {
		tp, err := TimeType(tests[i].timeStr)
		require.Equal(t, tp, tests[i].timeType)
		require.Equal(t, err, tests[i].err)
	}
}

// [Test_Check_TimeType_DateTime] is a unit test function to verify if the TimeValue function can calculate timestamps correctly.
// It includes multiple test cases for different time formats, such as date, time, and datetime formats.
// The test cases check whether the TimeType function returns the expected time type and error,
// and the require.Equal function is used to compare the actual result with the expected result.
func Test_Check_TimeType_DateTime(t *testing.T) {
	// test all time types
	tests := []struct {
		timeStr  string
		timeType uint
		err      error
	}{
		// valid cases
		{
			"2023-03-05 03:04:05",
			DatetimeFormat,
			nil,
		},
		{
			"2023-3-5 3:4:5",
			DatetimeFormat,
			nil,
		},
		{
			"2023-3-05 3:04:5",
			DatetimeFormat,
			nil,
		},
		{
			"2023-03-05 3:4:5",
			DatetimeFormat,
			nil,
		},
		{
			"2023-3-5 03:04:05",
			DatetimeFormat,
			nil,
		},
		// invalid cases
		{
			"02023-03-05 03:04:05",
			0,
			ErrUnSupportedTimeFormat,
		},
		{
			"02023-03-05  03:04:05",
			0,
			ErrUnSupportedTimeFormat,
		},
		{
			" 2023-03-05 03:04:05",
			0,
			ErrUnSupportedTimeFormat,
		},
		{
			"2023-03-05 03:04:05 ",
			0,
			ErrUnSupportedTimeFormat,
		},
	}
	// verify
	for i := 0; i < len(tests); i++ {
		tp, err := TimeType(tests[i].timeStr)
		require.Equal(t, tp, tests[i].timeType)
		require.Equal(t, err, tests[i].err)
	}
}

// [Test_Check_CheckOpts] is testing the CheckOpts function,
// which is a validation function for the input parameters of a time-related operation.
func Test_Check_CheckOpts_BaseTime(t *testing.T) {
	// test cases
	tests := []struct {
		opts Opts
		err  error
	}{
		// valid
		{
			opts: Opts{
				BaseTime:  "03:04:05", // valid
				Location:  "",
				Duration:  0 * time.Nanosecond,
				BaseList:  []string{},
				BeginTime: "",
				EndTime:   "",
			},
			err: nil,
		},
		{
			opts: Opts{
				BaseTime:  "2023-03-05 03:04:05", // valid
				Location:  "",
				Duration:  0 * time.Nanosecond,
				BaseList:  []string{},
				BeginTime: "",
				EndTime:   "",
			},
			err: nil,
		},
		{
			opts: Opts{
				BaseTime:  "2023-03-05          03:04:05", // valid but not recommend
				Location:  "",
				Duration:  0 * time.Nanosecond,
				BaseList:  []string{},
				BeginTime: "",
				EndTime:   "",
			},
			err: nil,
		},
		// invalid
		{
			opts: Opts{
				BaseTime:  "", // invalid
				Location:  "",
				Duration:  0 * time.Nanosecond,
				BaseList:  []string{},
				BeginTime: "",
				EndTime:   "",
			},
			err: ErrUnsupBasetime,
		},
		{
			opts: Opts{
				BaseTime:  "03:04:05 ", // invalid
				Location:  "",
				Duration:  0 * time.Nanosecond,
				BaseList:  []string{},
				BeginTime: "",
				EndTime:   "",
			},
			err: ErrUnSupportedTimeFormat,
		},
		{
			opts: Opts{
				BaseTime:  " 03:04:05", // invalid
				Location:  "",
				Duration:  0 * time.Nanosecond,
				BaseList:  []string{},
				BeginTime: "",
				EndTime:   "",
			},
			err: ErrUnSupportedTimeFormat,
		},
		{
			opts: Opts{
				BaseTime:  "2023-03-05", // invalid
				Location:  "",
				Duration:  0 * time.Nanosecond,
				BaseList:  []string{},
				BeginTime: "",
				EndTime:   "",
			},
			err: ErrUnsupBasetime,
		},
		{
			opts: Opts{
				BaseTime:  "2023-03-05 03:04:05 ", // invalid
				Location:  "",
				Duration:  0 * time.Nanosecond,
				BaseList:  []string{},
				BeginTime: "",
				EndTime:   "",
			},
			err: ErrUnSupportedTimeFormat,
		},
		{
			opts: Opts{
				BaseTime:  " 2023-03-05 03:04:05", // invalid
				Location:  "",
				Duration:  0 * time.Nanosecond,
				BaseList:  []string{},
				BeginTime: "",
				EndTime:   "",
			},
			err: ErrUnSupportedTimeFormat,
		},
	}

	// verify
	for i := 0; i < len(tests); i++ {
		require.Equal(t, tests[i].err, (tests[i].opts).CheckOpts())
	}
}

// [Test_Check_CheckOpts_Location] is testing the CheckOpts function,
// which is a validation function for the input parameters of a time-related operation.
func Test_Check_CheckOpts_Location(t *testing.T) {
	// test cases
	tests := []struct {
		opts Opts
		err  error
	}{
		// valid
		{
			opts: Opts{
				BaseTime:  "2023-03-05 03:04:05",
				Location:  "", // invalid
				Duration:  0 * time.Nanosecond,
				BaseList:  []string{},
				BeginTime: "",
				EndTime:   "",
			},
			err: nil,
		},
		// invalid
		{
			opts: Opts{
				BaseTime:  "2023-03-05 03:04:05",
				Location:  "heaven", // invalid
				Duration:  0 * time.Nanosecond,
				BaseList:  []string{},
				BeginTime: "",
				EndTime:   "",
			},
			err: ErrUnSupportedLocation,
		},
		{
			opts: Opts{
				BaseTime:  "2023-03-05 03:04:05",
				Location:  "",
				Duration:  -1 * time.Nanosecond, // negative duration, invalid
				BaseList:  []string{},
				BeginTime: "",
				EndTime:   "",
			},
			err: ErrNegativeDuration,
		},
		{
			opts: Opts{
				BaseTime:  "1970-01-01 01:01:01",
				Location:  "",
				Duration:  1 * time.Nanosecond,
				BaseList:  []string{"1970-01-01"}, // invalid
				BeginTime: "",
				EndTime:   "",
			},
			err: ErrUnSupportedBaseList,
		},
		{
			opts: Opts{
				BaseTime:  "1970-01-01 01:01:01", // supported
				Location:  "",                    // support empty
				Duration:  1 * time.Nanosecond,   // supported
				BaseList:  []string{"1970-1-1 1:1:1", "1970-01-01 01:01:01"},
				BeginTime: "",
				EndTime:   "",
			},
			err: ErrIncorrectBaseListOrder,
		},
		{
			opts: Opts{
				BaseTime:  "1970-01-01 01:01:01", // supported
				Location:  "",                    // support empty
				Duration:  1 * time.Nanosecond,   // supported
				BaseList:  []string{"1970-1-1 1:1:1", "1970-01-01 01:01:00"},
				BeginTime: "",
				EndTime:   "",
			},
			err: ErrIncorrectBaseListOrder,
		},
		{
			opts: Opts{
				BaseTime:  "1970-01-01 01:01:01", // supported
				Location:  "",                    // support empty
				Duration:  1 * time.Nanosecond,   // supported
				BaseList:  []string{"1970-1-1 1:1:1", "01:01:00"},
				BeginTime: "",
				EndTime:   "",
			},
			err: ErrBaseListDifferentTypes,
		},
		{
			opts: Opts{
				BaseTime:  "1970-01-01 01:01:01", // supported
				Location:  "",                    // support empty
				Duration:  0 * time.Nanosecond,   // supported
				BaseList:  []string{},
				BeginTime: "1970-01-01 01:01:01",
				EndTime:   "1970-01-01 01:01:02",
			},
			err: nil,
		},
		{
			opts: Opts{
				BaseTime:  "1970-01-01 01:01:01", // supported
				Location:  "",                    // support empty
				Duration:  0 * time.Nanosecond,   // supported
				BaseList:  []string{},
				BeginTime: "01:01:01",
				EndTime:   "01:01:02",
			},
			err: nil,
		},
		{
			opts: Opts{
				BaseTime:  "1970-01-01 01:01:01", // supported
				Location:  "",                    // support empty
				Duration:  0 * time.Nanosecond,   // supported
				BaseList:  []string{},
				BeginTime: "1970-01-01 01:01:01",
				EndTime:   "01:01:02",
			},
			err: ErrBeginEndTimeTypeNotEqual,
		},
		{
			opts: Opts{
				BaseTime:  "1970-01-01 01:01:01", // supported
				Location:  "",                    // support empty
				Duration:  0 * time.Nanosecond,   // supported
				BaseList:  []string{},
				BeginTime: "",
				EndTime:   "",
			},
			err: nil,
		},
		{
			opts: Opts{
				BaseTime:  "1970-01-01     01:01:01", // supported
				Location:  "",                        // support empty
				Duration:  0 * time.Nanosecond,       // supported
				BaseList:  []string{},
				BeginTime: "",
				EndTime:   "",
			},
			err: nil,
		},
		{
			opts: Opts{
				BaseTime:  "1970-01-01 01:01:01", // supported
				Location:  "",                    // support empty
				Duration:  0 * time.Nanosecond,   // supported
				BaseList:  []string{},
				BeginTime: "1970-01-01 01:01:02",
				EndTime:   "1970-01-01 01:01:01",
			},
			err: ErrIncorrectBeginEndTimeOrder,
		},
	}

	for i := 0; i < len(tests); i++ {
		require.Equal(t, tests[i].err, (tests[i].opts).CheckOpts())
	}
}

// [Test_Check_CheckOpts_Baselist] is testing the CheckOpts function,
// which is a validation function for the input parameters of a time-related operation.
func Test_Check_CheckOpts_Baselist(t *testing.T) {
	// test cases
	tests := []struct {
		opts Opts
		err  error
	}{
		// valid
		{
			opts: Opts{
				BaseTime:  "2023-03-05 03:04:05",
				Location:  "", // valid
				Duration:  0 * time.Nanosecond,
				BaseList:  []string{},
				BeginTime: "",
				EndTime:   "",
			},
			err: nil,
		},
		{
			opts: Opts{
				BaseTime:  "2023-03-05 03:04:05", // valid
				Location:  "",
				Duration:  1 * time.Nanosecond, // valid
				BaseList:  []string{"2023-03-06 03:04:06", "2023-03-07 03:04:07"},
				BeginTime: "",
				EndTime:   "",
			},
			err: nil,
		},
		// invalid
		{
			opts: Opts{
				BaseTime:  "2023-03-05 03:04:05",
				Location:  "heaven", // invalid
				Duration:  0 * time.Nanosecond,
				BaseList:  []string{},
				BeginTime: "",
				EndTime:   "",
			},
			err: ErrUnSupportedLocation,
		},
		{
			opts: Opts{
				BaseTime:  "2023-03-05 03:04:05",
				Location:  "",
				Duration:  -1 * time.Nanosecond, // negative duration, invalid
				BaseList:  []string{},
				BeginTime: "",
				EndTime:   "",
			},
			err: ErrNegativeDuration,
		},
		{
			opts: Opts{
				BaseTime:  "2023-03-05 03:04:05",
				Location:  "",
				Duration:  0 * time.Nanosecond,
				BaseList:  []string{"2023-03-06"}, // invalid
				BeginTime: "",
				EndTime:   "",
			},
			err: ErrUnSupportedBaseList,
		},
		{
			opts: Opts{
				BaseTime:  "2023-03-05 03:04:05",
				Location:  "",
				Duration:  0 * time.Nanosecond,
				BaseList:  []string{"2023-03-05 03:04:05", "2023-03-05 03:04:05"}, // invalid
				BeginTime: "",
				EndTime:   "",
			},
			err: ErrIncorrectBaseListOrder,
		},
		{
			opts: Opts{
				BaseTime:  "2023-03-05 03:04:05", // supported
				Location:  "",
				Duration:  0 * time.Nanosecond,
				BaseList:  []string{"2023-03-07 03:04:05", "2023-03-06 03:04:05"},
				BeginTime: "",
				EndTime:   "",
			},
			err: ErrIncorrectBaseListOrder,
		},
		{
			opts: Opts{
				BaseTime:  "2023-03-05 03:04:05", // supported
				Location:  "",                    // support empty
				Duration:  0 * time.Nanosecond,   // supported
				BaseList:  []string{"2023-03-05 03:04:07", "2023-03-05 03:04:06"},
				BeginTime: "",
				EndTime:   "",
			},
			err: ErrIncorrectBaseListOrder,
		},
		{
			opts: Opts{
				BaseTime:  "2023-03-05 03:04:05",
				Location:  "",
				Duration:  0 * time.Nanosecond,
				BaseList:  []string{"2023-03-05 03:04:06", "03:04:07"},
				BeginTime: "",
				EndTime:   "",
			},
			err: ErrBaseListDifferentTypes,
		},
	}

	// verify
	for i := 0; i < len(tests); i++ {
		require.Equal(t, tests[i].err, (tests[i].opts).CheckOpts())
	}
}

// [Test_Check_CheckOpts_BeginEndTime] is testing the CheckOpts function,
// which is a validation function for the input parameters of a time-related operation.
func Test_Check_CheckOpts_BeginEndTime(t *testing.T) {
	// test cases
	tests := []struct {
		opts Opts
		err  error
	}{
		// valid
		{
			opts: Opts{
				BaseTime:  "2023-03-05 03:04:05", // supported
				Location:  "",                    // support empty
				Duration:  0 * time.Nanosecond,   // supported
				BaseList:  []string{},
				BeginTime: "03:04:06",
				EndTime:   "03:04:07",
			},
			err: nil,
		},
		{
			opts: Opts{
				BaseTime:  "2023-03-05 03:04:05",
				Location:  "",
				Duration:  0 * time.Nanosecond,
				BaseList:  []string{},
				BeginTime: "2023-03-05 03:04:06",
				EndTime:   "2023-03-05 03:04:07",
			},
			err: nil,
		},
		{
			opts: Opts{
				BaseTime:  "2023-03-05 03:04:05",
				Location:  "",
				Duration:  0 * time.Nanosecond,
				BaseList:  []string{},
				BeginTime: "2023-03-05 03:04:06",
				EndTime:   "",
			},
			err: nil,
		},
		{
			opts: Opts{
				BaseTime:  "2023-03-05 03:04:05",
				Location:  "",
				Duration:  0 * time.Nanosecond,
				BaseList:  []string{},
				BeginTime: "",
				EndTime:   "2023-03-05 03:04:07",
			},
			err: nil,
		},
		// invalid
		{
			opts: Opts{
				BaseTime:  "2023-03-05 03:04:05", // supported
				Location:  "",                    // support empty
				Duration:  0 * time.Nanosecond,   // supported
				BaseList:  []string{},
				BeginTime: "1970-01-01 01:01:01",
				EndTime:   "01:01:02",
			},
			err: ErrBeginEndTimeTypeNotEqual,
		},
		{
			opts: Opts{
				BaseTime:  "2023-03-05 03:04:05", // supported
				Location:  "",                    // support empty
				Duration:  0 * time.Nanosecond,   // supported
				BaseList:  []string{},
				BeginTime: "1970-01-01 01:01:02",
				EndTime:   "1970-01-01 01:01:01",
			},
			err: ErrIncorrectBeginEndTimeOrder,
		},
	}

	// verify
	for i := 0; i < len(tests); i++ {
		require.Equal(t, tests[i].err, (tests[i].opts).CheckOpts())
	}
}
