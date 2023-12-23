package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

loopLabel:
	for {
		select {
		case <-ch:
			break loopLabel
		default:
			resp, err := http.Get("http://simple1.k8s-in-the-house.svc.cluster.local")
			if err != nil {
				log.Println(err)
			}

			if resp != nil {
				log.Println(resp.Status)
			}

			time.Sleep(1 * time.Second)
		}
	}
}
