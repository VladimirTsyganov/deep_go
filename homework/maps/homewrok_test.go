package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type TreeNode struct {
	key     int
	value   int
	lBranch *TreeNode
	rBranch *TreeNode
}

type OrderedMap struct {
	head *TreeNode
	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

func (m *OrderedMap) Insert(key, value int) {
	if m.head == nil {
		m.head = &TreeNode{
			key:   key,
			value: value,
		}
		m.size = 1
		return
	}

	var node *TreeNode

	if !m.Contains(key) {
		m.size += 1

		node = m.insertNode(key)
	} else {
		node = findNode(m.head, key)
	}

	node.value = value
}

func (m *OrderedMap) Erase(key int) {
	if m.empty() || !m.Contains(key) {
		return
	}

	if m.size == 1 {
		m.head = nil
		m.size = 0
		return
	}

	var replacement *TreeNode

	root := findRootNode(m.head, key)
	toDelete := findNode(root, key)

	if toDelete.lBranch == nil || toDelete.rBranch == nil {

		if toDelete.lBranch != nil {
			replacement = toDelete.lBranch
		} else {
			replacement = toDelete.rBranch
		}
	} else {
		replacement = findMinAndDetach(toDelete.rBranch)
		replacement.lBranch = toDelete.lBranch

		if replacement != toDelete.rBranch {
			replacement.rBranch = toDelete.rBranch
		}
	}

	if root.lBranch == toDelete {
		root.lBranch = replacement
	} else {
		root.rBranch = replacement
	}

	m.size -= 1
}

func (m *OrderedMap) Contains(key int) bool {
	if m.empty() {
		return false
	}

	node := findNode(m.head, key)

	return node != nil
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	if m.empty() {
		return
	}

	m.applyAction(m.head, action)
}

func (m *OrderedMap) insertNode(key int) *TreeNode {
	rootNode := findRootNode(m.head, key)
	node := &TreeNode{key: key}

	if key > rootNode.key {
		rootNode.rBranch = node
	} else {
		rootNode.lBranch = node
	}

	return node
}

func (m *OrderedMap) applyAction(from *TreeNode, action func(int, int)) {
	if from.lBranch != nil {
		m.applyAction(from.lBranch, action)
	}

	action(from.key, from.value)

	if from.rBranch != nil {
		m.applyAction(from.rBranch, action)
	}
}

func (m *OrderedMap) empty() bool {
	return m.Size() == 0
}

func findNode(from *TreeNode, key int) *TreeNode {
	if from == nil {
		return nil
	}

	if from.key == key {
		return from
	}

	if from.key > key {
		return findNode(from.lBranch, key)
	}

	return findNode(from.rBranch, key)
}

func findRootNode(from *TreeNode, key int) *TreeNode {
	var branch *TreeNode

	if from.key > key {
		branch = from.lBranch
	} else {
		branch = from.rBranch
	}

	if branch == nil || branch.key == key {
		return from
	}

	return findRootNode(branch, key)
}

func findMinAndDetach(from *TreeNode) *TreeNode {
	if from.lBranch == nil {
		return from
	}

	minParent := from
	minEl := minParent.lBranch

	for {
		if minEl.lBranch == nil {
			break
		}

		minParent = minEl
		minEl = minParent.lBranch
	}

	if minEl.rBranch != nil {
		minParent.lBranch = minEl.rBranch
	}

	return minEl
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

func TestCircularQueueWithCornerCases(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(12, 12)
	data.Insert(11, 11)
	data.Insert(20, 20)
	data.Insert(18, 18)
	data.Insert(19, 19)
	data.Insert(22, 22)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(10))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{10, 11, 12, 18, 19, 20, 22}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(12)

	assert.Equal(t, 6, data.Size())
	assert.False(t, data.Contains(12))
	assert.True(t, data.Contains(18))
	assert.True(t, data.Contains(19))
	assert.True(t, data.Contains(20))
	assert.True(t, data.Contains(22))

	keys = nil
	expectedKeys = []int{10, 11, 18, 19, 20, 22}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(19)
	data.Erase(18)

	keys = nil
	expectedKeys = []int{10, 11, 20, 22}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})
	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
