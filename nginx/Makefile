
dynamic_cache_key:
	docker build . -t testing_nginx
	docker run -it --rm -v $(pwd)/dynamic_cache_key:/etc/nginx testing_nginx sh 

.PHONY: dynamic_cache_key
