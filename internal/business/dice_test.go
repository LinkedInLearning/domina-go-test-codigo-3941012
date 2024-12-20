package business

import (
	"testing"
)

func TestSavageDice(t *testing.T) {
	t.Run("savage-dice-roll", func(t *testing.T) {
		savageDice := &SavageDice{
			BaseDice: BaseDice{
				Sides: 10,
			},
		}

		t.Run("roll", func(t *testing.T) {
			roll := savageDice.Roll()
			lowerBound := savageDice.Explosions * savageDice.Sides
			if roll <= lowerBound {
				t.Fatalf("expected roll to be greater than %d, got %d", lowerBound, roll)
			}
		})

		t.Run("explodes-to-the-max", func(t *testing.T) {
			oneSidedDice := &SavageDice{
				BaseDice: BaseDice{
					Sides: 1,
				},
				maxExplosions: 50,
			}
			// Porque el dado es de 1 cara, siempre explotará,
			// llegando al máximo de 50.
			roll := oneSidedDice.Roll()
			if roll != 50 {
				t.Fatalf("expected roll to be 50, got %d", roll)
			}
		})
	})

	t.Run("base-dice-roll", func(t *testing.T) {
		baseDice := &BaseDice{
			Sides: 10,
		}

		t.Run("roll", func(t *testing.T) {
			roll := baseDice.Roll()
			if roll < 1 {
				t.Fatalf("expected roll to be greater than or equal to 1, got %d", roll)
			}
			if roll > baseDice.Sides {
				t.Fatalf("expected roll to be less than or equal to %d, got %d", baseDice.Sides, roll)
			}
		})

		t.Run("roll-with-one-side", func(t *testing.T) {
			oneSidedDice := &BaseDice{
				Sides: 1,
			}

			roll := oneSidedDice.Roll()
			if roll != 1 {
				t.Fatalf("expected roll to be 1, got %d", roll)
			}
		})
	})
}
