package tlps_test

import (
	"bytes"
	"testing"

	"github.com/goropikari/tlps"
	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	runtime := tlps.NewRuntime()

	var tests = []struct {
		name     string
		expected tlps.TokenList
		code     string
	}{
		{
			name: "assign val",
			expected: tlps.TokenList{
				tlps.NewToken(tlps.Identifier, "x", nil, 1),
				tlps.NewToken(tlps.Equal, "=", nil, 1),
				tlps.NewToken(tlps.Number, "1", 1.0, 1),
				tlps.NewToken(tlps.EOF, "", nil, 1),
			},
			code: "x = 1",
		},
		{
			name: "if block",
			expected: tlps.TokenList{
				tlps.NewToken(tlps.If, "if", nil, 1),
				tlps.NewToken(tlps.Identifier, "hoge", nil, 1),
				tlps.NewToken(tlps.Colon, ":", nil, 1),
				tlps.NewToken(tlps.Newline, "\\n", nil, 1),
				tlps.NewToken(tlps.LeftBrace, "{", nil, 2),
				tlps.NewToken(tlps.Identifier, "x", nil, 2),
				tlps.NewToken(tlps.Newline, "\\n", nil, 2),
				tlps.NewToken(tlps.RightBrace, "}", nil, 3),
				tlps.NewToken(tlps.Else, "else", nil, 3),
				tlps.NewToken(tlps.Colon, ":", nil, 3),
				tlps.NewToken(tlps.Newline, "\\n", nil, 3),
				tlps.NewToken(tlps.LeftBrace, "{", nil, 4),
				tlps.NewToken(tlps.If, "if", nil, 4),
				tlps.NewToken(tlps.Identifier, "piyo", nil, 4),
				tlps.NewToken(tlps.Colon, ":", nil, 4),
				tlps.NewToken(tlps.Newline, "\\n", nil, 4),
				tlps.NewToken(tlps.LeftBrace, "{", nil, 5),
				tlps.NewToken(tlps.Identifier, "y", nil, 5),
				tlps.NewToken(tlps.Newline, "\\n", nil, 5),
				tlps.NewToken(tlps.RightBrace, "}", nil, 6),
				tlps.NewToken(tlps.Else, "else", nil, 6),
				tlps.NewToken(tlps.Colon, ":", nil, 6),
				tlps.NewToken(tlps.Newline, "\\n", nil, 6),
				tlps.NewToken(tlps.LeftBrace, "{", nil, 7),
				tlps.NewToken(tlps.Identifier, "z", nil, 7),
				tlps.NewToken(tlps.Newline, "\\n", nil, 7),
				tlps.NewToken(tlps.RightBrace, "}", nil, 8),
				tlps.NewToken(tlps.RightBrace, "}", nil, 8),
				tlps.NewToken(tlps.EOF, "", nil, 8),
			},
			code: "if hoge:\n  x\nelse:\n  if piyo:\n    y\n  else:\n    z\n",
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
			expected: tlps.TokenList{
				tlps.NewToken(tlps.Identifier, "x", nil, 1),
				tlps.NewToken(tlps.Equal, "=", nil, 1),
				tlps.NewToken(tlps.String, "\"hoge こんにちは\\\" piyo\"", "hoge こんにちは\" piyo", 1),
				tlps.NewToken(tlps.EOF, "", nil, 1),
			},
			code: "x = \"hoge こんにちは\\\" piyo\"",
		},
		{
			name: "useless newline",
			expected: tlps.TokenList{
				tlps.NewToken(tlps.String, "\"hoge\"", "hoge", 3),
				tlps.NewToken(tlps.Newline, "\\n", nil, 3),
				tlps.NewToken(tlps.Identifier, "piyo", nil, 5),
				tlps.NewToken(tlps.Newline, "\\n", nil, 5),
				tlps.NewToken(tlps.EOF, "", nil, 7),
			},
			code: "\n\n\"hoge\"\n\npiyo // hogehoge\n// piyopiyo\n   // fugafuga",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			buf.Write([]byte(tt.code))
			scanner := tlps.NewScanner(runtime, buf)
			actual := scanner.ScanTokens()
			assert.Equal(t, tt.expected, actual)
		})
	}
}
