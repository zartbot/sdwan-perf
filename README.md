# sdwan-perf
Sdwan-perf is based on golang and could support almost platform for performance
and policy validation.

```bash
   SDWAN Performance Test Report
+--------------+---------------------------+--------------------------+
|    Stats     |        Latency(ms)        |  Bandwidth(Per Session)  |
+--------------+---------------------------+--------------------------+
| mean         |               166.82ms    |                75.35Mbps |
| Jitter       |               770.53ms    |                          |
|              |                           |                          |
| Min          |                87.00ms    |                 0.80Mbps |
| p25          |               102.00ms    |                72.74Mbps |
| p75          |               110.00ms    |                78.44Mbps |
| p90          |               113.00ms    |                81.64Mbps |
| p95          |               116.00ms    |                83.34Mbps |
| p99          |               136.00ms    |                86.03Mbps |
| Max          |              9978.00ms    |                91.97Mbps |
+--------------+---------------------------+--------------------------+
| Count: 16102 | Error: 312 | Timeout: 300 | Total-BW:    7534.52Mbps |
+--------------+---------------------------+--------------------------+
```


## Server Mode

You could specify multi ports in port parameter.
This feature is useful when you want to compare different SDWAN policy

```bash
./sdwan-perf_linux -role=server -port=8000,8001,8002,8003
```

## Client Mode 

```bash
./sdwan-perf_linux -role=client  -duration=100 -server=127.0.0.1 -port=8001 -size=1000000 -num=100
```

You could use this tool to verify other site's performance like below
```bash
 ./sdwan-perf_linux -role=client -url=https://www.google.com -num=1 
```

## TODO
Will merge it to Ruta linkstate_probe node and provide a cloud native 
performance measurement framework.

## Usage
```bash
Usage of ./sdwan-perf_linux:
  -duration int
    	Test Duration (default 60)
  -fin
    	server mode close connection after send response
  -num int
    	Num of clients (default 10)
  -port string
    	Server Port (default "8000")
  -reqs int
    	Pipeline reqs per client (default 10)
  -role string
    	Role: client|server (default "client")
  -server string
    	Server IP address (default "127.0.0.1")
  -size int
    	bandwidth test block size (default 1)
  -timeout int
    	client timeout seconds (default 10)
  -url string
    	Testing URL
```


## Acknowlegement 
The Quantile stats lib is from, Thanks a lot for Beorn7 and bmizerany great work.
- https://github.com/beorn7/perks
- https://github.com/bmizerany/perks

