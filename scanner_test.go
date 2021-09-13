package golox_test

import (
	"bytes"
	"testing"

	"github.com/goropikari/golox"
	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	runtime := golox.NewRuntime()

	var tests = []struct {
		name     string
		expected golox.TokenList
		code     string
	}{
		{
			name: "assign val",
			expected: golox.TokenList{
				golox.NewToken(golox.IdentifierTT, "x", nil, 1),
				golox.NewToken(golox.EqualTT, "=", nil, 1),
				golox.NewToken(golox.NumberTT, "1", 1.0, 1),
				golox.NewToken(golox.EOFTT, "", nil, 1),
			},
			code: "x = 1",
		},
		{
			name: "if block",
			expected: golox.TokenList{
				golox.NewToken(golox.IfTT, "if", nil, 1),
				golox.NewToken(golox.IdentifierTT, "hoge", nil, 1),
				golox.NewToken(golox.LeftBraceTT, "{", nil, 1),
				golox.NewToken(golox.IdentifierTT, "x", nil, 1),
				golox.NewToken(golox.SemicolonTT, ";", nil, 1),
				golox.NewToken(golox.RightBraceTT, "}", nil, 1),
				golox.NewToken(golox.ElseTT, "else", nil, 1),
				golox.NewToken(golox.LeftBraceTT, "{", nil, 1),
				golox.NewToken(golox.IfTT, "if", nil, 1),
				golox.NewToken(golox.IdentifierTT, "piyo", nil, 1),
				golox.NewToken(golox.LeftBraceTT, "{", nil, 1),
				golox.NewToken(golox.IdentifierTT, "y", nil, 1),
				golox.NewToken(golox.SemicolonTT, ";", nil, 1),
				golox.NewToken(golox.RightBraceTT, "}", nil, 1),
				golox.NewToken(golox.ElseTT, "else", nil, 1),
				golox.NewToken(golox.LeftBraceTT, "{", nil, 1),
				golox.NewToken(golox.IdentifierTT, "z", nil, 1),
				golox.NewToken(golox.SemicolonTT, ";", nil, 1),
				golox.NewToken(golox.RightBraceTT, "}", nil, 1),
				golox.NewToken(golox.RightBraceTT, "}", nil, 1),
				golox.NewToken(golox.EOFTT, "", nil, 1),
			},
			code: "if hoge { x; } else { if piyo { y; } else { z; } }",
			// if hoge {
			//   x
			// } else {
			//   if piyo {
			//     y;
			//   } else {
			//     z;
			//   }
			// }
		},
		{
			name: "unicode string",
			expected: golox.TokenList{
				golox.NewToken(golox.IdentifierTT, "x", nil, 1),
				golox.NewToken(golox.EqualTT, "=", nil, 1),
				golox.NewToken(golox.StringTT, "\"hoge こんにちは\\\" piyo\"", "hoge こんにちは\" piyo", 1),
				golox.NewToken(golox.EOFTT, "", nil, 1),
			},
			code: "x = \"hoge こんにちは\\\" piyo\"",
		},
		{
			name: "useless newline",
			expected: golox.TokenList{
				golox.NewToken(golox.StringTT, "\"hoge\"", "hoge", 1),
				golox.NewToken(golox.IdentifierTT, "piyo", nil, 1),
				golox.NewToken(golox.EOFTT, "", nil, 1),
			},
			code: "\n\n\"hoge\"\n\npiyo // hogehoge\n// piyopiyo\n   // fugafuga",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			buf.Write([]byte(tt.code))
			scanner := golox.NewScanner(runtime, buf)
			actual := scanner.ScanTokens()
			assert.Equal(t, tt.expected, actual)
		})
	}
}
