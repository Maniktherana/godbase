package handler

import (
	"strconv"
	"sync"
	"time"

	"github.com/maniktherana/godbase/pkg/resp"
)

var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

func ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Typ: "string", Str: "PONG"}
	}

	return resp.Value{Typ: "string", Str: args[0].Bulk}
}

var SETs = map[string]resp.Value{}
var SETsMu = sync.RWMutex{}

func set(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.Value{Typ: "error", Str: "ERR syntax error"}
	}

	key := args[0].Bulk
	value := args[1].Bulk
	setter := ""
	ex := 0
	px := 0

	if len(args) == 3 {
		setter = args[2].Bulk
		if setter != "NX" && setter != "XX" {
			return resp.Value{Typ: "error", Str: "ERR syntax error"}
		}
	}

	if len(args) > 3 {
		switch args[2].Bulk {
		case "EX":
			ex, _ = strconv.Atoi(args[3].Bulk)
		case "PX":
			px, _ = strconv.Atoi(args[3].Bulk)
		default:
			setter = args[2].Bulk
		}
		switch args[3].Bulk {
		case "EX":
			ex, _ = strconv.Atoi(args[4].Bulk)
		case "PX":
			px, _ = strconv.Atoi(args[4].Bulk)
		default:
			setter = args[3].Bulk
		}
	}

	switch setter {
	case "NX":
		SETsMu.RLock()
		_, ok := SETs[key]
		SETsMu.RUnlock()
		if ok {
			return resp.Value{Typ: "null"}
		}
	case "XX":
		SETsMu.RLock()
		_, ok := SETs[key]
		SETsMu.RUnlock()
		if !ok {
			return resp.Value{Typ: "null"}
		}
	default:
		if setter != "" {
			return resp.Value{Typ: "error", Str: "ERR syntax error"}
		}
	}

	expiration := time.Now()
	val := resp.Value{Typ: "bulk", Bulk: value}
	if ex > 0 {
		expiration = expiration.Add(time.Duration(ex) * time.Second)
	} else if px > 0 {
		expiration = expiration.Add(time.Duration(px) * time.Millisecond)
	}
	val.Expires = expiration.Unix()

	SETsMu.Lock()
	SETs[key] = val
	SETsMu.Unlock()
	return resp.Value{Typ: "string", Str: "OK"}
}

func get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	if value.Expires > 0 && value.Expires < time.Now().Unix() {
		delete(SETs, key)
		return resp.Value{Typ: "null"}
	}
	SETsMu.RUnlock()

	if !ok {
		return resp.Value{Typ: "null"}
	}

	return value
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []resp.Value) resp.Value {
	if len(args) != 3 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return resp.Value{Typ: "string", Str: "OK"}
}

func hget(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return resp.Value{Typ: "null"}
	}

	return resp.Value{Typ: "bulk", Bulk: value}
}

func hgetall(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].Bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash]
	HSETsMu.RUnlock()

	if !ok {
		return resp.Value{Typ: "null"}
	}

	values := []resp.Value{}
	for key, value := range value {
		values = append(values, resp.Value{Typ: "bulk", Bulk: key})
		values = append(values, resp.Value{Typ: "bulk", Bulk: value})
	}

	return resp.Value{Typ: "bulk", Array: values}
}
