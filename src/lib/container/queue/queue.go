package queue

import "container/list"

type Queue struct {
	q *list.List
}

func (q *Queue) Enqueue(value interface{}) {
	q.q.PushBack(value)
}

func (q *Queue) Peek() interface{} {
	if q.q.Len() > 0 {
		return q.q.Front().Value
	}
	return nil
}

func (q *Queue) Dequeue() interface{} {
	if q.q.Len() > 0 {
		e := q.q.Front()
		q.q.Remove(e)
		return e.Value
	}
	return nil
}

func (q *Queue) Front() interface{} {
	if q.q.Len() > 0 {
		return q.q.Front().Value
	}
	return nil
}

func (q *Queue) Size() int   { return q.q.Len() }
func (q *Queue) Length() int { return q.q.Len() }
func (q *Queue) Len() int    { return q.q.Len() }
func (q *Queue) Empty() bool { return q.q.Len() == 0 }
func New() *Queue            { return &Queue{q: list.New()} }
