package handlers

import "testing"

func TestIsSupportedEmailDomain(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{name: "gmail", email: "user@gmail.com", want: true},
		{name: "yahoo", email: "user@yahoo.com", want: true},
		{name: "outlook", email: "user@outlook.com", want: true},
		{name: "unsupported", email: "user@example.com", want: false},
		{name: "missing domain", email: "user@", want: false},
		{name: "empty", email: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSupportedEmailDomain(tt.email); got != tt.want {
				t.Fatalf("isSupportedEmailDomain(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}
