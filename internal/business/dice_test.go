package business

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSavageDice(t *testing.T) {
	t.Run("roll", func(t *testing.T) {
		roll := savageDice(10)
		require.GreaterOrEqual(t, roll, 1)
	})

	t.Run("explodes-to-the-max", func(t *testing.T) {
		// Porque el dado es de 1 cara, siempre explotará,
		// llegando al máximo de 50.
		roll := savageDice(1)
		require.Equal(t, roll, 50)
	})
}
