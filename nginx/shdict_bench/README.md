
## Run

```
$ make shdict_bench
# nginx
# wrk -t 4 -c 100 http://127.0.0.1/set?key=foo\&value=hello
```


## Results

Tested on my Mackbook Air.

```
# wrk -t 4 -c 100 http://127.0.0.1/get?key=foo
Running 10s test @ http://127.0.0.1/get?key=foo
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    12.22ms   27.88ms 361.89ms   95.46%
    Req/Sec     3.62k     1.19k    8.22k    74.81%
  140864 requests in 10.09s, 28.61MB read
Requests/sec:  13958.58
Transfer/sec:      2.83MB
```

```
# wrk -t 4 -c 100 http://127.0.0.1/set?key=foo\&value=hello
Running 10s test @ http://127.0.0.1/set?key=foo&value=hello
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     6.19ms    4.54ms 149.07ms   95.03%
    Req/Sec     4.19k     0.97k   12.89k    89.75%
  166745 requests in 10.07s, 33.86MB read
Requests/sec:  16559.79
Transfer/sec:      3.36MB
```
