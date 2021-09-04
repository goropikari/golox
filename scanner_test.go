package mylang_test

import (
	"bytes"
	"testing"

	"github.com/goropikari/mylang"
	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	runtime := mylang.NewRuntime()

	var tests = []struct {
		name     string
		expected mylang.TokenList
		code     string
	}{
		{
			name: "assign val",
			expected: mylang.TokenList{
				mylang.NewToken(mylang.Identifier, "x", nil, 1),
				mylang.NewToken(mylang.Equal, "=", nil, 1),
				mylang.NewToken(mylang.Number, "1", 1.0, 1),
				mylang.NewToken(mylang.EOF, "", nil, 1),
			},
			code: "x = 1",
		},
		{
			name: "if block",
			expected: mylang.TokenList{
				mylang.NewToken(mylang.If, "if", nil, 1),
				mylang.NewToken(mylang.Identifier, "hoge", nil, 1),
				mylang.NewToken(mylang.Colon, ":", nil, 1),
				mylang.NewToken(mylang.LeftBrace, "{", nil, 2),
				mylang.NewToken(mylang.Identifier, "x", nil, 2),
				mylang.NewToken(mylang.RightBrace, "}", nil, 3),
				mylang.NewToken(mylang.Else, "else", nil, 3),
				mylang.NewToken(mylang.Colon, ":", nil, 3),
				mylang.NewToken(mylang.LeftBrace, "{", nil, 4),
				mylang.NewToken(mylang.If, "if", nil, 4),
				mylang.NewToken(mylang.Identifier, "piyo", nil, 4),
				mylang.NewToken(mylang.Colon, ":", nil, 4),
				mylang.NewToken(mylang.LeftBrace, "{", nil, 5),
				mylang.NewToken(mylang.Identifier, "y", nil, 5),
				mylang.NewToken(mylang.RightBrace, "}", nil, 6),
				mylang.NewToken(mylang.Else, "else", nil, 6),
				mylang.NewToken(mylang.Colon, ":", nil, 6),
				mylang.NewToken(mylang.LeftBrace, "{", nil, 7),
				mylang.NewToken(mylang.Identifier, "z", nil, 7),
				mylang.NewToken(mylang.RightBrace, "}", nil, 7),
				mylang.NewToken(mylang.RightBrace, "}", nil, 7),
				mylang.NewToken(mylang.EOF, "", nil, 7),
			},
			code: "if hoge:\n  x\nelse:\n  if piyo:\n    y\n  else:\n    z",
			// if hoge:
			//   x
			// else:
			//   if piyo:
			//     y
			//   else:
			//     z
		},
		{
			name: "unicode string",
			expected: mylang.TokenList{
				mylang.NewToken(mylang.Identifier, "x", nil, 1),
				mylang.NewToken(mylang.Equal, "=", nil, 1),
				mylang.NewToken(mylang.String, "\"hoge こんにちは piyo\"", []rune("hoge こんにちは piyo"), 1),
				mylang.NewToken(mylang.EOF, "", nil, 1),
			},
			code: "x = \"hoge こんにちは piyo\"",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			buf.Write([]byte(tt.code))
			scanner := mylang.NewScanner(runtime, buf)
			actual := scanner.ScanTokens()
			assert.Equal(t, tt.expected, actual)
		})
	}
}
