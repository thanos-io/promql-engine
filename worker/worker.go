// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package worker

import (
	"context"
	"sync"

	"github.com/thanos-community/promql-engine/executor/model"
)

type doneFunc func()

type Group []*Worker

func NewGroup(numWorkers int, task Task) Group {
	group := make([]*Worker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		group[i] = New(i, task)
	}
	return group
}

func (g Group) Start(ctx context.Context) {
	var wg sync.WaitGroup
	for _, w := range g {
		wg.Add(1)
		go w.start(wg.Done, ctx)
	}
	wg.Wait()
}

type Worker struct {
	ctx context.Context

	workerID int
	input    chan model.StepVector
	output   chan model.StepVector
	doWork   Task
}

type Task func(workerID int, in model.StepVector) model.StepVector

func New(workerID int, task Task) *Worker {
	input := make(chan model.StepVector, 1)
	output := make(chan model.StepVector, 1)

	return &Worker{
		workerID: workerID,
		input:    input,
		output:   output,
		doWork:   task,
	}
}

func (w *Worker) start(done doneFunc, ctx context.Context) {
	w.ctx = ctx
	done()
	for {
		select {
		case <-w.ctx.Done():
			close(w.output)
			return
		case task := <-w.input:
			w.output <- w.doWork(w.workerID, task)
		}
	}
}

func (w *Worker) Send(input model.StepVector) error {
	select {
	case <-w.ctx.Done():
		close(w.input)
		return w.ctx.Err()
	default:
		w.input <- input
		return nil
	}
}

func (w *Worker) GetOutput() (model.StepVector, error) {
	select {
	case <-w.ctx.Done():
		return model.StepVector{}, w.ctx.Err()
	default:
		return <-w.output, nil
	}
}
