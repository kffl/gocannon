# gocannon - HTTP benchmarking tool

Gocannon is a lightweight HTTP benchmarking tool, intended to measure changes in backend application performance over time. It keeps a detailed log of each request that is sent, not just the histogram of their latencies.

## Usage

```
usage: gocannon [<flags>] <target>

Flags:
      --help             Show context-sensitive help (also try --help-long and
                         --help-man).
  -d, --duration=10s     Load test duration
  -c, --connections=50   Maximum number of concurrent connections
  -t, --timeout=200ms    HTTP client timeout
  -o, --output=file.csv  File to save the CSV output (raw req data)
  -i, --interval=250ms   Interval for statistics calculation
      --version          Show application version.

Args:
  <target>  HTTP target URL
```

Below is an example of a load test conducted using gocannon against an Express.js server (notice the performance improvement over time under sustained load):

```
> gocannon http://localhost:3000/static -d 3s -c 50
Attacking http://localhost:3000/static with 50 connections over 3s
gocannon goes brr...
Total Req:    14990
Req/s:         4996.67
Interval stats: (interval = 250ms) 
|--REQS--|    |------------------------LATENCY-------------------------|
     Count           AVG         P50         P75         P90         P99
       673   18.004823ms 17.768106ms 18.810048ms 20.358396ms 31.970312ms
       777   16.090895ms 16.096573ms 17.475942ms 18.049402ms 19.867212ms
       922   13.630393ms 13.573315ms 14.279727ms 15.720507ms 17.725149ms
       987   12.715278ms 12.050508ms 14.378844ms 15.661408ms 20.669683ms
      1137   10.997768ms 10.201962ms 11.468307ms 13.736441ms 15.374815ms
      1354    9.186035ms  8.689116ms  9.831591ms 10.417963ms 11.393158ms
      1412    8.889394ms  8.444931ms  8.834144ms 10.480375ms 11.886929ms
      1428    8.792173ms  7.880262ms  9.580714ms 10.753837ms 14.997604ms
      1599    7.817586ms  7.578824ms  7.645575ms  9.083118ms 10.047772ms
      1505    8.244456ms  7.580401ms  8.307621ms  9.611815ms 17.449891ms
      1617    7.745439ms  7.422272ms  7.487905ms  9.132874ms 11.457255ms
      1579    7.901125ms  7.525814ms  7.655212ms  9.577895ms 10.203331ms
----------
     14990    9.986321ms  8.540818ms 11.252258ms 15.309979ms 19.632593ms
```

When used with the `--output=filename.csv` flag, raw request data (which includes the request HTTP code, starting timestamp and response timestamp both in nanoseconds since the beginning of the test) can be written into the specified csv file in the following format:

```
code;start;end;
200;509761;30592963;
200;694043;34837869;
200;963463;36094909;
200;841397;37319812;
[...]
```