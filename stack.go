package htmltotext

import (
	"errors"
	"strings"
)

//Stack -- A stack used when processing the html tags during the translation from html to text
type Stack []Tag

//Contains -- Checks to see if a Tag is in the Stack
func (s Stack) Contains(tags ...Tag) bool {
	result := false
	for _, t := range tags {
		for _, tag := range s {
			if strings.EqualFold(tag.String(), t.String()) {
				result = true
				break
			}
		}
	}
	return result
}

//OnTop -- Checks if the wanted Tag is on Top
func (s Stack) OnTop(t Tag) bool {
	result := false
	if len(s) > 0 {
		if strings.EqualFold(t.String(), s[len(s)-1].String()) {
			result = true
		}
	}
	return result
}

//Push -- Pushes a Tag onto the stack
func (s Stack) Push(t Tag) Stack {
	s = append(s, t)
	return s
}

//Pop -- Pops a Tag off of the stack
func (s Stack) Pop() (Stack, Tag, error) {
	if len(s) == 0 {
		return s, Tag(""), errors.New("Stack is Empty")
	}
	tag := s[len(s)-1]
	s = s[:len(s)-1]
	return s, tag, nil
}

//Peek -- Peeks at the Tag on top of the stack
func (s Stack) Peek() (t Tag, err error) {
	if len(s) == 0 {
		return t, errors.New("Stack is Empty")
	}
	return s[len(s)-1], nil
}

//NewStack -- Returns an empty Stack
func NewStack() (s Stack) {
	return
}
