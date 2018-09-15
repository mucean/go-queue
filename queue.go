package goqueue

import (
	"errors"
	"math"
)

func NewQueue() *Queue {
	return NewLimitQueue(0)
}

func NewLimitQueue(cap int) *Queue {
	return NewQueueWithELen(100, cap)
}

func NewQueueWithELen(eLen, cap int) *Queue {
	if cap <= 0 {
		cap = math.MaxInt32
	}
	lists := make([]*Elements, 0, 100)
	return &Queue{
		lists: append(lists, NewElements(eLen)),
		eLen:  eLen,
		len:   1,
		cap:   cap,
	}
}

var eqErr = errors.New("queue is empty")
var fqErr = errors.New("queue is full")

type Queue struct {
	lists []*Elements
	eLen  int
	count int
	len   int
	cap   int
}

func (q *Queue) Push(v interface{}) error {
	if err := q.pushable(); err != nil {
		return err
	}

	if err := q.tailElement().PushForce(v); err == nil {
		q.count++
		return err
	}

	if err := q.extend(); err != nil {
		return err
	}

	err := q.tailElement().PushForce(v)
	if err == nil {
		q.count++
	}

	return err
}

func (q *Queue) Pop() (v interface{}, err error) {
	if q.IsEmpty() {
		return nil, eqErr
	}

	e := q.lists[0]
	v, err = e.Pop()

	if err != nil {
		return
	}

	q.count--
	if q.len != 1 && e.IsEmpty() {
		q.lists = q.lists[1:]
		q.len--
	}

	return
}

func (q *Queue) tailElement() *Elements {
	return q.lists[q.len-1]
}

func (q *Queue) extend() error {
	if q.len == cap(q.lists) {
		lists := make([]*Elements, q.len, q.len+100)
		if copy(lists, q.lists) != q.len {
			return errors.New("extend queue failed")
		}

		q.lists = lists
	}
	q.len++
	q.lists = append(q.lists, NewElements(q.eLen))
	return nil
}

func (q *Queue) Pushable() bool {
	return q.pushable() == nil
}

func (q *Queue) pushable() error {
	if q.IsFull() {
		return fqErr
	}

	return nil
}

func (q *Queue) Len() int {
	return q.count
}

func (q *Queue) IsEmpty() bool {
	return q.count == 0
}

func (q *Queue) IsFull() bool {
	return q.cap <= q.count
}
