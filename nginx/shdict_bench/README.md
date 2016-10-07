
## Run

```
$ make shdict_bench
# nginx
# wrk -t 4 -c 100 http://127.0.0.1/set?key=foo\&value=hello
```


## Results

```
$ make shdict_bench
~ # nginx
~ # wrk -t 4 -c 100 http://127.0.0.1/set?key=foo\&value=hello
Running 10s test @ http://127.0.0.1/set?key=foo&value=hello
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.56ms    1.59ms  67.57ms   96.52%
    Req/Sec    10.17k     1.08k   13.94k    85.50%
  404842 requests in 10.03s, 82.22MB read
Requests/sec:  40346.08
Transfer/sec:      8.19MB
~ # wrk -t 4 -c 100 http://127.0.0.1/set?key=foo\&value=hello
Running 10s test @ http://127.0.0.1/set?key=foo&value=hello
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.76ms    4.21ms 173.24ms   99.17%
    Req/Sec     9.90k     1.33k   17.66k    86.75%
  394059 requests in 10.06s, 80.03MB read
Requests/sec:  39181.08
Transfer/sec:      7.96MB
~ # wrk -t 4 -c 100 http://127.0.0.1/set?key=foo\&value=hello
Running 10s test @ http://127.0.0.1/set?key=foo&value=hello
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.63ms    1.58ms  56.54ms   93.72%
    Req/Sec     9.88k     1.40k   15.44k    86.50%
  393195 requests in 10.06s, 79.85MB read
Requests/sec:  39081.77
Transfer/sec:      7.94MB
~ # wrk -t 4 -c 100 http://127.0.0.1/get?key=foo
Running 10s test @ http://127.0.0.1/get?key=foo
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.82ms    3.93ms 158.44ms   98.78%
    Req/Sec     9.67k     1.38k   14.36k    86.00%
  384912 requests in 10.05s, 78.17MB read
Requests/sec:  38317.24
Transfer/sec:      7.78MB
~ # wrk -t 4 -c 100 http://127.0.0.1/get?key=foo
Running 10s test @ http://127.0.0.1/get?key=foo
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.77ms    4.64ms 185.90ms   99.15%
    Req/Sec    10.01k     1.45k   23.95k    89.28%
  399165 requests in 10.10s, 81.06MB read
Requests/sec:  39525.47
Transfer/sec:      8.03MB
~ # wrk -t 4 -c 100 http://127.0.0.1/get?key=foo
Running 10s test @ http://127.0.0.1/get?key=foo
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.89ms    4.26ms 177.90ms   99.01%
    Req/Sec     9.40k     1.50k   19.51k    80.50%
  374562 requests in 10.07s, 76.07MB read
Requests/sec:  37194.31
Transfer/sec:      7.55MB
~ # curl localhost/baseline
{"key":"key","ok":"ok","value":"value"}
~ # wrk -t 4 -c 100 http://127.0.0.1/baseline
Running 10s test @ http://127.0.0.1/baseline
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.49ms    3.35ms 151.95ms   98.96%
    Req/Sec    10.80k     1.61k   22.37k    85.50%
  430070 requests in 10.08s, 87.34MB read
Requests/sec:  42670.26
Transfer/sec:      8.67MB
~ # wrk -t 4 -c 100 http://127.0.0.1/baseline
Running 10s test @ http://127.0.0.1/baseline
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.33ms    2.10ms  84.01ms   98.26%
    Req/Sec    11.27k     1.61k   25.82k    88.75%
  448788 requests in 10.07s, 91.14MB read
Requests/sec:  44582.43
Transfer/sec:      9.05MB
~ # wrk -t 4 -c 100 http://127.0.0.1/baseline
Running 10s test @ http://127.0.0.1/baseline
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.48ms    2.24ms  90.29ms   98.61%
    Req/Sec    10.61k     1.52k   26.32k    87.75%
  422446 requests in 10.06s, 85.79MB read
Requests/sec:  42008.75
Transfer/sec:      8.53MB
```

