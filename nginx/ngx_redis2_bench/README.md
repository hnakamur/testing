
## Run

```
$ make ngx_redis2_bench
# redis-server /etc/nginx/redis.conf &
# nginx
# wrk -t 4 -c 100 http://127.0.0.1/get?key=foo
```


## Results

Tested on my Mackbook Air.

### set
```
# curl http://127.0.0.1/set?key=foo\&value=bar
+OK
```

```
# wrk -t 4 -c 100 http://127.0.0.1/set?key=foo\&value=bar
Running 10s test @ http://127.0.0.1/set?key=foo&value=bar
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    12.43ms    6.66ms 208.27ms   84.41%
    Req/Sec     2.05k   357.10     4.48k    72.25%
  81685 requests in 10.06s, 13.79MB read
Requests/sec:   8120.20
Transfer/sec:      1.37MB
```

### get
```
# curl http://127.0.0.1/get?key=foo
$3
bar
```

```
# wrk -t 4 -c 100 http://127.0.0.1/get?key=foo
Running 10s test @ http://127.0.0.1/get?key=foo
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    11.31ms    5.11ms  56.43ms   73.69%
    Req/Sec     2.24k   498.14     6.55k    85.25%
  89125 requests in 10.05s, 15.38MB read
Requests/sec:   8866.00
Transfer/sec:      1.53MB
```
