package goqueue

import "testing"

var capLen int = 100

func TestNewElements(t *testing.T) {
	e := NewElements(capLen)
	if e.head != 0 {
		t.Error("the value of head in the new elements must equal to 0")
	}
	
	if e.Len() != 0 {
		t.Errorf("length of new elements must be 0")
	}
	
	if e.Cap() != capLen {
		t.Error("cap is not equal to argument")
	}
	
	if e.Cap() != len(e.values) {
		t.Error("values length is not equal to cap")
	}
	
	for _, v := range e.values {
		if v != nil {
			t.Error("each value of new elements must be nil")
		}
	}
}

func TestElements_Push(t *testing.T) {
	e := NewElements(capLen)
	vInt := 1
	if e.Push(vInt) != nil {
		t.Errorf("empty elements pust value must has no error")
	}
	
	if v, ok := e.values[0].(int); ok {
		if v != vInt {
			t.Errorf("the value in the elements by push method must be equal to the push value")
		}
	} else {
		t.Errorf("the value type in the elements by push method must be equal to the push value")
	}
}