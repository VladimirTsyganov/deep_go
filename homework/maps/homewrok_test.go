package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type TreeLeaf struct {
	key     int
	value   int
	lBranch *TreeLeaf
	rBranch *TreeLeaf
}

type OrderedMap struct {
	head *TreeLeaf
	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

func (m *OrderedMap) Insert(key, value int) {
	if m.head == nil {
		m.head = &TreeLeaf{
			key:   key,
			value: value,
		}
		m.size = 1
		return
	}

	el, built := findOrBuildLeaf(m.head, key)

	if built {
		m.size += 1
	}

	el.value = value
}

func (m *OrderedMap) Erase(key int) {
	if m.size == 0 {
		return
	}

	if m.head.key == key {
		m.head = nil
		m.size = 0
		return
	}

	root, toDelete := m.findRootLeafToDelete(m.head, key)
	if root == nil {
		return
	}

	var newBranch *TreeLeaf

	if toDelete.lBranch == nil && toDelete.rBranch == nil {
		newBranch = nil
	} else if toDelete.lBranch != nil && toDelete.rBranch != nil {
		base := toDelete
		minEl := toDelete.rBranch

		for {
			if minEl.lBranch == nil {
				break
			}

			base = minEl
			minEl = base.lBranch
		}

		base.lBranch = minEl.rBranch
		newBranch = minEl
	} else if toDelete.lBranch != nil {
		newBranch = toDelete.lBranch
	} else {
		newBranch = toDelete.rBranch
	}

	if root.key > key {
		root.lBranch = newBranch
	} else {
		root.rBranch = newBranch
	}

	m.size -= 1
}

func (m *OrderedMap) Contains(key int) bool {
	if m.size == 0 {
		return false
	}

	leaf := m.find(m.head, key)

	return leaf != nil
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	if m.size == 0 {
		return
	}

	m.applyAction(m.head, action)
}

func findOrBuildLeaf(from *TreeLeaf, key int) (*TreeLeaf, bool) {
	var leaf *TreeLeaf

	if from.key == key {
		return from, false
	}

	if from.key > key {
		leaf = from.lBranch
	} else {
		leaf = from.rBranch
	}

	if leaf != nil {
		return findOrBuildLeaf(leaf, key)
	}

	leaf = &TreeLeaf{
		key: key,
	}

	if from.key > key {
		from.lBranch = leaf
	} else {
		from.rBranch = leaf
	}

	return leaf, true
}

func (m *OrderedMap) find(from *TreeLeaf, key int) *TreeLeaf {
	var branch *TreeLeaf

	if from.key == key {
		return from
	}

	if from.key > key {
		branch = from.lBranch
	} else {
		branch = from.rBranch
	}

	if branch == nil {
		return nil
	}

	return m.find(branch, key)
}

func (m *OrderedMap) applyAction(from *TreeLeaf, action func(int, int)) {
	if from.lBranch != nil {
		m.applyAction(from.lBranch, action)
	}

	action(from.key, from.value)

	if from.rBranch != nil {
		m.applyAction(from.rBranch, action)
	}
}

func (m *OrderedMap) findRootLeafToDelete(from *TreeLeaf, key int) (*TreeLeaf, *TreeLeaf) {
	var branch *TreeLeaf

	if from.key < key {
		branch = from.rBranch
	} else {
		branch = from.lBranch
	}

	if branch == nil {
		return nil, nil
	} else if branch.key == key {
		return from, branch
	}

	return m.findRootLeafToDelete(branch, key)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
