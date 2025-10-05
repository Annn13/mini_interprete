package vm

import "fmt"

type Stack struct {
	data []any
}

// funiones basicas
func NewStack() *Stack { return &Stack{data: make([]any, 0, 32)} }

func (s *Stack) Push(v any) { s.data = append(s.data, v) }

func (s *Stack) Pop() (any, bool) {
	if len(s.data) == 0 {
		return nil, false
	}
	last := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return last, true
}

func (s *Stack) String() string {
	return fmt.Sprint(s.data)
}

func (s *Stack) Peek() (any, bool) {
	if len(s.data) == 0 {
		return nil, false
	}
	return s.data[len(s.data)-1], true
}

func (s *Stack) Len() int { return len(s.data) }
