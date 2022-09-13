package worker

import (
	"sync"

	"github.com/fpetkovski/promql-engine/operators/model"
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

	next *sync.Mutex
	done *sync.Mutex

	exit   bool
	doWork Task
}

type Task func(workerID int, in model.StepVector) model.StepVector

func New(workerID int, task Task) *Worker {
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
		w.next.Lock()
		w.output = w.doWork(w.workerID, w.input)
		w.done.Unlock()
	}
}

func (w *Worker) Send(input model.StepVector) {
	w.input = input
	w.next.Unlock()
}

func (w *Worker) Done() {
	w.done.Lock()
}

func (w *Worker) GetOutput() model.StepVector {
	return w.output
}

func (w *Worker) Shutdown() {
	w.exit = true
	w.next.Unlock()
}
