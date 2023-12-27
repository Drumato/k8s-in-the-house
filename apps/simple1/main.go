package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

const simple1Version = "0.1.0"

func getIndex(w http.ResponseWriter, _ *http.Request) {
	messages := []string{
		generateSimple1Message(),
	}

	resp, err := http.Get("http://simple2.k8s-in-the-house.svc.cluster.local/")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	simple2, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	messages = append(messages, string(simple2))
	io.WriteString(w, strings.Join(messages, ", "))
}

func generateSimple1Message() string {
	return fmt.Sprintf("Hello from Simple1(v%s)!", simple1Version)
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
