package main

import "testing"

func TestHello(t *testing.T) {
	want := "Hello, Mundial!"
	if got := Hello("Hello, Mundial!"); got != want {
		t.Errorf("Hello() got: %q, want: %q", got, want)

	}
}
