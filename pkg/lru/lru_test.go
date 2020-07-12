package lru

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLRU(t *testing.T) {
	l, err := NewLRU(128, nil)
	assert.NoError(t, err)

	for i := 0; i < 256; i++ {
		l.Add(i, i)
	}
	assert.Equal(t, 128, l.Len())
	assert.Equal(t, 128, l.size)

	for i, k := range l.Keys() {
		v, ok := l.Get(k)
		assert.True(t, ok)
		assert.Equal(t, v, k)
		assert.Equal(t, i+128, v)
	}
	for i := 0; i < 128; i++ {
		_, ok := l.Get(i)
		assert.False(t, ok)
	}
	for i := 128; i < 256; i++ {
		_, ok := l.Get(i)
		assert.True(t, ok)
	}
	for i := 128; i < 192; i++ {
		ok := l.Remove(i)
		assert.True(t, ok)
		ok = l.Remove(i)
		assert.False(t, ok)
		_, ok = l.Get(i)
		assert.False(t, ok)
	}

	l.Get(192) // expect 192 to be last key in l.Keys()

	l.Purge()
	assert.Equal(t, 0, l.Len())

	_, ok := l.Get(200)
	assert.False(t, ok)
}

func TestLRU_GetOldest_RemoveOldest(t *testing.T) {
	l, err := NewLRU(128, nil)
	assert.NoError(t, err)

	for i := 0; i < 256; i++ {
		l.Add(i, i)
	}
	k, _, ok := l.GetOldest()
	assert.True(t, ok)
	assert.Equal(t, 128, k.(int))

	k, _, ok = l.RemoveOldest()
	assert.True(t, ok)
	assert.Equal(t, 128, k.(int))

	k, _, ok = l.RemoveOldest()
	assert.True(t, ok)
	assert.Equal(t, 129, k.(int))
}

// Test that Add returns true/false if an eviction occurred.
func TestLRU_Add(t *testing.T) {
	l, err := NewLRU(1, nil)
	assert.NoError(t, err)
	assert.False(t, l.Add(1, 1))
	assert.True(t, l.Add(2, 2))
}

// Test that Contains doesn't update recent-ness.
func TestLRU_Contains(t *testing.T) {
	l, err := NewLRU(2, nil)
	assert.NoError(t, err)

	l.Add(1, 1)
	l.Add(2, 2)
	require.True(t, l.Contains(1))

	l.Add(3, 3)
	require.False(t, l.Contains(1))
}

// Test that Peek doesn't update recent-ness.
func TestLRU_Peek(t *testing.T) {
	l, err := NewLRU(2, nil)
	assert.NoError(t, err)

	l.Add(1, 1)
	l.Add(2, 2)
	v, ok := l.Peek(1)
	require.True(t, ok)
	require.Equal(t, 1, v)

	l.Add(3, 3)
	require.False(t, l.Contains(1))
}
