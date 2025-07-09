package main

import (
	"reflect"
	"runtime"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type COWBuffer struct {
	data   []byte
	refs   *int
	parent *COWBuffer
}

func NewCOWBuffer(data []byte) COWBuffer {
	buf := COWBuffer{
		data: data,
		refs: new(int),
	}

	runtime.SetFinalizer(&buf, func(b *COWBuffer) {
		b.Close()
	})

	return buf
}

func (b *COWBuffer) Clone() COWBuffer {
	b.addChild()

	buf := COWBuffer{
		data:   b.data,
		refs:   new(int),
		parent: b,
	}

	runtime.SetFinalizer(&buf, func(b *COWBuffer) {
		b.Close()
	})

	return buf
}

func (b *COWBuffer) Close() {
	if b.isClosed() {
		return
	}

	if b.hasParent() {
		b.parent.removeChild()
		b.parent = nil
	}

	runtime.SetFinalizer(&b, nil)

	b.refs = nil
	b.data = nil
}

func (b *COWBuffer) Update(index int, value byte) bool {
	if index < 0 || index >= len(b.data) {
		return false
	}

	if (b.refs != nil && *b.refs > 0) || b.hasParent() {
		cp := make([]byte, len(b.data), cap(b.data))
		copy(cp, b.data)
		b.data = cp
		b.refs = new(int)
		if b.hasParent() {
			b.parent.removeChild()
			b.parent = nil
		}
	}

	b.data[index] = value

	return true
}

func (b *COWBuffer) String() string {
	return unsafe.String(unsafe.SliceData(b.data), len(b.data))
}

func (b *COWBuffer) addChild() {
	if b.refs == nil {
		b.refs = new(int)
	}

	*b.refs += 1
}

func (b *COWBuffer) removeChild() {
	if b.isClosed() || *b.refs <= 0 {
		return
	}

	*b.refs -= 1
}

func (b *COWBuffer) hasParent() bool {
	if b.parent == nil || b.parent.isClosed() {
		return false
	}

	return unsafe.SliceData(b.parent.data) == unsafe.SliceData(b.data)
}

func (b *COWBuffer) isClosed() bool {
	return b.refs == nil
}

func TestCOWBuffer(t *testing.T) {
	data := []byte{'a', 'b', 'c', 'd'}
	buffer := NewCOWBuffer(data)
	defer buffer.Close()

	copy1 := buffer.Clone()
	copy2 := buffer.Clone()

	assert.Equal(t, unsafe.SliceData(data), unsafe.SliceData(buffer.data))
	assert.Equal(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	assert.True(t, (*byte)(unsafe.SliceData(data)) == unsafe.StringData(buffer.String()))
	assert.True(t, (*byte)(unsafe.StringData(buffer.String())) == unsafe.StringData(copy1.String()))
	assert.True(t, (*byte)(unsafe.StringData(copy1.String())) == unsafe.StringData(copy2.String()))

	assert.True(t, buffer.Update(0, 'g'))
	assert.False(t, buffer.Update(-1, 'g'))
	assert.False(t, buffer.Update(4, 'g'))

	assert.True(t, reflect.DeepEqual([]byte{'g', 'b', 'c', 'd'}, buffer.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy1.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy2.data))

	assert.NotEqual(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	copy1.Close()

	previous := copy2.data
	copy2.Update(0, 'f')
	current := copy2.data

	// 1 reference - don't need to copy buffer during update
	assert.Equal(t, unsafe.SliceData(previous), unsafe.SliceData(current))

	copy3 := copy2.Clone()
	defer copy3.Close()
	assert.Equal(t, unsafe.SliceData(copy2.data), unsafe.SliceData(copy3.data))

	copy3.Update(0, 'x')

	assert.True(t, reflect.DeepEqual([]byte{'f', 'b', 'c', 'd'}, copy2.data))
	assert.True(t, reflect.DeepEqual([]byte{'x', 'b', 'c', 'd'}, copy3.data))
	assert.NotEqual(t, unsafe.SliceData(copy2.data), unsafe.SliceData(copy3.data))

	copy2.Close()
}
