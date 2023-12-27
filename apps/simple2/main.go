package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
)

const version = "0.1.0"

func getIndex(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, generateSimple2Message())
}

func generateSimple2Message() string {
	return fmt.Sprintf("Hello from Simple2(v%s)!", version)
}

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	http.HandleFunc("/", getIndex)

	// Web サーバーを起動する
	log.Fatal(http.ListenAndServe(":12345", nil))

	if err := http.ListenAndServe(":12345", nil); err != nil {
		log.Fatalln(err)

	}
loopLabel:
	for {
		select {
		case <-ctx.Done():
			break loopLabel
		}
	}
}
