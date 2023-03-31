package gscall

import (
	"errors"
	"time"
)

type Command[T1 any, T2 any] struct {
	id      uint64
	value   T1
	resChan chan T2
}

func New[T1 any, T2 any](id uint64, value T1) *Command[T1, T2] {
	return &Command[T1, T2]{
		id:      id,
		resChan: make(chan T2, 1),
		value:   value,
	}
}

func (c *Command[T1, T2]) WaitForResult() (res T2, err error) {
	res, ok := <-c.resChan
	if !ok {
		return res, errors.New("channel occurs error")
	}

	return res, nil
}

func (c *Command[T1, T2]) WaitForResultTimeout(t time.Duration) (res T2, err error) {
	select {

	case <-time.After(t):
		return res, errors.New("wait for response timeout")

	case response, ok := <-c.resChan:
		if !ok {
			return res, errors.New("channel occurs error")
		}
		return response, nil

	}
}

func (c *Command[T1, T2]) GetID() uint64 {
	return c.id
}

func (c *Command[T1, T2]) GetValue() T1 {
	return c.value
}
