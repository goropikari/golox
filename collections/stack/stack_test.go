package stack_test

import (
	"testing"

	"github.com/goropikari/tlps/collections/stack"
	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	st := stack.NewStack()

	items := []int{123, 456, 789}
	revitems := make([]int, 0)
	for i := len(items) - 1; i >= 0; i-- {
		revitems = append(revitems, items[i])
	}

	for _, v := range items {
		st.Push(v)
	}

	assert.Equal(t, len(items), st.Size())

	for k, v := range revitems {
		g, _ := st.Get(k)
		assert.Equal(t, v, g)
	}

	for _, v := range revitems {
		g := st.Pop()
		assert.Equal(t, v, g)
	}
}
