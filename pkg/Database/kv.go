package Database

import (
	"github.com/maniktherana/godbase/pkg/resp"
	"net"
	"sync"
)

type Kv struct {
	SETs                 map[string]resp.Value
	SETsMu               sync.RWMutex
	HSETs                map[string]map[string]string
	HSETsMu              sync.RWMutex
	NumCommandsProcessed int
	Clients              map[string]net.Conn
}

func NewKv() *Kv {
	return &Kv{
		SETs:    map[string]resp.Value{},
		HSETs:   map[string]map[string]string{},
		Clients: map[string]net.Conn{},
	}
}
