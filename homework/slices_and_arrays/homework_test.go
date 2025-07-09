package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type CircularQueue[T int | int16 | int32 | int64] struct {
	values []T
	head   *int
	tail   int
}

func NewCircularQueue[T int | int16 | int32 | int64](size int) CircularQueue[T] {
	return CircularQueue[T]{
		values: make([]T, size),
	}
}

func (q *CircularQueue[T]) Push(value T) bool {
	if q.Full() {
		return false
	}

	q.values[q.tail] = value

	if q.head == nil {
		q.head = new(int)
		*q.head = q.tail
	}

	q.tail = q.nextIdx(q.tail)
	return true
}

func (q *CircularQueue[T]) Pop() bool {
	if q.Empty() {
		return false
	}

	nextHead := q.nextIdx(*q.head)

	if nextHead == q.tail {
		q.head = nil
	} else {
		*q.head = nextHead
	}

	return true
}

func (q *CircularQueue[T]) Front() T {
	if q.Empty() {
		return -1
	}

	return q.values[*q.head]
}

func (q *CircularQueue[T]) Back() T {
	if q.Empty() {
		return -1
	}

	backIdx := q.prevIdx(q.tail)

	return q.values[backIdx]
}

func (q *CircularQueue[T]) Empty() bool {
	return q.head == nil
}

func (q *CircularQueue[T]) Full() bool {
	return !q.Empty() && *q.head == q.tail
}

func (q *CircularQueue[T]) nextIdx(idx int) int {
	idx += 1

	if idx > len(q.values)-1 {
		return 0
	}

	return idx
}

func (q *CircularQueue[T]) prevIdx(idx int) int {
	idx -= 1

	if idx < 0 {
		return len(q.values) - 1
	}

	return idx
}

func TestCircularIntQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.True(t, queue.Push(1))
	assert.True(t, reflect.DeepEqual([]int{4, 1, 3}, queue.values))
	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 1, queue.Back())
}

func TestCircularInt32Queue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int32](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, int32(-1), queue.Front())
	assert.Equal(t, int32(-1), queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int32{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, int32(1), queue.Front())
	assert.Equal(t, int32(3), queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int32{4, 2, 3}, queue.values))

	assert.Equal(t, int32(2), queue.Front())
	assert.Equal(t, int32(4), queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.True(t, queue.Push(1))
	assert.True(t, reflect.DeepEqual([]int32{4, 1, 3}, queue.values))
	assert.Equal(t, int32(1), queue.Front())
	assert.Equal(t, int32(1), queue.Back())
}
