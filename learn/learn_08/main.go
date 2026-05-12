package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func main() {
	// Запрос с таймаутом 3 секунды
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", "<https://slow-api.example.com>", nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err) // context deadline exceeded
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
}
