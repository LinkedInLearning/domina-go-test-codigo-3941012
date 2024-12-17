package business

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRollDice(t *testing.T) {
	roll := rollDice(10)
	require.GreaterOrEqual(t, roll, 1)
	require.LessOrEqual(t, roll, 10)
}
