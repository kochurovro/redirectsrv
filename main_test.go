package main

import (
	"testing"
)

func TestIsUrl(t *testing.T) {
	type compfunc func(got bool) bool
	f := func(name string, compare compfunc, in string) {
		t.Run(name, func(t *testing.T) {
			if !compare(IsUrl(in)) {
				t.Error("wrong result")
			}
		})
	}
	f("empty url should return false", func(got bool) bool { return got == false }, "")
	f("wrong url should return false", func(got bool) bool { return got == false }, "google.com/")
	f("correct url should return true", func(got bool) bool { return got == true }, "http://google.com/")
}
