
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
    Latency    10.02ms    3.06ms  40.90ms   77.42%
    Req/Sec     2.50k   288.31     4.26k    78.00%
  99749 requests in 10.07s, 21.87MB read
Requests/sec:   9906.10
Transfer/sec:      2.17MB
```
