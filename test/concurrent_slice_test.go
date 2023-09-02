package test

import (
	"douyin/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"sync"
	"testing"
)

func TestConcurrentSlice(t *testing.T) {
	t.Run("Test ConcurrentSlice functionality", testConcurrentSliceFunctionality)
	t.Run("Test ConcurrentSlice concurrency", testConcurrentSliceConcurrency)
}

func testConcurrentSliceFunctionality(t *testing.T) {
	length := 30
	capacity := 50
	cs := util.NewConcurrentSlice[int](capacity, length)
	assert.Panics(t, func() {
		util.NewConcurrentSlice[int](length, capacity, 114514, 1919810)
	})

	// Test Cap and Len
	// 构造函数第一个参数为长度，第二个参数为容量
	// 期望是交换后的结果
	require.Equal(t, length, cs.Len())
	require.Equal(t, capacity, cs.Cap())

	// Test Append
	items := []int{1, 114514, 3}
	length = length + len(items)
	cs.Append(items...)
	require.Equal(t, length, cs.Len())

	// Test Get
	require.Equal(t, items[1], cs.Get(length-2))

	// Test Set
	index := 1
	item := 1919810
	cs.Set(index, item)
	require.Equal(t, item, cs.Get(index))

	// Test Copy
	cs2 := util.NewConcurrentSlice[int](16)
	for i := 0; i < cs2.Len(); i++ {
		cs2.Set(i, rand.Intn(256))
	}
	cs2.Copy(cs)
	require.Equal(t, cs.Len(), cs2.Len())
	for i := 0; i < cs.Len(); i++ {
		require.Equal(t, cs.Get(i), cs2.Get(i))
	}
	// 对 cs2 进行修改，不应该影响到 cs
	for i := 0; i < cs2.Len(); i++ {
		cs2.Set(i, rand.Intn(256)+256)
	}
	for i := 0; i < cs.Len(); i++ {
		require.NotEqual(t, cs.Get(i), cs2.Get(i))
	}

	// Test Slice
	beg := 3
	end := 10
	cs3 := cs.Slice(beg, end)
	require.Equal(t, end-beg, cs3.Len())
	for i := 0; i < cs3.Len(); i++ {
		require.Equal(t, cs.Get(beg+i), cs3.Get(i))
	}
}

func testConcurrentSliceConcurrency(t *testing.T) {
	cs := util.NewConcurrentSlice[int]()
	numGoroutines := 100
	numItemsPerGoroutine := 100

	var wg sync.WaitGroup

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			for j := 0; j < numItemsPerGoroutine; j++ {
				cs.Append(i*numItemsPerGoroutine + j)
			}
		}(i)
	}
	wg.Wait()

	// 验证长度是否正确
	expectedLen := numGoroutines * numItemsPerGoroutine
	require.Equal(t, expectedLen, cs.Len())

	// 验证元素无重复
	seen := make(map[int]struct{})
	for i := 0; i < expectedLen; i++ {
		item := cs.Get(i)
		_, ok := seen[item]
		require.False(t, ok, "Duplicate item found:", item)
		seen[item] = struct{}{}
	}
	for i := 0; i < expectedLen; i++ {
		item := cs.Get(i)
		_, ok := seen[item]
		require.True(t, ok, "Item not found:", item)
	}

	// 验证并发 Get 和 Set
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			for j := 0; j < numItemsPerGoroutine; j++ {
				index := i*numItemsPerGoroutine + j
				item := rand.Int()
				cs.Set(index, item)
				assert.Equal(t, item, cs.Get(index))
			}
		}(i)
	}
	wg.Wait()
}
