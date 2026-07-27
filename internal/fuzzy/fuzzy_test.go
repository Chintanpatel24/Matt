package fuzzy

import "testing"

func TestMatch(t *testing.T) {
	tests := []struct {
		pattern  string
		target   string
		expected bool
	}{
		{"", "main.go", true},
		{"go", "main.go", true},
		{"mgo", "main.go", true},
		{"xyz", "main.go", false},
		{"README", "readme.md", true},
	}

	for _, tt := range tests {
		result := Match(tt.pattern, tt.target)
		if result != tt.expected {
			t.Errorf("Match(%q, %q) = %v; want %v", tt.pattern, tt.target, result, tt.expected)
		}
	}
}
