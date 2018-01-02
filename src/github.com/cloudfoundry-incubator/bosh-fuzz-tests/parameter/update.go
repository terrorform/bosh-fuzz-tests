package parameter

import (
	"math/rand"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type update struct {
	canaries    []int
	maxInFlight []int
	// serial      []string
}

func NewUpdate(canaries []int, maxInFlight []int, serial []string) Parameter {
	return &update{
		canaries:    canaries,
		maxInFlight: maxInFlight,
		// serial:      serial,
	}
}

func (u *update) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	input.Update.Canaries = u.canaries[rand.Intn(len(u.canaries))]
	input.Update.MaxInFlight = u.maxInFlight[rand.Intn(len(u.maxInFlight))]

	seed := rand.Intn(3)
	if seed == 2 {
		input.Update.Serial = nil
	} else {
		b := !(0 == seed)
		input.Update.Serial = &b
	}

	return input
}
