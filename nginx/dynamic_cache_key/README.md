
## Results

```
$ make dynamic_cache_key
...
~ # nginx
~ # curl localhost/foo
1475907148.694
~ # curl localhost/foo
1475907148.694
~ # curl localhost/bar
1475907151.999
~ # curl localhost/bar
1475907151.999
~ # curl localhost/foo
1475907148.694
~ # curl 'localhost/internal/invalidate_cache?uri=/foo'
uri:/foo
newval:1
err:nil
forcible:false
~ # curl localhost/foo
1475907177.381
~ # curl localhost/foo
1475907177.381
~ # curl localhost/bar
1475907151.999
```
