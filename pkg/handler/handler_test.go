package handler

import (
	"testing"

	"github.com/maniktherana/godbase/pkg/resp"
	"github.com/stretchr/testify/assert"
)

func TestPingHandler(t *testing.T) {
	tt := []struct {
		name     string
		args     []resp.Value
		expected resp.Value
	}{
		{
			name:     "NoArgs",
			args:     []resp.Value{},
			expected: resp.Value{Typ: "string", Str: "PONG"},
		},
		{
			name:     "WithArgs",
			args:     []resp.Value{{Typ: "bulk", Bulk: "hello"}},
			expected: resp.Value{Typ: "string", Str: "hello"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := ping(tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSetHandler(t *testing.T) {
	tt := []struct {
		name     string
		args     []resp.Value
		expected resp.Value
	}{
		{
			name:     "NoArgs",
			args:     []resp.Value{{Typ: "bulk", Bulk: "key1"}, {Typ: "bulk", Bulk: "value1"}},
			expected: resp.Value{Typ: "string", Str: "OK"},
		},
		{
			name:     "EX",
			args:     []resp.Value{{Typ: "bulk", Bulk: "key2"}, {Typ: "bulk", Bulk: "value2"}, {Typ: "bulk", Bulk: "EX"}, {Typ: "bulk", Bulk: "10"}},
			expected: resp.Value{Typ: "string", Str: "OK"},
		},
		{
			name:     "PX",
			args:     []resp.Value{{Typ: "bulk", Bulk: "key3"}, {Typ: "bulk", Bulk: "value3"}, {Typ: "bulk", Bulk: "PX"}, {Typ: "bulk", Bulk: "10000"}},
			expected: resp.Value{Typ: "string", Str: "OK"},
		},
		{
			name:     "NX",
			args:     []resp.Value{{Typ: "bulk", Bulk: "key4"}, {Typ: "bulk", Bulk: "value4"}, {Typ: "bulk", Bulk: "NX"}},
			expected: resp.Value{Typ: "string", Str: "OK"},
		},
		{
			name:     "XX Error",
			args:     []resp.Value{{Typ: "bulk", Bulk: "key5"}, {Typ: "bulk", Bulk: "value5"}, {Typ: "bulk", Bulk: "XX"}},
			expected: resp.Value{Typ: "null", Str: "", Num: 0},
		},
		{
			name:     "GET",
			args:     []resp.Value{{Typ: "bulk", Bulk: "key6"}, {Typ: "bulk", Bulk: "value6"}, {Typ: "bulk", Bulk: "GET"}},
			expected: resp.Value{Typ: "string", Str: "value6"},
		},
		{
			name:     "KEEPTTL error",
			args:     []resp.Value{{Typ: "bulk", Bulk: "key7"}, {Typ: "bulk", Bulk: "value7"}, {Typ: "bulk", Bulk: "KEEPTTL"}},
			expected: resp.Value{Typ: "null", Str: "", Num: 0},
		},
		{
			name:     "Invalid Option",
			args:     []resp.Value{{Typ: "bulk", Bulk: "key8"}, {Typ: "bulk", Bulk: "value8"}, {Typ: "bulk", Bulk: "INVALID"}},
			expected: resp.Value{Typ: "error", Str: "ERR syntax error"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := set(tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetHandler(t *testing.T) {
	// Set up initial value
	SETsMu.Lock()
	SETs["mykey"] = resp.Value{Typ: "string", Str: "myvalue"}
	SETsMu.Unlock()

	tests := []struct {
		name     string
		args     []resp.Value
		expected resp.Value
	}{
		{
			name:     "Existing Key",
			args:     []resp.Value{{Typ: "bulk", Bulk: "mykey"}},
			expected: resp.Value{Typ: "string", Str: "myvalue"},
		},
		{
			name:     "Non Existing Key",
			args:     []resp.Value{{Typ: "bulk", Bulk: "nonexistingkey"}},
			expected: resp.Value{Typ: "null"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := get(test.args)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestHsetHandler(tt *testing.T) {
	tests := []struct {
		name     string
		args     []resp.Value
		expected resp.Value
	}{
		{
			name:     "Normal",
			args:     []resp.Value{{Typ: "bulk", Bulk: "hash"}, {Typ: "bulk", Bulk: "key"}, {Typ: "bulk", Bulk: "value"}},
			expected: resp.Value{Typ: "string", Str: "OK"},
		},
		{
			name:     "WrongNumberOfArguments",
			args:     []resp.Value{{Typ: "bulk", Bulk: "hash"}, {Typ: "bulk", Bulk: "key"}},
			expected: resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hset' command"},
		},
	}

	for _, tc := range tests {
		tt.Run(tc.name, func(t *testing.T) {
			result := hset(tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestHgetHandler(tt *testing.T) {
	tests := []struct {
		name     string
		args     []resp.Value
		setup    func()
		expected resp.Value
	}{
		{
			name: "ExistingKey",
			args: []resp.Value{{Typ: "bulk", Bulk: "hash"}, {Typ: "bulk", Bulk: "key"}},
			setup: func() {
				// Set up the initial key-value pair
				HSETsMu.Lock()
				HSETs["hash"] = map[string]string{"key": "value"}
				HSETsMu.Unlock()
			},
			expected: resp.Value{Typ: "bulk", Bulk: "value"},
		},
		{
			name:     "NonExistingKey",
			args:     []resp.Value{{Typ: "bulk", Bulk: "hash"}, {Typ: "bulk", Bulk: "nonexistent"}},
			expected: resp.Value{Typ: "null"},
		},
		{
			name:     "WrongNumberOfArguments",
			args:     []resp.Value{{Typ: "bulk", Bulk: "hash"}},
			expected: resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hget' command"},
		},
	}

	for _, tc := range tests {
		tt.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			result := hget(tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestHgetallHandler(tt *testing.T) {
	tests := []struct {
		name     string
		args     []resp.Value
		setup    func()
		expected resp.Value
	}{
		{
			name: "ExistingHash",
			args: []resp.Value{{Typ: "bulk", Bulk: "hash"}},
			setup: func() {
				// Set up the initial key-value pairs
				HSETsMu.Lock()
				HSETs["hash"] = map[string]string{"key1": "value1", "key2": "value2"}
				HSETsMu.Unlock()
			},
			expected: resp.Value{Typ: "bulk", Array: []resp.Value{
				{Typ: "bulk", Bulk: "key1"},
				{Typ: "bulk", Bulk: "value1"},
				{Typ: "bulk", Bulk: "key2"},
				{Typ: "bulk", Bulk: "value2"},
			}},
		},
		{
			name:     "NonExistingHash",
			args:     []resp.Value{{Typ: "bulk", Bulk: "nonexistent"}},
			expected: resp.Value{Typ: "null"},
		},
		{
			name:     "WrongNumberOfArguments",
			args:     []resp.Value{{Typ: "bulk", Bulk: "hash"}, {Typ: "bulk", Bulk: "this"}},
			expected: resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hget' command"},
		},
	}

	for _, tc := range tests {
		tt.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			result := hgetall(tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}
