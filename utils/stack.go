package utils

// create class Stack that inherits from the built-in type Array
type Stack struct {
	Elements [][4]uint32
}

// Add a new element to the end of the array
func (s *Stack) Append(element [4]uint32) {
	s.Elements = append(s.Elements, element)
}

// Return the last element in the array
func (s *Stack) Peek() [4]uint32 {
	return s.Elements[len(s.Elements)-1]
}

// Remove the last element from the array
func (s *Stack) Pop() [4]uint32 {
	lastIndex := len(s.Elements) - 1
	lastElement := s.Elements[lastIndex]
	s.Elements = s.Elements[:lastIndex]
	return lastElement
}

// Return the number of elements in the array
func (s *Stack) Length() int {
	return len(s.Elements)
}

