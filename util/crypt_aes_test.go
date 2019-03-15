package util

import "testing"

func TestAesDecryptEn(t *testing.T) {
	key := RandStringWordC(32)

	msg := RandomString(12)

	code, err := AesEncrypt([]byte(msg), key)
	if err != nil {
		t.Error(err)
	}
	tMsg, err := AesDecrypt(code, key)
	if err != nil {
		t.Error(err)
	}
	if tMsg != msg {
		t.Error("aes failed")
	}
}
