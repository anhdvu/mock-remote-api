package wtools

import (
	"testing"
)

func TestCalculateSHA1Checksum(t *testing.T) {
	SHA1Cases := []struct {
		inPk  string
		inSth string
		want  string
	}{
		{"Hello, World", "dlroW ,olleH", ""},
		{"Hello, 世界", "界世 ,olleH", ""},
		{"", "", ""},
	}

	for _, c := range SHA1Cases {
		got := CalculateSHA1Checksum(c.inPk, c.inSth)
		if got != c.want {
			t.Errorf("CalculateSHA1Checksum(%q, %q) == %q, want %q", c.inPk, c.inSth, got, c.want)
		}
	}
}

func TestCalculateSHA256Checksum(t *testing.T) {
	SHA256Cases := []struct {
		inPk  string
		inSth string
		want  string
	}{
		{},
		{},
		{},
	}

	for _, c := range SHA256Cases {
		got := CalculateSHA256Checksum(c.inPk, c.inSth)
		if got != c.want {
			t.Errorf("CalculateSHA1Checksum(%q, %q) == %q, want %q", c.inPk, c.inSth, got, c.want)
		}
	}
}
