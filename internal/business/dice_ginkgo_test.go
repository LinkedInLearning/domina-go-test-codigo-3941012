package business

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var savageDice *SavageDice
var baseDice *BaseDice

var _ = BeforeSuite(func() {
	// Inicializamos los dados antes de toda la suite de tests
	savageDice = &SavageDice{
		BaseDice: BaseDice{
			Sides: 6,
		},
	}

	baseDice = &BaseDice{
		Sides: 6,
	}
})

var _ = AfterSuite(func() {
	// Limpiamos los dados después de toda la suite de tests
	savageDice = nil
	baseDice = nil
})

var _ = Describe("Dices", func() {
	BeforeEach(func() {
		// Inicializamos los dados antes de cada test
		savageDice.Sides = 8
		savageDice.maxExplosions = 50

		baseDice.Sides = 6
	})

	Describe("Roll", func() {
		Context("with savage dice", func() {
			When("the dice has 1 side", func() {
				BeforeEach(func() {
					// inicializamos el dado con 1 lado
					savageDice.Sides = 1
					savageDice.maxExplosions = 50

					DeferCleanup(func() {
						// limpiamos el estado del dado al finalizar el test
						savageDice.Sides = 6
						savageDice.maxExplosions = 0
					})
				})

				It("always explodes to the maximum number of explosions", func() {
					roll := savageDice.Roll()
					Expect(roll).To(Equal(50))
					Expect(savageDice.Explosions).To(Equal(50))
				})
			})

			When("the dice has more than 1 side", func() {
				It("can explode", func() {
					roll := savageDice.Roll()
					Expect(roll).To(BeNumerically(">=", 1))
					Expect(savageDice.Explosions).To(BeNumerically(">=", 0))
					Expect(roll).To(BeNumerically("<=", savageDice.Explosions*savageDice.Sides+savageDice.Sides))
				})
			})
		})

		Context("with base dice", func() {
			It("cannot explode", func() {
				roll := baseDice.Roll()
				Expect(roll).To(BeNumerically(">=", 1))
				Expect(roll).To(BeNumerically("<=", baseDice.Sides))
			})
		})
	})
})

var _ = DescribeTable("All the dices",
	func(sides int, savage bool) {
		var d Dice
		d = &BaseDice{
			Sides: sides,
		}

		if savage {
			d = &SavageDice{
				BaseDice: BaseDice{
					Sides: sides,
				},
				maxExplosions: 50,
			}
		}

		roll := d.Roll()

		if savage {
			sd := d.(*SavageDice)
			Expect(roll).To(BeNumerically(">=", 1))
			Expect(roll).To(BeNumerically("<=", sd.Result()))
			Expect(roll).To(BeNumerically("<=", sd.Explosions*sd.Sides+sd.Sides))
			Expect(sd.Explosions).To(BeNumerically(">=", 0))
			Expect(sd.Explosions).To(Equal(len(sd.rolls) - 1))
		} else {
			Expect(roll).To(BeNumerically(">=", 1))
			Expect(roll).To(BeNumerically("<=", d.Result()))
		}
	},
	Entry("Dice 4", 4, false),
	Entry("Savage Dice 4", 4, true),
	Entry("Dice 6", 6, false),
	Entry("Savage Dice 6", 6, true),
	Entry("Dice 8", 8, false),
	Entry("Savage Dice 8", 8, true),
	Entry("Dice 10", 10, false),
	Entry("Savage Dice 10", 10, true),
	Entry("Dice 12", 12, false),
	Entry("Savage Dice 12", 12, true),
	Entry("Dice 20", 20, false),
	Entry("Savage Dice 20", 20, true),
	Entry("Dice 100", 100, false),
	Entry("Savage Dice 100", 100, true),
)

func TestDices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dices Suites")
}
