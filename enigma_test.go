package enigma

import (
	"testing"
)

const (
	PlainText  = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	CipherText = "BDZGOWCXLTKSBTMCDLPBMUQOFXYHCXTGYJFLINHNXSHIUNTHEORX"
)

func TestStandardEncrypt(t *testing.T) {
	m := NewStandardEnigma(RotorIII, RotorII, RotorI, ReflectorB)
	output := m.Encrypt([]byte(PlainText))
	t.Log(string(output))
	if string(output) != CipherText {
		t.Fatal("Encrypted text doesn't matched intended output:", CipherText)
	}
}

func TestStandardDecrypt(t *testing.T) {
	m := NewStandardEnigma(RotorIII, RotorII, RotorI, ReflectorB)
	output := m.Encrypt([]byte(CipherText))
	t.Log(string(output))
	if string(output) != PlainText {
		t.Fatal("Decrypted text doesn't matched intended output:", PlainText)
	}
}
