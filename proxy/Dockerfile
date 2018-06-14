FROM ubuntu:16.04
MAINTAINER sczyh30 <sczyh16@gmail.com>

ADD main /usr/local/bin/waffle-proxy-server

RUN mkdir /usr/local/bin/cert
ADD cert/cert.pem /usr/local/bin/cert/cert.pem
ADD cert/key.pem /usr/local/bin/cert/key.pem

RUN chgrp 2186 /usr/local/bin/waffle-proxy-server && \
 chmod 2755 /usr/local/bin/waffle-proxy-server

WORKDIR /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/waffle-proxy-server"]