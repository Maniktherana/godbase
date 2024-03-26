package handler

import (
	"fmt"
	"strconv"
	"strings"
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

	// Handling SETTER options
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
	}

	// Handling expiration
	expiration := time.Now()
	if keepttl {
		SETsMu.RLock()
		v, ok := SETs[key]
		SETsMu.RUnlock()
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
	if expiration.Unix() > 0 {
		val = resp.Value{Typ: "string", Str: value, Expires: expiration.Unix()}
	} else {
		val = resp.Value{Typ: "string", Str: value}
	}
	SETsMu.Lock()
	SETs[key] = val
	SETsMu.Unlock()

	if get {
		return val
	} else {
		return resp.Value{Typ: "string", Str: "OK"}
	}
}

func get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	if value.Expires > 0 && value.Expires > time.Now().Unix() {
		delete(SETs, key)
		return resp.Value{Typ: "null"}
	}
	SETsMu.RUnlock()

	if !ok {
		return resp.Value{Typ: "null"}
	}
	fmt.Println("getting value ", value)

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
