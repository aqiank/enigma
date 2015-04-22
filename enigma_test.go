package enigma

import (
	"testing"
)

const (
	PlainText  = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	CipherText = "BDZGOWCXLTKSBTMCDLPBMUQOFXYHCXTGYJFLINHNXSHIUNTHEORX"
)

func TestStandardEncrypt(t *testing.T) {
	m := NewEnigma(RotorIII, RotorII, RotorI, ReflectorB)
	output := m.Encrypt(PlainText)
	t.Log(output)
	if output != CipherText {
		t.Fatal("Encrypted text doesn't matched intended output:", CipherText)
	}
}

func TestStandardDecrypt(t *testing.T) {
	m := NewEnigma(RotorIII, RotorII, RotorI, ReflectorB)
	output := m.Encrypt(CipherText)
	t.Log(output)
	if output != PlainText {
		t.Fatal("Decrypted text doesn't matched intended output:", PlainText)
	}
}
