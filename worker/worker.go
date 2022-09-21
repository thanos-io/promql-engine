// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package worker

import (
	"sync"

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

func (g Group) Start() {
	for _, w := range g {
		go w.Start()
	}
}

func (g Group) Shutdown() {
	for _, w := range g {
		w.Shutdown()
	}
}

type Worker struct {
	workerID int
	input    model.StepVector
	output   model.StepVector

	// next and done help in starting/stopping Workers
	// and tracking completion.
	next *sync.Mutex
	done *sync.Mutex

	exit   bool
	doWork Task
}

type Task func(workerID int, in model.StepVector) model.StepVector

func New(workerID int, task Task) *Worker {
	// next and done are both locked when a new Worker is spawned.
	next := &sync.Mutex{}
	next.Lock()
	done := &sync.Mutex{}
	done.Lock()

	return &Worker{
		workerID: workerID,
		next:     next,
		done:     done,
		doWork:   task,
	}
}

func (w *Worker) Start() {
	for {
		// Wait for next to be unlocked by Send
		// before starting work by acquiring lock.
		w.next.Lock()
		// Shutdown can also unlock next, so in this case
		// do not run doWork.
		if !w.exit {
			w.output = w.doWork(w.workerID, w.input)
			w.done.Unlock()
		}
	}
}

func (w *Worker) Send(input model.StepVector) {
	// input is received for Worker, so run Start
	// by unlocking next.
	w.input = input
	w.next.Unlock()
}

func (w *Worker) Done() {
	// We can only acquire done lock once Start
	// has finished executing doWork, i.e the Worker
	// has finished its task.
	// If this function returns, it indicates completion.
	w.done.Lock()
}

func (w *Worker) GetOutput() model.StepVector {
	return w.output
}

func (w *Worker) Shutdown() {
	w.exit = true
	w.next.Unlock()
}
