# :boom: gocannon - HTTP benchmarking tool

[![CI Workflow](https://github.com/kffl/gocannon/workflows/CI/badge.svg)](https://github.com/kffl/gocannon/actions) [![Go Report Card](https://goreportcard.com/badge/github.com/kffl/gocannon)](https://goreportcard.com/report/github.com/kffl/gocannon)

Gocannon is a lightweight HTTP benchmarking tool, intended to measure changes in backend application performance over time. It keeps a detailed log of each request that is sent, not just the histogram of their latencies.

## Usage

```
usage: gocannon [<flags>] <target>

Flags:
      --help              Show context-sensitive help (also try --help-long and
                          --help-man).
  -d, --duration=10s      Load test duration
  -c, --connections=50    Maximum number of concurrent connections
  -t, --timeout=200ms     HTTP client timeout
  -m, --mode="reqlog"     Statistics collection mode: reqlog (logs each request) or hist
                          (stores histogram of completed requests latencies)
  -o, --output=file.csv   File to save the request log in CSV format (reqlog mode) or a
                          text file with raw histogram data (hist mode)
  -i, --interval=250ms    Interval for statistics calculation (reqlog mode)
      --preallocate=1000  Number of requests in req log to preallocate memory for per
                          connection (reqlog mode)
      --method=GET        The HTTP request method (GET, POST, PUT, PATCH or DELETE)
  -b, --body="{data..."   HTTP request body
  -h, --header="k:v" ...  HTTP request header(s). You can set more than one header by
                          repeating this flag.
      --trust-all         Omit SSL certificate validation
      --plugin=/to/p.so   Plugin to run Gocannon with
      --version           Show application version.

Args:
  <target>  HTTP target URL with port (i.e. http://localhost:80/test or https://host:443/x)
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
Responses by HTTP status code:
  200  ->  14990
Requests ended with timeout/socket error: 1
```

### Saving request log to a CSV file

When used with the `--output=filename.csv` flag, raw request data (which includes the request HTTP code, starting timestamp, response timestamp both in nanoseconds since the beginning of the test as well as the ID of the connection which was used for that request) can be written into the specified csv file in the following format:

```
code;start;end;connection;
200;509761;30592963;0;
200;694043;34837869;0;
200;963463;36094909;0;
200;841397;37319812;0;
[...]
```

In the CSV output file, the requests are sorted primarily by connection ID and secondarily by response timestamp (`end` column).

### Preallocating memory for request log

Gocannon stores data of each completed request in a slice (one per each connection). If a given slice runs out of its initial capacity, the underlying array is re-allocated to a newly selected memory address with double the capacity. In order to minimize the performance impact of slice resizing during the load test, you can specify the number of requests in log per connection to preallocate memory for using the `--preallocate` flag.

### Histogram mode

There are some use cases in which request log mode is not desirable. If you don't need a detailed log of each request that was completed during the load test, you can use histogram mode using the `--mode=hist` flag. When using histogram mode, upon each completed request, gocannon will only save its latency in a histogram (in a manner similar to how wrk or bombardier does that). That way, memory usage of the stats collector will remain constant over the entire duration of the load test. This results in gocannon having a total memory footprint of about 3MB when conducting a load test in histogram mode (when the timeout is set to a default 200ms).

Below is an example of a load test conducted with gocannon in histogram mode:

```
> gocannon http://localhost:3000/static --duration=2m --mode=hist
Attacking http://localhost:3000/static with 50 connections over 2m0s
gocannon goes brr...
Total Req:  13715707
Req/s:        114297.56
|------------------------LATENCY (μs)-----------------------|
          AVG         P50         P75         P90         P99
       436.45         278         508         925        2681
Responses by HTTP status code:
  200  ->  13715707
```

Similarly to saving CSV output in request log mode, you can write the histogram data to a text file using the `--output=filename` flag. The output contains the number of hits in each latency histogram bin (having width of 1μs) truncated up until the last non-zero value:

```
0
0
1
1
15
45
[...]
```
### Custom plugins

Gocannon supports user-provided plugins, which can customize the requests sent during the load test. A custom plugin has to satisfy the `GocannonPlugin` interface defined in the `common` package, which also contains types required for plugin development. An example implementation with additional comments is provided in `_example_plugin` folder.

In order to build a plugin, use the following command inside its directory:

```
go build -buildmode=plugin -o plugin.so plugin.go
```

Once you obtain a shared object (`.so`) file, you can provide a path to it via `--plugin` flag. Bare in mind that a custom plugin `.so` file and the Gocannon binary using it must both be compiled using the same Go version.

## Load testing recommendations

-   **Killing non-essential background processes**: Since the activity of background processes running alongside the SUT may introduce anomalies in the obtained experiment results, it is advised to disable non-essential services for the load test duration. Additionally, when conducting a load test against a SUT running on the same machine as gocannon, you may want to assign the SUT and gocannon processes to separate sets of logical cores (i.e. via `taskset`) and update the `GOMAXPROCS` env variable accordingly (so as to reflect the number of cores available to gocannon).
-   **Disabling Turbo Boost, EIST and ensuring a constant CPU clock**: Changes of the CPU clock over the duration of the load test may introduce inconsistencies in the collected metrics, especially given that the detection of a suitable workload after the test start and subsequent increase of the CPU clock may cause observable improvement in application throughput shortly after the start of the load test. If your goal is to measure how the application performance changes over time under sustained load, you may want to disable Turbo Boost and EIST (or their equivalents in your system) as well as ensure that the CPU clock is not being reduced at idle (i.e. by setting the CPU mode to performance via the `cpupower` linux package).
-   **Repeating the test runs**: Performing multiple test runs improves the statistical significance of the obtained results and minimizes impact of potential one-off anomalies (i.e. background process activity during the load test) in the obtained data on the overall result. When combining data from multiple test runs, remember not to average the percentiles. Instead, place the obtained request latencies of all runs in a single dataset and calculate the desired percentiles on that dataset.
