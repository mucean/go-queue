package goqueue

import (
	"errors"
)

func NewQueue() Queue {
	eLen := 100
	return Queue{
		lists: []*Elements{NewElements(eLen)},
		eLen:  eLen,
	}
}

type Queue struct {
	lists []*Elements
	eLen  int
	count int
	head  int
	tail  int
}

func (q *Queue) Len() int {
	return q.count
}

func (q *Queue) Pop() (interface{}, error) {
	if q.count == 0 {
		return nil, errors.New("queue is empty")
	}
	return nil, errors.New("queue is empty")
}