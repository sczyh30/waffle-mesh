FROM ubuntu:16.04

RUN apt-get update && apt-get install -y \
    iproute2 \
    iptables \
 && rm -rf /var/lib/apt/lists/*

ADD iptables.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/iptables.sh
ENTRYPOINT ["/usr/local/bin/iptables.sh"]