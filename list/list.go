package list


// create class List that inherits from the built-in type Array
type List struct {
    Elements []string
}

// Add a new element to the end of the array
func (l *List) Append(element string) {
    l.Elements = append(l.Elements, element)
}

// Add a new element to the beginning of the array
func (l *List) Prepend(element string) {
    l.Elements = append([]string{element}, l.Elements...)
}

// Add a new element at the specified index
func (l *List) Insert(index int, element string) {
    l.Elements = append(l.Elements[:index], append([]string{element}, l.Elements[index:]...)...)
}

// Remove the first element from the array
func (l *List) Shift() string {
    firstElement := l.Elements[0]
    l.Elements = l.Elements[1:]
    return firstElement
}

// Remove the last element from the array
func (l *List) Pop() string {
	lastIndex := len(l.Elements) - 1
	lastElement := l.Elements[lastIndex]
	l.Elements = l.Elements[:lastIndex]
	return lastElement
}

// Return the number of elements in the array
func (l *List) Length() int {
    return len(l.Elements)
}

// Return the element at the specified index
func (l *List) Get(index int) string {
    return l.Elements[index]
}

// Return the index of the first occurrence of the specified element
func (l *List) IndexOf(element string) int {
	for i, e := range l.Elements {
		if e == element {
			return i
		}
	}
	return -1
}


