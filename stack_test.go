package htmltotext

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

func TestContainsMultiple(t *testing.T) {
	stack := Stack{OpenPTag, OpenH1Tag, OpenTableTag, OpenOLTag}

	if stack.Contains(OpenH2Tag, OpenH3Tag, OpenH6Tag) {
		t.Errorf("TestContainsMultiple Case A failed")
	}

	if !stack.Contains(OpenH2Tag, OpenH1Tag, OpenH6Tag) {
		t.Errorf("TestContainsMultiple Case B failed")
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		Stack  Stack
		Tag    Tag
		Answer bool
	}{
		{Stack{}, OpenPTag, false},
		{Stack{OpenPTag}, OpenPTag, true},
		{Stack{OpenH1Tag, OpenH2Tag, OpenH3Tag}, OpenPTag, false},
		{Stack{OpenH1Tag, OpenH2Tag, OpenH3Tag, OpenPTag}, OpenPTag, true},
	}

	for _, test := range tests {
		if test.Stack.Contains(test.Tag) != test.Answer {
			t.Errorf("stack.Contains(%s): Expected %t, but got %t", test.Tag.String(), test.Answer, test.Stack.Contains(test.Tag))
		}
	}
}

func TestOnTop(t *testing.T) {
	tests := []struct {
		Stack  Stack
		Tag    Tag
		Answer bool
	}{
		{Stack{}, OpenPTag, false},
		{Stack{OpenPTag}, OpenPTag, true},
		{Stack{OpenH1Tag, OpenPTag, OpenH3Tag}, OpenPTag, false},
		{Stack{OpenH1Tag, OpenPTag, OpenH3Tag, OpenPTag}, OpenPTag, true},
	}

	for _, test := range tests {
		if test.Stack.OnTop(test.Tag) != test.Answer {
			t.Errorf("stack.OnTop(%s): Expected %t, but got %t\n%v", test.Tag.String(), test.Answer, test.Stack.Contains(test.Tag), test.Stack)
		}
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
