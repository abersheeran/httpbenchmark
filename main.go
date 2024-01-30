package main

import (
	"flag"

	"github.com/abersheeran/httpbenchmark/core"
	. "github.com/abersheeran/rgo-error"
	"github.com/go-resty/resty/v2"
)

func main() {
	url := flag.String("url", "", "Target URL")
	concurrency := flag.Int("concurrency", 100, "Concurrency")
	requests := flag.Int("requests", 10000, "Requests")

	flag.Parse()
	if *url == "" {
		panic("URL is required")
	}

	res := core.Benchmark{
		Concurrency: *concurrency,
		Requests:    *requests,
		CallRequest: func(request *resty.Request) Result[*resty.Response] {
			return AsResult(request.Get(*url))
		},
	}.Run()
	res.Print()
}
