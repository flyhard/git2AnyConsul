FROM debian:jessie

RUN apt-get update
RUN apt-get install -y git

COPY bin/gitProperties2Consul /
COPY run.sh /
RUN chmod +x run.sh

## Uncomment next lines if you want to use SSH to access
## a repo. The key in id_rsa needs to be on the list of authorized
## keys on the server
#COPY id_rsa /root/.ssh/id_rsa
#RUN chmod 700 /root/.ssh


#VOLUME /data

CMD /run.sh