package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContains(t *testing.T) {
	slice := []string{"A", "B"}

	require.True(t, Contains(slice, "A"))
	require.False(t, Contains(slice, "C"))
}
