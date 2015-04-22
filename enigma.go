package enigma

import (
	"bytes"

	"enigma/bytesutil"
	"enigma/stringutil"
)

// Types of Components
const (
	Plugboard = iota
	Rotor
	Reflector
)

const (
	RotorI = iota
	RotorII
	RotorIII
)

const (
	ReflectorA = iota
	ReflectorB
	ReflectorC
)

var rotorWiring = []string{
	"EKMFLGDQVZNTOWYHXUSPAIBRCJ",
	"AJDKSIRUXBLHWTMCQGZNPYFVOE",
	"BDFHJLCPRTXVZNYEIWGAKMUSQO",
}

var reflectorWiring = []string {
	"EJMZALYXVBWFCRQUONTSPIKHGD",
	"YRUHQSLDPXNGOKMIEBFZCWVJAT",
	"FVPJIAOYEDRZXWGCTKUQSBNMHL",
}

var notch = []int {
	'Q',
	'E',
	'V',
}

const (
	NumAlphabets = 26
)

// Enigma machine is made of several components with similar functionalities,
// namely the plugboard, rotors, and reflector. Each has an input and output
// "sockets" and they are all chained to each other. Therefore, the components
// are implemented as double-linked lists that have input and output character
// maps. The maps are implemented in simple arrays of bytes instead of Go's map
// for simplicity when used in two-directional way e.g. key[value] and
// value[key].
type Component struct {
	right  [NumAlphabets]byte
	left   [NumAlphabets]byte
	offset int // offset is used by Rotors and ignored by other components
	notch  int
	next   *Component
	prev   *Component
	type_  int
}

func NewComponent(type_ int) *Component {
	c := &Component{}
	c.type_ = type_
	for i := byte(0); i < NumAlphabets; i++ {
		c.left[i] = i
	}
	return c
}

func NewPlugboard() *Component {
	c := NewComponent(Plugboard)
	c.right = c.left
	return c
}

func NewRotor(rotorType int) *Component {
	c := NewComponent(Rotor)
	c.notch = notch[rotorType] - 'A'
	if rotorType >= 0 && rotorType < len(rotorWiring) {
		for i := 0; i < NumAlphabets; i++ {
			c.right[i] = rotorWiring[rotorType][i] - 'A'
		}
	}
	return c
}

func NewReflector(reflectorType int) *Component {
	c := NewComponent(Reflector)
	if reflectorType >= 0 && reflectorType < len(reflectorWiring) {
		for i := 0; i < NumAlphabets; i++ {
			c.right[i] = reflectorWiring[reflectorType][i] - 'A'
		}
	}
	return c
}

// Convenient function to create Enigma in standard configurations
// e.g. a plugboard, three rotors, and a reflector
func NewStandardEnigma(rotorType1, rotorType2, rotorType3, reflectorType int) *Component {
	pb := NewPlugboard()
	r1 := NewRotor(rotorType1)
	r2 := NewRotor(rotorType2)
	r3 := NewRotor(rotorType3)
	rfl := NewReflector(reflectorType)
	Connect(pb, r1, r2, r3, rfl)
	return pb
}

// Connect a set of Enigma components together
func Connect(comps ...*Component) *Component {
	if len(comps) <= 0 {
		return nil
	}

	for i := 0; i < len(comps)-1; i++ {
		comps[i].next = comps[i+1]
		comps[i+1].prev = comps[i]
	}
	return comps[0]
}

// Encrypt a message using a chain of Enigma components. Because of the way
// Enigma works, the process of decrypting is the same as encrypting. So
// use this for decrypting as well!
func (comp *Component) Encrypt(msg []byte) []byte {
	b := stringutil.Sanitize(msg)
	emsg := make([]byte, len(b))
	for i, _ := range b {
		comp.Step(1)
		emsg[i] = comp.encryptChar(b[i])
	}
	return emsg
}

// Set initial settings for the Enigma component
func (comp *Component) SetCharacterMap(right string) {
	for i := 0; i < NumAlphabets; i++ {
		comp.left[i] = byte(i)
		comp.right[i] = right[i] - 'A'
	}
}

// Encrypt a single character
func (comp *Component) encryptChar(c byte) byte {
	r := comp
	i := byte(0)
	c -= 'A'

	for ; r != nil; r = r.next {
		switch r.type_ {
		case Plugboard:
			i = r.lIndex(c)
			c = r.right[i]
		case Rotor:
			c = r.right[i]
			i = r.lIndex(c)
		case Reflector:
			c = r.right[i]
			i = r.lIndex(c)
			goto out
		}
	}

out:
	for r = r.prev; r != nil; r = r.prev {
		c = r.left[i]
		i = r.rIndex(c)
	}

	return 'A' + c
}

func (comp *Component) lIndex(c byte) byte {
	return byte(bytes.IndexByte(comp.left[:], c))
}

func (comp *Component) rIndex(c byte) byte {
	return byte(bytes.IndexByte(comp.right[:], c))
}

// Step all rotors that are forward-linked to this component.
func (comp *Component) Step(steps int) {
	if steps <= 0 {
		return
	}

	// Number of times the notch is encountered
	revs := comp.countNotchRevs(steps)

	// Step the rotor
	comp.step(steps)

	// Rotate the next component on the condition that the current rotor has
	// done a full revolution, or if the current component is not a rotor.
	if comp.next != nil && (revs > 0 || comp.type_ != Rotor) {
		comp.next.Step(revs)
	}
}

// Count the number of times the rotor has passed the notch. Has to take care of
// a case when the number of steps is a multiple of NumAlphabets (e.g. for
// initial settings).
func (comp *Component) countNotchRevs(steps int) int {
	// Returns steps if not a Rotor.
	if comp.type_ != Rotor {
		return steps
	}

	// Return 0 if it is determined that the notch won't be reached in the steps
	if comp.offset > comp.notch && steps < comp.notch + (NumAlphabets - comp.notch) {
		return 0
	}

	return reallyCountNotchRevs(comp.offset, comp.notch, NumAlphabets, steps)
}

func reallyCountNotchRevs(current, notch, max, steps int) int {
	revs := 0
	steps -= notch - current
	if steps > 0 {
		revs++
	}
	revs += steps / max
	return revs
}

// Step only current rotor component by n position.
func (comp *Component) step(n int) {
	if comp.type_ != Rotor {
		return
	}

	comp.offset = (comp.offset + n) % NumAlphabets
	bytesutil.Shift(comp.right[:], 1)
	bytesutil.Shift(comp.left[:], 1)
}
