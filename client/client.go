package client

import (
	"fmt"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type ClientRespMetric struct {
	Status  int
	Latency float64
	Size    int
}

type PerfClient struct {
	timeout  time.Duration
	num      int32
	pipeline int32
	uri      string
	RespChan chan *ClientRespMetric
	ErrChan  chan error
	Client   *fasthttp.PipelineClient
}

//Prinf Nothing, I hate the logs during performance test
func (p *PerfClient) Printf(format string, args ...interface{}) {
}

func New(clients int32, pipeline int32, timeout time.Duration, uri string) (*PerfClient, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return &PerfClient{}, err
	}

	address := fmt.Sprintf("%v:%v", u.Hostname(), u.Port())
	if u.Port() == "" {
		if u.Scheme == "https" {
			address += "443"
		} else {
			address += "80"
		}
	}

	c := &PerfClient{
		num:      clients,
		pipeline: pipeline,
		timeout:  timeout,
		uri:      uri,
		RespChan: make(chan *ClientRespMetric, 2*clients*pipeline),
		ErrChan:  make(chan error, 2*clients*pipeline),
		Client: &fasthttp.PipelineClient{
			Addr:               address,
			IsTLS:              u.Scheme == "https",
			MaxPendingRequests: int(pipeline),
		},
	}
	logrus.Info("Start testing with: ", uri)
	c.Client.Logger = c
	return c, nil
}

func (p *PerfClient) Run() {
	for j := 0; j < int(p.num); j++ {
		go func() {
			req := fasthttp.AcquireRequest()
			req.SetRequestURI(p.uri)
			res := fasthttp.AcquireResponse()

			for {
				startTime := time.Now()
				if err := p.Client.DoTimeout(req, res, p.timeout); err != nil {
					p.ErrChan <- err
				} else {
					size := len(res.Body()) + 2
					res.Header.VisitAll(func(key, value []byte) {
						size += len(key) + len(value) + 2
					})
					p.RespChan <- &ClientRespMetric{
						Status:  res.Header.StatusCode(),
						Latency: float64(time.Since(startTime).Milliseconds()),
						Size:    size,
					}
					res.Reset()
				}
			}
		}()
	}
}
