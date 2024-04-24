package handler

import (
	"github.com/maniktherana/godbase/pkg/Database"
	"strconv"
	"strings"
	"time"

	"github.com/maniktherana/godbase/pkg/resp"
)

var Handlers = map[string]func(
	[]resp.Value,
	*Database.Kv,
) resp.Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

func ping(args []resp.Value, kv *Database.Kv) resp.Value {
	if len(args) == 0 {
		return resp.Value{Typ: "string", Str: "PONG"}
	}

	return resp.Value{Typ: "string", Str: args[0].Bulk}
}

func set(args []resp.Value, kv *Database.Kv) resp.Value {
	if len(args) < 2 {
		return resp.Value{Typ: "error", Str: "ERR syntax error"}
	}

	key := args[0].Bulk
	value := args[1].Bulk
	var setter string
	keepttl := false
	var get bool
	var ex, px int

	// Parsing command options
	for i := 2; i < len(args); i++ {
		switch strings.ToUpper(args[i].Bulk) {
		case "NX", "XX":
			setter = strings.ToUpper(args[i].Bulk)
		case "KEEPTTL":
			keepttl = true
		case "GET":
			get = true
		case "EX":
			if keepttl {
				return resp.Value{Typ: "error", Str: "ERR syntax error"}
			}
			if i+1 < len(args) {
				ex, _ = strconv.Atoi(args[i+1].Bulk)
				i++
			} else {
				return resp.Value{Typ: "error", Str: "ERR syntax error"}
			}
		case "PX":
			if keepttl {
				return resp.Value{Typ: "error", Str: "ERR syntax error"}
			}
			if i+1 < len(args) {
				px, _ = strconv.Atoi(args[i+1].Bulk)
				i++
			} else {
				return resp.Value{Typ: "error", Str: "ERR syntax error"}
			}
		default:
			return resp.Value{Typ: "error", Str: "ERR syntax error"}
		}
	}

	// Handling SETTER (XX/NX) options
	switch setter {
	case "NX":
		kv.SETsMu.RLock()
		_, ok := kv.SETs[key]
		kv.SETsMu.RUnlock()
		if ok {
			return resp.Value{Typ: "null"}
		}
	case "XX":
		kv.SETsMu.RLock()
		_, ok := kv.SETs[key]
		kv.SETsMu.RUnlock()
		if !ok {
			return resp.Value{Typ: "null"}
		}
	}

	// Handling expiration
	expiration := time.Now()
	if keepttl {
		kv.SETsMu.RLock()
		v, ok := kv.SETs[key]
		kv.SETsMu.RUnlock()
		if !ok {
			return resp.Value{Typ: "null"}
		}
		expiration = time.Unix(v.Expires, 0)
	} else {
		if ex > 0 {
			expiration = expiration.Add(time.Duration(ex) * time.Second)
		} else if px > 0 {
			expiration = expiration.Add(time.Duration(px) * time.Millisecond)
		}
	}

	// Setting the value
	var val resp.Value
	if expiration.UnixMilli() > time.Now().UnixMilli() {
		val = resp.Value{Typ: "string", Str: value, Expires: expiration.UnixMilli()}
	} else {
		val = resp.Value{Typ: "string", Str: value}
	}
	kv.SETsMu.Lock()
	kv.SETs[key] = val
	kv.SETsMu.Unlock()

	if get {
		return val
	} else {
		return resp.Value{Typ: "string", Str: "OK"}
	}
}

func get(args []resp.Value, kv *Database.Kv) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk

	kv.SETsMu.RLock()
	value, ok := kv.SETs[key]
	kv.SETsMu.RUnlock()

	if !ok {
		return resp.Value{Typ: "null"}
	}

	if value.Expires > 0 && value.Expires < time.Now().UnixMilli() {
		kv.SETsMu.Lock()
		delete(kv.SETs, key)
		kv.SETsMu.Unlock()
		return resp.Value{Typ: "null"}
	}

	return value
}

func hset(args []resp.Value, kv *Database.Kv) resp.Value {
	if len(args) != 3 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	kv.HSETsMu.Lock()
	if _, ok := kv.HSETs[hash]; !ok {
		kv.HSETs[hash] = map[string]string{}
	}
	kv.HSETs[hash][key] = value
	kv.HSETsMu.Unlock()

	return resp.Value{Typ: "string", Str: "OK"}
}

func hget(args []resp.Value, kv *Database.Kv) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	kv.HSETsMu.RLock()
	value, ok := kv.HSETs[hash][key]
	kv.HSETsMu.RUnlock()

	if !ok {
		return resp.Value{Typ: "null"}
	}

	return resp.Value{Typ: "bulk", Bulk: value}
}

func hgetall(args []resp.Value, kv *Database.Kv) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].Bulk

	kv.HSETsMu.RLock()
	value, ok := kv.HSETs[hash]
	kv.HSETsMu.RUnlock()

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
