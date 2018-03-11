package htmlToText

import (
	"strings"
	"testing"
)

func TestNewStack(t *testing.T) {
	stack := NewStack()

	if len(stack) != 0 {
		t.Errorf("NewStack return a non-empty stack...")
	}
}

func TestPush(t *testing.T) {
	tests := []struct {
		Stack    Stack
		NumOfElm int
	}{
		{Stack{}, 1},
		{Stack{OpenPTag}, 2},
		{Stack{OpenPTag, OpenPTag}, 3},
	}

	for _, test := range tests {
		s := test.Stack.Push(OpenPTag)

		if len(s) != test.NumOfElm {
			t.Errorf("Stack.Push: Expected %d elements, but got %d", test.NumOfElm, len(test.Stack))
		}
	}
}

func TestPop(t *testing.T) {
	tests := []struct {
		Stack    Stack
		NumOfElm int
		Error    bool
	}{
		{Stack{}, 0, true},
		{Stack{OpenPTag}, 0, false},
		{Stack{OpenPTag, OpenPTag}, 1, false},
	}

	for _, test := range tests {
		s, _, err := test.Stack.Pop()

		if test.Error && err != nil {
			continue
		} else if test.Error && err == nil {
			t.Errorf("expected error, but received none")
		} else if !test.Error && err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}

		if len(s) != test.NumOfElm {
			t.Errorf("Expected %d of elements, but got %d", test.NumOfElm, len(s))
		}
	}

}

func TestPeek(t *testing.T) {
	tests := []struct {
		Stack Stack
		Tag   Tag
		Error bool
	}{
		{Stack{}, OpenPTag, true},
		{Stack{OpenPTag}, OpenPTag, false},
	}

	for _, test := range tests {
		tag, err := test.Stack.Peek()

		if test.Error && err != nil {
			continue
		} else if test.Error && err == nil {
			t.Errorf("Expected an error, but received none")
		} else if !test.Error && err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}

		if !strings.EqualFold(tag.String(), test.Tag.String()) {
			t.Errorf("Expected %s tag, but got %s instead", test.Tag, tag)
		}
	}
}
