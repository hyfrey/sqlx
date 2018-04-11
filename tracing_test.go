package sqlx

import "testing"

func TestFirstWord(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty query", args{""}, ""},
		{"one word", args{"select"}, "select"},
		{"one word with space", args{"select "}, "select"},
		{"select from", args{"SELECT * FROM t"}, "select"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := firstWord(tt.args.query); got != tt.want {
				t.Errorf("firstWord() = %v, want %v", got, tt.want)
			}
		})
	}
}
