package main

import (
	"reflect"
	"testing"
)

func TestCliArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		fixture string
	}{
		{"no arguments", []string{}, "no-args.golden"},
		{"help", []string{"help"}, "help.golden"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runBinary(tt.args)
			if err != nil {
				t.Fatal(err)
			}

			if *update {
				writeFixture(t, tt.fixture, output)
			}

			actual := string(output)

			expected := loadFixture(t, tt.fixture)

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("actual = %s, expected = %s", actual, expected)
			}
		})
	}
}
