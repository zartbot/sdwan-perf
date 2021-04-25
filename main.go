package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"

	"github.com/zartbot/sdwan-perf/client"
	"github.com/zartbot/sdwan-perf/server"
	"github.com/zartbot/sdwan-perf/stats/describe"
	"github.com/zartbot/sdwan-perf/stats/quantile"
)

type Config struct {
	Role            string
	URL             string
	Duration        int
	ServerIP        string
	ServerPort      string
	ServerConnClose bool
	ClientNum       int
	ClientPipeline  int
	ClientTimeout   int
	BWTestSize      int
}

var cli = &Config{
	Role:            "client",
	URL:             "",
	Duration:        60,
	ServerIP:        "127.0.0.1",
	ServerPort:      "8000",
	ServerConnClose: false,
	ClientNum:       1,
	ClientPipeline:  10,
	ClientTimeout:   10,
	BWTestSize:      100,
}

func init() {
	flag.StringVar(&cli.Role, "role", cli.Role, "Role: client|server")
	flag.StringVar(&cli.URL, "url", cli.URL, "Testing URL")
	flag.IntVar(&cli.Duration, "duration", cli.Duration, "Test Duration")
	flag.StringVar(&cli.ServerIP, "server", cli.ServerIP, "Server IP address")
	flag.StringVar(&cli.ServerPort, "port", cli.ServerPort, "Server Port")
	flag.IntVar(&cli.ClientNum, "num", cli.ClientNum, "Num of clients")
	flag.IntVar(&cli.ClientPipeline, "reqs", cli.ClientPipeline, "Pipeline reqs per client")
	flag.IntVar(&cli.ClientTimeout, "timeout", cli.ClientTimeout, "client timeout seconds")
	flag.IntVar(&cli.BWTestSize, "size", cli.BWTestSize, "bandwidth test block size")
	flag.BoolVar(&cli.ServerConnClose, "fin", cli.ServerConnClose, "server mode close connection after send response")
	flag.Parse()
}

func main() {

	switch cli.Role {
	case "server":
		logrus.Info("sdwan_perf is running in server mode.")
		ports := strings.Split(cli.ServerPort, ",")
		for _, port := range ports {
			addr := "0.0.0.0:" + port
			go server.Run(addr, cli.ServerConnClose)
		}

		var wg sync.WaitGroup
		wg.Add(1)
		wg.Wait()

	//case : add controller mode later
	default:
		//client mode
		logrus.Info("sdwan_perf is running in client mode.")
		//if url is not configured, use the server ip and port build url
		if cli.URL == "" {
			if cli.ServerIP == "" || cli.ServerPort == "" {
				logrus.Fatal("Test require specific server address or 3rd party URL")
			}
			if cli.BWTestSize <= 0 {
				logrus.Fatal("Invalid testing block size: ", cli.BWTestSize)
			}
			cli.URL = fmt.Sprintf("http://%s:%s/speed?size=%d", cli.ServerIP, cli.ServerPort, cli.BWTestSize)
		}

		PerfromanceTest(cli)
	}

}

func PerfromanceTest(cli *Config) {
	//create quantile stats
	LatencyQuantile := quantile.NewTargeted(map[float64]float64{
		0.50: 0.005,
		0.90: 0.001,
		0.99: 0.0001,
	})

	BWQuantile := quantile.NewTargeted(map[float64]float64{
		0.50: 0.005,
		0.90: 0.001,
		0.99: 0.0001,
	})
	//create describe stats
	LatencyStats := describe.New()
	BWStats := describe.New()

	cc, err := client.New(int32(cli.ClientNum), int32(cli.ClientPipeline), time.Second*time.Duration(cli.ClientTimeout), cli.URL)
	if err != nil {
		logrus.Fatal("Failed to create client:", err)
	}
	ticker := time.NewTicker(1 * time.Second)
	runTimeout := time.NewTimer(time.Second * time.Duration(cli.Duration))

	errors := 0
	timeouts := 0

	cc.Run()

	for {
		select {
		case err := <-cc.ErrChan:
			errors++
			if err == fasthttp.ErrTimeout {
				timeouts++
			}
		case res := <-cc.RespChan:
			LatencyQuantile.Insert(float64(res.Latency))
			LatencyStats.Append(float64(res.Latency), 2)

			//per-session
			bw := (float64(res.Size) / res.Latency) * 0.008
			if res.Latency == 0 {
				bw = 0
			}
			BWQuantile.Insert(bw)
			BWStats.Append(bw, 2)
		case <-ticker.C:
			TableRender(LatencyStats, LatencyQuantile, BWStats, BWQuantile, errors, timeouts)

		case <-runTimeout.C:
			TableRender(LatencyStats, LatencyQuantile, BWStats, BWQuantile, errors, timeouts)
			//TODO: Add csv output ?
			os.Exit(0)
		}
	}
}

func TableRender(LatencyStats *describe.Item, LatencyQuantile *quantile.Stream, BWStats *describe.Item, BWQuantile *quantile.Stream, errors int, timeouts int) {
	fmt.Printf("\033[H\033[2J")
	fmt.Printf("\n   SDWAN Performance Test Report\n\n")

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Stats ", "Latency(ms)", "Bandwidth(Per Session)"})
	table.SetAutoFormatHeaders(false)
	table.Append([]string{"mean", fmt.Sprintf("%20.2fms", LatencyStats.Mean), fmt.Sprintf("%20.2fMbps", BWStats.Mean)})
	table.Append([]string{"Jitter", fmt.Sprintf("%20.2fms", LatencyStats.Std()), ""})
	table.Append([]string{"", "", ""})
	table.Append([]string{"Min", fmt.Sprintf("%20.2fms", LatencyStats.Min), fmt.Sprintf("%20.2fMbps", BWStats.Min)})
	table.Append([]string{"p25", fmt.Sprintf("%20.2fms", LatencyQuantile.Query(0.25)), fmt.Sprintf("%20.2fMbps", BWQuantile.Query(0.25))})
	table.Append([]string{"p75", fmt.Sprintf("%20.2fms", LatencyQuantile.Query(0.75)), fmt.Sprintf("%20.2fMbps", BWQuantile.Query(0.75))})
	table.Append([]string{"p90", fmt.Sprintf("%20.2fms", LatencyQuantile.Query(0.90)), fmt.Sprintf("%20.2fMbps", BWQuantile.Query(0.90))})
	table.Append([]string{"p95", fmt.Sprintf("%20.2fms", LatencyQuantile.Query(0.95)), fmt.Sprintf("%20.2fMbps", BWQuantile.Query(0.95))})
	table.Append([]string{"p99", fmt.Sprintf("%20.2fms", LatencyQuantile.Query(0.99)), fmt.Sprintf("%20.2fMbps", BWQuantile.Query(0.99))})
	table.Append([]string{"Max", fmt.Sprintf("%20.2fms", LatencyStats.Max), fmt.Sprintf("%20.2fMbps", BWStats.Max)})
	table.SetFooter([]string{fmt.Sprintf("Count: %d", LatencyQuantile.Count()), fmt.Sprintf("Error: %d | Timeout: %d", errors, timeouts), fmt.Sprintf("Total-BW: %10.2fMbps", BWStats.Mean*float64(cli.ClientNum))})

	table.Render()
}
