package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	defaultPort = 8080

	defaultDelayInSeconds = 1
)

func hello(w http.ResponseWriter, req *http.Request) {
	time.Sleep(defaultDelayInSeconds * time.Second)

	if _, err := fmt.Fprintf(w, "hello\n"); err != nil {
		fmt.Printf("err: %v\n", err)
	}
}

func listenAndServe() {
	http.HandleFunc("/hello", hello)

	address := fmt.Sprintf("0.0.0.0:%v", defaultPort)
	fmt.Printf("listening on address: %v\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("can't listen on port %v: %v", defaultPort, err)
	}
}

func client(requestsNum int, infinityRequests bool, goroutinesNum int) {
	var wg sync.WaitGroup
	wg.Add(goroutinesNum)

	for i := 0; i < goroutinesNum; i++ {
		go func(i int) {
			defer wg.Done()

			for j := 0; j < requestsNum || infinityRequests; j++ {
				address := fmt.Sprintf("http://localhost:%v/hello", defaultPort)
				resp, err := http.Get(address)
				if err != nil {
					fmt.Printf("err: %v\n", err)
					continue
				}

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("err: %v\n", err)
					continue
				}
				_ = body
				//fmt.Printf("body: %s\n", body)
				//fmt.Printf("goroutine %v, request %v is done\n", i, j)
			}

			fmt.Printf("goroutine %v is done\n", i)
		}(i)
	}

	wg.Wait()
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "server",
				Aliases: []string{"s"},
				Action: func(cCtx *cli.Context) error {
					listenAndServe()
					return nil
				},
			},
			{
				Name:    "client",
				Aliases: []string{"c"},
				Action: func(cCtx *cli.Context) error {
					client(10, true, 10)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
