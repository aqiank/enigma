package enigma

import (
	"strings"
)

const (
	Plugboard = iota
	Rotor
	Reflector
)

const (
	NumAlphabets = 26
	A = 65
	Z = 90
)

type Component struct {
	out [NumAlphabets]byte
	in [NumAlphabets]byte
	offset byte
	next *Component
	prev *Component
	type_ int
}

// Create an Enigma component with default settings
func NewComponent(type_ int) *Component {
	comp := &Component{}
	comp.type_ = type_
	if comp.type_ == Reflector {
		for i := 0; i < NumAlphabets; i++ {
			comp.in[i] = byte(A + i)
			comp.out[i] = byte(Z - i)
		}
	} else {
		for i := 0; i < NumAlphabets; i++ {
			comp.in[i] = byte(A + i)
			comp.out[i] = byte(A + ((i + 1) % NumAlphabets))
		}
	}
	return comp
}

// Connect an enigma component to another
func (comp *Component) Connect(oth *Component) *Component {
	comp.next = oth
	oth.prev = comp
	return oth
}

// Encrypt a message using a chain of Enigma components
func (comp *Component) Encrypt(msg []byte) []byte {
	b := clean(msg)
	emsg := make([]byte, len(b))
	for i, _ := range b {
		comp.step()
		emsg[i] = comp.encryptChar(b[i])
	}
	return emsg
}

// Set initial settings for the Enigma component
func (comp *Component) Set(in, out string) {
	for i := 0; i < NumAlphabets; i++ {
		inc := in[i] - A // Input Character Index
		outc := out[i] - A // Output Character Index
		comp.in[outc] = inc
		comp.out[inc] = outc
	}
}

func (comp *Component) encryptChar(c byte) byte {
	r := comp
	j := c - A

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

	return A + j
}

func (comp *Component) step() {
	// Only step current component it's a rotor
	if comp.type_ == Rotor {
		comp.offset = (comp.offset + 1) % NumAlphabets
		for i := byte(0); i < NumAlphabets; i++ {
			j := (comp.out[i] + 1) % NumAlphabets
			comp.out[i] = j
			comp.in[j] = i
		}
	}
	if comp.next != nil && comp.offset == 0 {
		comp.next.step()
	}
}

func clean(msg []byte) []byte {
	s := string(msg)
	s = strings.ToUpper(s)
	s = strings.TrimSpace(s)
	s = stripChars(s, " `1234567890-=~!@#$%^&*()_+[]\\;',./{}|:\"<>?")
	return []byte(s)
}

func stripChars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}
