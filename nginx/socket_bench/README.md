
## Run

```
$ make socket_bench
# redis-server /etc/nginx/redis.conf &
# nginx
# wrk -t 4 -c 100 http://127.0.0.1/set?key=foo\&value=hello
```


## Results

Tested on my Mackbook Air.

```
# wrk -t 4 -c 100 http://127.0.0.1/set?key=foo\&value=hello
Running 10s test @ http://127.0.0.1/set?key=foo&value=hello
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     9.72ms    3.39ms  77.03ms   89.67%
    Req/Sec     2.60k   361.84     3.63k    78.75%
  103739 requests in 10.06s, 22.65MB read
Requests/sec:  10312.96
Transfer/sec:      2.25MB
```
