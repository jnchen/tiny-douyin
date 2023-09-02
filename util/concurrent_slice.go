package util

import (
	"log"
	"sync"
)

type ConcurrentSlice[T any] struct {
	mu   sync.RWMutex
	data []T
}

// NewConcurrentSlice 构造一个新的并发切片，参数分别为长度和容量，可省略
func NewConcurrentSlice[T any](size ...int) *ConcurrentSlice[T] {
	l := 0
	c := 0
	if len(size) == 0 {
	} else if len(size) == 1 {
		l = size[0]
		c = size[0]
	} else if len(size) == 2 {
		l = size[0]
		c = size[1]
	} else {
		log.Panicln("NewConcurrentSlice: Too many arguments!")
	}

	if l > c {
		l, c = c, l
		log.Println("NewConcurrentSlice: Len > Cap, swapping!")
	}

	return &ConcurrentSlice[T]{
		data: make([]T, l, c),
	}
}

func (cs *ConcurrentSlice[T]) Cap() int {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	return cap(cs.data)
}

func (cs *ConcurrentSlice[T]) Len() int {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	return len(cs.data)
}

func (cs *ConcurrentSlice[T]) Slice(beg, end int) *ConcurrentSlice[T] {
	newSlice := NewConcurrentSlice[T]()

	cs.mu.RLock()
	defer cs.mu.RUnlock()

	newSlice.data = cs.data[beg:end]
	return newSlice
}

func (cs *ConcurrentSlice[T]) RawSlice() []T {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	return cs.data
}

func (cs *ConcurrentSlice[T]) Copy(src *ConcurrentSlice[T]) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.data = make([]T, len(src.data))
	copy(cs.data, src.data)
}

func (cs *ConcurrentSlice[T]) Append(item ...T) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.data = append(cs.data, item...)
}

func (cs *ConcurrentSlice[T]) Get(index int) T {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	return cs.data[index]
}

func (cs *ConcurrentSlice[T]) Set(index int, item T) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.data[index] = item
}
