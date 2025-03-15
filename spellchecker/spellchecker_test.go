package spellchecker

import "testing"

func TestFindScore(t *testing.T) {
	type args struct {
		s1 string
		s2 string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "equal_eqsize", args: args{"apple", "apple"}, want: 0},
		{name: "equal_dfsize", args: args{"appletree", "apple"}, want: 4},
		{name: "empty_eqsize", args: args{"", ""}, want: 0},
		{name: "empty_dfsize", args: args{"empty", ""}, want: 5},
		{name: "dfany_eqsize", args: args{"ignition", "position"}, want: 3},
		{name: "dfany_dfsize", args: args{"situation", "position"}, want: 5},
		{name: "dfmax_eqsize", args: args{"aaaaa", "bbbbb"}, want: 5},
		{name: "dfmax_dfsize", args: args{"aaaaaaa", "bbbbb"}, want: 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindScore(tt.args.s1, tt.args.s2); got != tt.want {
				t.Errorf("FindScore() = %v, want %v", got, tt.want)
			}
		})
	}
}
