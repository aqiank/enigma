package main

import (
	"fmt"

	"enigma"
)

const (
	//text = `Lorem ipsum dolor sit amet, consectetur adipiscing elit.`
	text = `AAAAAAAAAAAAAAAAAAAAAAAAAAAA`
)

func main() {
	m := enigma.NewStandardEnigma(enigma.RotorIII, enigma.RotorII, enigma.RotorI, enigma.ReflectorB)
	output := m.Encrypt([]byte(text))
	fmt.Println("Encrypted", text, "as", string(output))
}
