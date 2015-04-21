package enigma

import (
	"strings"
)

// Types of Components
const (
	Plugboard = iota
	Rotor
	Reflector
)

const (
	NumAlphabets = 26
)

// Enigma machine is made of several components with similar functionalities,
// namely the plugboard, rotors, and reflector. Each has an input and output
// "sockets" and they are all chained to each other. Therefore, the components
// are implemented as a double-linked lists that have input and output
// character maps. The maps are implemented in simple arrays of bytes instead
// of Go's map for simplicity.
type Component struct {
	out    [NumAlphabets]byte // in->out map of characters e.g. out[in]
	in     [NumAlphabets]byte // out->in map of characters e.g. in[out]
	offset byte               // offset is used by Rotors and ignored by other components
	next   *Component
	prev   *Component
	type_  int
}

// Create an Enigma component with default settings
func NewComponent(type_ int) *Component {
	comp := &Component{}
	comp.type_ = type_
	if comp.type_ == Reflector {
		for i := 0; i < NumAlphabets; i++ {
			comp.in[i] = byte(i)
			comp.out[i] = byte(NumAlphabets - (i + 1))
		}
	} else {
		for i := 0; i < NumAlphabets; i++ {
			comp.in[i] = byte(i)
			comp.out[i] = byte(i)
		}
	}
	return comp
}

// Connect an enigma component to another
func Connect(comps ...*Component) *Component{
	if len(comps) <= 0 {
		return nil
	}

	for i := 0; i < len(comps) - 1; i++ {
		comps[i].next = comps[i + 1]
		comps[i + 1].prev = comps[i]
	}
	return comps[0]
}

// Encrypt a message using a chain of Enigma components. Because of the way
// Enigma works, the process of decrypting is the same as encrypting. So
// use this for decrypting as well!
func (comp *Component) Encrypt(msg []byte) []byte {
	b := sanitizeString(msg)
	emsg := make([]byte, len(b))
	for i, _ := range b {
		comp.step()
		emsg[i] = comp.encryptChar(b[i])
	}
	return emsg
}

// Set initial settings for the Enigma component
func (comp *Component) SetCharacterMap(in, out string) {
	for i := 0; i < NumAlphabets; i++ {
		inc := in[i] - 'A'   // Input Character Index
		outc := out[i] - 'A' // Output Character Index
		comp.in[outc] = inc
		comp.out[inc] = outc
	}
}

// Offset current rotor component by n position. Should only be used before
// encryption (with the exception of an internal function that calls it).
func (comp *Component) OffsetBy(n byte) {
	if comp.type_ != Rotor {
		return
	}

	comp.offset = (comp.offset + n) % NumAlphabets
	for i := byte(0); i < NumAlphabets; i++ {
		j := (comp.out[i] + n) % NumAlphabets
		comp.out[i] = j
		comp.in[j] = i
	}
}

// Encrypt a single character
func (comp *Component) encryptChar(c byte) byte {
	r := comp
	j := c - 'A'

	// Run character through the rotors
	for ; r != nil; r = r.next {
		j = r.out[j]
		if r.type_ == Reflector {
			break
		}
	}

	// Reflecting
	for r = r.prev; r != nil; r = r.prev {
		j = r.in[j]
	}

	return 'A' + j
}

// Step all rotors that are forward-linked to this component.
func (comp *Component) step() {
	comp.OffsetBy(1)

	// Rotate the next component on the condition that the current rotor has
	// done a full revolution, or if the current component is not a rotor.
	// It knows that the component is not a rotor when the offset is zero,
	// and even if it is a rotor, the next component should step anyway.
	if comp.next != nil && comp.offset == 0 {
		comp.next.step()
	}
}

// Remove unacceptable characters from message
func sanitizeString(msg []byte) []byte {
	s := string(msg)
	s = strings.ToUpper(s)
	s = strings.TrimSpace(s)
	s = stripChars(s, " `1234567890-=~!@#$%^&*()_+[]\\;',./{}|:\"<>?")
	return []byte(s)
}

// Strip a set of characters from string
func stripChars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}
