package core

import (
	"fmt"
	"sync"
	"time"

	. "github.com/abersheeran/rgo-error"
	"github.com/go-resty/resty/v2"
)

type Benchmark struct {
	Concurrency int
	Requests    int
	CallRequest func(*resty.Request) Result[*resty.Response]
}

type BenchmarkResult struct {
	TotalRequests     int
	TotalTime         time.Duration
	RequestsPerSecond float64

	MaxRequestTime time.Duration
	MinRequestTime time.Duration

	SuccessRequests int
	BadRequests     int
	FailRequests    int
}

func (b Benchmark) Run() BenchmarkResult {
	client := resty.New()
	channel := make(chan Result[*resty.Response], b.Requests)
	control := make(chan bool, b.Concurrency)

	var start_time, end_time time.Time
	var res BenchmarkResult

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		start_time = time.Now()
		for i := 0; i < b.Requests; i++ {
			control <- true
			go func() {
				defer func() { <-control }()
				r := b.CallRequest(client.R().EnableTrace())
				channel <- r
			}()
		}
	}()

	go func() {
		defer wg.Done()

		res = BenchmarkResult{}
		for i := 0; i < b.Requests; i++ {
			r := <-channel
			res.TotalRequests++
			if r.IsErr() {
				res.FailRequests++
				continue
			}
			response := r.Unwrap()
			switch response.StatusCode() / 100 {
			case 2:
				res.SuccessRequests++
			case 4:
				res.BadRequests++
			case 5:
				res.FailRequests++
			}

			ti := response.Request.TraceInfo()
			if res.MaxRequestTime < ti.TotalTime {
				res.MaxRequestTime = ti.TotalTime
			}
			if res.MinRequestTime == 0 || res.MinRequestTime > ti.TotalTime {
				res.MinRequestTime = ti.TotalTime
			}
		}

		end_time = time.Now()

		res.TotalTime = end_time.Sub(start_time)
		res.RequestsPerSecond = float64(res.TotalRequests) / res.TotalTime.Seconds()
	}()

	wg.Wait()

	return res
}

func (res BenchmarkResult) Print() {
	fmt.Println("Total Requests:", res.TotalRequests)
	fmt.Println("Total Time:", res.TotalTime)
	fmt.Println("Requests Per Second:", res.RequestsPerSecond)
	fmt.Println("Max Request Time:", res.MaxRequestTime)
	fmt.Println("Min Request Time:", res.MinRequestTime)
	fmt.Println("Success Requests:", res.SuccessRequests)
	fmt.Println("Bad Requests:", res.BadRequests)
	fmt.Println("Fail Requests:", res.FailRequests)
}
