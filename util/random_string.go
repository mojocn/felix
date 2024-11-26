package util

import (
	"fmt"
	"math/rand"
)

var (
	stringMisc = []byte(".$#@&*_")

	stringDigit    = []byte("1234567890")
	stringDigitLen = len(stringDigit)

	stringLword    = []byte("abcdefghijklmnopqrstuvwxyz")
	stringLwordLen = len(stringLword)

	stringUWord    = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	stringUWordLen = len(stringUWord)

	stringCWord    = []byte(fmt.Sprintf("%s%s", stringLword, stringUWord))
	stringCWordLen = len(stringCWord)

	letterRunes    = []byte(fmt.Sprintf("%s%s%s%s", stringDigit, stringLword, stringMisc, stringUWord))
	letterRunesLen = len(letterRunes)
)

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(letterRunesLen)]
	}
	return string(b)
}

func RandStringWordL(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = stringLword[rand.Intn(stringLwordLen)]
	}
	return string(b)
}

func RandStringWordU(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = stringUWord[rand.Intn(stringUWordLen)]
	}
	return string(b)
}
func RandStringWordC(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = stringCWord[rand.Intn(stringCWordLen)]
	}
	return string(b)
}
func RandStringDigit(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = stringDigit[rand.Intn(stringDigitLen)]
	}
	return string(b)
}
