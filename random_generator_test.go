package main

import "testing"

func TestGenerate(t *testing.T) {
	var body = "${title}"

	if generate(body) == "${title}" {
		t.Error("random generator not working.")
	}
}