```
$ git diff master
diff --git a/nginx/shdict_bench/nginx.conf b/nginx/shdict_bench/nginx.conf
index 2f77838..b02a0d6 100644
--- a/nginx/shdict_bench/nginx.conf
+++ b/nginx/shdict_bench/nginx.conf
@@ -86,5 +86,19 @@ http {
         })
       }
     }
+
+    # NOTE: url is /baselin2 not /baseline2 so that the length
+    # of url become the same as /baseline.
+    location = /baselin2 {
+      content_by_lua_block {
+        ngx.say('{"key":"key","ok":"ok","value":"value"}')
+      }
+    }
+
+    # NOTE: url is /baselin3 not /baseline3 so that the length
+    # of url become the same as /baseline.
+    location = /baselin3 {
+      echo '{"key":"key","ok":"ok","value":"value"}';
+    }
   }
 }
```

```
$ make shdict_bench
~ # nginx
~ # curl localhost/baselin2
{"key":"key","ok":"ok","value":"value"}
~ # wrk -t 4 -c 100 http://127.0.0.1/baselin2
Running 10s test @ http://127.0.0.1/baselin2
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.32ms    2.84ms 135.60ms   99.22%
    Req/Sec    11.47k     1.30k   18.92k    89.50%
  456278 requests in 10.05s, 92.66MB read
Requests/sec:  45390.13
Transfer/sec:      9.22MB
~ # wrk -t 4 -c 100 http://127.0.0.1/baselin2
Running 10s test @ http://127.0.0.1/baselin2
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.51ms    4.05ms 139.08ms   98.91%
    Req/Sec    11.22k     1.48k   17.49k    87.00%
  446571 requests in 10.05s, 90.69MB read
Requests/sec:  44436.44
Transfer/sec:      9.02MB
~ # wrk -t 4 -c 100 http://127.0.0.1/baselin2
Running 10s test @ http://127.0.0.1/baselin2
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.48ms    1.97ms  83.32ms   94.27%
    Req/Sec    10.73k     1.82k   17.59k    83.75%
  427186 requests in 10.06s, 86.76MB read
Requests/sec:  42480.19
Transfer/sec:      8.63MB
```


```
$ make shdict_bench
~ # nginx
~ # curl localhost/baselin3
{"key":"key","ok":"ok","value":"value"}
~ # wrk -t 4 -c 100 http://127.0.0.1/baselin3
Running 10s test @ http://127.0.0.1/baselin3
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.92ms    2.98ms 122.43ms   98.64%
    Req/Sec    14.44k     1.79k   19.38k    85.75%
  574775 requests in 10.06s, 116.73MB read
Requests/sec:  57146.55
Transfer/sec:     11.61MB
~ # wrk -t 4 -c 100 http://127.0.0.1/baselin3
Running 10s test @ http://127.0.0.1/baselin3
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.07ms    2.45ms  70.98ms   98.37%
    Req/Sec    13.38k     2.15k   22.59k    81.75%
  532605 requests in 10.05s, 108.16MB read
Requests/sec:  53011.89
Transfer/sec:     10.77MB
~ # wrk -t 4 -c 100 http://127.0.0.1/baselin3
Running 10s test @ http://127.0.0.1/baselin3
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.85ms  708.58us  36.91ms   93.00%
    Req/Sec    13.81k     1.08k   18.83k    79.25%
  549533 requests in 10.03s, 111.60MB read
Requests/sec:  54773.61
Transfer/sec:     11.12MB
```

```
$ docker run -it --rm testing_nginx sh
~ # nginx
~ # curl localhost
<!DOCTYPE html>
<html>
<head>
<title>Welcome to OpenResty!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to OpenResty!</h1>
<p>If you see this page, the OpenResty web platform is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="https://openresty.org/">openresty.org</a>.<br/>

<p><em>Thank you for flying OpenResty.</em></p>
</body>
</html>
~ # wrk -t 4 -c 100 http://127.0.0.1/
Running 10s test @ http://127.0.0.1/
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     9.10ms   16.32ms 372.06ms   98.70%
    Req/Sec     3.29k   408.68     7.68k    90.73%
  130731 requests in 10.07s, 99.98MB read
Requests/sec:  12983.11
Transfer/sec:      9.93MB
~ # wrk -t 4 -c 100 http://127.0.0.1/
Running 10s test @ http://127.0.0.1/
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     8.82ms   12.29ms 257.04ms   98.35%
    Req/Sec     3.30k   460.26     8.78k    89.50%
  131415 requests in 10.06s, 100.51MB read
Requests/sec:  13057.95
Transfer/sec:      9.99MB
~ # wrk -t 4 -c 100 http://127.0.0.1/
Running 10s test @ http://127.0.0.1/
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     8.99ms   14.82ms 299.39ms   98.33%
    Req/Sec     3.35k   501.13     7.30k    92.75%
  133586 requests in 10.07s, 102.17MB read
Requests/sec:  13264.39
Transfer/sec:     10.14MB
```

