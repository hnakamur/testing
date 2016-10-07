FROM alpine:3.4

RUN apk add --no-cache build-base linux-headers openssl-dev pcre-dev ca-certificates wget perl curl vim && update-ca-certificates

ENV WORKDIR=/root
WORKDIR $WORKDIR

RUN wget -O- https://openresty.org/download/openresty-1.11.2.1.tar.gz  | tar xz

RUN cd ./openresty-1.11.2.1 && ./configure \
	--sbin-path=/usr/sbin/nginx \
	--conf-path=/etc/nginx/nginx.conf \
	--error-log-path=/var/log/nginx/error.log \
	--pid-path=/var/run/nginx.pid \
	--lock-path=/var/lock/nginx.lock \
	--http-log-path=/var/log/nginx/access.log \
	--http-client-body-temp-path=/var/lib/nginx/client_body_temp \
	--http-proxy-temp-path=/var/lib/nginx/proxy_temp \
	--http-fastcgi-temp-path=/var/lib/nginx/fastcgi_temp \
	--http-uwsgi-temp-path=/var/lib/nginx/uwsgi_temp \
	--http-scgi-temp-path=/var/lib/nginx/scgi_temp \
	--with-debug \
	-j4 \
	--with-luajit \
	--with-threads \
	--with-file-aio \
	--with-http_dav_module \
	--with-http_gzip_static_module \
	--with-http_auth_request_module \
	--with-http_secure_link_module \
	--with-http_degradation_module \
	--with-http_stub_status_module \
	--with-pcre \
	--with-pcre-jit \
	&& make install -j4

RUN mkdir -p /var/lib/nginx

RUN apk add --no-cache redis

RUN apk add --no-cache git
RUN git clone https://github.com/wg/wrk \
	&& cd wrk \
	&& make \
	&& cp ./wrk /usr/local/bin/wrk
