package main

import (
	"fmt"

	"enigma"
)

const (
	text = `AAAAAAAAAAAAAAAAAAAAAAAAAA`
)

func main() {
	m := enigma.NewEnigma(enigma.RotorIII, enigma.RotorII, enigma.RotorI, enigma.ReflectorB)
	m.SetStartingPositions('A', 'Z', 'A')
	output := m.Encrypt(text)
	fmt.Println("Encrypted", text, "as", output)

}
