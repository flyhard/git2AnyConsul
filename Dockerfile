FROM debian:jessie

RUN apt-get update
RUN apt-get install -y git

COPY bin/gitProperties2Consul /
COPY run.sh /
RUN chmod +x run.sh

#VOLUME /data

CMD /run.sh