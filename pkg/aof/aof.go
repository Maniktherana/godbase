package aof

import (
	"bufio"
	"os"
	"sync"
	"time"
)

type Aof struct {
	File *os.File
	Rd   *bufio.Reader
	Mu   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		File: f,
		Rd:   bufio.NewReader(f),
	}

	// Start a goroutine to sync AOF to disk every 1 second
	go func() {
		for {
			aof.Mu.Lock()

			aof.File.Sync()

			aof.Mu.Unlock()

			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func (aof *Aof) Close() error {
	aof.Mu.Lock()
	defer aof.Mu.Unlock()

	return aof.File.Close()
}
