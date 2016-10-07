
## Run

```
$ make redis_bench
# redis-server /etc/nginx/redis.conf &
# nginx
# wrk -t 4 -c 100 http://127.0.0.1/baseline
```


## Results

Tested on my Mackbook Air.

```
# wrk -t 4 -c 100 http://127.0.0.1/baseline
Running 10s test @ http://127.0.0.1/baseline
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     6.30ms    2.60ms 109.22ms   93.58%
    Req/Sec     4.03k   568.44     8.37k    84.50%
  160426 requests in 10.05s, 32.58MB read
Requests/sec:  15958.35
Transfer/sec:      3.24MB
```

```
# wrk -t 4 -c 100 http://127.0.0.1/get?key=test
Running 10s test @ http://127.0.0.1/get?key=test
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    10.56ms    4.29ms 127.31ms   79.76%
    Req/Sec     2.40k   435.09     6.62k    78.50%
  95630 requests in 10.04s, 19.42MB read
Requests/sec:   9520.65
Transfer/sec:      1.93MB
```

```
# wrk -t 4 -c 100 http://127.0.0.1/set?key=test\&value=hello
Running 10s test @ http://127.0.0.1/set?key=test&value=hello
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    11.12ms   11.73ms 301.95ms   98.52%
    Req/Sec     2.44k   510.88     6.98k    87.53%
  97276 requests in 10.09s, 19.85MB read
Requests/sec:   9636.47
Transfer/sec:      1.97MB
```

```
# wrk -t 4 -c 100 -s ./wrk/bench.lua  http://127.0.0.1/
Running 10s test @ http://127.0.0.1/
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    11.58ms    9.23ms 292.13ms   97.37%
    Req/Sec     2.26k   476.26     7.55k    89.78%
  90293 requests in 10.10s, 17.82MB read
Requests/sec:   8940.18
Transfer/sec:      1.76MB
```
