package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/abersheeran/httpbenchmark/core"
	. "github.com/abersheeran/rgo-error"
	"github.com/go-resty/resty/v2"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	concurrency := flag.Int("c", 100, "Concurrency")
	requests := flag.Int("r", 10000, "Requests")
	method := flag.String("X", "GET", "Method")
	var headers arrayFlags
	flag.Var(&headers, "H", "Header")
	body := flag.String("d", "", "Request body")
	file := flag.String("f", "", "Request body file path")

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
	if *body != "" && *file != "" {
		panic("Please specify only one of `-b` and `-f`")
	}

	res := core.Benchmark{
		Concurrency: *concurrency,
		Requests:    *requests,
		CallRequest: func(request *resty.Request) Result[*resty.Response] {
			for _, header := range headers {
				parts := strings.SplitN(header, ":", 2)
				if len(parts) != 2 {
					panic(fmt.Errorf("invalid header: %s", header))
				}
				request.SetHeader(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
			if *file != "" {
				file, err := os.Open(*file)
				if err != nil {
					panic(err)
				}
				defer file.Close()

				bytes, err := io.ReadAll(file)
				if err != nil {
					panic(err)
				}
				request.SetBody(bytes)
			}
			if *body != "" {
				request.SetBody(*body)
			}

			return AsResult(request.Execute(*method, url))
		},
	}.Run()
	res.Print()
}
