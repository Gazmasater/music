package utils

import (
	"music/internal/utils"
	"testing"
)

func TestNormalizeSongName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Song Name\u00A0-\u00A0Artist", "Song Name-Artist"},
		{" Song  Name  -  Artist ", "Song Name-Artist"},
		{"Song - Name, Artist ", "Song-Name,Artist"},
		{"  Song  Name  ", "Song Name"},
		{"Song Name -  Artist , Some Other Artist", "Song Name-Artist,Some Other Artist"},
		{"За     тебя, Родина   -   мать", "За тебя,Родина-мать"},
		{"Test    Song", "Test Song"},
	}

	for _, tt := range tests {
		result := utils.NormalizeSongName(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeSongName(%q) = %q; expected %q", tt.input, result, tt.expected)
		}
	}
}
