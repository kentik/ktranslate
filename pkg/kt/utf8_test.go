package kt

import (
	"testing"
)

func TestSanitizeUTF8(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: ""},
		{name: "ascii", in: "ifAlias-ok", want: "ifAlias-ok"},
		{name: "valid unicode", in: "café", want: "café"},
		{name: "invalid byte", in: "bad\xffname", want: "626164ff6e616d65"},
		{name: "only invalid", in: "\x80\xff", want: "80ff"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeUTF8(tt.in); got != tt.want {
				t.Fatalf("SanitizeUTF8(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestSanitizeUTF8WithHint(t *testing.T) {
	got, changed := SanitizeUTF8WithHint("ok")
	if got != "ok" || changed {
		t.Fatalf("valid input: got=%q changed=%v", got, changed)
	}
	got, changed = SanitizeUTF8WithHint("x\xff")
	if !changed || got != "78ff" {
		t.Fatalf("invalid input: got=%q changed=%v", got, changed)
	}
}
