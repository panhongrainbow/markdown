package md2json

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

// Test_Check_StringSlicePool tests the behavior of StringSlicePool's Get and Put methods.
func Test_Check_StringSlicePool(t *testing.T) {
	// Create a new StringSlicePool with a specified capacity.
	capacity := 5
	stringPool := NewStringSlicePool(capacity)

	// Test the initial state after Get
	strSlice := stringPool.Get()
	assert.Empty(t, strSlice, "Expected an empty slice from the pool")

	// Perform operations on the retrieved slice
	strSlice = append(strSlice, "Hello", "World")

	// Assert the length and content of the modified slice
	assert.Len(t, strSlice, 2, "Expected the slice to have 2 elements")
	assert.Equal(t, "Hello", strSlice[0], "Expected 'Hello' in the first element")
	assert.Equal(t, "World", strSlice[1], "Expected 'World' in the second element")

	// Put the modified slice back to the pool
	stringPool.Put(strSlice)

	// Test the state after Put
	strSlice = stringPool.Get()
	assert.Empty(t, strSlice, "Expected an empty slice from the pool")
}

// Test_Race_StringSlicePool runs the test with data race detection using the specified regular expression.
// go test -run='^\QTest_Race_StringSlicePool\E$' -race
func Test_Race_StringSlicePool(t *testing.T) {
	// Enable concurrent execution of the test.
	t.Parallel()

	// Set the capacity of the string slice pool.
	capacity := 5
	stringPool := NewStringSlicePool(capacity)

	// Define the number of goroutines and iterations for testing.
	const goroutines = 100
	const iterations = 100

	// Use a WaitGroup to wait for all goroutines to finish.
	var wg sync.WaitGroup

	// Launch multiple goroutines.
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			// Mark the completion of the current goroutine when it finishes.
			defer wg.Done()

			// Perform iterations on each goroutine.
			for j := 0; j < iterations; j++ {
				// Get a string slice from the pool.
				strSlice := stringPool.Get()

				// Modify the string slice by appending "Hello" and "World".
				strSlice = append(strSlice, "Hello", "World")

				// Put the modified string slice back into the pool.
				stringPool.Put(strSlice)
			}
		}()
	}

	// Wait for all goroutines to finish.
	wg.Wait()
}

// Benchmark_Comparator_StringSlicePool compares pool-based and non-pool-based string slice management in Go benchmarks.
func Benchmark_Comparator_StringSlicePool(b *testing.B) {
	// Set the capacity of the string slice pool.
	capacity := 10
	stringPool := NewStringSlicePool(capacity)

	// Reset the timer before starting the benchmark.
	b.ResetTimer()

	// Run a sub-benchmark with pool usage.
	b.Run("WithPool", func(b *testing.B) {
		// Perform benchmarking for 'b.N' iterations.
		for i := 0; i < b.N; i++ {
			// Acquire a string slice from the pool.
			strSlice := stringPool.Get()
			// Modify the string slice by appending "Hello" and "World".
			strSlice = append(strSlice, "Hello", "World")
			// Return the modified string slice to the pool.
			stringPool.Put(strSlice)
		}
	})

	// Reset the timer before starting the next sub-benchmark.
	b.ResetTimer()

	// Run a sub-benchmark without pool usage.
	b.Run("WithoutPool", func(b *testing.B) {
		// Perform benchmarking for 'b.N' iterations.
		for i := 0; i < b.N; i++ {
			// Create a new string slice without using a pool.
			strSlice := make([]string, 0, capacity)
			// Modify the string slice by appending "Hello" and "World".
			strSlice = append(strSlice, "Hello", "World")
		}
	})
}
