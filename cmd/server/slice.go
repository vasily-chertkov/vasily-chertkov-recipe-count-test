package main

import (
	"sort"
	"sync"
)

func NewStringSlice(n ...string) *Slice { return NewSlice(sort.StringSlice(n)) }

func NewSlice(n sort.StringSlice) *Slice {
	s := &Slice{StringSlice: n, pos: make(map[string][2]int, n.Len())}
	sort.Sort(s)
	s.calculate()
	return s
}

type Slice struct {
	sort.StringSlice
	sync.RWMutex
	pos   map[string][2]int
	bytes []byte
}

func (s *Slice) calculate() {
	pos := 0 // [ in the beginning
	last := len(s.StringSlice) - 1
	for i, c := range s.StringSlice {
		if i < last {
			pos += 1 // double quote
		}
		start := pos
		// if the string is the only - move start by one
		if i == 0 && i == last {
			start += 1
		}

		// string length, double quote, comma or last square bracket
		pos += len(c) + 1 + 1
		end := pos + 1 // points to the double quote of the next string

		s.Lock()
		s.pos[c] = [2]int{start, end}
		s.Unlock()
	}
	b, _ := jsonfast.Marshal(s.StringSlice)
	s.bytes = b
}

func (s *Slice) Bytes(id string) []byte {
	s.RLock()
	pos := s.pos[id]
	s.RUnlock()
	start, end := pos[0], pos[1]
	b := make([]byte, len(s.bytes)-(end-start))
	copy(b[:start], s.bytes[:start])
	copy(b[start:], s.bytes[end:])
	return b
}
