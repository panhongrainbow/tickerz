package goTicker

import (
	tickerBase "github.com/panhongrainbow/tickerz/base"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// [Test_Check_NowDate] is a Golang unit test example that includes two sub-tests,
// both testing the ReloadLocation and UpdateNowDate functions of the Base structure.
//
// * In the first sub-test,
// there are three assertion statements that respectively test whether the UpdateNowDate function returns the expected error when bs.NowDate is empty,
// whether the ReloadLocation function returns the expected error when bs.Opts.Location is empty,
// and whether the ReloadLocation function can successfully load the time zone when bs.Opts.Location is not empty.
//
// * In the second sub-test, there are also three assertion statements that respectively test whether the TimeType function can successfully identify the format of bs.NowDate
// after it is updated, and whether bs.NowDate is in the DateFormat format.
//
// * The purpose of the entire test is to ensure that the ReloadLocation and UpdateNowDate functions of the Base structure can work correctly
// and that the format of bs.NowDate conforms to the expected date format.
func Test_Check_NowDate(t *testing.T) {
	t.Run("reload location", func(t *testing.T) {
		// empty bs time
		bs := GoTicker{}
		err := bs.UpdateNowDate()
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
		err = gt.UpdateNowDate()
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
			Repeat:    false,
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
			Repeat:    false,
			Duration:  0 * time.Nanosecond,
			BaseList:  []string{"2023-03-05 03:04:06", "2023-03-05 03:04:07", "2023-03-05 03:04:08"},
			BeginTime: "2023-03-05 03:04:06",
			EndTime:   "2023-03-05 03:04:07",
		}
		offOpts := tickerBase.OffOpts{}
		gtk, err := New(opts, offOpts)
		require.NoError(t, err)
		// check base stamp
		require.Equal(t, int64(1677956645), gtk.BaseStamp)
		// check location
		require.Equal(t, "Asia/Shanghai", gtk.BaseLocation.String())
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
