package main

import (
	"fmt"

	"enigma"
)

const (
	text = `AAAAAAAAAAAAAAAAAAAAAAAAAA`
)

func main() {
	m := enigma.CreateEnigma(enigma.RotorIII, enigma.RotorII, enigma.RotorI, enigma.ReflectorB)
	m.SetOffsets('A', 'Z', 'A')
	output := m.Encrypt(text)
	fmt.Println("Encrypted", text, "as", output)
}