## Test environment

Tested on Mackbook Pro (Retina, 2012 Mid).
CPU: Intel Core i7 2.6GHz Quad Cores
RAM: 16GB 1600MHz DDR3

```
~ # nginx -V
nginx version: openresty/1.11.2.1
built by gcc 5.3.0 (Alpine 5.3.0)
built with OpenSSL 1.0.2j  26 Sep 2016
TLS SNI support enabled
configure arguments: --prefix=/usr/local/openresty/nginx --with-debug --with-cc-opt='-DNGX_LUA_USE_ASSERT -DNGX_LUA_ABORT_AT_PANIC -O2' --add-module=../ngx_devel_kit-0.3.0 --add-module=../echo-nginx-module-0.60 --add-module=../xss-nginx-module-0.05 --add-module=../ngx_coolkit-0.2rc3 --add-module=../set-misc-nginx-module-0.31 --add-module=../form-input-nginx-module-0.12 --add-module=../encrypted-session-nginx-module-0.06 --add-module=../srcache-nginx-module-0.31 --add-module=../ngx_lua-0.10.6 --add-module=../ngx_lua_upstream-0.06 --add-module=../headers-more-nginx-module-0.31 --add-module=../array-var-nginx-module-0.05 --add-module=../memc-nginx-module-0.17 --add-module=../redis2-nginx-module-0.13 --add-module=../redis-nginx-module-0.3.7 --add-module=../rds-json-nginx-module-0.14 --add-module=../rds-csv-nginx-module-0.07 --with-ld-opt=-Wl,-rpath,/usr/local/openresty/luajit/lib --sbin-path=/usr/sbin/nginx --conf-path=/etc/nginx/nginx.conf --error-log-path=/var/log/nginx/error.log --pid-path=/var/run/nginx.pid --lock-path=/var/lock/nginx.lock --http-log-path=/var/log/nginx/access.log --http-client-body-temp-path=/var/lib/nginx/client_body_temp --http-proxy-temp-path=/var/lib/nginx/proxy_temp --http-fastcgi-temp-path=/var/lib/nginx/fastcgi_temp --http-uwsgi-temp-path=/var/lib/nginx/uwsgi_temp --http-scgi-temp-path=/var/lib/nginx/scgi_temp --with-threads --with-file-aio --with-http_dav_module --with-http_gzip_static_module --with-http_auth_request_module --with-http_secure_link_module --with-http_degradation_module --with-http_stub_status_module --with-pcre --with-pcre-jit --with-http_ssl_module
~ # wrk --version
wrk 4.0.2 [epoll] Copyright (C) 2012 Will Glozer
Usage: wrk <options> <url>
  Options:
    -c, --connections <N>  Connections to keep open
    -d, --duration    <T>  Duration of test
    -t, --threads     <N>  Number of threads to use

    -s, --script      <S>  Load Lua script file
    -H, --header      <H>  Add header to request
        --latency          Print latency statistics
        --timeout     <T>  Socket/request timeout
    -v, --version          Print version details

  Numeric arguments may include a SI unit (1k, 1M, 1G)
  Time arguments may include a time unit (2s, 2m, 2h)
$ docker version
Client:
 Version:      1.12.1
 API version:  1.24
 Go version:   go1.7.1
 Git commit:   6f9534c
 Built:        Thu Sep  8 10:31:18 2016
 OS/Arch:      darwin/amd64

Server:
 Version:      1.12.1
 API version:  1.24
 Go version:   go1.6.3
 Git commit:   23cf638
 Built:        Thu Aug 18 17:52:38 2016
 OS/Arch:      linux/amd64
```
