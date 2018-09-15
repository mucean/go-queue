package goqueue

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestNewQueue(t *testing.T) {
	q := NewQueue()
	if q == nil {
		t.Fatalf("type of return value by NewQueue function is nil, but Queue pointer expected")
	}
}

func TestNewLimitQueue(t *testing.T) {
	qs := []*Queue{
		NewLimitQueue(-1),
		NewLimitQueue(0),
		NewLimitQueue(100),
	}

	for _, q := range qs {
		if q == nil {
			t.Fatalf("type of return value by NewLimitQueue function is nil, but Queue pointer expected")
		}
	}
}

func TestNewQueueWithELen(t *testing.T) {
	tests := []struct {
		q     *Queue
		eLen  int
		len   int
		cap   int
		count int
	}{
		{NewQueueWithELen(100, -1), 100, 1, math.MaxInt32, 0},
		{NewQueueWithELen(200, 0), 200, 1, math.MaxInt32, 0},
		{NewQueueWithELen(100, 101), 100, 1, 101, 0},
	}

	for i, test := range tests {
		if test.eLen != test.q.eLen {
			fatalf(t, i, "length of elements expect %d, %d gives", test.eLen, test.q.eLen)
		}

		if test.len != test.q.len {
			fatalf(t, i, "length of elements slice expect %d, %d gives", test.len, test.q.len)
		}

		if test.q.len != len(test.q.lists) {
			fatalf(t, i, "length of elements slice expect %d, %d gives", test.q.len, len(test.q.lists))
		}

		if test.cap != test.q.cap {
			fatalf(t, i, "capacity of Queue expect %d, %d gives", test.cap, test.q.cap)
		}

		if test.count != test.q.count {
			fatalf(t, i, "the count of value in Queue expect %d, %d gives", test.count, test.q.count)
		}
	}
}

func TestQueue_Push(t *testing.T) {
	type check struct {
		t   int
		err error
	}
	tests := []struct {
		q  *Queue
		ck []check
	}{
		{emptyQueue(), []check{
			{50, nil},
			{capacity - 50, nil},
			{1, fqErr},
		}},
		{popQueue(fullQueue(), 50), []check{
			{10, nil},
			{39, nil},
			{1, nil},
			{10, fqErr},
		}},
		{fullQueue(), []check{
			{50, fqErr},
		}},
	}

	printf := func(j, n int, format string, args ...interface{}) string {
		return fmt.Sprintf(fmt.Sprintf("check index %d, %d times, %s", j, n+1, format), args...)
	}
	for i, test := range tests {
		for j, ck := range test.ck {
			for n := 0; n < ck.t; n++ {
				if err := test.q.Push(n); err != ck.err {
					fatalf(t, i, printf(j, n, "the return err of Push method expect `%s`, `%s` gives", ck.err, err))
				}
			}
		}

		for j := len(test.ck) - 1; j >= 0; j-- {
			ck := test.ck[j]
			if ck.err != nil {
				continue
			}
			for n := ck.t - 1; n >= 0; n-- {
				if v, err := test.q.Pop(); err != nil {
					fatalf(t, i, printf(j, n, "Pop method must return `nil` error, `%s` gives", err))
				} else if !reflect.DeepEqual(n, v) {
					fatalf(t, i, printf(j, n, "Push into Queue value expect `%v`, `%v` gives", n, v))
				}
			}
		}
	}
}

func TestQueue_Pop(t *testing.T) {
	tests := []struct {
		q   *Queue
		v   interface{}
		err error
	}{
		{emptyQueue(), nil, eqErr},
		{fullQueue(), capacity - 1, nil},
		{popQueue(fullQueue(), 5), capacity - 1 - 5, nil},
		{popQueue(fullQueue(), capacity), nil, eqErr},
		{popQueue(pushQueue(emptyQueue(), 10), 3), 10 - 1 - 3, nil},
	}

	for i, test := range tests {
		if v, err := test.q.Pop(); err != test.err {
			fatalf(t, i, "Pop method of Queue expect `%s` error, `%s` gives", test.err, err)
		} else if !reflect.DeepEqual(v, test.v) {
			fatalf(t, i, "Pop method of Queue expect `%v` value, `%v` gives", test.v, v)
		}
	}
}

