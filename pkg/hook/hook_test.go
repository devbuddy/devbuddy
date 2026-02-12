package hook

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShellEscapeDoubleQuoted(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "plain string",
			input: "hello world",
			want:  "hello world",
		},
		{
			name:  "double quotes",
			input: `say "hi"`,
			want:  `say \"hi\"`,
		},
		{
			name:  "backslash",
			input: `path\to\file`,
			want:  `path\\to\\file`,
		},
		{
			name:  "json with nested quotes",
			input: `{"key":"value"}`,
			want:  `{\"key\":\"value\"}`,
		},
		{
			name:  "json with escaped inner quotes (marshaled nested JSON)",
			input: `{"param":"{\"FOO\":\"bar\"}"}`,
			want:  `{\"param\":\"{\\\"FOO\\\":\\\"bar\\\"}\"}`,
		},
		{
			name:  "dollar sign",
			input: `$HOME`,
			want:  `\$HOME`,
		},
		{
			name:  "backtick",
			input: "run `cmd`",
			want:  "run \\`cmd\\`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellEscapeDoubleQuoted(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}
