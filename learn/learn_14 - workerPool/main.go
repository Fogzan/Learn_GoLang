package main

import (
	"context"
	"fmt"
	workerPool "main/workerpool"
	"math/rand"
	"time"
)

type Messages struct {
	id       int
	Messages string
	creator  string
}

const maxMessagesCount = 10

func getMessages() []Messages {
	countMessages := rand.Intn(maxMessagesCount) + 10
	result := make([]Messages, countMessages)
	for i := range countMessages {
		result[i].id = i
	}
	return result
}

func processMessage(workerId int, msg Messages) {
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("Worker: %v завершил обрабатывать сообщение: %v\n", workerId, msg.id)
}

func main() {
	pool := workerPool.New(processMessage)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

mainfor:
	for {
		select {
		case <-ctx.Done():
			break mainfor
		default:
		}
		pool.Create()
		messages := getMessages()

		for _, message := range messages {
			pool.Execute(message)
		}

		pool.Wait()
	}

	pool.ShowStats()
}
