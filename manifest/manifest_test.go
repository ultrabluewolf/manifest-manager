package manifest

import (
	"testing"

	"io/ioutil"
	"strings"
)

func TestParseManifestFile(t *testing.T) {
	cases := []struct {
		in   string
		want Manifest
	}{
		{
			"../mocks/manifest.1.txt",
			Manifest{
				Filename: "../mocks/manifest.1.txt",
				Files: map[string]bool{
					"../mocks/files/one.txt":   true,
					"../mocks/files/two.out":   true,
					"../mocks/files/three.log": true,
				},
			},
		},
		{
			"../mocks/manifest.2.txt",
			Manifest{
				Filename: "../mocks/manifest.2.txt",
				Files:    map[string]bool{},
			},
		},
	}

	for _, c := range cases {
		got, err := ParseManifestFile(c.in)
		if err != nil || got.Filename != c.in || len(got.Files) != len(c.want.Files) {
			t.Errorf("ParseManifestFile(%q) == %q, want %q", c.in, got, c.want)
		}
		for k, v := range got.Files {
			if v2, ok := c.want.Files[k]; !ok || v != v2 {
				t.Errorf("ParseManifestFile(%q) == %q, want %q", c.in, got, c.want)
			}
		}
		for k, v := range c.want.Files {
			if v2, ok := got.Files[k]; !ok || v != v2 {
				t.Errorf("ParseManifestFile(%q) == %q, want %q", c.in, got, c.want)
			}
		}
	}
}

func TestFileList(t *testing.T) {
	cases := []struct {
		in   Manifest
		want []string
	}{
		{
			Manifest{
				Filename: "../mocks/manifest.1.txt",
				Files: map[string]bool{
					"one.txt":   true,
					"two.out":   true,
					"three.log": true,
				},
			},
			[]string{
				"one.txt",
				"two.out",
				"three.log",
			},
		},
		{
			Manifest{
				Filename: "../mocks/manifest.2.txt",
				Files:    map[string]bool{},
			},
			[]string{},
		},
		{
			Manifest{
				Filename: "../mocks/manifest.x.txt",
				Files: map[string]bool{
					"one.txt":   true,
					"two.out":   true,
					"dne.bin":   false,
					"three.log": true,
				}},
			[]string{
				"one.txt",
				"two.out",
				"three.log",
			},
		},
	}

	for _, c := range cases {
		got := c.in.FileList()
		for _, k := range got {
			if v, ok := c.in.Files[k]; !ok || v != true {
				t.Errorf("%q.FileList() == %q, want %q", c.in, got, c.want)
			}
		}
		for k1, isActive := range c.in.Files {
			found := false
			for _, k2 := range got {
				if k1 == k2 {
					found = true
					break
				}
			}
			if !found && isActive {
				t.Errorf("%q.FileList() == %q, want %q", c.in, got, c.want)
			}
		}
	}
}

