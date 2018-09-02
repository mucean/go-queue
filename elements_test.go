package goqueue

import "testing"

var capLen = 100

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

func TestElements_Push_One(t *testing.T) {
	e := NewElements(capLen)
	vInt := 1
	if err := e.Push(vInt); err != nil {
		t.Errorf(err.Error())
	}
	
	if v, ok := e.values[0].(int); ok {
		if v != vInt {
			t.Errorf("the value in the elements by push method must be equal to the push value")
		}
	} else {
		t.Errorf("the value type in the elements by push method must be equal to the push value")
	}
	
	if e.len != 1 {
		t.Errorf("the length of elements is not equal to 1 when empty elements push one value")
	}
	
	if e.head != 0 {
		t.Errorf("the head value must equal to 0 when empty elements push one value")
	}
}

func TestElements_Push_Full(t *testing.T) {
	e := NewElements(capLen)
	for i := 0; i < capLen; i++ {
		if err := e.Push(i); err != nil {
			t.Errorf(err.Error())
		}
		if e.values[i].(int) != i {
			t.Errorf("not all value equal to the push one")
		}
	}
	
	if e.len != e.cap {
		t.Errorf("length of one full elements is not equal to it's capcity")
	}
	
	if e.head != 0 {
		t.Errorf("head index on full elements must equal to zero")
	}
	
	if err := e.Push("test"); err != fullErr {
		t.Error("full elements push value must return full error")
	}
}

func TestElements_PushForce(t *testing.T) {
	ee := emptyElements()
	fe := fullElements()
	if err := ee.PushForce(1); err != nil {
		t.Fatalf("PushForce method must return nil error when the elements empty")
	}
	
	if err := fe.PushForce(1); err != fullErr {
		t.Fatalf("PushForce method must return fullErr error when the elements full")
	}
	
	fe.Pop()
	if err := fe.PushForce(1); err != nil {
		t.Fatalf("PushForce method must return nil error, but error is: %s", err.Error())
	}

	// check PushForce is use Rebuild and Push method
	if err := fe.PushForce(1); err != fullErr {
		t.Fatalf("PushForce method must return fullErr error when the elements full")
	}
}

func TestElements_Pop(t *testing.T) {
	e := fullElements()
	
	length := e.len
	for i := 0; i < capLen; i++ {
		if v, e := e.Pop(); e != nil {
			t.Error(e)
		} else {
			if vi, ok := v.(int); !ok {
				t.Errorf("the value type from pop method is not the push one")
			} else if vi != i {
				t.Errorf("left value: %d, right value: %d", vi, i)
			}
		}
		length--
		if length != e.len {
			t.Errorf("pop method reduce the length of elements, the length must be %d, %d gives", length, e.len)
		}

		if i + 1 != e.head {
			t.Errorf("pop method add the head index, head must be %d, %d gives", i + 1, e.head)
		}
	}
}

func TestElements_pushable(t *testing.T) {
	e := emptyElements()
	for i := 0; i < capLen; i++ {
		if err := e.pushable(); err != nil {
			t.Errorf("push %d times failed: %s", i + 1, err.Error())
		}
		if !e.Pushable() {
			t.Fatalf("Pushable method return false when push %d times", i + 1)
		}
		e.Push(i)
	}
	
	if e.pushable() != fullErr {
		t.Errorf("pushable of the full elements must return `fullErr`")
	}
	
	if e.Pushable() {
		t.Fatalf("Pushable method return true when elements is full")
	}
	
	for i := 0; i < capLen; i++ {
		e.Pop()
		if err := e.pushable(); err != pushErr {
			t.Error(err)
		}
		
		if e.Pushable() {
			t.Fatalf("Pushable method must return true when index of tail in the end")
		}
	}

	if err := e.pushable(); err != pushErr {
		t.Error(err)
	}
	
	if e.Pushable() {
		t.Fatalf("Pushable method must return true when index of tail in the end")
	}
}

func TestElements_tail(t *testing.T) {
	e := emptyElements()
	for i := 0; i < capLen; i++ {
		e.Push(i)
		if e.tail() != i {
			t.Errorf("tail must equal to %d, %d gives", i, e.tail())
		}
	}
	if e.tail() != capLen - 1 {
		t.Errorf("tail must equal to %d, %d gives", capLen - 1, e.tail())
	}
	
	for i := 0; i < capLen; i++ {
		e.Pop()
		if e.tail() != capLen - 1 {
			t.Errorf("pop method can not change tail index")
		}
	}
}

func emptyElements() *Elements {
	return NewElements(capLen)
}

func fullElements() *Elements {
	e := emptyElements()
	for i := 0; i < capLen; i++ {
		e.Push(i)
	}
	
	return e
}

func tailEndNotFullElements(popTime int) *Elements {
	e := fullElements()
	for i := 0; i < popTime; i++ {
		e.Pop()
	}
	return e
}