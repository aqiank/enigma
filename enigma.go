package enigma

import (
	"strings"

	"enigma/stringutil"
)

// These information were obtained from
// https://en.wikipedia.org/wiki/Enigma_rotor_details

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

var rotorCharMap = []string{
	"EKMFLGDQVZNTOWYHXUSPAIBRCJ",
	"AJDKSIRUXBLHWTMCQGZNPYFVOE",
	"BDFHJLCPRTXVZNYEIWGAKMUSQO",
}

var reflectorCharMap = []string {
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
	Alphabets = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NumAlphabets = 26
)

// Enigma machine is made of several components with similar functionalities,
// namely the plugboard, rotors, and reflector. Therefore, the components
// are implemented as interface that have methods for encrypting characters and
// stepping (rotors only, others will just pass the stepping count).
type Component interface {
	processChar(ci int, reflecting bool) (int, bool)
	step(int) int
}

// A plugboard is simply a character scrambler, but does not step (or rotate)
// or change its output character after several characters like the rotors do.
type Plugboard struct {
	charMap string
}

// Maps the input character to specified output on the plugboard
func (p *Plugboard) processChar(ci int, reflecting bool) (int, bool) {
	c := rune(p.charMap[ci])
	return strings.IndexRune(Alphabets, c), reflecting
}

// Returns original step count since it's not a rotor
func (p *Plugboard) step(n int) int {
	return n
}

// Rotor always scrambles character differently every time a character enters.
// When the rotor offset reaches its notch, its next step will cause the next
// rotor to also step.
type Rotor struct {
	type_ int
	offset int
}

func CreateRotor(type_ int) Rotor {
	return Rotor{
		type_: type_,
		offset: 0,
	}
}

// Set the initial offset (or position) of the rotor. Also known as
// grundstellung.
func (r *Rotor) SetStartingPosition(pos int) {
	r.offset = pos
}

// Scrambles a character depending on its current rotor position.
func (r *Rotor) processChar(ci int, reflecting bool) (int, bool) {
	if reflecting {
		idx := (ci + r.offset) % NumAlphabets
		lc := rune(Alphabets[idx])
		ri := strings.IndexRune(rotorCharMap[r.type_], lc)
		ri -= r.offset
		if ri < 0 {
			ri += NumAlphabets
		} else {
			ri %= NumAlphabets
		}
		return ri, reflecting
	} else {
		idx := (ci + r.offset) % NumAlphabets
		rc := rune(rotorCharMap[r.type_][idx])
		li := strings.IndexRune(Alphabets, rc)
		li -= r.offset
		if li < 0 {
			li += NumAlphabets
		} else {
			li %= NumAlphabets
		}
		return li, reflecting
	}
}

// Step (or rotate) the rotor which causes the characters to scramble
// differently
func (r *Rotor) step(n int) int {
	revs := r.countNotchRevs(n)
	r.offset = (r.offset + n) % NumAlphabets
	return revs
}

// Count how many times the rotor has encountered the notch, in case it has
// a large step count
func (r *Rotor) countNotchRevs(steps int) int {
	// Return 0 if it is determined that the notch won't be reached in the steps
	nch := notch[r.type_] - 'A'
	if r.offset > nch && steps < nch + (NumAlphabets - nch) {
		return 0
	}
	return reallyCountNotchRevs(r.offset, nch, NumAlphabets, steps)
}

// Actual counting function of the notch encounters
func reallyCountNotchRevs(current, notch, max, steps int) int {
	revs := 0
	steps -= notch - current
	if steps > 0 {
		revs++
	}
	revs += steps / max
	return revs
}

// Reflector mirrors the character that enters it (e.g. if A refers to Y, then
// Y also refers mirrors to A. It is what makes the decrypting process in an
// Enigma possible, even though it means that a character can never refer to
// itself (which is a disadvantage).
type Reflector struct {
	type_ int
	charMap string
}

func CreateReflector(type_ int) Reflector {
	return Reflector{
		type_: type_,
		charMap: reflectorCharMap[type_],
	}
}

// Mirrors a character with its partner
func (r *Reflector) processChar(ci int, reflecting bool) (int, bool) {
	c := rune(reflectorCharMap[r.type_][ci])
	return strings.IndexRune(Alphabets, c), !reflecting
}

// Returns original step count since it's not a rotor
func (r *Reflector) step(n int) int {
	return n
}

type Enigma struct {
	components []Component
}

// Convenient function to create Enigma in standard configurations
// e.g. a plugboard, three rotors, and a reflector
func NewEnigma(rotorType1, rotorType2, rotorType3, reflectorType int) *Enigma {
	e := &Enigma{}
	pb := Plugboard{Alphabets}
	r1 := CreateRotor(rotorType1)
	r2 := CreateRotor(rotorType2)
	r3 := CreateRotor(rotorType3)
	rfl := CreateReflector(reflectorType)
	e.connect(&pb, &r1, &r2, &r3, &rfl)
	return e
}

// Set starting positions of the rotors (also known as Grundstellung).
func (e *Enigma) SetStartingPositions(a... int) {
	for i, v := range a {
		if v >= 65 && v <= 90 {
			v -= 'A'
		} else {
			return
		}
		e.components[i + 1].(*Rotor).SetStartingPosition(v)
	}
}

// Connect a set of Enigma components together
func (e *Enigma) connect(comps ...Component) {
	if len(comps) <= 0 {
		return
	}
	e.components = append(e.components, comps...)
}

// Encrypt a message using a chain of Enigma components. Because of the way
// Enigma works, the process of decrypting is the same as encrypting. So
// use this for decrypting as well!
func (e *Enigma) Encrypt(msg string) string {
	b := stringutil.Sanitize(msg)
	emsg := make([]rune, len(b))
	for i := range b {
		emsg[i] = e.encryptChar(rune(b[i]))
	}
	return string(emsg)
}

// Encrypt a single character
func (e *Enigma) encryptChar(c rune) rune {
	e.Step(1)

	ci := int(c - 'A')
	reflecting := false
	for i := 0; i < len(e.components); i++ {
		if ci, reflecting = e.components[i].processChar(ci, reflecting); reflecting {
			break
		}
	}

	if reflecting {
		for i := len(e.components) - 2; i >= 0; i-- {
			ci, _ = e.components[i].processChar(ci, reflecting)
		}
	}

	return rune(ci + 'A')
}

// Step all rotors that are forward-linked to this component.
func (e *Enigma) Step(steps int) {
	if steps <= 0 {
		return
	}

	for i := range e.components {
		steps = e.components[i].step(steps)
		if steps <= 0 {
			break
		}
	}
}
