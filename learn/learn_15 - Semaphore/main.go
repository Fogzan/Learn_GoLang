package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var countWriteToCache = 0
var mute = sync.Mutex{}

func writeToCache() {
	time.Sleep(10 * time.Millisecond)
	mute.Lock()
	countWriteToCache++
	mute.Unlock()
}

var countWorkData = 0
var muteWork = sync.Mutex{}

func workData() {

	semaphore(func() {
		writeToCache()
	})

	muteWork.Lock()
	countWorkData++
	muteWork.Unlock()
}

const maxGorutine = 10

var sem = make(chan struct{}, maxGorutine)

func semaphore(f func()) {
	select {
	case sem <- struct{}{}:
	default:
		return
	}

	go func() {
		f()
		<-sem
	}()
}

func main() {
	ticker := time.NewTicker(1 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

mainFor:
	for {
		select {
		case <-ticker.C:
		case <-ctx.Done():
			break mainFor
		}
		workData()
	}
	fmt.Println("Work Data: ", countWorkData)
	fmt.Println("Write Cache: ", countWriteToCache)

}
