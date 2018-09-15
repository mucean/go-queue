package goqueue

import (
	"errors"
	"sort"
)

func NewElements(cap int) *Elements {
	return &Elements{
		values: make([]interface{}, cap),
		cap:    cap,
	}
}

var emptyErr = errors.New("elements is empty")
var fullErr = errors.New("elements is full")
var pushErr = errors.New("push index is in the end")
var notFoundErr = errors.New("not found")

type Elements struct {
	values []interface{}
	len    int
	cap    int
	head   int
}

func (e *Elements) Pop() (interface{}, error) {
	if e.IsEmpty() {
		return nil, emptyErr
	}

	tail := e.tail()
	e.len -= 1
	value := e.values[tail]
	e.values[tail] = nil
	return value, nil
}

func (e *Elements) Push(v interface{}) error {
	err := e.pushable()
	if err != nil {
		return err
	}

	e.len += 1
	e.values[e.tail()] = v

	return nil
}

func (e *Elements) PushForce(v interface{}) error {
	err := e.Push(v)

	if err != pushErr || e.head == 0 {
		return err
	}

	if err = e.Rebuild(); err == nil {
		return e.Push(v)
	} else {
		return err
	}
}

func (e *Elements) Find(f func(v interface{}) bool) (int, error) {
	if e.IsEmpty() {
		return 0, emptyErr
	}

	for start := e.tail(); start >= e.head; start-- {
		if f(e.values[start]) {
			return start, nil
		}
	}

	return 0, notFoundErr
}

func (e *Elements) FindAll(f func(v interface{}) bool) ([]int, error) {
	if e.IsEmpty() {
		return nil, emptyErr
	}

	index := make([]int, 0)
	for start := e.tail(); start >= e.head; start-- {
		if f(e.values[start]) {
			index = append(index, start)
		}
	}

	if len(index) > 0 {
		return index, nil
	}
	return nil, notFoundErr
}

func (e *Elements) Rebuild() error {
	if e.head == 0 {
		return nil
	}

	if e.IsEmpty() {
		e.head = 0
		return nil
	}

	old := e.values
	e.values = make([]interface{}, e.cap)
	length := copy(e.values, old[e.head:e.tail()+1])
	if length != e.len {
		e.values = old
		return errors.New("rebuild failed")
	}

	e.head = 0

	return nil
}

func (e *Elements) MoveHead(head int) int {
	if e.head > head || e.IsEmpty() {
		return 0
	}

	if tail := e.tail(); tail < head {
		head = tail
	}

	index := make([]int, 0, head-e.head+1)
	for i := e.head; i <= head; i++ {
		index = append(index, i)
	}

	return e.eraseByIndex(index)
}

func (e *Elements) MoveTail(t int) int {
	if e.IsEmpty() {
		return 0
	}

	tail := e.tail()

	if t > tail {
		return 0
	}

	if t < e.head {
		t = e.head
	}

	index := make([]int, 0, tail-t+1)
	for i := t; i <= tail; i++ {
		index = append(index, i)
	}

	return e.eraseByIndex(index)
}

func (e *Elements) eraseByIndex(index []int) int {
	eCount := 0
	if index == nil {
		return eCount
	}

	if e.IsEmpty() {
		return eCount
	}

	sort.Ints(index)

	doIndex := make([]int, 0, len(index))
	tail := e.tail()
	for _, i := range index {
		if i < e.head || i > tail {
			continue
		}

		e.values[i] = nil
		doIndex = append(doIndex, i)
		eCount++
	}
	e.len -= eCount

	if e.IsEmpty() {
		e.head += eCount
		return eCount
	}

	var ii int
	for len(doIndex) != 0 {
		ii = doIndex[0]
		if doIndex[0] != e.head {
			break
		}
		e.head++
		doIndex = doIndex[1:]
	}

	for len(doIndex) != 0 {
		ui := doIndex[0]
		doIndex = doIndex[1:]

		for len(doIndex) != 0 {
			if doIndex[0] != ui+1 {
				break
			}
			ui = doIndex[0]
			doIndex = doIndex[1:]
		}
		var start, end int
		start = ui + 1
		if len(doIndex) == 0 {
			end = e.tail()
		} else {
			end = doIndex[0] - 1
		}

		for ; start <= end; start += 1 {
			e.values[ii] = e.values[start]
			e.values[start] = nil
			ii++
		}
	}

	return eCount
}

func (e *Elements) Tail() (int, error) {
	if e.IsEmpty() {
		return 0, emptyErr
	}

	return e.tail(), nil
}

func (e *Elements) tail() int {
	return e.head + e.len - 1
}

func (e *Elements) Pushable() bool {
	return e.pushable() == nil
}

func (e *Elements) pushable() error {
	if e.IsFull() {
		return fullErr
	}

	if e.tail()+1 >= e.cap {
		return pushErr
	}

	return nil
}

func (e *Elements) Len() int {
	return e.len
}

func (e *Elements) Cap() int {
	return e.cap
}

func (e *Elements) IsFull() bool {
	return e.len >= e.cap
}

func (e *Elements) IsEmpty() bool {
	return e.len == 0
}
