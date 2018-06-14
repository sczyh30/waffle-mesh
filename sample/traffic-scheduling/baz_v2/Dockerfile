FROM nginx:1.13.12
MAINTAINER sczyh30 <sczyh16@gmail.com>

RUN apt-get update && \
    apt-get install --no-install-recommends -y \
      curl \
      iptables \
      iproute2 \
      iputils-ping \
      dnsutils \
      tcpdump \
      net-tools \
      sudo &&  apt-get upgrade -y && \
    rm -rf /var/lib/apt/lists/*

COPY main /usr/local/bin/baz-v2
RUN chmod +x /usr/local/bin/baz-v2

EXPOSE 5763

ENTRYPOINT ["/usr/local/bin/baz-v2"]