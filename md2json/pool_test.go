package md2json

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func Test_test(t *testing.T) {
	stringPool := NewStringSlicePool(5)

	// 檢查 Inited 值是否已設置
	if stringPool.Inited.Load() {
		fmt.Println("StringSlicePool is initialized.")
	} else {
		fmt.Println("StringSlicePool is not initialized.")
	}
}

func TestStringSlicePool(t *testing.T) {
	t.Run("Test Get and Put", func(t *testing.T) {
		capacity := 5
		stringPool := NewStringSlicePool(capacity)

		strSlice, err := stringPool.Get()
		assert.NoError(t, err, "Expected no error from Get")

		strSlice = append(strSlice, "Hello", "World")
		stringPool.Put(strSlice)

		strSlice2, _ := stringPool.Get()
		assert.Empty(t, strSlice2, "Expected an empty slice from the pool")
	})

	t.Run("Test Get on uninitialized pool", func(t *testing.T) {
		uninitializedPool := &StringSlicePool{} // Create an uninitialized pool

		_, err := uninitializedPool.Get()
		assert.Error(t, err, "Expected an error for getting from uninitialized pool")
	})
}

// go test -run=TestStringSlicePool_Concurrent -race
func Test_Race_StringSlicePool(t *testing.T) {
	// Enable data race detection
	t.Parallel()

	capacity := 5
	stringPool := NewStringSlicePool(capacity)

	const goroutines = 100
	const iterations = 100

	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				strSlice, err := stringPool.Get()
				if err == nil {
					strSlice = append(strSlice, "Hello", "World")
					stringPool.Put(strSlice)
				}
			}
		}()
	}

	wg.Wait()
}