func TestPrune(t *testing.T) {
	cases := []struct {
		in   Manifest
		want Manifest
	}{
		{
			Manifest{
				Filename: "../mocks/manifest.1.txt",
				Files: map[string]bool{
					"../mocks/files/one.txt":   true,
					"../mocks/files/two.out":   true,
					"../mocks/files/three.log": true,
				},
			},
			Manifest{
				Filename: "../mocks/manifest.1.txt",
				Files: map[string]bool{
					"../mocks/files/one.txt":   true,
					"../mocks/files/two.out":   true,
					"../mocks/files/three.log": true,
				},
			},
		},
		{
			Manifest{
				Filename: "../mocks/manifest.2.txt",
				Files:    map[string]bool{},
			},
			Manifest{
				Filename: "../mocks/manifest.2.txt",
				Files:    map[string]bool{},
			},
		},
		{
			Manifest{
				Filename: "../mocks/manifest.x.txt",
				Files: map[string]bool{
					"../mocks/files/one.txt":   true,
					"../mocks/files/two.out":   true,
					"../mocks/files/dne.bin":   true,
					"../mocks/files/three.log": true,
					"../mocks/files/a-path":    true,
				},
			},
			Manifest{
				Filename: "../mocks/manifest.x.txt",
				Files: map[string]bool{
					"../mocks/files/one.txt":   true,
					"../mocks/files/two.out":   true,
					"../mocks/files/dne.bin":   false,
					"../mocks/files/three.log": true,
					"../mocks/files/a-path":    false,
				},
			},
		},
	}

	for _, c := range cases {
		got := *(&c.in)
		err := got.Prune()
		if err != nil || got.Filename != c.want.Filename || len(got.Files) != len(c.want.Files) {
			t.Errorf("%q.Prune() => %q, want %q", c.in, got, c.want)
		}
		for k, v := range got.Files {
			if v2, ok := c.want.Files[k]; !ok || v != v2 {
				t.Errorf("%q.Prune() => %q, want %q", c.in, got, c.want)
			}
		}
		for k, v := range c.want.Files {
			if v2, ok := got.Files[k]; !ok || v != v2 {
				t.Errorf("%q.Prune() => %q, want %q", c.in, got, c.want)
			}
		}
	}
}

func TestSave(t *testing.T) {
	cases := []struct {
		in Manifest
	}{
		{
			Manifest{
				Filename: "../mocks/manifest.1.txt",
				Files: map[string]bool{
					"../mocks/files/one.txt":   true,
					"../mocks/files/two.out":   true,
					"../mocks/files/three.log": true,
				},
			},
		},
		{
			Manifest{
				Filename: "../mocks/manifest.2.txt",
				Files:    map[string]bool{},
			},
		},
		{
			Manifest{
				Filename: "../mocks/manifest.x.txt",
				Files: map[string]bool{
					"../mocks/files/one.txt":     true,
					"../mocks/files/two.out":     true,
					"../mocks/files/dne.bin":     false,
					"../mocks/files/dne/dne.bin": false,
					"../mocks/files/three.log":   true,
					"../mocks/files/a-path":      false,
				},
			},
		},
	}

	for _, c := range cases {
		if err := c.in.Save(); err != nil {
			t.Error(err.Error())
		}

		result, err := ioutil.ReadFile(c.in.Filename)
		files := c.in.FileList()
		expected := strings.Join(files, "\n") + "\n"
		if len(files) == 0 {
			expected = ""
		}
		if err != nil || string(result) != expected {
			t.Errorf("%q.Save() => %q, want %q", c.in, string(result), expected)
		}
	}
}

func TestAdd(t *testing.T) {
	cases := []struct {
		in    Manifest
		addIn string
		want  Manifest
	}{
		{
			Manifest{
				Filename: "../mocks/manifest.1.txt",
				Files: map[string]bool{
					"../mocks/files/one.txt":   true,
					"../mocks/files/two.out":   true,
					"../mocks/files/three.log": true,
					"../mocks/files/a-path":    false,
				},
			},
			"../mocks/files/*",
			Manifest{
				Filename: "../mocks/manifest.1.txt",
				Files: map[string]bool{
					"../mocks/files/one.txt":   true,
					"../mocks/files/two.out":   true,
					"../mocks/files/three.log": true,
					"../mocks/files/a-path":    false,
				},
			},
		},
		{
			Manifest{
				Filename: "../mocks/manifest.2.txt",
				Files:    map[string]bool{},
			},
			"../mocks/files/*",
			Manifest{
				Filename: "../mocks/manifest.2.txt",
				Files: map[string]bool{
					"../mocks/files/one.txt":   true,
					"../mocks/files/two.out":   true,
					"../mocks/files/three.log": true,
					"../mocks/files/a-path":    false,
				},
			},
		},
		{
			Manifest{
				Filename: "../mocks/manifest.x.txt",
				Files: map[string]bool{
					"../mocks/files/one.txt":   true,
					"../mocks/files/two.out":   true,
					"../mocks/files/dne.bin":   true,
					"../mocks/files/three.log": true,
					"../mocks/files/a-path":    false,
				},
			},
			"../mocks/files/*",
			Manifest{
				Filename: "../mocks/manifest.x.txt",
				Files: map[string]bool{
					"../mocks/files/one.txt":   true,
					"../mocks/files/two.out":   true,
					"../mocks/files/dne.bin":   false,
					"../mocks/files/three.log": true,
					"../mocks/files/a-path":    false,
				},
			},
		},
	}

	for _, c := range cases {
		got := *(&c.in)
		err := got.Add(c.addIn)
		if err != nil || got.Filename != c.want.Filename || len(got.Files) != len(c.want.Files) {
			t.Errorf("%q.Add(%q) => %q, want %q", c.in, c.addIn, got, c.want)
		}
		for k, v := range got.Files {
			if v2, ok := c.want.Files[k]; !ok || v != v2 {
				t.Errorf("%q.Add(%q) => %q, want %q", c.in, c.addIn, got, c.want)
			}
		}
		for k, v := range c.want.Files {
			if v2, ok := got.Files[k]; !ok || v != v2 {
				t.Errorf("%q.Add(%q) => %q, want %q", c.in, c.addIn, got, c.want)
			}
		}
	}
}

