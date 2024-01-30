package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/abersheeran/httpbenchmark/core"
	. "github.com/abersheeran/rgo-error"
	"github.com/go-resty/resty/v2"
)

func main() {
	concurrency := flag.Int("c", 100, "Concurrency")
	requests := flag.Int("r", 10000, "Requests")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] url\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		panic("Please specify url")
	}
	if len(args) > 1 {
		panic("Please use: `[options] url`, not `url [options]`")
	}
	url := args[0]

	res := core.Benchmark{
		Concurrency: *concurrency,
		Requests:    *requests,
		CallRequest: func(request *resty.Request) Result[*resty.Response] {
			return AsResult(request.Get(url))
		},
	}.Run()
	res.Print()
}
