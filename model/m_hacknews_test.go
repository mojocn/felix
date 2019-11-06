package model

import "testing"

func TestTranslateEn2Ch(t *testing.T) {

	got, err := TranslateEn2Ch("Worldwide observations confirm nearby ‘lensing’ exoplanet")
	if err != nil {
		t.Log(err)
	}
	t.Log(got)

}
