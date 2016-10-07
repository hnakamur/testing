
dynamic_cache_key:
	docker build . -t testing_nginx
	docker run -it --rm -v $(CURDIR)/dynamic_cache_key:/etc/nginx testing_nginx sh 

redis:
	docker build . -t testing_nginx
	docker run -it --rm -v $(CURDIR)/redis:/etc/nginx testing_nginx sh 

redis_bench:
	docker build . -t testing_nginx
	docker run -it --rm -v $(CURDIR)/redis_bench:/etc/nginx testing_nginx sh 

ngx_redis2_bench:
	docker build . -t testing_nginx
	docker run -it --rm -v $(CURDIR)/ngx_redis2_bench:/etc/nginx testing_nginx sh 

socket_bench:
	docker build . -t testing_nginx
	docker run -it --rm -v $(CURDIR)/$@:/etc/nginx testing_nginx sh 

.PHONY: dynamic_cache_key \
	redis \
	redis_bench \
	ngx_redis2_bench \
	socket_bench
