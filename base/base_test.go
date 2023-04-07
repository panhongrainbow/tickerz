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

/*
Test_Check_checkOptsBaseList checks if the checkOptsBaseList function correctly validates the format of baseList.
It includes test cases for valid and invalid baseList formats and orders.
*/
func Test_Check_checkOptsBaseList(t *testing.T) {
	// Define test cases
	tests := []struct {
		name          string
		baseList      []string
		date          string
		expectedError error
	}{
		// Test case: Valid baseList with TimeFormat
		{
			name:          "Valid baseList with TimeFormat",
			baseList:      []string{"09:00:00", "10:00:00", "11:00:00"},
			date:          "2022-01-01",
			expectedError: nil,
		},
		// Test case: Valid baseList with DatetimeFormat
		{
			name:          "Valid baseList with DatetimeFormat",
			baseList:      []string{"2022-01-01 09:00:00", "2022-01-01 10:00:00", "2022-01-01 11:00:00"},
			date:          "",
			expectedError: nil,
		},
		// Test case: Invalid baseList with unsupported format
		{
			name:          "Invalid baseList with unsupported format",
			baseList:      []string{"2023-09-01", "2023-10-01", "2023-11-01"},
			date:          "2023-01-01",
			expectedError: ErrUnSupportedBaseList,
		},
		// Test case: Invalid baseList with different time formats
		{
			name:          "Invalid baseList with different time formats",
			baseList:      []string{"09:00:00", "2022-01-01 10:00:00", "11:00:00"},
			date:          "2022-01-01",
			expectedError: ErrBaseListDifferentTypes,
		},
		// Test case: Invalid baseList with incorrect order
		{
			name:          "Invalid baseList with incorrect order",
			baseList:      []string{"11:00:00", "10:00:00", "09:00:00"},
			date:          "2022-01-01",
			expectedError: ErrIncorrectBaseListOrder,
		},
	}
	// Run test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Call function with test inputs
			err := checkOptsBaseList(test.baseList, test.date)
			// Check if returned error matches expected error
			require.Equal(t, test.expectedError, err)
		})
	}
}

// Test_Check_checkOptsBeginEndTimeToTime is to check if the checkOptsBeginEndTimeToTime function returns the correct time.Time object for different inputs.
func Test_Check_checkOptsBeginEndTimeToTime(t *testing.T) {
	// Define test cases
	tests := []struct {
		name            string
		beginEndTimeStr string
		date            string
		expectedType    uint
		expectedTime    time.Time
		expectedError   error
	}{
		// Test case for valid TimeFormat input
		{
			name:            "Valid TimeFormat",
			beginEndTimeStr: "12:34:56",
			date:            "2022-01-01",
			expectedType:    TimeFormat,
			expectedTime:    time.Date(2022, time.January, 1, 12, 34, 56, 0, time.UTC),
			expectedError:   nil,
		},
		// Test case for valid DatetimeFormat input
		{
			name:            "Valid DatetimeFormat",
			beginEndTimeStr: "2022-01-01 12:34:56",
			date:            "",
			expectedType:    DatetimeFormat,
			expectedTime:    time.Date(2022, time.January, 1, 12, 34, 56, 0, time.UTC),
			expectedError:   nil,
		},
		// Test case for unsupported format input
		{
			name:            "Unsupported format",
			beginEndTimeStr: "01/01/2022",
			date:            "",
			expectedType:    0,
			expectedTime:    time.Time{},
			expectedError:   ErrUnSupportedBeginOrEndTimeFormat,
		},
	}

	// Run test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Call function with test inputs
			beginEndTimeType, beginEndTime, err := checkOptsBeginEndTimeToTime(test.beginEndTimeStr, test.date)
			// Check if returned type matches expected type
			require.Equal(t, test.expectedType, beginEndTimeType)
			// Check if returned time matches expected time
			require.Equal(t, test.expectedTime, beginEndTime)
			// Check if returned error matches expected error
			require.Equal(t, test.expectedError, err)
		})
	}
}
