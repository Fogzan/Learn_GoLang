package workerPool

import "fmt"

type workers struct {
	id                int
	countComlitedWork int
}

var arrayWorkers = []*workers{
	&workers{id: 0},
	&workers{id: 1},
	&workers{id: 2},
	&workers{id: 3},
	&workers{id: 4},
	&workers{id: 5},
	&workers{id: 6},
	&workers{id: 7},
	&workers{id: 8},
	&workers{id: 9},
}

type Pool[Data any] struct {
	handler func(int, Data)
	chPoll  chan *workers
}

func New[Data any](handler func(int, Data)) *Pool[Data] {
	return &Pool[Data]{
		handler: handler,
		chPoll:  make(chan *workers, len(arrayWorkers)),
	}
}

func (p *Pool[Data]) Create() {
	for i := range len(arrayWorkers) {
		p.chPoll <- arrayWorkers[i]
	}
}

func (p *Pool[Data]) Execute(data Data) {
	w := <-p.chPoll
	go func() {
		p.handler(w.id, data)
		w.countComlitedWork++
		// fmt.Println("----------------------> ", w.id)s
		p.chPoll <- w
	}()

}

func (p *Pool[Data]) Wait() {
	for range arrayWorkers {
		<-p.chPoll
	}
}

func (p *Pool[Data]) ShowStats() {
	for _, i := range arrayWorkers {
		fmt.Printf("Worker %v отработал %v раз\n", i.id, i.countComlitedWork)
	}
}