func TestQueue_extend(t *testing.T) {
	tests := []struct {
		q *Queue
		l int
		c int
	}{
		{emptyQueue(), 100, 0},
		{pushQueue(emptyQueue(), eLen*100), 200, eLen * 100},
		{fullQueue(), 200, capacity},
		{popQueue(pushQueue(emptyQueue(), eLen*100), eLen*2), 100, eLen*100 - eLen*2},
		{popQueue(pushQueue(emptyQueue(), eLen*100), eLen), 100, eLen*100 - eLen},
		{pushQueue(popQueue(pushQueue(emptyQueue(), eLen*100), eLen), 1), 200, eLen*100 - eLen + 1},
	}

	for i, test := range tests {
		if test.q.extend() != nil {
			fatalf(t, i, "extend method has failed!!!!!!!!!!!!!!!!!")
		}

		if c := test.q.Len(); c != test.c {
			fatalf(t, i, "count of values in the queue must be %d, %d gives", test.c, c)
		}

		if l := cap(test.q.lists); l != test.l {
			fatalf(t, i, "list length must be %d, %d gives when extend method called", test.l, l)
		}
	}
}

func TestQueue_Pushable(t *testing.T) {
	tests := []struct {
		q   *Queue
		err error
		b   bool
	}{
		{emptyQueue(), nil, true},
		{popQueue(fullQueue(), capacity), nil, true},
		{fullQueue(), fqErr, false},
		{popQueue(fullQueue(), 5), nil, true},
		{pushQueue(popQueue(fullQueue(), capacity), 10), nil, true},
		{popQueue(pushQueue(popQueue(fullQueue(), capacity), capacity), capacity), nil, true},
		{pushQueue(popQueue(pushQueue(popQueue(fullQueue(), capacity), capacity), capacity), capacity), fqErr, false},
	}

	for i, test := range tests {
		if err := test.q.pushable(); err != test.err {
			fatalf(t, i, "index of %d, error returned by pushable of Queue expect `%s`, `%s` gives", test.err, err)
		}

		if b := test.q.Pushable(); b != test.b {
			fatalf(t, i, "index of %d, bool value returned by Pushable of Queue expect %v, %v gives", test.b, b)
		}
	}
}

func TestQueue_Len(t *testing.T) {
	tests := []struct {
		q     *Queue
		count int
	}{
		{emptyQueue(), 0},
		{popQueue(fullQueue(), capacity), 0},
		{fullQueue(), capacity},
		{popQueue(fullQueue(), 5), capacity - 5},
		{pushQueue(popQueue(fullQueue(), capacity), 10), 10},
		{popQueue(pushQueue(popQueue(fullQueue(), capacity), capacity), capacity), 0},
		{pushQueue(popQueue(pushQueue(popQueue(fullQueue(), capacity), capacity), capacity), 10), 10},
	}

	for i, test := range tests {
		if count := test.q.Len(); count != test.count {
			fatalf(t, i, "length of Queue expect %d, %d gives", test.count, count)
		}
	}
}

func TestQueue_IsEmpty(t *testing.T) {
	tests := []struct {
		q *Queue
		b bool
	}{
		{emptyQueue(), true},
		{popQueue(fullQueue(), capacity), true},
		{fullQueue(), false},
		{popQueue(fullQueue(), 5), false},
		{pushQueue(popQueue(fullQueue(), capacity), 10), false},
		{popQueue(pushQueue(popQueue(fullQueue(), capacity), capacity), capacity), true},
		{pushQueue(popQueue(pushQueue(popQueue(fullQueue(), capacity), capacity), capacity), 10), false},
	}

	for i, test := range tests {
		if b := test.q.IsEmpty(); b != test.b {
			fatalf(t, i, "the value return by IsEmpty of Queue expect %v, %v gives", test.b, b)
		}
	}
}

func TestQueue_IsFull(t *testing.T) {
	tests := []struct {
		q *Queue
		b bool
	}{
		{emptyQueue(), false},
		{popQueue(fullQueue(), capacity), false},
		{fullQueue(), true},
		{popQueue(fullQueue(), 5), false},
		{pushQueue(popQueue(fullQueue(), capacity), capacity), true},
		{popQueue(pushQueue(popQueue(fullQueue(), capacity), capacity), capacity), false},
		{pushQueue(popQueue(pushQueue(popQueue(fullQueue(), capacity), capacity), capacity), capacity), true},
	}

	for i, test := range tests {
		if b := test.q.IsFull(); b != test.b {
			fatalf(t, i, "index of %d, the value return by IsFull of Queue expect %v, %v gives", test.b, b)
		}
	}
}

var capacity = 500
var eLen = 4

func emptyQueue() *Queue {
	return NewQueueWithELen(eLen, capacity)
}

func fullQueue() *Queue {
	q := emptyQueue()

	for i := 0; i < capacity; i++ {
		q.Push(i)
	}

	return q
}

func pushQueue(q *Queue, t int) *Queue {
	for i := 0; i < t; i++ {
		q.Push(i)
	}

	return q
}

func popQueue(q *Queue, t int) *Queue {
	for i := 0; i < t; i++ {
		q.Pop()
	}

	return q
}

func fatalf(t *testing.T, i int, format string, args ...interface{}) {
	format = fmt.Sprintf("index %d, %s", i, format)
	t.Fatalf(format, args...)
}
