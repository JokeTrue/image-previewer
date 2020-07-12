package cache

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestLRU(t *testing.T) {
	l, err := NewCache(128)
	assert.NoError(t, err)

	for i := 0; i < 256; i++ {
		l.Add(i, i)
	}
	require.Equal(t, 128, l.Len())

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
		l.Remove(i)
		_, ok := l.Get(i)
		assert.False(t, ok)
	}

	l.Get(192) // expect 192 to be last key in l.Keys()

	l.Purge()
	assert.Equal(t, 0, l.Len())
	_, ok := l.Get(200)
	assert.False(t, ok)
}

// test that Add returns true/false if an eviction occurred.
func TestLRUAdd(t *testing.T) {
	l, err := NewCache(1)
	assert.NoError(t, err)

	require.False(t, l.Add(1, 1))
	assert.False(t, l.Add(1, 1))
	assert.True(t, l.Add(2, 2))
}

// test that Contains doesn't update recent-ness.
func TestLRUContains(t *testing.T) {
	l, err := NewCache(2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	l.Add(1, 1)
	l.Add(2, 2)
	assert.True(t, l.Contains(1))

	l.Add(3, 3)
	assert.False(t, l.Contains(1))
}

// test that ContainsOrAdd doesn't update recent-ness.
func TestLRUContainsOrAdd(t *testing.T) {
	l, err := NewCache(2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	l.Add(1, 1)
	l.Add(2, 2)
	contains, evict := l.ContainsOrAdd(1, 1)
	assert.True(t, contains)
	assert.False(t, evict)

	l.Add(3, 3)
	contains, evict = l.ContainsOrAdd(1, 1)
	assert.False(t, contains)
	assert.True(t, evict)
	assert.True(t, l.Contains(1))
}

// test that PeekOrAdd doesn't update recent-ness.
func TestLRUPeekOrAdd(t *testing.T) {
	l, err := NewCache(2)
	assert.NoError(t, err)

	l.Add(1, 1)
	l.Add(2, 2)
	previous, contains, evict := l.PeekOrAdd(1, 1)
	assert.True(t, contains)
	assert.False(t, evict)
	assert.Equal(t, 1, previous)

	l.Add(3, 3)
	contains, evict = l.ContainsOrAdd(1, 1)
	assert.False(t, contains)
	assert.True(t, evict)
	assert.True(t, l.Contains(1))
}

// test that Peek doesn't update recent-ness.
func TestLRUPeek(t *testing.T) {
	l, err := NewCache(2)
	assert.NoError(t, err)

	l.Add(1, 1)
	l.Add(2, 2)

	v, ok := l.Peek(1)
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	l.Add(3, 3)
	assert.False(t, l.Contains(1))
}
