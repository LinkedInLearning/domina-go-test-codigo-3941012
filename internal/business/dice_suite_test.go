package business

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// DiceTestSuite is a suite of tests for the Dice struct.
type DiceTestSuite struct {
	// embed the suite.Suite type
	suite.Suite

	// sides es el número de caras del dado a testear.
	sides int
}

// Todos los métodos que comiencen por `Test` son tests de la suite.
func (s *DiceTestSuite) TestSavageRoll() {
	testSavageRoll(s.T(), s.sides)
}

func (s *DiceTestSuite) TestBaseRoll() {
	baseDice := &BaseDice{
		Sides: s.sides,
	}

	roll := baseDice.Roll()
	s.Require().Greater(roll, 0)
	s.Require().LessOrEqual(roll, baseDice.Sides)
}

// Para que `go test` ejecute la suite, debemos crear una función
// de test habitual, p.e. `TestSavageDiceSuite`, que invoque a
// `suite.Run` pasando como parámetro una instancia de nuestra suite.
func TestSavageDiceSuite(t *testing.T) {
	sides := []int{4, 6, 8, 10, 12, 20, 100}

	for _, side := range sides {
		// creamos una suite para cada lado.
		suite.Run(t, &DiceTestSuite{sides: side})
	}
}
