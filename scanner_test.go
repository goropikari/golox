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
				tlps.NewToken(tlps.IdentifierTT, "x", nil, 1),
				tlps.NewToken(tlps.EqualTT, "=", nil, 1),
				tlps.NewToken(tlps.NumberTT, "1", 1.0, 1),
				tlps.NewToken(tlps.EOFTT, "", nil, 1),
			},
			code: "x = 1",
		},
		{
			name: "if block",
			expected: tlps.TokenList{
				tlps.NewToken(tlps.IfTT, "if", nil, 1),
				tlps.NewToken(tlps.IdentifierTT, "hoge", nil, 1),
				tlps.NewToken(tlps.ColonTT, ":", nil, 1),
				tlps.NewToken(tlps.NewlineTT, "\\n", nil, 1),
				tlps.NewToken(tlps.LeftBraceTT, "{", nil, 2),
				tlps.NewToken(tlps.IdentifierTT, "x", nil, 2),
				tlps.NewToken(tlps.NewlineTT, "\\n", nil, 2),
				tlps.NewToken(tlps.RightBraceTT, "}", nil, 3),
				tlps.NewToken(tlps.ElseTT, "else", nil, 3),
				tlps.NewToken(tlps.ColonTT, ":", nil, 3),
				tlps.NewToken(tlps.NewlineTT, "\\n", nil, 3),
				tlps.NewToken(tlps.LeftBraceTT, "{", nil, 4),
				tlps.NewToken(tlps.IfTT, "if", nil, 4),
				tlps.NewToken(tlps.IdentifierTT, "piyo", nil, 4),
				tlps.NewToken(tlps.ColonTT, ":", nil, 4),
				tlps.NewToken(tlps.NewlineTT, "\\n", nil, 4),
				tlps.NewToken(tlps.LeftBraceTT, "{", nil, 5),
				tlps.NewToken(tlps.IdentifierTT, "y", nil, 5),
				tlps.NewToken(tlps.NewlineTT, "\\n", nil, 5),
				tlps.NewToken(tlps.RightBraceTT, "}", nil, 6),
				tlps.NewToken(tlps.ElseTT, "else", nil, 6),
				tlps.NewToken(tlps.ColonTT, ":", nil, 6),
				tlps.NewToken(tlps.NewlineTT, "\\n", nil, 6),
				tlps.NewToken(tlps.LeftBraceTT, "{", nil, 7),
				tlps.NewToken(tlps.IdentifierTT, "z", nil, 7),
				tlps.NewToken(tlps.NewlineTT, "\\n", nil, 7),
				tlps.NewToken(tlps.RightBraceTT, "}", nil, 8),
				tlps.NewToken(tlps.RightBraceTT, "}", nil, 8),
				tlps.NewToken(tlps.EOFTT, "", nil, 8),
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
				tlps.NewToken(tlps.IdentifierTT, "x", nil, 1),
				tlps.NewToken(tlps.EqualTT, "=", nil, 1),
				tlps.NewToken(tlps.StringTT, "\"hoge こんにちは\\\" piyo\"", "hoge こんにちは\" piyo", 1),
				tlps.NewToken(tlps.EOFTT, "", nil, 1),
			},
			code: "x = \"hoge こんにちは\\\" piyo\"",
		},
		{
			name: "useless newline",
			expected: tlps.TokenList{
				tlps.NewToken(tlps.StringTT, "\"hoge\"", "hoge", 3),
				tlps.NewToken(tlps.NewlineTT, "\\n", nil, 3),
				tlps.NewToken(tlps.IdentifierTT, "piyo", nil, 5),
				tlps.NewToken(tlps.NewlineTT, "\\n", nil, 5),
				tlps.NewToken(tlps.EOFTT, "", nil, 7),
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
