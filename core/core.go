package core

import (
	"fmt"
	"os"
	"time"

	. "github.com/abersheeran/rgo-error"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

type R = Result[*resty.Response]

type Benchmark struct {
	Concurrency int
	Requests    int
	CallRequest func(*resty.Request) R
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
	channel := make(chan R, b.Requests)
	control := make(chan bool, b.Concurrency)

	var startTime, endTime time.Time
	res := BenchmarkResult{}

	p := tea.NewProgram(counter[R]{
		sub:     channel,
		spinner: spinner.New(),
		listenForActivity: func(sub chan R) tea.Cmd {
			return func() tea.Msg {
				for i := 0; i < b.Requests; i++ {
					control <- true
					go func() {
						defer func() { <-control }()
						r := b.CallRequest(client.R().EnableTrace())
						channel <- r
					}()
				}
				return nil
			}
		},
		waitForActivity: func(sub chan R) tea.Cmd {
			return func() tea.Msg {
				for {
					select {
					case r := <-channel:
						if r.IsErr() {
							res.FailRequests++
						} else {
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
						res.TotalRequests++
						return r
					default:
						if res.TotalRequests == b.Requests {
							return tea.Quit()
						}
					}
				}
			}
		},
	})

	startTime = time.Now()
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
	endTime = time.Now()

	res.TotalTime = endTime.Sub(startTime)
	res.RequestsPerSecond = float64(res.TotalRequests) / res.TotalTime.Seconds()

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
