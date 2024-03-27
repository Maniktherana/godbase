package writer

import (
	"bytes"
	"testing"

	"github.com/maniktherana/godbase/pkg/resp"
	"github.com/stretchr/testify/assert"
)

func TestWriterWrite(t *testing.T) {
	tt := []struct {
		name       string
		inputValue resp.Value
		expected   string
	}{
		{
			name: "String",
			inputValue: resp.Value{
				Typ: "string",
				Str: "hello",
			},
			expected: "+hello\r\n",
		},
		{
			name: "Error",
			inputValue: resp.Value{
				Typ: "error",
				Str: "ERR Something went wrong",
			},
			expected: "-ERR Something went wrong\r\n",
		},

		{
			name: "Bulk",
			inputValue: resp.Value{
				Typ:  "bulk",
				Bulk: "value",
			},
			expected: "$5\r\nvalue\r\n",
		},
		{
			name: "Array",
			inputValue: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "string", Str: "hello"},
					{Typ: "string", Str: "world"},
				},
			},
			expected: "*2\r\n+hello\r\n+world\r\n",
		},
		{
			name: "Integer",
			inputValue: resp.Value{
				Typ: "integer",
				Num: 123,
			},
			expected: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := NewWriter(&buf)
			err := writer.Write(tc.inputValue)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, buf.String())
		})
	}
}
