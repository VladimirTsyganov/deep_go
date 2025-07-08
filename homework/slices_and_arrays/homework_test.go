package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type CircularQueue[T int | int16 | int32 | int64] struct {
	values []T
	head   unsafe.Pointer
	tail   unsafe.Pointer
}

func NewCircularQueue[T int | int16 | int32 | int64](size int) CircularQueue[T] {
	s := make([]T, size)
	return CircularQueue[T]{
		values: s,
		tail:   unsafe.Pointer(&s[0]),
	}
}

func (q *CircularQueue[T]) Push(value T) bool {
	if q.head != nil && q.head == q.tail {
		return false
	}

	*(*T)(q.tail) = value

	if q.head == nil {
		q.head = q.tail
	}

	q.tail = q.nextPtr(q.tail)
	return true
}

func (q *CircularQueue[T]) Pop() bool {
	if q.head == nil {
		return false
	}

	nextHead := q.nextPtr(q.head)

	if nextHead == q.tail {
		q.head = nil
	} else {
		q.head = nextHead
	}

	return true
}

func (q *CircularQueue[T]) Front() T {
	if q.Empty() {
		return -1
	}

	return *(*T)(q.head)
}

func (q *CircularQueue[T]) Back() T {
	if q.Empty() {
		return -1
	}

	return *(*T)(q.prevPtr(q.tail))
}

func (q *CircularQueue[T]) Empty() bool {
	return q.head == nil
}

func (q *CircularQueue[T]) Full() bool {
	return !q.Empty() && q.head == q.tail
}

func (q *CircularQueue[T]) nextPtr(ptr unsafe.Pointer) unsafe.Pointer {
	if &q.values[len(q.values)-1] == (*T)(ptr) {
		return unsafe.Pointer(&q.values[0])
	}

	return unsafe.Add(ptr, unsafe.Sizeof(q.values[0]))
}

func (q *CircularQueue[T]) prevPtr(ptr unsafe.Pointer) unsafe.Pointer {
	if &q.values[0] == (*T)(ptr) {
		return unsafe.Pointer(&q.values[len(q.values)-1])
	}

	return unsafe.Add(ptr, -unsafe.Sizeof(q.values[0]))
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
