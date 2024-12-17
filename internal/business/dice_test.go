package business

import (
	"testing"
)

func TestSavageDice(t *testing.T) {
	t.Run("roll", func(t *testing.T) {
		roll := savageDice(10)
		if roll < 1 {
			t.Fatalf("expected roll to be greater than or equal to 1, got %d", roll)
		}
	})

	t.Run("explodes-to-the-max", func(t *testing.T) {
		// Porque el dado es de 1 cara, siempre explotará,
		// llegando al máximo de 50.
		roll := savageDice(1)
		if roll != 50 {
			t.Fatalf("expected roll to be 50, got %d", roll)
		}
	})
}
