package resp

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadLine(t *testing.T) {
	tt := []struct {
		input    string
		expected string
	}{
		{"hello\r\n", "hello"},
		{"world\r\n", "world"},
		{"testing 123\r\n", "testing 123"},
	}

	for _, tc := range tt {
		reader := strings.NewReader(tc.input)
		respReader := NewResp(reader)
		line, _, err := respReader.readLine()
		require.NoError(t, err)
		assert.Equal(t, tc.expected, string(line))
	}
}

func TestReadInteger(t *testing.T) {
	tt := []struct {
		input    string
		expected int
	}{
		{"123\r\n", 123},
		{"456\r\n", 456},
		{"789\r\n", 789},
	}

	for _, tc := range tt {
		reader := strings.NewReader(tc.input)
		respReader := NewResp(reader)
		val, _, err := respReader.readInteger()
		require.NoError(t, err)
		assert.Equal(t, tc.expected, val)
	}
}

func TestReadArray(t *testing.T) {
	tt := []struct {
		input    string
		expected []Value
	}{
		{"3\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$3\r\nbaz\r\n", []Value{
			{Typ: "bulk", Bulk: "foo"},
			{Typ: "bulk", Bulk: "bar"},
			{Typ: "bulk", Bulk: "baz"},
		}},
		{"2\r\n$5\r\nhello\r\n$5\r\nworld\r\n", []Value{
			{Typ: "bulk", Bulk: "hello"},
			{Typ: "bulk", Bulk: "world"},
		}},
	}

	for _, tc := range tt {
		reader := strings.NewReader(tc.input)
		respReader := NewResp(reader)
		val, err := respReader.readArray()
		require.NoError(t, err)
		assert.Equal(t, tc.expected, val.Array)
	}
}

func TestReadBulk(t *testing.T) {
	tt := []struct {
		input    string
		expected Value
	}{
		{"5\r\nhello\r\n", Value{Typ: "bulk", Bulk: "hello"}},
		{"5\r\nworld\r\n", Value{Typ: "bulk", Bulk: "world"}},
	}

	for _, tc := range tt {
		reader := strings.NewReader(tc.input)
		respReader := NewResp(reader)
		val, err := respReader.readBulk()
		require.NoError(t, err)
		assert.Equal(t, tc.expected, val)
	}
}

func TestMarshal(t *testing.T) {
	tt := []struct {
		input    Value
		expected []byte
	}{
		{Value{Typ: "string", Str: "hello"}, []byte("+hello\r\n")},
		{Value{Typ: "bulk", Bulk: "world"}, []byte("$5\r\nworld\r\n")},
		{Value{Typ: "array", Array: []Value{{Typ: "string", Str: "foo"}, {Typ: "string", Str: "bar"}}}, []byte("*2\r\n+foo\r\n+bar\r\n")},
		{Value{Typ: "error", Str: "oops"}, []byte("-oops\r\n")},
		{Value{Typ: "null"}, []byte("$-1\r\n")},
	}

	for _, tc := range tt {
		val := tc.input
		result := val.Marshal()
		assert.Equal(t, tc.expected, result)
	}
}

func TestRead(t *testing.T) {
	tt := []struct {
		input    string
		expected Value
	}{
		{"$5\r\nworld\r\n", Value{Typ: "bulk", Bulk: "world"}},
		{"*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n", Value{Typ: "array", Array: []Value{
			{Typ: "bulk", Bulk: "foo"},
			{Typ: "bulk", Bulk: "bar"},
		}}},
	}

	for _, tc := range tt {
		reader := strings.NewReader(tc.input)
		respReader := NewResp(reader)
		val, err := respReader.Read()
		require.NoError(t, err)
		assert.Equal(t, tc.expected, val)
	}
}
