package files

import (
	"testing"
)

func TestExists(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"./files.go", true},
		{"./dne/dne.txt", false},
	}

	for _, c := range cases {
		if got := Exists(c.in); got != c.want {
			t.Errorf("Exists(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
