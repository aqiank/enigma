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
	Connect(pb, r1, r2, r3, rf)
	output := pb.Encrypt([]byte(Input))
	t.Log(string(output))
	if string(output) != Output {
		t.Fatal("Encrypted text doesn't matched intended output:", Output)
	}
}
