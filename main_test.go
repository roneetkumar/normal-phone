package main

import "testing"

func TestNormalizer(t *testing.T) {

	testCases := []struct {
		input string
		want  string
	}{
		{"1234567890", "1234567890"},
		{"123 456 7890", "1234567890"},
		{"(123) 456 7890", "1234567890"},
		{"(123) 456 - 7890", "1234567890"},
		{"123-456-7890", "1234567890"},
		{"(123)456-7890", "1234567890"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := normalizer(tc.input)
			if actual != tc.want {
				t.Errorf("got %s; want %s", actual, tc.want)
			}
		})
	}
}
