package common

import (
	"testing"
)

func TestDownload(t *testing.T) {
	got, err := Download("https://wos.develop.meetwhale.com/Pk_P1rgUyRkq-Qfu2ARs7")
	if err != nil {
		t.Error(err)
	}

	if len(got) != 263143 {
		t.Error("length error")
	}
}

func TestUnzip(t *testing.T) {
	got, err := Download("https://wos.develop.meetwhale.com/Pk_P1rgUyRkq-Qfu2ARs7")
	if err != nil {
		t.Error(err)
	}

	resources, err := Unzip(got)
	if err != nil {
		t.Error(err)
	}

	if len(resources) != 215 {
		t.Error("length error")
	}
}
