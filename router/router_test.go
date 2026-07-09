package router

import (
	"reflect"
	"testing"
)

// testRouter returns a Router registered with the bot's real command set.
func testRouter() *Router {
	m := New()
	for _, p := range []string{"help", "pp", "simp", "gayrate", "waifu", "stankrate", "epicgamer"} {
		m.Route(p, "desc", nil)
	}
	return m
}

func TestFuzzyMatch(t *testing.T) {
	m := testRouter()

	tests := []struct {
		name string
		msg  string
		want string // matched pattern, "" for no match
		args []string
	}{
		{"exact", "help", "help", nil},
		{"exact with args", "pp @someone now", "pp", []string{"@someone", "now"}},
		{"case insensitive", "HeLp Me", "help", []string{"me"}},
		{"unique prefix", "stank", "stankrate", nil},
		{"unique prefix epic", "epic", "epicgamer", nil},
		{"unique single letter prefix", "p", "pp", nil},
		{"ambiguous prefix", "s", "", nil},
		{"typo transposition", "hlep", "help", nil},
		{"typo transposition waifu", "wiafu", "waifu", nil},
		{"typo missing letter", "stankrat", "stankrate", nil},
		{"typo extra letter", "simps", "simp", nil},
		{"typo gayrate", "gayrat", "gayrate", nil},
		{"too different", "zzzzz", "", nil},
		{"short garbage skips typo tier", "xy", "", nil},
		{"mid sentence exact", "can you help me", "help", []string{"me"}},
		{"mid sentence exact last word", "show me pp", "pp", nil},
		{"mid sentence near miss ignored", "tell me a story", "", nil},
		{"empty", "", "", nil},
		{"whitespace only", "   ", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, args := m.FuzzyMatch(tt.msg)
			got := ""
			if r != nil {
				got = r.Pattern
			}
			if got != tt.want {
				t.Fatalf("FuzzyMatch(%q) matched %q, want %q", tt.msg, got, tt.want)
			}
			if len(args) != len(tt.args) {
				t.Fatalf("FuzzyMatch(%q) args = %v, want %v", tt.msg, args, tt.args)
			}
			for i := range args {
				if args[i] != tt.args[i] {
					t.Fatalf("FuzzyMatch(%q) args = %v, want %v", tt.msg, args, tt.args)
				}
			}
		})
	}
}

func TestSuggest(t *testing.T) {
	m := testRouter()

	tests := []struct {
		word string
		want []string
	}{
		{"s", []string{"simp", "stankrate"}}, // ambiguous prefix: list both
		{"hallp", []string{"help"}},          // two edits away: too far to run, close enough to suggest
		{"zzzzz", nil},
		{"hi", nil}, // short words must not suggest their whole edit ball
	}

	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			if got := m.Suggest(tt.word); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Suggest(%q) = %v, want %v", tt.word, got, tt.want)
			}
		})
	}
}

func TestEditDistance(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"", "", 0},
		{"abc", "", 3},
		{"", "abc", 3},
		{"help", "help", 0},
		{"hlep", "help", 1},  // adjacent transposition counts as one edit
		{"halp", "help", 1},  // substitution
		{"hel", "help", 1},   // insertion
		{"helpp", "help", 1}, // deletion
		{"stankrate", "epicgamer", 8},
	}

	for _, tt := range tests {
		if got := editDistance(tt.a, tt.b); got != tt.want {
			t.Errorf("editDistance(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}
