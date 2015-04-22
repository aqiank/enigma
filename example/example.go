package main

import (
	"fmt"

	"enigma"
)

const (
	text = `AAAAAAAAAAAAAAAAAAAAAAAAAAAA`
)

func main() {
	m := enigma.NewEnigma(enigma.RotorIII, enigma.RotorII, enigma.RotorI, enigma.ReflectorB)
	output := m.Encrypt(text)
	fmt.Println("Encrypted", text, "as", output)

}
