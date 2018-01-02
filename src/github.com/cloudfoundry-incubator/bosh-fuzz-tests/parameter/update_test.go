package parameter_test

import (
	"math/rand"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

var _ = Describe("Update", func() {
	var (
		update Parameter
	)

	BeforeEach(func() {
		update = NewUpdate([]int{1, 5}, []int{1, 3}, []string{"not_specified", "true", "false"})
	})

	Describe("Apply", func() {
		It("returns true in some instances", func() {
			rand.Seed(4)
			input := bftinput.Input{}
			result := update.Apply(input, bftinput.Input{})
			truthy := true

			Expect(result).To(Equal(bftinput.Input{
				Update: bftinput.UpdateConfig{
					Canaries:    5,
					MaxInFlight: 1,
					Serial:      &truthy,
				},
			}))
		})

		It("returns false in some instances", func() {
			rand.Seed(2)
			input := bftinput.Input{}
			result := update.Apply(input, bftinput.Input{})

			falsy := false

			Expect(result).To(Equal(bftinput.Input{
				Update: bftinput.UpdateConfig{
					Canaries:    1,
					MaxInFlight: 1,
					Serial:      &falsy,
				},
			}))
		})
	})

	It("returns nil in some instances", func() {
		rand.Seed(64)
		input := bftinput.Input{}
		result := update.Apply(input, bftinput.Input{})

		Expect(result).To(Equal(bftinput.Input{
			Update: bftinput.UpdateConfig{
				Canaries:    5,
				MaxInFlight: 1,
				Serial:      nil,
			},
		}))
	})
})
