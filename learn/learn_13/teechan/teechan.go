package teechan

import (
	"context"
	"sync"
)

type WaitAll = sync.WaitGroup

type WaitFast struct{}

func (w *WaitFast) Add(_ int) {}
func (w *WaitFast) Done()     {}
func (w *WaitFast) Wait()     {}

type WaitG interface {
	Add(int)
	Done()
	Wait()
}

type Teechan struct {
	out       []chan int
	countChan int
	wgFast    []WaitG
	wgAll     WaitG
}

const (
	Fast = iota
	All
)

func New(countChan int, ttype int) *Teechan {

	out := make([]chan int, countChan)
	wgFast := make([]WaitG, countChan)

	for i := range countChan {
		out[i] = make(chan int)
	}

	var wgAll WaitG
	if ttype == Fast {
		wgAll = &WaitFast{}
		for i := range countChan {
			wgFast[i] = &WaitAll{}
		}
	}
	if ttype == All {
		wgAll = &WaitAll{}
		for i := range countChan {
			wgFast[i] = &WaitFast{}
		}
	}

	return &Teechan{
		out:       out,
		wgFast:    wgFast,
		wgAll:     wgAll,
		countChan: countChan,
	}
}

func (t *Teechan) Execute(ctx context.Context, in chan int) []chan int {
	go func() {
		defer func() {
			for i := range t.countChan {
				go func() {
					// time.Sleep(1 * time.Millisecond)
					t.wgFast[i].Wait()
					close(t.out[i])
				}()
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				for i := range t.countChan {
					t.wgAll.Add(1)
					t.wgFast[i].Add(1)
					go func() {
						defer t.wgAll.Done()
						defer t.wgFast[i].Done()
						select {
						case <-ctx.Done():
							return
						case t.out[i] <- val:
						}
					}()
				}
				t.wgAll.Wait()
			}
		}
	}()

	return t.out
}
