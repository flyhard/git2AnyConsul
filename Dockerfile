FROM debian:jessie

RUN apt-get update
RUN apt-get install -y git

COPY bin/gitProperties2Consul /

#VOLUME /data

CMD /gitProperties2Consul -dataDir /data -host $CONSUL_PORT_8500_TCP_ADDR -port $CONSUL_PORT_8500_TCP_PORT