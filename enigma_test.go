package enigma

import (
	"testing"
)

const (
	Input = "Hello World."
	Output = "QRIGBRXSWC"
)

func TestDefaultEncrypt(t *testing.T) {
	pb := NewComponent(Plugboard)
	r1 := NewComponent(Rotor)
	r2 := NewComponent(Rotor)
	r3 := NewComponent(Rotor)
	rf := NewComponent(Reflector)
	pb.Connect(r1).Connect(r2).Connect(r3).Connect(rf)
	output := pb.Encrypt([]byte(Input))
	if string(output) != Output {
		t.Fatal("Encrypted text doesn't matched intended output")
	}
}
