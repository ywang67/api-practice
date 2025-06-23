package stringutil

import (
	"log"
	"os"
	"testing"

	"github.comcast.com/genome/pkg/tests"
)

func TestSliceEquals(t *testing.T) {
	testCases := []struct {
		name     string
		a, b     []string
		expected bool
	}{
		{
			name:     "identical",
			a:        []string{"a", "b", "c"},
			b:        []string{"a", "b", "c"},
			expected: true,
		}, {
			name:     "jumbled",
			a:        []string{"a", "b", "c"},
			b:        []string{"c", "a", "b"},
			expected: true,
		}, {
			name:     "diffElements",
			a:        []string{"a", "b", "c"},
			b:        []string{"d", "e", "f"},
			expected: false,
		}, {
			name:     "diffLen",
			a:        []string{"a", "b", "c"},
			b:        []string{"a", "b", "c", "d"},
			expected: false,
		}, {
			name:     "duplicatesTrue",
			a:        []string{"a", "a", "b"},
			b:        []string{"b", "a", "a"},
			expected: true,
		}, {
			name:     "duplicatesFalse",
			a:        []string{"a", "b", "b"},
			b:        []string{"b", "a", "a"},
			expected: false,
		},
	}

	for _, tt := range testCases {
		tt := tt // shadow to avoid concurrency issues
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tl := tests.Tlogger{
				Logger: log.New(os.Stderr, t.Name()+": ", 0),
				T:      t,
			}

			if tt.expected != SliceSetEquals(tt.a, tt.b) {
				tl.Fatalf("expected does not match: want %t, got %t", tt.expected, !tt.expected)
			}
		})
	}
}

func TestTrimSpecialCharter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"EmptyString", "", ""},
		{"OnlyMacAddress", "80:da:c2:e9:f5:0a", "80:da:c2:e9:f5:0a"},
		{"MacAddressWithSpace", "  80:da:c2:e9:f5:0a  ", "80:da:c2:e9:f5:0a"},
		{"MacAddressWithSpecialCharacters", "80:da:c2:e9:f5:0a:", "80:da:c2:e9:f5:0a"},
		{"Fqdn", "cbr01.roseville.mi.michigan.comcast.net", "cbr01.roseville.mi.michigan.comcast.net"},
		{"FqdnWithSpecialCharacters", "   cbr01.roseville.mi.michigan.comcast.net:::", "cbr01.roseville.mi.michigan.comcast.net"},
		{"PpodName", "NJABPP101", "NJABPP101"},
		{"PpodNameWithSpecialCharacters", ">.<NJABPP101   ", "NJABPP101"},
		{"OnlySpecialCharacters", "!@#$%^&*", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := TrimSpecialCharter(test.input)
			if result != test.expected {
				t.Errorf("Test %s failed. Expected '%s', got '%s'", test.name, test.expected, result)
			}
		})
	}
}

func BenchmarkTrimSpecialCharter(b *testing.B) {
	testString := "  80:da:c2:e9:f5:0a  "

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = TrimSpecialCharter(testString)
	}
}