func TestRemove(t *testing.T) {
	cases := []struct {
		in       Manifest
		removeIn string
		want     Manifest
	}{
		{
			Manifest{
				Filename: "../mocks/manifest.1.txt",
				Files: map[string]bool{
					"../mocks/files/one.txt":   true,
					"../mocks/files/two.out":   true,
					"../mocks/files/three.log": true,
				},
			},
			"../mocks/files/*",
			Manifest{
				Filename: "../mocks/manifest.1.txt",
				Files: map[string]bool{
					"../mocks/files/one.txt":   false,
					"../mocks/files/two.out":   false,
					"../mocks/files/three.log": false,
					"../mocks/files/a-path":    false,
				},
			},
		},
		{
			Manifest{
				Filename: "../mocks/manifest.2.txt",
				Files:    map[string]bool{},
			},
			"../mocks/files/*",
			Manifest{
				Filename: "../mocks/manifest.2.txt",
				Files: map[string]bool{
					"../mocks/files/one.txt":   false,
					"../mocks/files/two.out":   false,
					"../mocks/files/three.log": false,
					"../mocks/files/a-path":    false,
				},
			},
		},
		{
			Manifest{
				Filename: "../mocks/manifest.x.txt",
				Files: map[string]bool{
					"../mocks/files/one.txt":     true,
					"../mocks/files/two.out":     true,
					"../mocks/files/dne.bin":     false,
					"../mocks/files/dne/dne.bin": true,
					"../mocks/files/three.log":   true,
				},
			},
			"../mocks/files/*",
			Manifest{
				Filename: "../mocks/manifest.x.txt",
				Files: map[string]bool{
					"../mocks/files/one.txt":     false,
					"../mocks/files/two.out":     false,
					"../mocks/files/dne.bin":     false,
					"../mocks/files/dne/dne.bin": false,
					"../mocks/files/three.log":   false,
					"../mocks/files/a-path":      false,
				},
			},
		},
	}

	for _, c := range cases {
		got := *(&c.in)
		err := got.Remove(c.removeIn)
		if err != nil || got.Filename != c.want.Filename || len(got.Files) != len(c.want.Files) {
			t.Errorf("%+v.Remove(%q) => %q, want %+v", c.in, c.removeIn, got, c.want)
		}
		for k, v := range got.Files {
			if v2, ok := c.want.Files[k]; !ok || v != v2 {
				t.Errorf("%+v.Remove(%q) => %q, want %+v", c.in, c.removeIn, got, c.want)
			}
		}
		for k, v := range c.want.Files {
			if v2, ok := got.Files[k]; !ok || v != v2 {
				t.Errorf("%+v.Remove(%q) => %q, want %+v", c.in, c.removeIn, got, c.want)
			}
		}
	}
}
