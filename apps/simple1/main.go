package main

import (
	"fmt"
	"net/http"

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
	e := echo.New()
	e.GET("/", getIndex)
	e.Logger.Fatal(e.Start(":1323"))
}
