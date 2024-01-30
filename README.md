# http benchmark

A programmable HTTP stress testing tool. Written in Golang.

## Install

```bash
git clone https://github.com/abersheeran/httpbenchmark.git
cd httpbenchmark
```

And then, if use `go install` to install it, you can use `httpbenchmark` command in your terminal.

```bash
httpbenchmark -h
```

If you want to use it without installing, you can use `go run` to run it.

```bash
go run main.go -h
```

## Usage

```bash
httpbenchmark https://example.com
```

```bash
httpbenchmark -X POST -H "Content-Type: application/json" -d '{"key": "value"}' https://httpbin.org/post
```

```bash
httpbenchmark -c 300 -r 100000 https://example.com
```
