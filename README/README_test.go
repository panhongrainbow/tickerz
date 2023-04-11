package README

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// BenchmarkTickerImitator benchmarks the performance of using time.After to imitate a tickerz.
func BenchmarkTickerImitator(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-time.After(time.Millisecond * 10)
	}
}

// BenchmarkLongSleep benchmarks the performance of using time.Sleep to pause the execution for 10 milliseconds.
func BenchmarkLongSleep(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		time.Sleep(time.Millisecond * 10)
	}
}

// BenchmarkFragmentedSleep benchmarks the performance of using time.Sleep to pause the execution for 1 millisecond in a loop that runs 10 times.
func BenchmarkFragmentedSleep(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			time.Sleep(time.Millisecond * 1)
		}
	}
}

// BenchmarkCompareCPU runs all the other benchmark functions and compares their results.
func Benchmark_Compare_CPU(b *testing.B) {
	// Run the benchmark for using time.Sleep to pause the execution for 10 milliseconds
	longSleep := b.Run("LongSleep", BenchmarkLongSleep)
	require.Equal(b, longSleep, true)

	// Run the benchmark for using time.After to imitate a ticker
	tickerImitator := b.Run("TickerImitator", BenchmarkTickerImitator)
	require.Equal(b, tickerImitator, true)

	// Run the benchmark for using time.Sleep to pause the execution
	fragmentedSleep := b.Run("FragmentedSleep", BenchmarkFragmentedSleep)
	require.Equal(b, fragmentedSleep, true)
}
