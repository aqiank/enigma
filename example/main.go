package main

import (
	"fmt"

	"enigma"
)

const (
	text = `Lorem ipsum dolor sit amet, consectetur adipiscing elit.`
)

var m1, m2 *enigma.Component

func main() {
	var err error
	m1, err = enigma.FromJSONFile("enigma.json")
	if err != nil {
		fmt.Println(err)
	}
	m2, err = enigma.FromJSONFile("enigma.json")
	if err != nil {
		fmt.Println(err)
	}
	emsg := encode([]byte(text))
	decode(emsg)
}

func encode(msg []byte) []byte {
	emsg := m1.Encrypt(msg)
	println("Encoding " + string(msg) + " as " + string(emsg))
	return emsg
}

func decode(msg []byte) []byte {
	dmsg := m2.Encrypt(msg)
	println("Decoding " + string(msg) + " as " + string(dmsg))
	return dmsg
}
