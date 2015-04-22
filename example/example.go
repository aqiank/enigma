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
	m.SetStartingPositions(0, 2, 0)
	output := m.Encrypt(text)
	fmt.Println("Encrypted", text, "as", output)

}
