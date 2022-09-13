package worker

import (
	"github.com/fpetkovski/promql-engine/operators/model"
)

type Worker struct {
	input  chan model.StepVector
	output chan model.StepVector

	exit   bool
	doWork func(in model.StepVector) model.StepVector
}

func New(doWork func(in model.StepVector) model.StepVector) *Worker {
	return &Worker{
		input:  make(chan model.StepVector),
		output: make(chan model.StepVector),
		doWork: doWork,
	}
}

func (w *Worker) Send(input model.StepVector) {
	w.input <- input
}

func (w *Worker) GetOutput() model.StepVector {
	return <-w.output
}

func (w *Worker) Shutdown() {
	w.exit = true
	w.input <- model.StepVector{}
}

func (w *Worker) Start() {
	for {
		select {
		case in := <-w.input:
			if w.exit {
				return
			}

			w.output <- w.doWork(in)
		}
	}
}
