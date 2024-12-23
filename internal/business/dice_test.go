package business

import (
	"testing"
)

func testSavageRoll(t *testing.T, sides int) {
	t.Helper()

	if sides <= 2 {
		t.Skip("skipping test for dice with less than 2 sides")
	}

	savageDice := &SavageDice{
		BaseDice: BaseDice{
			Sides: sides,
		},
		maxExplosions: 50,
	}

	roll := savageDice.Roll()
	lowerBound := savageDice.Explosions * savageDice.Sides
	upperBound := (savageDice.Explosions + 1) * savageDice.Sides
	if roll < lowerBound {
		t.Fatalf("expected roll to be greater than %d, got %d", lowerBound, roll)
	}
	if roll > upperBound {
		t.Fatalf("expected roll to be less than %d, got %d", upperBound, roll)
	}
	t.Logf("sides: %d, roll: %d, rolls: %v, explosions: %d, lowerBound: %d, upperBound: %d", sides, roll, savageDice.rolls, savageDice.Explosions, lowerBound, upperBound)
}

func TestSavageDice(t *testing.T) {
	t.Run("savage-dice-roll", func(t *testing.T) {
		t.Run("sides", func(t *testing.T) {
			t.Run("2", func(t *testing.T) {
				testSavageRoll(t, 2)
			})
			t.Run("4", func(t *testing.T) {
				testSavageRoll(t, 4)
			})
			t.Run("6", func(t *testing.T) {
				testSavageRoll(t, 6)
			})
			t.Run("8", func(t *testing.T) {
				testSavageRoll(t, 8)
			})
			t.Run("10", func(t *testing.T) {
				testSavageRoll(t, 10)
			})
			t.Run("12", func(t *testing.T) {
				testSavageRoll(t, 12)
			})
			t.Run("20", func(t *testing.T) {
				testSavageRoll(t, 20)
			})
			t.Run("100", func(t *testing.T) {
				testSavageRoll(t, 100)
			})
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
		testCases := []struct {
			name  string
			sides int
		}{
			{name: "100", sides: 100},
			{name: "20", sides: 20},
			{name: "12", sides: 12},
			{name: "10", sides: 10},
			{name: "8", sides: 8},
			{name: "6", sides: 6},
			{name: "4", sides: 4},
			{name: "2", sides: 2},
			{name: "1", sides: 1},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				testCase := testCase
				t.Parallel()

				baseDice := &BaseDice{
					Sides: testCase.sides,
				}

				roll := baseDice.Roll()
				if roll <= 0 {
					t.Fatalf("expected roll to be greater than 0, got %d", roll)
				}
				if roll > baseDice.Sides {
					t.Fatalf("expected roll to be less than or equal to %d, got %d", baseDice.Sides, roll)
				}
			})
		}
	})
}
