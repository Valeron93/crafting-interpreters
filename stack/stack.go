package stack

import "fmt"

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(elem T) {
	s.items = append(s.items, elem)
}

func (s *Stack[T]) Count() int {
	return len(s.items)
}

func (s *Stack[T]) Empty() bool {
	return s.Count() == 0
}

func (s *Stack[T]) GetIdx(idx int) T {
	return s.items[idx]
}

func (s *Stack[T]) MustPeek() T {
	return s.items[s.Count()-1]
}

func (s *Stack[T]) MustPop() T {

	last := s.Count() - 1
	elem := s.items[last]
	s.items = s.items[:last]
	return elem
}

func (s *Stack[T]) Pop() (T, bool) {
	if s.Count() < 1 {
		var zero T
		return zero, false
	}
	return s.MustPop(), true
}

func (s *Stack[T]) Peek() (T, bool) {
	if s.Count() < 1 {
		var zero T
		return zero, false
	}
	return s.MustPeek(), true
}

func (s *Stack[T]) String() string {

	return fmt.Sprintf("Stack %v", s.items)

}
