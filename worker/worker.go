// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package worker

import (
	"context"
	"github.com/thanos-community/promql-engine/physicalplan/model"
)

type Group []*Worker

func NewGroup(numWorkers int, task Task) Group {
	group := make([]*Worker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		group[i] = New(i, task)
	}
	return group
}

func (g Group) Start(ctx context.Context) {
	for _, w := range g {
		go w.Start(ctx)
	}
}

type Worker struct {
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

func (w *Worker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(w.input)
			close(w.output)
			return
		case task := <-w.input:
			w.output <- w.doWork(w.workerID, task)
		}
	}
}

func (w *Worker) Send(input model.StepVector) {
	w.input <- input
}

func (w *Worker) GetOutput() model.StepVector {
	return <-w.output
}
