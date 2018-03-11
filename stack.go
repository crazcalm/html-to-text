package htmlToText

import (
	"errors"
)

//Tag -- A struct for HTML tags
type Tag string

func (t Tag) String() string {
	return string(t)
}

const (
	//OpenPTag -- Open html p tag
	OpenPTag = Tag("<p>")
	//ClosePTag -- Closed html p tag
	ClosePTag = Tag("</p>")
)

//Stack -- A stack used when processing the html tags during the translation from html to text
type Stack []Tag

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
