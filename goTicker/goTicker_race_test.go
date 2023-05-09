package goTicker

import (
	tickerBase "github.com/panhongrainbow/tickerz/base"
	"testing"
	"time"
)

// Test_Race_waitForNextDay is to check for data race for that method.
func Test_Race_waitForNextDay(t *testing.T) {
	// Get the current time
	now := time.Now()

	// Create a new GoTicker and set its properties
	gt := &GoTicker{
		BaseStamp: now.Unix(),
		BaseList: []int64{
			now.Add(1 * time.Second).Unix(),
		},
		BeginStamp: now.Add(0 * time.Second).Unix(),
		EndStamp:   now.Add(2 * time.Second).Unix(),
	}

	// Create a channel to receive signals from the ticker
	gt.SignalChan = make(chan tickerBase.TickerSignal)

	// Reload the location information for the ticker
	err := gt.ReloadLocation()
	if err != nil {
		panic(err)
	}

	// Update the ticker's current date
	err = gt.UpdateNowDateOrMock("")
	if err != nil {
		panic(err)
	}

	// Create 1000 goroutines to call waitForNextDay in order to check for data race
	for i := 0; i < 1000; i++ {
		go func() {
			err = gt.waitForNextDay()
			if err != nil {
				panic(err)
			}
		}()
	}
}
