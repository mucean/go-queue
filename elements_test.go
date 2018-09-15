package goqueue

import (
	"fmt"
	"testing"
)

var capLen = 100

func TestNewElements(t *testing.T) {
	e := NewElements(capLen)
	if e.head != 0 {
		t.Error("the value of head in the new elements must equal to 0")
	}

	if e.Len() != 0 {
		t.Fatalf("length of new elements must be 0")
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
		t.Fatalf(err.Error())
	}

	if v, ok := e.values[0].(int); ok {
		if v != vInt {
			t.Fatalf("the value in the elements by push method must be equal to the push value")
		}
	} else {
		t.Fatalf("the value type in the elements by push method must be equal to the push value")
	}

	if e.len != 1 {
		t.Fatalf("the length of elements is not equal to 1 when empty elements push one value")
	}

	if e.head != 0 {
		t.Fatalf("the head value must equal to 0 when empty elements push one value")
	}
}

func TestElements_Push_Full(t *testing.T) {
	e := NewElements(capLen)
	for i := 0; i < capLen; i++ {
		if err := e.Push(i); err != nil {
			t.Fatalf(err.Error())
		}
		if e.values[i].(int) != i {
			t.Fatalf("not all value equal to the push one")
		}
	}

	if e.len != e.cap {
		t.Fatalf("length of one full elements is not equal to it's capcity")
	}

	if e.head != 0 {
		t.Fatalf("head index on full elements must equal to zero")
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
	for i := capLen - 1; i >= 0; i-- {
		if v, e := e.Pop(); e != nil {
			t.Error(e)
		} else {
			if vi, ok := v.(int); !ok {
				t.Fatalf("the value type from pop method is not the push one")
			} else if vi != i {
				t.Fatalf("left value: %d, right value: %d", vi, i)
			}
		}
		length--
		if length != e.len {
			t.Fatalf("pop method reduce the length of elements, the length must be %d, %d gives", length, e.len)
		}

		if 0 != e.head {
			t.Fatalf("pop method add the head index, head must be %d, %d gives", 0, e.head)
		}
	}
}

func TestElements_Find(t *testing.T) {
	findTests := []struct {
		e   *Elements
		i   int
		err error
		f   func(v interface{}) bool
	}{
		{emptyElements(), 0, emptyErr, func(v interface{}) bool {
			return true
		}},
		{fullElements(), 0, nil, func(v interface{}) bool {
			return v.(int) == 0
		}},
		{fullElements(), capLen - 1, nil, func(v interface{}) bool {
			return v.(int) == capLen-1
		}},
		{fullElements(), 0, notFoundErr, func(v interface{}) bool {
			return v.(int) == capLen
		}},
		{fullElements(), 2, nil, func(v interface{}) bool {
			return v.(int) < 3
		}},
	}

	for j, test := range findTests {
		if i, err := test.e.Find(test.f); i != test.i || err != test.err {
			t.Fatalf("index %d, Find method must return index: %d, err: %v, but index result: %d, err result: %v", j, test.i, test.err, i, err)
		}
	}
}

func TestElements_FindAll(t *testing.T) {
	rem := make([]int, 0, 50)
	for i := capLen - 1; i >= 0; i-- {
		if i%2 == 0 {
			rem = append(rem, i)
		}
	}
	findTests := []struct {
		e   *Elements
		ins []int
		err error
		f   func(v interface{}) bool
	}{
		{emptyElements(), nil, emptyErr, func(v interface{}) bool {
			return true
		}},
		{fullElements(), nil, notFoundErr, func(v interface{}) bool {
			return v.(int) < 0
		}},
		{fullElements(), []int{2, 1, 0}, nil, func(v interface{}) bool {
			return v.(int) < 3
		}},
		{fullElements(), rem, nil, func(v interface{}) bool {
			return v.(int)%2 == 0
		}},
	}

	for j, test := range findTests {
		if ins, err := test.e.FindAll(test.f); err != test.err {
			t.Fatalf("index %d, FindAll method must return err: %s, but err result: %s", j, test.err, err)
		} else {
			if len(ins) != len(test.ins) {
				t.Fatalf("index %d, FindAll method must return index: %v, but resutl index is %v", j, test.ins, ins)
			}
			for key, i := range ins {
				if i != test.ins[key] {
					t.Fatalf("index %d, FindAll method must return index: %v, but resutl index is %v", j, test.ins, ins)
				}
			}
		}
	}
}

func TestElements_Rebuild(t *testing.T) {
	tests := []struct {
		e   *Elements
		oe  *Elements
		err error
		len int
	}{
		{emptyElements(), emptyElements(), nil, 0},
		{popTimes(5, emptyElementsPush(5)), popTimes(5, emptyElementsPush(5)), nil, 0},
		{emptyElementsPush(5), emptyElementsPush(5), nil, 5},
		{popTimes(2, emptyElementsPush(5)), popTimes(2, emptyElementsPush(5)), nil, 3},
	}
	for i, test := range tests {
		if err := test.e.Rebuild(); err != test.err {
			t.Fatalf("Rebuild method must return `%s`, but now return `%s`", test.err, err)
		}
		if test.e.head != 0 {
			t.Fatalf("the head must equal to zero when Rebuild method called")
		}
		if test.e.Len() != test.len {
			t.Fatalf("the tail method must return %d when Rebuild method called, not tail method return %d", test.len, test.e.Len())
		}

		if test.len > 0 {
			elementsEqual(test.e, test.oe, i, t)
		}
	}
}

func TestElements_MoveHead(t *testing.T) {
	e1 := fullElements()
	e1.head = 2
	tests := []struct {
		e     *Elements
		head  int
		count int
	}{
		{emptyElements(), -1, 0},
		{emptyElements(), 0, 0},
		{emptyElements(), 2, 0},
		{e1, 1, 0},
		{popTimes(2, emptyElementsPush(2)), 2, 0},
		{emptyElementsPush(1), 0, 1},
		{fullElements(), capLen - 1, capLen},
		{fullElements(), capLen, capLen},
		{emptyElementsPush(2), 0, 1},
	}

	for i, test := range tests {
		if count := test.e.MoveHead(test.head); count != test.count {
			t.Fatalf("index %d, MoveHead method must return %d, now return %d", i, test.head, count)
		}
	}
}

func TestElements_MoveTail(t *testing.T) {
	tests := []struct {
		e     *Elements
		t     int
		count int
	}{
		{emptyElements(), 0, 0},
		{emptyElementsPush(1), 2, 0},
		{fullElements(), capLen - 1, 1},
		{fullElements(), capLen - 49, 49},
	}

	for i, test := range tests {
		if count := test.e.MoveTail(test.t); count != test.count {
			t.Fatalf("index %d, MoveTail method must return %d, now return %d", i, test.t, count)
		}
	}
}

func TestElements_eraseByIndex(t *testing.T) {
	e1 := popTimes(3, emptyElementsPush(5))
	e1.Push(6)
	e1.Push(8)
	tests := []struct {
		e   *Elements
		ins []int
		c   int
		l   int
		ae  *Elements
	}{
		{emptyElements(), nil, 0, 0, nil},
		{emptyElements(), []int{1}, 0, 0, emptyElements()},
		{emptyElementsPush(10), []int{0, 1, 2}, 3, 7, moveHead(emptyElementsPush(10), 2)},
		{emptyElementsPush(10), []int{7, 8, 9}, 3, 7, emptyElementsPush(7)},
		{emptyElementsPush(6), []int{0, 1, 2, 3, 4, 5}, 6, 0, emptyElements()},
		{moveHead(emptyElementsPush(10), 1), []int{0, 20, 1, 2, 9, 7, -1, 5}, 4, 4, e1},
		{fullElements(), []int{100, 101, 99, 0}, 2, 98, moveHead(emptyElementsPush(99), 0)},
	}

	for i, test := range tests {
		if c := test.e.eraseByIndex(test.ins); c != test.c {
			t.Fatalf("index %d, eraseByIndex method must return %d, now return %d", i, test.c, c)
		}

		if test.e.Len() != test.l {
			t.Fatalf("index %d, the length of element must equal to %d after eraseByIndex method, now is %d", i, test.l, test.e.Len())
		}

		if test.ae != nil && test.ae.Len() > 0 {
			elementsEqual(test.e, test.ae, i, t)
		}
	}
}

func TestElements_Pushable(t *testing.T) {
	tests := []struct {
		e *Elements
		b bool
	}{
		{emptyElements(), true},
		{fullElements(), false},
		{popTimes(capLen, fullElements()), true},
		{moveHead(fullElements(), capLen-1), false},
		{moveHead(fullElements(), capLen), false},
	}

	for i, test := range tests {
		if b := test.e.Pushable(); b != test.b {
			fmt.Println(test.e)
			t.Fatalf("index %d, Pushable method must return bool %v, now bool %v gives", i, test.b, b)
		}
	}
}

func TestElements_pushable(t *testing.T) {
	e := emptyElements()
	for i := 0; i < capLen; i++ {
		if err := e.pushable(); err != nil {
			t.Fatalf("push %d times failed: %s", i+1, err.Error())
		}
		if !e.Pushable() {
			t.Fatalf("Pushable method return false when push %d times", i+1)
		}
		e.Push(i)
	}

	if e.pushable() != fullErr {
		t.Fatalf("pushable of the full elements must return `fullErr`")
	}

	if e.Pushable() {
		t.Fatalf("Pushable method return true when elements is full")
	}

	for i := 0; i < capLen+1; i++ {
		e.Pop()
		if err := e.pushable(); err != nil {
			t.Error(err)
		}

		if !e.Pushable() {
			t.Fatalf("Pushable method must return false")
		}
	}
}

func TestElements_Tail(t *testing.T) {
	tests := []struct {
		e   *Elements
		t   int
		err error
	}{
		{emptyElements(), 0, emptyErr},
		{moveHead(emptyElementsPush(2), 1), 0, emptyErr},
		{emptyElementsPush(2), 1, nil},
		{fullElements(), capLen - 1, nil},
		{moveHead(fullElements(), 4), capLen - 1, nil},
	}

	for i, test := range tests {
		if tail, err := test.e.Tail(); tail != test.t {
			t.Fatalf("index %d, Tail method must return %d, %d gives", i, test.t, tail)
		} else if err != test.err {
			t.Fatalf("index %d, Tail method must return `%s` error, now `%s` error gives", i, test.err, err)
		}
	}
}

func TestElements_tail(t *testing.T) {
	e := emptyElements()
	for i := 0; i < capLen; i++ {
		e.Push(i)
		if e.tail() != i {
			t.Fatalf("tail must equal to %d, %d gives", i, e.tail())
		}
	}
	if e.tail() != capLen-1 {
		t.Fatalf("tail must equal to %d, %d gives", capLen-1, e.tail())
	}

	for i := 0; i < capLen; i++ {
		if e.tail() != capLen-(i+1) {
			t.Fatalf("left is %d, right %d", e.tail(), capLen-(i+1))
		}
		e.Pop()
	}
}

func TestElements_Len(t *testing.T) {
	tests := []struct {
		e *Elements
		l int
	}{
		{emptyElements(), 0},
		{emptyElementsPush(2), 2},
		{popTimes(2, emptyElementsPush(5)), 5 - 2},
		{fullElements(), capLen},
		{popTimes(50, fullElements()), capLen - 50},
	}

	for _, test := range tests {
		if l := test.e.Len(); l != test.l {
			t.Fatalf("Len method must return %d, %d gives", test.l, l)
		}
	}
}

func TestElements_Cap(t *testing.T) {
	tests := []struct {
		e *Elements
		c int
	}{
		{emptyElements(), capLen},
		{emptyElementsPush(2), capLen},
		{fullElements(), capLen},
		{popTimes(5, fullElements()), capLen},
	}

	for _, test := range tests {
		if c := test.e.Cap(); c != test.c {
			t.Fatalf("Cap method must return %d, %d gives", test.c, c)
		}
	}
}

func TestElements_IsFull(t *testing.T) {
	tests := []struct {
		e *Elements
		b bool
	}{
		{emptyElements(), false},
		{fullElements(), true},
		{popTimes(capLen, fullElements()), false},
		{popTimes(5, fullElements()), false},
	}

	for _, test := range tests {
		if b := test.e.IsFull(); b != test.b {
			t.Fatalf("IsFull method must return bool %v, now bool %v gives", test.b, b)
		}
	}
}

func TestElements_IsEmpty(t *testing.T) {
	tests := []struct {
		e *Elements
		b bool
	}{
		{emptyElements(), true},
		{fullElements(), false},
		{popTimes(capLen, fullElements()), true},
		{popTimes(5, fullElements()), false},
	}

	for _, test := range tests {
		if b := test.e.IsEmpty(); b != test.b {
			t.Fatalf("IsEmpty method must return bool %v, now bool %v gives", test.b, b)
		}
	}
}

func emptyElements() *Elements {
	return NewElements(capLen)
}

func fullElements() *Elements {
	return emptyElementsPush(capLen)
}

func emptyElementsPush(pushTime int) *Elements {
	e := emptyElements()
	if pushTime > capLen {
		pushTime = capLen
	}
	for i := 0; i < pushTime; i++ {
		e.Push(i)
	}

	return e
}

func popTimes(times int, e *Elements) *Elements {
	if times > e.Len() {
		times = e.Len()
	}
	for i := 0; i < times; i++ {
		e.Pop()
	}
	return e
}

func moveHead(e *Elements, i int) *Elements {
	e.MoveHead(i)

	return e
}

func fullElementsPop(popTime int) *Elements {
	e := fullElements()
	if popTime > capLen {
		popTime = capLen
	}
	for i := 0; i < popTime; i++ {
		e.Pop()
	}
	return e
}

func elementsEqual(e1, e2 *Elements, i int, t *testing.T) {
	if e1.Len() != e2.Len() {
		t.Fatalf("index %d, left length: %d, right length: %d", i, e1.Len(), e2.Len())
	}
	var err1, err2 error
	var v1, v2 interface{}
	for j := 0; j < e1.len; j++ {
		v1, err1 = e1.Pop()
		v2, err2 = e2.Pop()
		if v1.(int) != v2.(int) {
			t.Fatalf("index %d, left is %v, right is %v", i, v1, v2)
		}

		if err1 != err2 {
			t.Fatalf("index %d, left error is %s, right error is %s", i, err1, err2)
		}
	}
}
