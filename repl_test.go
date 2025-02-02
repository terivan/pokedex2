package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {

	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "This is a test",
			expected: []string{"this", "is", "a", "test"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		length := len(actual)
		t.Log(actual)
		t.Log(c.expected)
		if length != len(c.expected) {
			t.Errorf("Output slice length doesn't match")
			t.Fatal()
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			if expectedWord != word {
				t.Errorf("Wrong word!")
				t.Fatal()
			}
		}
	}
}
