package main

import (
	"strings"
	"sync"
)

type LogBuilders struct {
	mu    sync.RWMutex
	items map[int]*strings.Builder
}

func NewLogBuilders(count int) *LogBuilders {
	logBuilders := &LogBuilders{
		items: make(map[int]*strings.Builder),
	}

	logBuilders.mu.Lock()
	for i := 0; i < count; i++ {
		logBuilders.items[i] = &strings.Builder{}
	}
	logBuilders.mu.Unlock()

	return logBuilders
}

func (l *LogBuilders) Get(index int) *strings.Builder {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.items[index]
}
