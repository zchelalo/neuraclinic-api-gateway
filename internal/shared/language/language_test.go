package language

import "testing"

func TestResolveHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "spanish", value: "es", want: Spanish},
		{name: "english", value: "en", want: English},
		{name: "regional with quality", value: "es-MX,en;q=0.8", want: Spanish},
		{name: "empty", value: "", want: English},
		{name: "invalid", value: "fr-CA,pt-BR;q=0.8", want: English},
		{name: "skip invalid and use later supported", value: "fr, en-US;q=0.7", want: English},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := ResolveHeader(tt.value); got != tt.want {
				t.Fatalf("ResolveHeader(%q) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}
