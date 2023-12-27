package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
)

const simple1Version = "0.1.0"

type Response struct {
	Messages []string `json:"messages"`
}

func getIndex(c echo.Context) error {
	messages := []string{
		generateSimple1Message(),
	}

	return c.JSON(http.StatusOK, Response{Messages: messages})
}

func generateSimple1Message() string {
	return fmt.Sprintf("Hello from Simple1(v%s)!", simple1Version)
}

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	r := echo.New()
	r.GET("/", getIndex)

	go func() {
		if err := r.Start(":12345"); err != nil {
			log.Println(err)
		}
	}()

loopLabel:
	for {
		select {
		case <-ctx.Done():
			if err := r.Shutdown(ctx); err != nil {
				log.Println(err)
			}
			break loopLabel
		}
	}
}
